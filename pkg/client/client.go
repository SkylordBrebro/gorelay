package client

import (
	"fmt"
	"net"
	"sync"
	"time"

	"gorelay/pkg/account"
	"gorelay/pkg/crypto"
	"gorelay/pkg/events"
	"gorelay/pkg/logger"
	"gorelay/pkg/packets"
	"gorelay/pkg/resources"
)

// Client represents a connected RotMG client
type Client struct {
	// Connection info
	conn       net.Conn
	connected  bool
	serverAddr string
	serverPort int
	mu         sync.Mutex
	rc4        *crypto.RC4Manager

	// Game state
	objectID    int32
	worldPos    *packets.WorldPosData
	playerData  *packets.PlayerData
	accountInfo *account.Account
	buildVer    string
	gameID      int32
	nextPos     []*packets.WorldPosData

	// Resources
	resources *resources.ResourceManager
	logger    *logger.Logger

	// Packet handling
	packetHandler *packets.PacketHandler
	versionMgr    *packets.VersionManager

	// State tracking
	lastFrameTime int64
	projectiles   []*Projectile
	enemies       map[int32]*Enemy
	players       map[int32]*Player

	// Event handling
	events *events.EventEmitter
}

// NewClient creates a new RotMG client instance
func NewClient(acc *account.Account, res *resources.ResourceManager, log *logger.Logger) *Client {
	client := &Client{
		accountInfo:   acc,
		resources:     res,
		logger:        log,
		serverAddr:    "server.realmofthemadgod.com",
		serverPort:    2050,
		packetHandler: packets.NewPacketHandler(),
		versionMgr:    packets.NewVersionManager(),
		enemies:       make(map[int32]*Enemy),
		players:       make(map[int32]*Player),
		nextPos:       make([]*packets.WorldPosData, 0),
		events:        events.NewEventEmitter(),
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

	addr := fmt.Sprintf("%s:%d", c.serverAddr, c.serverPort)
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return fmt.Errorf("failed to connect: %v", err)
	}

	c.conn = conn
	c.connected = true

	// Start packet handling goroutine
	go c.handlePackets()

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

// registerPacketHandlers sets up handlers for different packet types
func (c *Client) registerPacketHandlers() {
	// Handle AoE packets
	c.packetHandler.RegisterHandler(int((&packets.AoePacket{}).ID()), func(data []byte) error {
		packet := &packets.AoePacket{}
		// TODO: Implement packet decoding
		if packet.Pos.SquareDistanceTo(c.worldPos) < packet.Radius*packet.Radius {
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
		c.lastFrameTime = time.Now().UnixNano() / int64(time.Millisecond)

		// Process updates
		for _, update := range packet.Updates {
			if update.ObjectID == c.objectID {
				c.worldPos = update.Pos
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
			if obj.Status.ObjectID == c.objectID {
				c.worldPos = obj.Status.Pos
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
		if packet.Recipient == c.playerData.Name {
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
			c.buildVer = packet.ErrorDescription
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
			Time: int32(c.lastFrameTime),
		}
		c.send(ack)

		if packet.ObjectID == c.objectID {
			c.worldPos = packet.Position
			c.emit(events.EventPlayerMove, packet, &events.PlayerEventData{
				PlayerData: c.playerData,
				Position:   c.worldPos,
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

// GetPosition returns the client's current position
func (c *Client) GetPosition() *packets.WorldPosData {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.worldPos
}

// SetPosition updates the client's position
func (c *Client) SetPosition(pos *packets.WorldPosData) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.worldPos = pos
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
		n, err := c.conn.Read(buffer)
		if err != nil {
			c.logger.Error("Client", "Error reading packet: %v", err)
			return
		}

		// Process the packet
		if err := c.packetHandler.HandlePacket(int(buffer[0]), buffer[1:n]); err != nil {
			c.logger.Error("Client", "Error handling packet: %v", err)
		}
	}
}

// Add send method
func (c *Client) send(packet packets.Packet) error {
	if !c.connected {
		return fmt.Errorf("not connected")
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
