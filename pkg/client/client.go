package client

import (
	"fmt"
	"io"
	"net"
	"strings"
	"sync"
	"time"

	"gorelay/pkg/account"
	"gorelay/pkg/crypto"
	"gorelay/pkg/events"
	"gorelay/pkg/logger"
	"gorelay/pkg/models"
	"gorelay/pkg/packets"
	"gorelay/pkg/resources"
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

	// Resources
	resources *resources.ResourceManager
	logger    *logger.Logger

	// Packet handling
	packetHandler *packets.PacketHandler
	versionMgr    *packets.VersionManager

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
}

// NewClient creates a new RotMG client instance
func NewClient(acc *account.Account, res *resources.ResourceManager, log *logger.Logger) *Client {
	// Fetch server list using account credentials
	servers, err := models.FetchServers(acc.Email, acc.Password)
	if err != nil {
		log.Warning("Client", "Failed to fetch servers: %v. Using default server.", err)
		// Use the default server instead of trying to fetch the list
		server := models.DefaultServer
		return createClient(acc, res, log, server)
	}

	// Get server from account preference or pick first available
	var server *models.Server
	if pref := acc.ServerPref; pref != "" {
		if s, ok := servers[pref]; ok {
			server = s
		} else {
			// If preferred server not found, pick first available
			for _, s := range servers {
				server = s
				break
			}
		}
	} else {
		// If no preference, pick first available
		for _, s := range servers {
			server = s
			break
		}
	}

	// If somehow we still don't have a server, use default
	if server == nil {
		server = models.DefaultServer
	}

	return createClient(acc, res, log, server)
}

// createClient creates a new client instance with the given server
func createClient(acc *account.Account, res *resources.ResourceManager, log *logger.Logger, server *models.Server) *Client {
	client := &Client{
		accountInfo:   acc,
		resources:     res,
		logger:        log,
		server:        server,
		packetHandler: packets.NewPacketHandler(),
		versionMgr:    packets.NewVersionManager(),

		// Initialize game state
		state:       &GameState{},
		enemies:     make(map[int32]*Enemy),
		players:     make(map[int32]*Player),
		projectiles: make(map[int32]*Projectile),
		events:      events.NewEventEmitter(),

		// Initialize connection management
		maxReconnectAttempts: 3,
		reconnectDelay:       5 * time.Second,
		readTimeout:          30 * time.Second,
		writeTimeout:         10 * time.Second,
	}

	// Register packet handlers
	client.registerPacketHandlers()

	return client
}

// emit dispatches an event to all subscribed handlers
func (c *Client) emit(eventType events.EventType, packet packets.Packet, data interface{}) {
	c.events.Emit(&events.Event{
		Type:   eventType,
		Client: c,
		Packet: packet,
		Data:   data,
	})
}

