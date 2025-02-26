package client

import (
	"encoding/binary"
	"fmt"
	"io"
	"math"
	"net"
	"strings"
	"sync"
	"time"

	"gorelay/pkg/account"
	"gorelay/pkg/config"
	"gorelay/pkg/crypto"
	"gorelay/pkg/events"
	"gorelay/pkg/logger"
	"gorelay/pkg/models"
	"gorelay/pkg/packets"
	"gorelay/pkg/packets/client"
	"gorelay/pkg/packets/dataobjects"
	"gorelay/pkg/packets/interfaces"
	"gorelay/pkg/packets/server"
)

// Client represents a connected RotMG client
type Client struct {
	// Connection info
	conn      net.Conn
	connected bool
	server    *models.Server
	mu        sync.Mutex
	rc4       *crypto.RC4Manager

	// Game state
	state       *GameState
	accountInfo *account.Account
	config      *config.Config

	// Packet handling
	packetHandler      *packets.PacketHandler
	versionMgr         *packets.VersionManager
	handlersRegistered bool

	// Game tracking
	enemies     map[int32]*Enemy
	players     map[int32]*Player
	projectiles map[int32]*Projectile
	currentMap  *Map

	// Event handling
	events *events.EventEmitter

	// Connection management
	reconnectAttempts    int
	maxReconnectAttempts int
	reconnectDelay       time.Duration
	readTimeout          time.Duration
	writeTimeout         time.Duration

	// Logging
	logger *logger.Logger

	// Movement management
	nextPositions []*WorldPosData
	moveSpeed     float32
	lastMoveTime  time.Time
}

// NewClient creates a new RotMG client instance
func NewClient(acc *account.Account, cfg *config.Config, log *logger.Logger) *Client {
	// First verify the account if needed
	if acc.NeedAccountVerify() {
		log.Info("Client", "Verifying account %s...", acc.Alias)
		if err := acc.VerifyAccount(cfg.HWIDToken); err != nil {
			log.Error("Client", "Failed to verify account %s: %v", acc.Alias, err)
			return nil
		}
	}

	// Then fetch character list if needed
	if acc.NeedCharList() {
		log.Info("Client", "Fetching character list for %s...", acc.Alias)
		if err := acc.GetCharList(); err != nil {
			log.Error("Client", "Failed to get character list for %s: %v", acc.Alias, err)
			return nil
		}
	}

	// Fetch server list using account credentials
	servers, err := models.FetchServers(acc.Email, acc.Password)
	if err != nil {
		log.Warning("Client", "Failed to fetch servers: %v. Using default server.", err)
		// Use the default server instead of trying to fetch the list
		server := models.DefaultServer
		return createClient(acc, cfg, log, server)
	}

	// Get server from account preference or pick first available
	var server *models.Server
	if pref := acc.ServerPref; pref != "" {
		if s, ok := servers[pref]; ok {
			server = s
		} else {
			// If preferred server not found, pick first available
			foundServer := false
			for _, s := range servers {
				server = s
				foundServer = true
				break
			}

			if foundServer {
				log.Warning("Client", "Preferred server %s not found. Using %s instead.", pref, server.Name)
			} else {
				// No servers available in the map, use default
				server = models.DefaultServer
				log.Warning("Client", "Preferred server %s not found and no servers available. Using default server %s.",
					pref, server.Name)
			}
		}
	} else {
		// If no preference, pick first available
		foundServer := false
		for _, s := range servers {
			server = s
			foundServer = true
			break
		}

		if !foundServer {
			// No servers available in the map, use default
			server = models.DefaultServer
			log.Warning("Client", "No servers available. Using default server %s.", server.Name)
		}
	}

	client := createClient(acc, cfg, log, server)
	if client == nil {
		return nil
	}

	client.state.BuildVer = cfg.BuildVersion
	log.Info("Client", "Successfully created client for %s on server %s", acc.Alias, server.Name)
	return client
}

// createClient creates a new client instance with the given server
func createClient(acc *account.Account, cfg *config.Config, log *logger.Logger, server *models.Server) *Client {
	client := &Client{
		accountInfo:   acc,
		config:        cfg,
		logger:        log,
		server:        server,
		packetHandler: packets.NewPacketHandler(),
		versionMgr:    packets.NewVersionManager(),

		// Initialize game state
		state: &GameState{
			WorldPos:      &WorldPosData{X: 0, Y: 0},
			PlayerData:    &PlayerData{},
			LastUpdate:    time.Now(),
			LastFrameTime: time.Now().UnixNano() / int64(time.Millisecond),
		},
		enemies:     make(map[int32]*Enemy),
		players:     make(map[int32]*Player),
		projectiles: make(map[int32]*Projectile),
		events:      events.NewEventEmitter(),

		// Initialize movement management
		nextPositions: make([]*WorldPosData, 0),
		moveSpeed:     4.0, // Default speed in tiles per second
		lastMoveTime:  time.Now(),

		// Initialize connection management
		maxReconnectAttempts: 3,
		reconnectDelay:       time.Duration(cfg.ReconnectDelay) * time.Millisecond,
		readTimeout:          30 * time.Second,
		writeTimeout:         10 * time.Second,
	}

	// Register packet handlers
	client.registerPacketHandlers()
	client.handlersRegistered = true

	return client
}

