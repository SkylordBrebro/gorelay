package models

import "time"

// Entity represents a base game entity
type Entity struct {
	ObjectID   int32
	ObjectType int32
	Position   struct {
		X float32
		Y float32
	}
	Size      float32
	Condition int32 // Bitmask of condition effects
	LastMove  time.Time
}

// Object represents a static game object
type Object struct {
	Entity
	Properties map[string]interface{}
	Class      string
	DisplayID  int32
	Tex1       int32
	Tex2       int32
}

// Enemy represents an enemy entity
type Enemy struct {
	Entity
	HP          int32
	MaxHP       int32
	Defense     int32
	Experience  int32
	Dead        bool
	God         bool
	Quest       bool
	NoKnockback bool
}

// Player represents a player entity
type Player struct {
	Entity
	Name        string
	Level       int32
	Experience  int32
	Fame        int32
	Stars       int32
	AccountID   string
	Guild       string
	GuildRank   GuildRank
	HP          int32
	MaxHP       int32
	MP          int32
	MaxMP       int32
	Stats       map[string]int32
	Inventory   []int32
	HasBackpack bool
	Class       CharacterClass
	Tex1        int32
	Tex2        int32
}

// Container represents an inventory container
type Container struct {
	Entity
	Items     []int32
	Slots     int32
	OwnerId   int32
	UpdatedAt time.Time
}

// Projectile represents a projectile entity
type Projectile struct {
	Entity
	OwnerID     int32
	BulletID    int32
	BulletType  int32
	Damage      int32
	Angle       float32
	StartTime   time.Time
	MultiHit    bool
	PassesCover bool
}

// Portal represents a portal entity
type Portal struct {
	Entity
	Name        string
	Nexus       bool
	Locked      bool
	Active      bool
	UsableIn    int32
	DungeonName string
}
