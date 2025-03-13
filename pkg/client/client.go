package client

import (
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"math"
	"net"
	"reflect"
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
	if true || acc.NeedAccountVerify() {
		log.Info("Client", "Verifying account %s (token %s)", acc.Alias, acc.HwidToken)
		if err := acc.VerifyAccount(acc.HwidToken); err != nil {
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
	if c.rc4 == nil {
		inKey, err := hex.DecodeString("c91d9eec420160730d825604e0")
		if err != nil {
			return fmt.Errorf("failed to decode inKey: %v", err)
		}
		outKey, err := hex.DecodeString("5a4d2016bc16dc64883194ffd9")
		if err != nil {
			return fmt.Errorf("failed to decode outKey: %v", err)
		}
		rc4Manager, err := crypto.NewRC4Manager(inKey, outKey)
		if err != nil {
			return fmt.Errorf("failed to initialize RC4: %v", err)
		}
		c.rc4 = rc4Manager
	} else {
		// Reset RC4 ciphers for new connection
		if err := c.rc4.Reset(); err != nil {
			return fmt.Errorf("failed to reset RC4: %v", err)
		}
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

		// Register packet handlers if not already done
		if !c.handlersRegistered {
			c.registerPacketHandlers()
			c.handlersRegistered = true
		}

		// Create and send Hello packet
		hello := client.NewHello()
		hello.GameID = -2
		hello.BuildVersion = c.config.BuildVersion
		hello.AccessToken = c.accountInfo.AccessToken
		hello.KeyTime = -1
		hello.Key = []byte{}
		hello.GameNet = "rotmg"
		hello.PlayPlatform = "rotmg"
		hello.PlatformToken = ""
		hello.ClientToken = c.accountInfo.HwidToken
		hello.ClientIdentification = "XQpu8CWkMehb5rLVP3DG47FcafExRUvg"

		c.logger.Info("Client", "Sending Hello")
		//c.logger.Info("Client", "Sending Hello packet: %s", hello.ToString())

		if err := c.sendPacket(hello); err != nil {
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

// sendPacket sends a packet to the server with proper RC4 encryption and header construction
func (c *Client) sendPacket(p packets.Packet) error {
	// Create packet writer and write packet contents
	writer := packets.NewPacketWriter()
	p.Write(writer)

	// Create header with packet size and ID
	header := packets.NewPacketWriter()
	header.WriteInt32(int32(5 + len(writer.Bytes())))
	header.WriteByte(byte(p.ID()))
	header.WriteBytes(writer.Bytes())

	// Get final encoded bytes
	data := header.Bytes()

	// Make a copy for encryption
	encryptedData := make([]byte, len(data))
	copy(encryptedData, data)

	// Encrypt payload (skip header)
	c.rc4.Encrypt(encryptedData)

	// Send to server
	if _, err := c.conn.Write(encryptedData); err != nil {
		c.logger.Error("Client", "Failed to send packet: %v", err)
		return err
	}

	return nil
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

var packetTypes = map[interfaces.PacketType]packets.Packet{
	interfaces.AccountList:                    &server.AccountList{},
	interfaces.ActivePet:                      &server.ActivePet{},
	interfaces.AllyShoot:                      &server.AllyShoot{},
	interfaces.AOE:                            &server.AOE{},
	interfaces.BoostBPMilestoneResult:         &server.BoostBPMilestoneResult{},
	interfaces.BuyItemResult:                  &server.BuyItemResult{},
	interfaces.BuyResult:                      &server.BuyResult{},
	interfaces.ClaimBPMilestoneResult:         &server.ClaimBPMilestoneResult{},
	interfaces.ClaimMissionResult:             &server.ClaimMissionResult{},
	interfaces.CreateSuccess:                  &server.CreateSuccess{},
	interfaces.CrucibleResult:                 &server.CrucibleResult{},
	interfaces.Damage:                         &server.Damage{},
	interfaces.Death:                          &server.Death{},
	interfaces.DeletePet:                      &server.DeletePet{},
	interfaces.DrawDebugArrow:                 &server.DrawDebugArrow{},
	interfaces.DrawDebugShape:                 &server.DrawDebugShape{},
	interfaces.EnemyShoot:                     &server.EnemyShoot{},
	interfaces.EvolvedPet:                     &server.EvolvedPet{},
	interfaces.ExaltationBonusChanged:         &server.ExaltationBonusChanged{},
	interfaces.Failure:                        &server.Failure{},
	interfaces.File:                           &server.File{},
	interfaces.ForgeResult:                    &server.ForgeResult{},
	interfaces.ForgeUnlockedBlueprints:        &server.ForgeUnlockedBlueprints{},
	interfaces.Goto:                           &server.Goto{},
	interfaces.GuildResult:                    &server.GuildResult{},
	interfaces.HatchPet:                       &server.HatchPet{},
	interfaces.HeroLeft:                       &server.HeroLeft{},
	interfaces.IncomingPartyInvite:            &server.IncomingPartyInvite{},
	interfaces.IncomingPartyMemberInfo:        &server.IncomingPartyMemberInfo{},
	interfaces.InventoryResult:                &server.InventoryResult{},
	interfaces.InvitedToGuild:                 &server.InvitedToGuild{},
	interfaces.KeyInfoResponse:                &server.KeyInfoResponse{},
	interfaces.MapInfo:                        &server.MapInfo{},
	interfaces.MissionProgressUpdate:          &server.MissionProgressUpdate{},
	interfaces.MultipleMissionsProgressUpdate: &server.MultipleMissionsProgressUpdate{},
	interfaces.NameResult:                     &server.NameResult{},
	interfaces.NewAbility:                     &server.NewAbility{},
	interfaces.NewCharacterInformation:        &server.NewCharacterInformation{},
	interfaces.NewTick:                        &server.NewTick{},
	interfaces.Notification:                   &server.Notification{},
	interfaces.PartyAction:                    &server.PartyAction{},
	interfaces.PartyJoinRequestResponse:       &server.PartyJoinRequestResponse{},
	interfaces.PartyJoinResponse:              &server.PartyJoinResponse{},
	interfaces.PartyList:                      &server.PartyList{},
	interfaces.PartyMemberAdded:               &server.PartyMemberAdded{},
	interfaces.PasswordPrompt:                 &server.PasswordPrompt{},
	interfaces.PetYardUpdate:                  &server.PetYardUpdate{},
	interfaces.Pic:                            &server.Pic{},
	interfaces.Ping:                           &server.Ping{},
	interfaces.PlayersList:                    &server.PlayersList{},
	interfaces.PlaySound:                      &server.PlaySound{},
	interfaces.QuestFetchResponse:             &server.QuestFetchResponse{},
	interfaces.QuestObjectId:                  &server.QuestObjectId{},
	interfaces.QuestRedeemResponse:            &server.QuestRedeemResponse{},
	interfaces.Queue:                          &server.Queue{},
	interfaces.RealmScoreUpdate:               &server.RealmScoreUpdate{},
	interfaces.Reconnect:                      &server.Reconnect{},
	interfaces.RefineResult:                   &server.RefineResult{},
	interfaces.ResetDailyQuests:               &server.ResetDailyQuests{},
	interfaces.ServerPlayerShoot:              &server.ServerPlayerShoot{},
	interfaces.ShowEffect:                     &server.ShowEffect{},
	interfaces.SkinRecycleResponse:            &server.SkinRecycleResponse{},
	interfaces.Text:                           &server.Text{},
	interfaces.TradeAccepted:                  &server.TradeAccepted{},
	interfaces.TradeChanged:                   &server.TradeChanged{},
	interfaces.TradeDone:                      &server.TradeDone{},
	interfaces.TradeRequested:                 &server.TradeRequested{},
	interfaces.TradeStart:                     &server.TradeStart{},
	interfaces.UnlockCustomization:            &server.UnlockCustomization{},
	interfaces.UnlockNewSlot:                  &server.UnlockNewSlot{},
	interfaces.Update:                         &server.Update{},
	interfaces.VaultContent:                   &server.VaultContent{},
}

// registerPacketHandlers sets up handlers for different packet types
func (c *Client) registerPacketHandlers() {
	c.packetHandler.RegisterHandler(int(interfaces.MapInfo), func(packet packets.Packet) error {
		mapInfo := packet.(*server.MapInfo)
		c.logger.Info("Client", "MapInfo: %v", mapInfo)

		// First check if we have a character ID in the account config
		if c.accountInfo != nil && c.accountInfo.CharInfo != nil && c.accountInfo.CharInfo.CharID > 0 {
			c.logger.Info("Client", "Loading character %d from config", c.accountInfo.CharInfo.CharID)
			load := &client.Load{
				CharacterID: c.accountInfo.CharInfo.CharID,
			}
			if err := c.send(load); err != nil {
				c.logger.Error("Client", "Failed to send Load packet: %v", err)
			}
			return nil
		}

		// If no character ID in config, check character list
		needsNewChar := true
		if c.accountInfo != nil && c.accountInfo.Chars != nil && len(c.accountInfo.Chars.Characters) > 0 {
			needsNewChar = false
			// Update character ID in config
			if c.accountInfo.CharInfo == nil {
				c.accountInfo.CharInfo = &account.CharInfo{}
			}
			c.accountInfo.CharInfo.CharID = int32(c.accountInfo.Chars.Characters[0].ID)
		}

		if needsNewChar {
			c.logger.Info("Client", "Creating new character")
			create := &client.Create{
				ClassType:    768, //wizard
				SkinType:     0,
				IsChallenger: false,
				IsSeasonal:   false,
			}
			if err := c.send(create); err != nil {
				c.logger.Error("Client", "Failed to send Create packet: %v", err)
			}
		} else {
			load := &client.Load{
				CharacterID: int32(c.accountInfo.Chars.Characters[0].ID),
			}
			if err := c.send(load); err != nil {
				c.logger.Error("Client", "Failed to send Load packet: %v", err)
			}
			c.logger.Info("Client", "Loading character %d", load.CharacterID)
		}

		return nil
	})

	// Handle AoE packets
	c.packetHandler.RegisterHandler(int(interfaces.AOE), func(packet packets.Packet) error {
		aoe := packet.(*server.AOE)
		if aoe.Location != nil && c.state.WorldPos != nil {
			// Convert Location to WorldPosData for distance calculation
			packetPos := &WorldPosData{X: float32(aoe.Location.X), Y: float32(aoe.Location.Y)}
			if packetPos.SquareDistanceTo(c.state.WorldPos) < aoe.Radius*aoe.Radius {
				// Apply AoE damage
				c.applyDamage(int32(aoe.Damage), aoe.ArmorPierce)
			}
		}
		return nil
	})

	// Handle enemy shoot packets
	c.packetHandler.RegisterHandler(int(interfaces.EnemyShoot), func(packet packets.Packet) error {
		enemyShoot := packet.(*server.EnemyShoot)

		// Send ShootAck as keep-alive response
		shootAck := &client.ShootAckCounter{
			Time:   int32(time.Now().UnixNano() / int64(time.Millisecond)),
			Amount: 1,
		}

		if err := c.send(shootAck); err != nil {
			c.logger.Error("Client", "Failed to send ShootAck: %v", err)
		}

		if enemy, ok := c.enemies[enemyShoot.OwnerId]; ok && !enemy.IsDead() {
			startPos := &WorldPosData{X: float32(enemyShoot.Location.X), Y: float32(enemyShoot.Location.Y)}
			for i := byte(0); i < enemyShoot.NumShots; i++ {
				angle := enemyShoot.Angle + float32(i)*enemyShoot.AngleInc
				c.addProjectile(int32(enemyShoot.BulletType), enemyShoot.OwnerId, int32(enemyShoot.BulletId)+int32(i), angle, startPos)
			}
		}
		return nil
	})

	// Handle ping packets
	c.packetHandler.RegisterHandler(int(interfaces.Ping), func(packet packets.Packet) error {
		ping := packet.(*server.Ping)

		// Create and send pong response with proper time calculation
		currentTime := time.Now().UnixMilli()
		pong := &client.Pong{
			Serial: ping.Serial,
			Time:   int32(currentTime % (1 << 31)), // Ensure time fits in int32
		}

		if err := c.send(pong); err != nil {
			c.logger.Error("Client", "Failed to send Pong: %v", err)
		}
		return nil
	})

	// Handle update packets
	c.packetHandler.RegisterHandler(int(interfaces.Update), func(packet packets.Packet) error {
		update := packet.(*server.Update)

		// Update player position if provided and non-zero
		if update.PlayerPosition != nil && (update.PlayerPosition.X != 0 || update.PlayerPosition.Y != 0) {
			c.state.WorldPos = &WorldPosData{
				X: float32(update.PlayerPosition.X),
				Y: float32(update.PlayerPosition.Y),
			}
			c.logger.Debug("Client", "Updated position to X=%f, Y=%f", c.state.WorldPos.X, c.state.WorldPos.Y)
		}

		// Send UpdateAck as keep-alive response
		updateAck := &client.UpdateAck{}
		if err := c.send(updateAck); err != nil {
			c.logger.Error("Client", "Failed to send UpdateAck: %v", err)
		}

		// Process new objects
		for _, entity := range update.NewObjs {
			c.handleNewObject(entity)
		}

		// Process dropped objects
		for _, objID := range update.Drops {
			delete(c.enemies, objID)
			delete(c.players, objID)
		}
		return nil
	})

	// Handle new tick packets
	c.packetHandler.RegisterHandler(int(interfaces.NewTick), func(packet packets.Packet) error {
		newTick := packet.(*server.NewTick)

		// Update last frame time from server's tick time
		c.state.LastFrameTime = int64(newTick.ServerRealTimeMs)

		// Create and send move packet with correct timing
		movePacket := client.NewMove()
		movePacket.TickID = newTick.TickId // Use server's tick ID
		movePacket.Time = int32(newTick.ServerRealTimeMs)

		// Add current position
		record := dataobjects.NewLocationRecord()
		record.Time = int32(newTick.ServerRealTimeMs)

		// Only send position if we have a valid one
		if c.state.WorldPos != nil && (c.state.WorldPos.X != 0 || c.state.WorldPos.Y != 0) {
			record.Position = dataobjects.NewLocationWithCoords(float64(c.state.WorldPos.X), float64(c.state.WorldPos.Y))
			c.logger.Debug("Client", "Sending Move with position X=%f, Y=%f", c.state.WorldPos.X, c.state.WorldPos.Y)
		} else {
			// If we don't have a valid position, don't send the move packet
			c.logger.Debug("Client", "Skipping Move packet - no valid position")
			return nil
		}

		movePacket.Records = append(movePacket.Records, record)

		if err := c.send(movePacket); err != nil {
			c.logger.Error("Client", "Failed to send Move response to NewTick: %v", err)
		}

		// Process statuses
		for _, status := range newTick.Statuses {
			if int32(status.ObjectID) == c.state.ObjectID {
				if status.Position != nil && (status.Position.X != 0 || status.Position.Y != 0) {
					c.state.WorldPos = &WorldPosData{X: float32(status.Position.X), Y: float32(status.Position.Y)}
					c.logger.Debug("Client", "Updated position from status to X=%f, Y=%f", c.state.WorldPos.X, c.state.WorldPos.Y)
				}
				for _, stat := range status.Data {
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

	// Handle text packets
	c.packetHandler.RegisterHandler(int(interfaces.Text), func(packet packets.Packet) error {
		text := packet.(*server.Text)

		// Determine message type and format appropriately
		switch {
		case text.Recipient != "": // Private message
			if text.Name == "" {
				// System tell (from game)
				c.logger.Info("Client", "From: %s> %s", text.Name, text.RawText)
			} else {
				// Player tell
				c.logger.Info("Client", "From %s: %s", text.Name, text.RawText)
			}
		case strings.HasPrefix(text.Name, "#"): // Oryx/Admin message
			c.logger.Info("Client", "[Announcement] %s: %s", strings.TrimPrefix(text.Name, "#"), text.RawText)
		case strings.HasPrefix(text.Name, "*"): // Guild message
			c.logger.Info("Client", "[Guild] %s: %s", strings.TrimPrefix(text.Name, "*"), text.RawText)
		case strings.HasPrefix(text.Name, "@"): // Party message
			c.logger.Info("Client", "[Party] %s: %s", strings.TrimPrefix(text.Name, "@"), text.RawText)
		case text.Name == "": // Pure server message
			c.logger.Info("Client", "[Server] %s", text.RawText)
		default: // Normal chat
			c.logger.Info("Client", "<%s> %s", text.Name, text.RawText)
		}

		return nil
	})

	// Handle notification packets
	c.packetHandler.RegisterHandler(int(interfaces.Notification), func(packet packets.Packet) error {
		return nil
	})

	// Handle ClientStat packets
	c.packetHandler.RegisterHandler(int(interfaces.ClientStat), func(packet packets.Packet) error {
		stat := packet.(*server.ClientStat)
		c.logger.Debug("Client", "Received client stat: %s = %d", stat.Name, stat.Value)
		return nil
	})

	// Handle ServerPlayerShoot packets
	c.packetHandler.RegisterHandler(int(interfaces.ServerPlayerShoot), func(packet packets.Packet) error {
		shoot := packet.(*server.ServerPlayerShoot)
		c.logger.Debug("Client", "Server player shoot: BulletId=%d, OwnerId=%d, ContainerType=%d, Pos=(%f,%f), Angle=%f, Damage=%d",
			shoot.BulletId, shoot.OwnerId, shoot.ContainerType,
			shoot.StartingPos.X, shoot.StartingPos.Y,
			shoot.Angle, shoot.Damage)
		return nil
	})

	// Handle ShowEffect packets (ID 11)
	c.packetHandler.RegisterHandler(int(interfaces.ShowEffect), func(packet packets.Packet) error {
		effect := packet.(*server.ShowEffect)

		// Log effect details at debug level
		c.logger.Debug("Client", "ShowEffect: Type=%d, Value=%d, TargetId=%d, Pos=(%f,%f)",
			effect.EffectType, effect.EffectValue, effect.TargetId,
			effect.PosA.X, effect.PosA.Y)

		// Handle different effect types
		switch effect.EffectType {
		case 1: // Heal
			if int32(effect.TargetId) == c.state.ObjectID {
				c.logger.Debug("Client", "Received heal effect")
			}
		case 2: // Teleport
			if int32(effect.TargetId) == c.state.ObjectID {
				c.logger.Debug("Client", "Received teleport effect to (%f,%f)",
					effect.PosA.X, effect.PosA.Y)
			}
		case 3: // Stream
			c.logger.Debug("Client", "Received stream effect")
		case 4: // Throw
			c.logger.Debug("Client", "Received throw effect")
		case 5: // Nova
			c.logger.Debug("Client", "Received nova effect")
		case 6: // Poison
			c.logger.Debug("Client", "Received poison effect")
		case 7: // Line
			c.logger.Debug("Client", "Received line effect")
		case 8: // Burst
			c.logger.Debug("Client", "Received burst effect")
		case 9: // Flow
			c.logger.Debug("Client", "Received flow effect")
		case 10: // Ring
			c.logger.Debug("Client", "Received ring effect")
		case 11: // Lightning
			c.logger.Debug("Client", "Received lightning effect")
		case 12: // Collapse
			c.logger.Debug("Client", "Received collapse effect")
		case 13: // Coneblast
			c.logger.Debug("Client", "Received coneblast effect")
		default:
			c.logger.Debug("Client", "Received unknown effect type: %d", effect.EffectType)
		}

		return nil
	})

	// Handle ReskinUnlock packets (ID 114)
	c.packetHandler.RegisterHandler(114, func(packet packets.Packet) error {
		c.logger.Debug("Client", "Received ReskinUnlock packet")
		// Just acknowledge the packet since we don't need to process it
		return nil
	})

	// Handle unknown packet type 120 (possibly server status)
	c.packetHandler.RegisterHandler(120, func(packet packets.Packet) error {
		c.logger.Debug("Client", "Received server status packet (type 120)")
		return nil
	})

	// Handle failure packets with improved logging and keep-alive detection
	c.packetHandler.RegisterHandler(int(interfaces.Failure), func(packet packets.Packet) error {
		failure := packet.(*server.Failure)

		// Handle keep-alive packets (empty failure packets)
		if failure.ErrorId == 0 && failure.ErrorMessage == "" {
			// Just log at debug level and continue - no need to send response
			c.logger.Debug("Client", "Received keep-alive packet")
			return nil
		}

		switch failure.ErrorId {
		case int32(4): // IncorrectVersion
			c.logger.Info("Client", "Build version out of date. Updating and reconnecting...")
			// Update build version in config and state
			c.config.BuildVersion = failure.ErrorMessage
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
		default:
			c.logger.Error("Client", "Received failure %d: %s", failure.ErrorId, failure.ErrorMessage)
		}
		return nil
	})

	// Handle CreateSuccess packets
	c.packetHandler.RegisterHandler(int(interfaces.CreateSuccess), func(packet packets.Packet) error {
		createSuccess := packet.(*server.CreateSuccess)
		c.logger.Info("Client", "Character loaded successfully - ObjectId: %d, CharId: %d",
			createSuccess.ObjectId, createSuccess.CharId)

		// Update our state with the character info
		c.state.ObjectID = createSuccess.ObjectId
		if c.accountInfo != nil && c.accountInfo.CharInfo != nil {
			c.accountInfo.CharInfo.CharID = createSuccess.CharId
		}

		return nil
	})

	// Handle goto packets
	c.packetHandler.RegisterHandler(int(interfaces.Goto), func(packet packets.Packet) error {
		gotoPacket := &server.Goto{}
		// TODO: Implement packet decoding

		// Create and send acknowledgment
		gotoAck := client.NewGotoAck()
		gotoAck.Time = int32(c.state.LastFrameTime)
		gotoAck.Unknown = false

		if err := c.send(gotoAck); err != nil {
			c.logger.Error("Client", "Failed to send GotoAck: %v", err)
		}

		if gotoPacket.ObjectId == c.state.ObjectID {
			// Convert Location to WorldPosData
			pos := &WorldPosData{X: float32(gotoPacket.Location.X), Y: float32(gotoPacket.Location.Y)}
			c.state.WorldPos = pos
			c.emit(events.EventPlayerMove, gotoPacket, &events.PlayerEventData{
				PlayerData: c.state.PlayerData,
				Position:   c.state.WorldPos,
			})
		}
		return nil
	})

	// Handle MultipleMissionsProgressUpdate packets
	c.packetHandler.RegisterHandler(int(interfaces.MultipleMissionsProgressUpdate), func(packet packets.Packet) error {
		missionUpdate := packet.(*server.MultipleMissionsProgressUpdate)
		c.logger.Debug("Client", "Received mission progress update: %s", missionUpdate.UnknownString)
		return nil
	})

	// Handle Trade packets
	c.packetHandler.RegisterHandler(int(interfaces.TradeRequested), func(packet packets.Packet) error {
		trade := packet.(*server.TradeRequested)
		c.logger.Debug("Client", "Received trade request from: %s", trade.Name)

		// For now, we'll just log the trade request
		// You can implement trade acceptance/rejection logic here if needed
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
		if !c.connected {
			return
		}

		// Set read deadline for each packet
		if err := c.conn.SetReadDeadline(time.Now().Add(c.readTimeout)); err != nil {
			c.logger.Warning("Client", "Failed to set read deadline: %v", err)
			continue
		}

		// Read packet header (not encrypted)
		header := make([]byte, 5)
		bytesRead := 0
		for bytesRead < 5 {
			n, err := c.conn.Read(header[bytesRead:])
			if err != nil {
				if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
					continue
				}
				c.logger.Error("Client", "Failed to read packet header: %v", err)
				return
			}
			bytesRead += n
		}

		// Convert 4 bytes to int32 length as big endian
		packetLength := int32(binary.BigEndian.Uint32(header[0:4]))
		if packetLength < 5 || packetLength > 16384 { // Add reasonable size limit
			c.logger.Warning("Client", "Invalid packet length: %d", packetLength)
			continue
		}

		packetId := header[4]

		// Read packet data in chunks
		packetData := make([]byte, 0, packetLength-5)
		remaining := packetLength - 5

		for remaining > 0 {
			// Read in chunks of up to 8192 bytes
			chunkSize := remaining
			if chunkSize > 8192 {
				chunkSize = 8192
			}

			chunk := make([]byte, chunkSize)
			n, err := c.conn.Read(chunk)
			if err != nil {
				if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
					continue
				}
				c.logger.Error("Client", "Failed to read packet data: %v", err)
				return
			}

			// Decrypt the chunk
			c.rc4.Decrypt(chunk[:n])

			// Append decrypted chunk
			packetData = append(packetData, chunk[:n]...)
			remaining -= int32(n)
		}

		packet, ok := packetTypes[interfaces.PacketType(packetId)]
		if !ok {
			c.logger.Warning("Client", "Unknown packet type: %d", packetId)
			continue
		}

		// Create a new instance of the packet type
		newPacket := reflect.New(reflect.TypeOf(packet).Elem()).Interface().(packets.Packet)

		reader := packets.NewPacketReader(packetData)
		if err := newPacket.Read(reader); err != nil {
			c.logger.Warning("Client", "Failed to read packet: %v", err)
			continue
		}

		// Log received packet
		c.logger.Debug("Client", "RECV [%s] Type: %d, Length: %d, Data: %+v",
			interfaces.PacketType(packetId), packetId, packetLength, newPacket)

		// Process the decrypted packet
		if err := c.packetHandler.HandlePacket(int(packetId), newPacket); err != nil {
			c.logger.Warning("Client", "Error handling packet: %v", err)
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
	c.state = &GameState{}
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
	var packetType interfaces.PacketType

	// Handle different packet types
	switch p := packet.(type) {
	case *client.Load:
		// Create packet header
		writer := packets.NewPacketWriter()
		// Write packet size (4 bytes) - will update after writing data
		writer.WriteInt32(0)
		// Write packet type (1 byte)
		writer.WriteByte(byte(interfaces.Load))
		// Write packet data
		if err := p.Write(writer); err != nil {
			return fmt.Errorf("failed to write Load packet: %v", err)
		}
		// Update packet size
		data = writer.Bytes()
		binary.BigEndian.PutUint32(data[0:4], uint32(len(data)))
		packetType = interfaces.Load

	case *client.Create:
		writer := packets.NewPacketWriter()
		writer.WriteInt32(0)
		writer.WriteByte(byte(interfaces.Create))
		if err := p.Write(writer); err != nil {
			return fmt.Errorf("failed to write Create packet: %v", err)
		}
		data = writer.Bytes()
		binary.BigEndian.PutUint32(data[0:4], uint32(len(data)))
		packetType = interfaces.Create

	case *client.Pong:
		writer := packets.NewPacketWriter()
		writer.WriteInt32(0)
		writer.WriteByte(byte(interfaces.Pong))
		if err := p.Write(writer); err != nil {
			return fmt.Errorf("failed to write Pong packet: %v", err)
		}
		data = writer.Bytes()
		binary.BigEndian.PutUint32(data[0:4], uint32(len(data)))
		packetType = interfaces.Pong

	case *client.UpdateAck:
		writer := packets.NewPacketWriter()
		writer.WriteInt32(0)
		writer.WriteByte(byte(interfaces.UpdateAck))
		if err := p.Write(writer); err != nil {
			return fmt.Errorf("failed to write UpdateAck packet: %v", err)
		}
		data = writer.Bytes()
		binary.BigEndian.PutUint32(data[0:4], uint32(len(data)))
		packetType = interfaces.UpdateAck

	case *client.ShootAckCounter:
		writer := packets.NewPacketWriter()
		writer.WriteInt32(0)
		writer.WriteByte(byte(interfaces.ShootAckCounter))
		if err := p.Write(writer); err != nil {
			return fmt.Errorf("failed to write ShootAckCounter packet: %v", err)
		}
		data = writer.Bytes()
		binary.BigEndian.PutUint32(data[0:4], uint32(len(data)))
		packetType = interfaces.ShootAckCounter

	case *client.Move:
		writer := packets.NewPacketWriter()
		writer.WriteInt32(0)
		writer.WriteByte(byte(interfaces.Move))
		if err := p.Write(writer); err != nil {
			return fmt.Errorf("failed to write Move packet: %v", err)
		}
		data = writer.Bytes()
		binary.BigEndian.PutUint32(data[0:4], uint32(len(data)))
		packetType = interfaces.Move

	case *client.GotoAck:
		writer := packets.NewPacketWriter()
		writer.WriteInt32(0)
		writer.WriteByte(byte(interfaces.GotoAck))
		if err := p.Write(writer); err != nil {
			return fmt.Errorf("failed to write GotoAck packet: %v", err)
		}
		data = writer.Bytes()
		binary.BigEndian.PutUint32(data[0:4], uint32(len(data)))
		packetType = interfaces.GotoAck

	case packets.Packet:
		writer := packets.NewPacketWriter()
		writer.WriteInt32(0)
		writer.WriteByte(byte(p.Type()))
		if err := p.Write(writer); err != nil {
			return fmt.Errorf("failed to write packet: %v", err)
		}
		data = writer.Bytes()
		binary.BigEndian.PutUint32(data[0:4], uint32(len(data)))
		packetType = p.Type()

	default:
		return fmt.Errorf("unsupported packet type: %T", packet)
	}

	if err != nil {
		return fmt.Errorf("failed to encode packet: %v", err)
	}

	// Log outgoing packet before encryption
	c.logger.Debug("Client", "SEND [%d] Type: %d, Length: %d, Data: %#v",
		packetType, int(packetType), len(data), packet)

	// Make a copy of the data for encryption
	encryptedData := make([]byte, len(data))
	copy(encryptedData, data)

	// Encrypt if RC4 is initialized
	if c.rc4 != nil {
		c.rc4.Encrypt(encryptedData)
	}

	// Send the encrypted packet
	_, err = c.conn.Write(encryptedData)
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