// emit dispatches an event to all subscribed handlers
func (c *Client) emit(eventType events.EventType, packet interface{}, data interface{}) {
	// Create an event with the packet
	event := &events.Event{
		Type:   eventType,
		Client: c,
		Data:   data,
	}

	// Handle different packet types
	switch p := packet.(type) {
	case *server.Goto:
		// Create a wrapper for Goto that adapts the Write method
		wrapper := &packetWrapper{
			id: int32(interfaces.Goto),
			writeFunc: func(w *packets.PacketWriter) error {
				return p.Write(w)
			},
			hasNulls:   func() bool { return false },
			packetType: interfaces.Goto,
		}
		event.Packet = wrapper
	case *client.PlayerShoot:
		// Create a wrapper for PlayerShoot that adapts the Write method
		wrapper := &packetWrapper{
			id: int32(interfaces.PlayerShoot),
			writeFunc: func(w *packets.PacketWriter) error {
				return p.Write(w)
			},
			hasNulls:   func() bool { return false },
			packetType: interfaces.PlayerShoot,
		}
		event.Packet = wrapper
	case packets.Packet:
		// Use the Packet interface directly if implemented
		event.Packet = p
	default:
		// For other types, just set the packet to nil
		event.Packet = nil
	}

	// Emit the event
	c.events.Emit(event)
}

// Connect establishes a connection to the game server
func (c *Client) Connect() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.connected {
		return fmt.Errorf("client already connected")
	}

	// Fetch Unity build version from init endpoint
	buildVersion, err := c.fetchUnityBuildVersion()
	if err != nil {
		return fmt.Errorf("failed to fetch Unity build version: %v", err)
	}
	c.state.BuildVer = buildVersion
	c.config.BuildVersion = buildVersion

	// Ensure required components are initialized
	if c.packetHandler == nil {
		c.packetHandler = packets.NewPacketHandler()
	}
	if c.logger == nil {
		return fmt.Errorf("logger not initialized")
	}
	if c.state == nil {
		c.state = &GameState{
			WorldPos:      &WorldPosData{X: 0, Y: 0},
			PlayerData:    &PlayerData{},
			LastUpdate:    time.Now(),
			LastFrameTime: time.Now().UnixNano() / int64(time.Millisecond),
		}
	}

	// Initialize RC4 encryption
	inKey := []byte("c91d9eec420160730d825604e0")
	outKey := []byte("5a4d2016bc16dc64883194ffd9")
	rc4Manager, err := crypto.NewRC4Manager(inKey, outKey)
	if err != nil {
		return fmt.Errorf("failed to initialize RC4: %v", err)
	}
	c.rc4 = rc4Manager

	var lastErr error
	for attempt := 0; attempt <= c.maxReconnectAttempts; attempt++ {
		if attempt > 0 {
			c.logger.Info("Client", "Reconnection attempt %d/%d in %v...",
				attempt, c.maxReconnectAttempts, c.reconnectDelay)
			time.Sleep(c.reconnectDelay)
		}

		addr := fmt.Sprintf("%s:%d", c.server.Address, 2050)
		conn, err := net.Dial("tcp", addr)
		if err != nil {
			lastErr = fmt.Errorf("failed to connect: %v", err)
			continue
		}

		// Set connection timeouts
		tcpConn := conn.(*net.TCPConn)
		tcpConn.SetKeepAlive(true)
		tcpConn.SetKeepAlivePeriod(60 * time.Second)
		tcpConn.SetReadBuffer(8192)
		tcpConn.SetWriteBuffer(8192)

		c.conn = conn
		c.connected = true
		c.reconnectAttempts = 0

		// Register packet handlers if not already done
		if !c.handlersRegistered {
			c.registerPacketHandlers()
			c.handlersRegistered = true
		}

		// Create and send Hello packet
		// We'll need to implement a proper send method that works with the pclient.Hello type
		// For now, we'll create a simple struct that matches what we need
		type HelloPacket struct {
			BuildVersion  string
			GameNet       string
			PlayPlatform  string
			PlatformToken string
			AccessToken   string
			KeyTime       int32
			Key           []byte
			GameID        int32
			ClientToken   string
		}

		// Create the hello packet data
		helloData := HelloPacket{
			BuildVersion:  c.state.BuildVer,
			GameNet:       "Unity",
			PlayPlatform:  "Unity",
			PlatformToken: "",
			AccessToken:   c.accountInfo.AccessToken,
			KeyTime:       int32(time.Now().Unix()),
			Key:           []byte{}, // Empty for now
			GameID:        -2,
			ClientToken:   c.config.HWIDToken,
		}

		// TODO: Implement proper packet encoding and sending
		// This is a placeholder until we implement the proper packet handling
		c.logger.Info("Client", "Sending Hello packet with build version: %s", helloData.BuildVersion)

		// For now, we'll just create a simple byte array to send
		// In the future, this should use the proper packet encoding
		helloBytes := []byte{0x00} // Packet ID 0 for Hello
		// Append other data...

		if _, err := c.conn.Write(helloBytes); err != nil {
			c.logger.Error("Client", "Failed to send Hello packet: %v", err)
			c.conn.Close()
			continue
		}

		// Start packet handling goroutine
		go c.handlePackets()

		return nil
	}

	return fmt.Errorf("failed to connect after %d attempts: %v",
		c.maxReconnectAttempts, lastErr)
}