// Connect establishes a connection to the game server
func (c *Client) Connect() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.connected {
		return fmt.Errorf("client already connected")
	}

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
	// Handle AoE packets
	c.packetHandler.RegisterHandler(int((&packets.AoePacket{}).ID()), func(data []byte) error {
		packet := &packets.AoePacket{}
		// TODO: Implement packet decoding
		if packet.Pos.SquareDistanceTo(c.state.WorldPos) < packet.Radius*packet.Radius {
			// Apply AoE damage
			c.applyDamage(packet.Damage, packet.ArmorPiercing)
		}
		return nil
	})

	// Handle enemy shoot packets
	c.packetHandler.RegisterHandler(int((&packets.EnemyShootPacket{}).ID()), func(data []byte) error {
		packet := &packets.EnemyShootPacket{}
		// TODO: Implement packet decoding
		if enemy, ok := c.enemies[packet.OwnerID]; ok && !enemy.IsDead() {
			for i := int32(0); i < packet.NumShots; i++ {
				angle := packet.Angle + float32(i)*packet.AngleInc
				c.addProjectile(packet.BulletType, packet.OwnerID, packet.BulletID+i, angle, packet.StartingPos)
			}
		}
		return nil
	})

	// Handle new tick packets
	c.packetHandler.RegisterHandler(int((&packets.NewTickPacket{}).ID()), func(data []byte) error {
		packet := &packets.NewTickPacket{}
		// TODO: Implement packet decoding
		c.state.LastFrameTime = time.Now().UnixNano() / int64(time.Millisecond)

		// Process updates
		for _, update := range packet.Updates {
			if update.ObjectID == c.state.ObjectID {
				c.state.WorldPos = update.Pos
				// Update player stats
				for _, stat := range update.Stats {
					c.updateStat(stat)
				}
			}
		}
		return nil
	})

	// Handle update packets
	c.packetHandler.RegisterHandler(int((&packets.UpdatePacket{}).ID()), func(data []byte) error {
		packet := &packets.UpdatePacket{}
		// TODO: Implement packet decoding

		// Process new objects
		for _, obj := range packet.NewObjects {
			if obj.Status.ObjectID == c.state.ObjectID {
				c.state.WorldPos = obj.Status.Pos
				// Update player stats
				for _, stat := range obj.Status.Stats {
					c.updateStat(stat)
				}
				continue
			}

			// Handle other objects based on type
			c.handleNewObject(obj)
		}

		// Process dropped objects
		for _, objID := range packet.Drops {
			delete(c.enemies, objID)
			delete(c.players, objID)
		}
		return nil
	})

	// Handle text packets
	c.packetHandler.RegisterHandler(int((&packets.TextPacket{}).ID()), func(data []byte) error {
		packet := &packets.TextPacket{}
		// TODO: Implement packet decoding

		// Handle chat messages
		if packet.Recipient == c.state.PlayerData.Name {
			c.handlePrivateMessage(packet)
		} else {
			c.handleChatMessage(packet)
		}
		return nil
	})

	// Handle failure packets
	c.packetHandler.RegisterHandler(int((&packets.FailurePacket{}).ID()), func(data []byte) error {
		packet := &packets.FailurePacket{}
		// TODO: Implement packet decoding
		switch packet.ErrorID {
		case packets.IncorrectVersion:
			c.logger.Info("Client", "Build version out of date. Updating and reconnecting...")
			c.state.BuildVer = packet.ErrorDescription
			// TODO: Update build version
		case packets.InvalidTeleportTarget:
			c.logger.Warning("Client", "Invalid teleport target")
		case packets.EmailVerificationNeeded:
			c.logger.Error("Client", "Email verification required")
		case packets.BadKey:
			c.logger.Error("Client", "Invalid key used")
			// Reset key info
		case packets.InvalidCharacter:
			c.logger.Info("Client", "Character not found. Creating new character...")
			// TODO: Handle character creation
		default:
			c.logger.Error("Client", "Received failure %d: %s", packet.ErrorID, packet.ErrorDescription)
		}
		return nil
	})

	// Handle goto packets
	c.packetHandler.RegisterHandler(int((&packets.GotoPacket{}).ID()), func(data []byte) error {
		packet := &packets.GotoPacket{}
		// TODO: Implement packet decoding

		// Send acknowledgment
		ack := &packets.GotoAckPacket{
			Time: int32(c.state.LastFrameTime),
		}
		c.send(ack)

		if packet.ObjectID == c.state.ObjectID {
			c.state.WorldPos = packet.Position
			c.emit(events.EventPlayerMove, packet, &events.PlayerEventData{
				PlayerData: c.state.PlayerData,
				Position:   c.state.WorldPos,
			})
		}
		return nil
	})

	// Handle player shoot
	c.packetHandler.RegisterHandler(int((&packets.PlayerShootPacket{}).ID()), func(data []byte) error {
		packet := &packets.PlayerShootPacket{}
		// TODO: Implement packet decoding

		c.emit(events.EventPlayerShoot, packet, nil)
		return nil
	})
}

// Helper methods

func (c *Client) applyDamage(damage int32, armorPiercing bool) {
	// TODO: Implement damage calculation and application
}

func (c *Client) addProjectile(bulletType, ownerID, bulletID int32, angle float32, startPos *packets.WorldPosData) {
	// TODO: Implement projectile tracking
}

func (c *Client) updateStat(stat packets.StatData) {
	// TODO: Implement stat updates
}

func (c *Client) handleNewObject(obj packets.ObjectData) {
	// TODO: Implement object handling based on type
}

func (c *Client) handlePrivateMessage(packet *packets.TextPacket) {
	// TODO: Implement private message handling
}

func (c *Client) handleChatMessage(packet *packets.TextPacket) {
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
func (c *Client) GetPosition() *packets.WorldPosData {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.state.WorldPos
}

// SetPosition updates the client's position
func (c *Client) SetPosition(pos *packets.WorldPosData) {
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

	buffer := make([]byte, 8192)
	for {
		// Set read deadline for each packet
		if err := c.conn.SetReadDeadline(time.Now().Add(c.readTimeout)); err != nil {
			c.logger.Error("Client", "Failed to set read deadline: %v", err)
			return
		}

		n, err := c.conn.Read(buffer)
		if err != nil {
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

			c.logger.Error("Client", "Error reading packet: %v", err)
			return
		}

		// Process the packet
		if err := c.packetHandler.HandlePacket(int(buffer[0]), buffer[1:n]); err != nil {
			c.logger.Error("Client", "Error handling packet: %v", err)
		}
	}
}

func (c *Client) reconnect() {
	c.mu.Lock()
	defer c.mu.Unlock()

	if !c.connected {
		return
	}

	c.conn.Close()
	c.connected = false

	if c.reconnectAttempts < c.maxReconnectAttempts {
		c.reconnectAttempts++
		go func() {
			if err := c.Connect(); err != nil {
				c.logger.Error("Client", "Failed to reconnect: %v", err)
			}
		}()
	}
}

// Add send method
func (c *Client) send(packet packets.Packet) error {
	if !c.connected {
		return fmt.Errorf("not connected")
	}

	// Set write deadline for sending packet
	if err := c.conn.SetWriteDeadline(time.Now().Add(c.writeTimeout)); err != nil {
		return fmt.Errorf("failed to set write deadline: %v", err)
	}

	data, err := packets.Encode(packet)
	if err != nil {
		return fmt.Errorf("failed to encode packet: %v", err)
	}

	if c.rc4 != nil {
		c.rc4.Encrypt(data)
	}

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
