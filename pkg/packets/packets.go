package packets

// Packet is the interface that all packets must implement
type Packet interface {
	ID() int32
}

// WorldPosData represents a position in the game world
type WorldPosData struct {
	X float32
	Y float32
}

// SquareDistanceTo returns the square of the distance to another position
func (w *WorldPosData) SquareDistanceTo(other *WorldPosData) float32 {
	dx := w.X - other.X
	dy := w.Y - other.Y
	return dx*dx + dy*dy
}

// PlayerData contains player information
type PlayerData struct {
	Name      string
	ObjectID  int32
	Pos       *WorldPosData
	HP        int32
	MaxHP     int32
	MP        int32
	MaxMP     int32
	Level     int32
	Exp       int32
	Fame      int32
	Stats     map[string]int32
	Inventory []int32
}

// AoePacket represents an area of effect attack
type AoePacket struct {
	Pos           *WorldPosData
	Radius        float32
	Damage        int32
	Effect        int32
	Duration      float32
	OrigType      int32
	Color         int32
	ArmorPiercing bool
}

func (p *AoePacket) ID() int32 { return 21 }

// EnemyShootPacket represents an enemy shooting projectiles
type EnemyShootPacket struct {
	BulletID    int32
	OwnerID     int32
	BulletType  int32
	StartingPos *WorldPosData
	Angle       float32
	Damage      int32
	NumShots    int32
	AngleInc    float32
}

func (p *EnemyShootPacket) ID() int32 { return 42 }

// ServerPlayerShootPacket represents another player shooting
type ServerPlayerShootPacket struct {
	BulletID      int32
	OwnerID       int32
	ContainerType int32
	StartingPos   *WorldPosData
	Angle         float32
	Damage        int32
}

func (p *ServerPlayerShootPacket) ID() int32 { return 43 }

// NewTickPacket contains game tick information
type NewTickPacket struct {
	TickID                  int32
	TickTime                int32
	ServerRealTimeMS        int32
	ServerLastRealTimeRTTMS int32
	Updates                 []UpdateData
}

func (p *NewTickPacket) ID() int32 { return 44 }

// UpdatePacket contains entity updates
type UpdatePacket struct {
	Tiles      []GroundTileData
	NewObjects []ObjectData
	Drops      []int32
}

func (p *UpdatePacket) ID() int32 { return 45 }

// TextPacket represents a chat message
type TextPacket struct {
	Name       string
	ObjectID   int32
	NumStars   int32
	BubbleTime int32
	Recipient  string
	Text       string
	CleanText  string
}

func (p *TextPacket) ID() int32 { return 46 }

// StatData represents a single stat value
type StatData struct {
	StatType    int32
	StatValue   int32
	StringValue string
}

// UpdateData contains entity update information
type UpdateData struct {
	ObjectID int32
	Pos      *WorldPosData
	Stats    []StatData
}

// GroundTileData represents map tile information
type GroundTileData struct {
	X    int32
	Y    int32
	Type int32
}

// ObjectData contains information about game objects
type ObjectData struct {
	ObjectType int32
	Status     ObjectStatusData
}

// ObjectStatusData contains object status information
type ObjectStatusData struct {
	ObjectID int32
	Pos      *WorldPosData
	Stats    []StatData
}

// MapInfoPacket contains information about the current map
type MapInfoPacket struct {
	Width               int32
	Height              int32
	Name                string
	Seed                int32
	Background          int32
	Difficulty          int32
	AllowPlayerTeleport bool
	ShowDisplays        bool
	ClientXML           []byte
	ExtraXML            []byte
}

func (p *MapInfoPacket) ID() int32 { return 83 }

// FailureCode represents different types of connection failures
type FailureCode int32

const (
	IncorrectVersion FailureCode = iota
	BadKey
	InvalidTeleportTarget
	EmailVerificationNeeded
	InvalidCharacter
)

// FailurePacket indicates a connection or game error
type FailurePacket struct {
	ErrorID          FailureCode
	ErrorDescription string
}

func (p *FailurePacket) ID() int32 { return 0 }

// UpdateAckPacket acknowledges an UpdatePacket
type UpdateAckPacket struct{}

func (p *UpdateAckPacket) ID() int32 { return 47 }

// GotoAckPacket acknowledges a GotoPacket
type GotoAckPacket struct {
	Time int32
}

func (p *GotoAckPacket) ID() int32 { return 48 }

// PlayerShootPacket is sent when the player shoots
type PlayerShootPacket struct {
	Time          int32
	BulletID      int32
	ContainerType int32
	StartingPos   *WorldPosData
	Angle         float32
}

func (p *PlayerShootPacket) ID() int32 { return 49 }

// MovePacket relays player position to server
type MovePacket struct {
	TickID      int32
	Time        int32
	NewPosition *WorldPosData
	Records     []MoveRecord
}

func (p *MovePacket) ID() int32 { return 50 }

// MoveRecord represents a movement record
type MoveRecord struct {
	Time int32
	X    float32
	Y    float32
}

// UseItemPacket is sent to use an inventory item
type UseItemPacket struct {
	Time       int32
	SlotObject *SlotObjectData
	ItemUsePos *WorldPosData
	UseType    int32
}

func (p *UseItemPacket) ID() int32 { return 51 }

// SlotObjectData represents an inventory slot
type SlotObjectData struct {
	ObjectID   int32
	SlotID     int32
	ObjectType int32
}

// GotoPacket indicates an object moving to a new position
type GotoPacket struct {
	ObjectID int32
	Position *WorldPosData
}

func (p *GotoPacket) ID() int32 { return 52 }

// Encode serializes a packet for network transmission
func Encode(packet Packet) ([]byte, error) {
	// TODO: Implement proper packet encoding based on RotMG protocol
	return nil, nil
}