// Disconnect closes the connection to the game server
func (c *Client) Disconnect() {
	c.mu.Lock()
	defer c.mu.Unlock()

	if !c.connected {
		return
	}

	if c.conn != nil {
		c.conn.Close()
	}
	c.connected = false
}

// registerPacketHandlers sets up handlers for different packet types
func (c *Client) registerPacketHandlers() {
	// Handle Hello packet
	c.packetHandler.RegisterHandler(0, func(data []byte) error {
		c.logger.Info("Client", "Handling Hello packet")
		c.logger.Info("Client", "Hello packet data: % x", data)

		// Log connection details
		c.logger.Info("Client", "Connection details:")
		c.logger.Info("Client", "  Build Version: %s", c.state.BuildVer)
		c.logger.Info("Client", "  Account: %s", c.accountInfo.Alias)
		c.logger.Info("Client", "  Server: %s", c.server.Name)

		return nil
	})

	// Handle AoE packets
	c.packetHandler.RegisterHandler(int(interfaces.AOE), func(data []byte) error {
		packet := &server.AOE{}
		// TODO: Implement packet decoding
		if packet.Location != nil && c.state.WorldPos != nil {
			// Convert Location to WorldPosData for distance calculation
			packetPos := &WorldPosData{X: float32(packet.Location.X), Y: float32(packet.Location.Y)}
			if packetPos.SquareDistanceTo(c.state.WorldPos) < packet.Radius*packet.Radius {
				// Apply AoE damage
				c.applyDamage(int32(packet.Damage), packet.ArmorPierce)
			}
		}
		return nil
	})

	// Handle enemy shoot packets
	c.packetHandler.RegisterHandler(int(interfaces.EnemyShoot), func(data []byte) error {
		packet := &server.EnemyShoot{}
		// TODO: Implement packet decoding
		if enemy, ok := c.enemies[packet.OwnerId]; ok && !enemy.IsDead() {
			// Convert Location to WorldPosData
			startPos := &WorldPosData{X: float32(packet.Location.X), Y: float32(packet.Location.Y)}
			for i := byte(0); i < packet.NumShots; i++ {
				angle := packet.Angle + float32(i)*packet.AngleInc
				c.addProjectile(int32(packet.BulletType), packet.OwnerId, int32(packet.BulletId)+int32(i), angle, startPos)
			}
		}
		return nil
	})

	// Handle new tick packets
	c.packetHandler.RegisterHandler(int(interfaces.NewTick), func(data []byte) error {
		packet := &server.NewTick{}
		// TODO: Implement packet decoding
		c.state.LastFrameTime = time.Now().UnixNano() / int64(time.Millisecond)

		// Process statuses
		for _, status := range packet.Statuses {
			if int32(status.ObjectID) == c.state.ObjectID {
				if status.Position != nil {
					// Convert Position to WorldPosData
					c.state.WorldPos = &WorldPosData{X: float32(status.Position.X), Y: float32(status.Position.Y)}
				}
				// Update player stats
				for _, stat := range status.Data {
					// Convert StatsType to int32 and handle string vs int values
					if stat.IsStringData() {
						c.updateStat(int32(stat.ID), 0, stat.StringValue)
					} else {
						c.updateStat(int32(stat.ID), int32(stat.IntValue), "")
					}
				}
			}
		}
		return nil
	})

	// Handle update packets
	c.packetHandler.RegisterHandler(int(interfaces.Update), func(data []byte) error {
		packet := &server.Update{}
		// TODO: Implement packet decoding

		// Process new objects
		for _, entity := range packet.NewObjs {
			// Handle entity based on its properties
			// This will need to be adjusted based on the actual structure of Entity
			c.handleNewObject(entity)
		}

		// Process dropped objects
		for _, objID := range packet.Drops {
			delete(c.enemies, objID)
			delete(c.players, objID)
		}
		return nil
	})

	// Handle text packets
	c.packetHandler.RegisterHandler(int(interfaces.Text), func(data []byte) error {
		packet := &server.Text{}
		// TODO: Implement packet decoding

		// Handle chat messages based on the Text packet
		c.handleChatMessage(packet)
		return nil
	})

	// Handle failure packets
	c.packetHandler.RegisterHandler(int(interfaces.Failure), func(data []byte) error {
		packet := &server.Failure{}
		// TODO: Implement packet decoding
		switch packet.ErrorId {
		case int32(4): // IncorrectVersion
			c.logger.Info("Client", "Build version out of date. Updating and reconnecting...")
			// Update build version in config and state
			c.config.BuildVersion = packet.ErrorMessage
			c.state.BuildVer = packet.ErrorMessage
			// Save updated config
			if err := config.SaveConfig("config.json", c.config); err != nil {
				c.logger.Error("Client", "Failed to save updated build version: %v", err)
			}
			// Reconnect with new version
			c.reconnect()
		case int32(5): // InvalidTeleportTarget
			c.logger.Warning("Client", "Invalid teleport target")
		case int32(7): // EmailVerificationNeeded
			c.logger.Error("Client", "Email verification required")
		case int32(8): // BadKey
			c.logger.Error("Client", "Invalid key used")
		case int32(11): // InvalidCharacter
			c.logger.Info("Client", "Character not found. Creating new character...")
			// TODO: Handle character creation
		default:
			c.logger.Error("Client", "Received failure %d: %s", packet.ErrorId, packet.ErrorMessage)
		}
		return nil
	})

	// Handle goto packets
	c.packetHandler.RegisterHandler(int(interfaces.Goto), func(data []byte) error {
		packet := &server.Goto{}
		// TODO: Implement packet decoding

		// Create and send acknowledgment
		ack := client.NewGotoAck()
		ack.Time = int32(c.state.LastFrameTime)
		ack.Unknown = false

		if err := c.send(ack); err != nil {
			c.logger.Error("Client", "Failed to send GotoAck: %v", err)
		}

		if packet.ObjectId == c.state.ObjectID {
			// Convert Location to WorldPosData
			pos := &WorldPosData{X: float32(packet.Location.X), Y: float32(packet.Location.Y)}
			c.state.WorldPos = pos
			c.emit(events.EventPlayerMove, packet, &events.PlayerEventData{
				PlayerData: c.state.PlayerData,
				Position:   c.state.WorldPos,
			})
		}
		return nil
	})

	// Handle player shoot
	c.packetHandler.RegisterHandler(int(interfaces.PlayerShoot), func(data []byte) error {
		packet := client.NewPlayerShoot()
		// TODO: Implement packet decoding
		c.emit(events.EventPlayerShoot, packet, nil)
		return nil
	})
}

// Helper methods

func (c *Client) applyDamage(damage int32, armorPiercing bool) {
	// TODO: Implement damage calculation and application
}

func (c *Client) addProjectile(bulletType, ownerID, bulletID int32, angle float32, startPos *WorldPosData) {
	// TODO: Implement projectile tracking
}

func (c *Client) updateStat(statType int32, statValue int32, stringValue string) {
	if c.state.PlayerData == nil {
		c.state.PlayerData = &PlayerData{
			Stats:     make(map[string]int32),
			Inventory: make([]int32, 20), // 12 inventory + 8 backpack slots
		}
	}

	switch models.StatType(statType) {
	case models.MAXHPSTAT:
		c.state.PlayerData.MaxHP = statValue
	case models.HPSTAT:
		c.state.PlayerData.HP = statValue
	case models.MAXMPSTAT:
		c.state.PlayerData.MaxMP = statValue
	case models.MPSTAT:
		c.state.PlayerData.MP = statValue
	case models.NEXTLEVELEXPSTAT:
		c.state.PlayerData.NextLevelExp = statValue
	case models.EXPSTAT:
		c.state.PlayerData.Exp = statValue
	case models.LEVELSTAT:
		c.state.PlayerData.Level = statValue
	case models.NAMESTAT:
		c.state.PlayerData.Name = stringValue
	case models.ATTACKSTAT:
		c.state.PlayerData.Stats["atk"] = statValue
	case models.DEFENSESTAT:
		c.state.PlayerData.Stats["def"] = statValue
	case models.SPEEDSTAT:
		c.state.PlayerData.Stats["spd"] = statValue
	case models.DEXTERITYSTAT:
		c.state.PlayerData.Stats["dex"] = statValue
	case models.VITALITYSTAT:
		c.state.PlayerData.Stats["vit"] = statValue
	case models.WISDOMSTAT:
		c.state.PlayerData.Stats["wis"] = statValue
	case models.FAMESTAT:
		c.state.PlayerData.Fame = statValue
	case models.CURRFAMESTAT:
		c.state.PlayerData.CurrentFame = statValue
	case models.NUMSTARSSTAT:
		c.state.PlayerData.Stars = statValue
	case models.ACCOUNTIDSTAT:
		c.state.PlayerData.AccountID = stringValue
	case models.GUILDNAMESTAT:
		c.state.PlayerData.GuildName = stringValue
	case models.GUILDRANKSTAT:
		c.state.PlayerData.GuildRank = statValue
		/* 	case models.HEALTHPOTIONSTACKSTAT:
		   		c.state.PlayerData.HPPots = statValue
		   	case models.MAGICPOTIONSTACKSTAT:
		   		c.state.PlayerData.MPPots = statValue
		   	case models.HASBACKPACKSTAT:
		   		c.state.PlayerData.HasBackpack = statValue == 1 */
	default:
		// Handle inventory slots
		statTypeEnum := models.StatType(statType)
		if statTypeEnum >= models.INVENTORY0STAT && statTypeEnum <= models.INVENTORY11STAT {
			slot := int(statTypeEnum - models.INVENTORY0STAT)
			if slot >= 0 && slot < len(c.state.PlayerData.Inventory) {
				c.state.PlayerData.Inventory[slot] = statValue
			}
		} else if statTypeEnum >= models.BACKPACK0STAT && statTypeEnum <= models.BACKPACK7STAT {
			slot := int(statTypeEnum - models.BACKPACK0STAT + 12) // Offset by 12 inventory slots
			if slot >= 0 && slot < len(c.state.PlayerData.Inventory) {
				c.state.PlayerData.Inventory[slot] = statValue
			}
		}
	}
}

func (c *Client) handleNewObject(obj interface{}) {
	// TODO: Implement object handling based on type
}

func (c *Client) handlePrivateMessage(packet interface{}) {
	// TODO: Implement private message handling
}

func (c *Client) handleChatMessage(packet interface{}) {
	// TODO: Implement chat message handling
}

// GetState returns the current game state
func (c *Client) GetState() *GameState {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.state
}

// GetEnemy returns an enemy by ID
func (c *Client) GetEnemy(id int32) *Enemy {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.enemies[id]
}

// GetPlayer returns a player by ID
func (c *Client) GetPlayer(id int32) *Player {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.players[id]
}

// GetProjectile returns a projectile by ID
func (c *Client) GetProjectile(id int32) *Projectile {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.projectiles[id]
}

// GetMap returns the current map
func (c *Client) GetMap() *Map {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.currentMap
}

// GetPosition returns the client's current position
func (c *Client) GetPosition() *WorldPosData {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.state.WorldPos
}

// SetPosition updates the client's position
func (c *Client) SetPosition(pos *WorldPosData) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.state.WorldPos = pos
}

// IsConnected returns whether the client is connected
func (c *Client) IsConnected() bool {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.connected
}

// GetLogger returns the client's logger
func (c *Client) GetLogger() *logger.Logger {
	return c.logger
}

// handlePackets processes incoming packets
func (c *Client) handlePackets() {
	defer c.Disconnect()

	for {
		// Set read deadline for each packet
		if err := c.conn.SetReadDeadline(time.Now().Add(c.readTimeout)); err != nil {
			c.logger.Error("Client", "Failed to set read deadline: %v", err)
			return
		}

		// First read the packet length (4 bytes) and packet ID (1 byte)
		header := make([]byte, 5)
		if _, err := io.ReadFull(c.conn, header); err != nil {
			if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
				c.logger.Warning("Client", "Read timeout, attempting to reconnect...")
				c.reconnect()
				return
			}

			if err == io.EOF || strings.Contains(err.Error(), "connection reset by peer") ||
				strings.Contains(err.Error(), "forcibly closed") {
				c.logger.Warning("Client", "Connection closed by server, attempting to reconnect...")
				c.reconnect()
				return
			}

			c.logger.Error("Client", "Error reading packet header: %v", err)
			return
		}

		// Extract packet length and ID
		packetLength := int(binary.BigEndian.Uint32(header[:4]))
		packetID := header[4]

		// Validate packet length
		//todo: dynamically adjust buffer size if too big
		if packetLength <= 5 || packetLength > 100000 {
			c.logger.Error("Client", "Invalid packet length: %d", packetLength)
			continue
		}

		// Read the rest of the packet
		payloadLength := packetLength - 5
		payload := make([]byte, payloadLength)
		if _, err := io.ReadFull(c.conn, payload); err != nil {
			c.logger.Error("Client", "Error reading packet payload: %v", err)
			return
		}

		// Make a copy of the payload for decryption
		decryptedPayload := make([]byte, payloadLength)
		copy(decryptedPayload, payload)

		// Decrypt only the payload if RC4 is initialized
		if c.rc4 != nil {
			c.rc4.Decrypt(decryptedPayload)
		} else {
			c.logger.Warning("Client", "RC4 not initialized, processing raw data")
		}

		// Log packet details
		c.logger.Debug("Client", "Received packet - ID: %d, Length: %d bytes", packetID, packetLength)
		if payloadLength > 0 {
			maxBytes := 16
			if payloadLength < maxBytes {
				maxBytes = payloadLength
			}
			c.logger.Debug("Client", "Encrypted payload (first %d bytes): % x", maxBytes, payload[:maxBytes])
			c.logger.Debug("Client", "Decrypted payload (first %d bytes): % x", maxBytes, decryptedPayload[:maxBytes])
		}

		// Special handling for Hello packet (ID 0)
		if packetID == 0 {
			c.logger.Info("Client", "Received Hello packet response")
			c.logger.Debug("Client", "Full Hello packet payload: % x", decryptedPayload)

			if len(decryptedPayload) > 0 {
				c.logger.Info("Client", "Hello packet details:")
				c.logger.Info("Client", "  Build Version: %s", c.state.BuildVer)
				c.logger.Info("Client", "  Account: %s", c.accountInfo.Alias)
				c.logger.Info("Client", "  Connected to: %s", c.server.Name)
			}
		}

		// If we have a version manager, try to get the packet name
		if c.versionMgr != nil {
			if packetName, err := c.versionMgr.GetPacketName(int(packetID)); err == nil {
				c.logger.Debug("Client", "Packet type: %s", packetName)
			}
		}

		// Process the decrypted packet
		if err := c.packetHandler.HandlePacket(int(packetID), decryptedPayload); err != nil {
			c.logger.Error("Client", "Error handling packet: %v", err)
			// Don't return on packet handling errors, continue processing other packets
		}
	}
}

func (c *Client) reconnect() {
	c.mu.Lock()
	defer c.mu.Unlock()

	// If already disconnected, no need to proceed
	if !c.connected {
		return
	}

	// Log reconnection attempt
	c.logger.Info("Client", "Initiating reconnection sequence...")

	// Properly close existing connection
	if c.conn != nil {
		c.conn.Close()
		c.conn = nil
	}
	c.connected = false

	// Reset game state
	c.state = &GameState{
		BuildVer: c.state.BuildVer, // Preserve build version
	}
	c.enemies = make(map[int32]*Enemy)
	c.players = make(map[int32]*Player)
	c.projectiles = make(map[int32]*Projectile)

	// Check if we should attempt reconnection
	if c.reconnectAttempts >= c.maxReconnectAttempts {
		c.logger.Error("Client", "Max reconnection attempts (%d) reached", c.maxReconnectAttempts)
		return
	}

	c.reconnectAttempts++
	attemptNum := c.reconnectAttempts

	// Start reconnection attempt in a goroutine
	go func() {
		// Wait for the configured delay
		c.logger.Info("Client", "Waiting %v before reconnection attempt %d/%d...",
			c.reconnectDelay, attemptNum, c.maxReconnectAttempts)
		time.Sleep(c.reconnectDelay)

		// Check if this was a manual reconnection request
		if c.accountInfo != nil && c.accountInfo.Reconnect {
			c.accountInfo.Reconnect = false // Reset the flag
			c.reconnectAttempts = 0         // Reset attempts for manual reconnection
		}

		// Attempt to reconnect
		if err := c.Connect(); err != nil {
			c.logger.Error("Client", "Reconnection attempt %d failed: %v", attemptNum, err)
		} else {
			c.logger.Info("Client", "Successfully reconnected on attempt %d", attemptNum)
			c.reconnectAttempts = 0 // Reset counter on successful connection
		}
	}()
}

// send sends a packet to the server
func (c *Client) send(packet interface{}) error {
	if !c.connected {
		return fmt.Errorf("not connected")
	}

	// Set write deadline for sending packet
	if err := c.conn.SetWriteDeadline(time.Now().Add(c.writeTimeout)); err != nil {
		return fmt.Errorf("failed to set write deadline: %v", err)
	}

	var data []byte
	var err error

	// Handle different packet types
	switch p := packet.(type) {
	case *client.GotoAck:
		// Create a wrapper for GotoAck
		wrapper := &packetWrapper{
			id:         int32(interfaces.GotoAck),
			writeFunc:  p.Write,
			hasNulls:   func() bool { return false },
			packetType: interfaces.GotoAck,
		}
		data, err = packets.EncodePacket(wrapper)
	case *client.PlayerShoot:
		// Create a wrapper for PlayerShoot
		wrapper := &packetWrapper{
			id:         int32(interfaces.PlayerShoot),
			writeFunc:  p.Write,
			hasNulls:   func() bool { return false },
			packetType: interfaces.PlayerShoot,
		}
		data, err = packets.EncodePacket(wrapper)
	case *client.Move:
		// Create a wrapper for Move
		wrapper := &packetWrapper{
			id:         int32(interfaces.Move),
			writeFunc:  p.Write,
			hasNulls:   func() bool { return false },
			packetType: interfaces.Move,
		}
		data, err = packets.EncodePacket(wrapper)
	case packets.Packet:
		// Use the Packet interface directly if implemented
		data, err = packets.EncodePacket(p)
	default:
		return fmt.Errorf("unsupported packet type: %T", packet)
	}

	if err != nil {
		return fmt.Errorf("failed to encode packet: %v", err)
	}

	// Log outgoing packet for debugging
	c.logger.Debug("Client", "Sending packet, data: % x", data)

	// Encrypt if RC4 is initialized
	if c.rc4 != nil {
		c.rc4.Encrypt(data)
		c.logger.Debug("Client", "Encrypted data: % x", data)
	}

	// Send the packet
	_, err = c.conn.Write(data)
	return err
}

// SwitchServer changes the client's server and attempts to connect to it
func (c *Client) SwitchServer(serverName string) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	server := models.GetServer(serverName)
	if server == nil {
		return fmt.Errorf("unknown server: %s", serverName)
	}

	// Update server info
	c.server = server

	// Disconnect from current server if connected
	if c.connected {
		c.Disconnect()
	}

	// Connect to new server
	return c.Connect()
}

// GetCurrentServer returns the current server configuration
func (c *Client) GetCurrentServer() *models.Server {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.server
}

// moveTo updates the client's position for smooth movement
func (c *Client) moveTo(target *WorldPosData) bool {
	if target == nil {
		return false
	}

	now := time.Now()
	elapsed := now.Sub(c.lastMoveTime).Seconds()
	c.lastMoveTime = now

	step := float32(elapsed) * c.moveSpeed

	// Calculate distance to target
	dx := target.X - c.state.WorldPos.X
	dy := target.Y - c.state.WorldPos.Y
	distSq := dx*dx + dy*dy

	// If we can reach target in this step
	if distSq <= step*step {
		c.state.WorldPos.X = target.X
		c.state.WorldPos.Y = target.Y
		if len(c.nextPositions) > 0 {
			c.nextPositions = c.nextPositions[1:]
		}
		return true
	}

	// Move towards target
	angle := float32(math.Atan2(float64(dy), float64(dx)))
	c.state.WorldPos.X += float32(math.Cos(float64(angle))) * step
	c.state.WorldPos.Y += float32(math.Sin(float64(angle))) * step

	return true
}

// AddPath adds a path of positions to move through
func (c *Client) AddPath(path []*WorldPosData) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.nextPositions = append(c.nextPositions, path...)
}

// ClearPath clears the movement queue
func (c *Client) ClearPath() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.nextPositions = c.nextPositions[:0]
}

// HasNextPosition returns whether there are more positions to move to
func (c *Client) HasNextPosition() bool {
	c.mu.Lock()
	defer c.mu.Unlock()
	return len(c.nextPositions) > 0
}

// GetNextPosition returns the next target position, or nil if none
func (c *Client) GetNextPosition() *WorldPosData {
	c.mu.Lock()
	defer c.mu.Unlock()
	if len(c.nextPositions) == 0 {
		return nil
	}
	return c.nextPositions[0]
}

// Update handles client updates including movement
func (c *Client) Update() {
	c.mu.Lock()
	defer c.mu.Unlock()

	if len(c.nextPositions) > 0 {
		if moved := c.moveTo(c.nextPositions[0]); moved {
			// Create a new Move packet
			movePacket := client.NewMove()
			movePacket.TickID = int32(time.Now().UnixNano() / int64(time.Millisecond))
			movePacket.Time = int32(c.state.LastFrameTime)

			// Create a location record with the current position
			record := dataobjects.NewLocationRecord()
			record.Time = int32(c.state.LastFrameTime)
			record.Position = dataobjects.NewLocationWithCoords(float64(c.state.WorldPos.X), float64(c.state.WorldPos.Y))

			// Add the record to the packet
			movePacket.Records = append(movePacket.Records, record)

			// Encode and send the packet
			if c.connected {
				// Create a wrapper that implements the packets.Packet interface
				wrapper := &packetWrapper{
					id:         int32(interfaces.Move),
					writeFunc:  movePacket.Write,
					hasNulls:   func() bool { return false },
					packetType: interfaces.Move,
				}

				// Encode the packet
				data, err := packets.EncodePacket(wrapper)
				if err != nil {
					c.logger.Error("Client", "Failed to encode Move packet: %v", err)
					return
				}

				// Log outgoing packet for debugging
				c.logger.Debug("Client", "Sending Move packet, data: % x", data)

				// Encrypt if RC4 is initialized
				if c.rc4 != nil {
					c.rc4.Encrypt(data)
					c.logger.Debug("Client", "Encrypted data: % x", data)
				}

				// Send the packet
				if _, err := c.conn.Write(data); err != nil {
					c.logger.Error("Client", "Failed to send Move packet: %v", err)
				}
			}
		}
	}
}

// packetWrapper is a helper struct to adapt client/server packets to the packets.Packet interface
type packetWrapper struct {
	id         int32
	writeFunc  func(w *packets.PacketWriter) error
	hasNulls   func() bool
	packetType interfaces.PacketType
}

func (p *packetWrapper) ID() int32 {
	return p.id
}

func (p *packetWrapper) Write(w interfaces.Writer) error {
	if pw, ok := w.(*packets.PacketWriter); ok {
		return p.writeFunc(pw)
	}
	return fmt.Errorf("expected *packets.PacketWriter, got %T", w)
}

func (p *packetWrapper) Read(r interfaces.Reader) error {
	return fmt.Errorf("read not implemented for wrapper")
}

func (p *packetWrapper) String() string {
	return fmt.Sprintf("PacketWrapper(ID=%d, Type=%v)", p.id, p.packetType)
}

func (p *packetWrapper) HasNulls() bool {
	return p.hasNulls()
}

func (p *packetWrapper) Structure() string {
	return fmt.Sprintf("PacketWrapper(ID=%d, Type=%v)", p.id, p.packetType)
}

func (p *packetWrapper) Type() interfaces.PacketType {
	return p.packetType
}

// Add new function to fetch Unity build version
func (c *Client) fetchUnityBuildVersion() (string, error) {
	// Hardcoded build version as requested
	return "5.8.0.0.0", nil

	// Original implementation commented out below
	/*
		baseURL := "https://www.realmofthemadgod.com/app/init"

		// Create HTTP client with appropriate headers
		client := &http.Client{}

		// Prepare form data
		data := make(url.Values)
		data.Set("platform", "standalonewindows64")
		data.Set("key", "9KnJFxtTvLu2frXv")
		data.Set("game_net", "Unity")
		data.Set("play_platform", "Unity")
		data.Set("game_net_user_id", "")

		// Create request
		req, err := http.NewRequest("POST", baseURL, strings.NewReader(data.Encode()))
		if err != nil {
			return "", fmt.Errorf("failed to create request: %v", err)
		}

		// Set Unity-specific headers
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		req.Header.Set("X-Unity-Version", "2021.3.16f1")
		req.Header.Set("User-Agent", "UnityPlayer/2021.3.16f1 (UnityWebRequest/1.0, libcurl/7.84.0-DEV)")

		// Send request
		resp, err := client.Do(req)
		if err != nil {
			return "", fmt.Errorf("failed to send request: %v", err)
		}
		defer resp.Body.Close()

		// Check response status
		if resp.StatusCode != http.StatusOK {
			body, _ := io.ReadAll(resp.Body)
			return "", fmt.Errorf("server returned status %d: %s", resp.StatusCode, string(body))
		}

		// Read response
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return "", fmt.Errorf("failed to read response: %v", err)
		}

		// Log response for debugging
		c.logger.Debug("Client", "Init response: %s", string(body))

		// Parse XML response
		doc := &struct {
			XMLName      xml.Name `xml:"AppSettings"`
			BuildHash    string   `xml:"BuildHash"`
			BuildCDN     string   `xml:"BuildCDN"`
			BuildVersion string   `xml:"BuildVersion"`
		}{}

		if err := xml.Unmarshal(body, doc); err != nil {
			return "", fmt.Errorf("failed to parse XML response: %v", err)
		}

		// Return BuildVersion if available, otherwise use BuildHash
		if doc.BuildVersion != "" {
			return doc.BuildVersion, nil
		}
		if doc.BuildHash != "" {
			return doc.BuildHash, nil
		}

		return "", fmt.Errorf("neither BuildVersion nor BuildHash found in response")
	*/
}
