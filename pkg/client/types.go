package client

import (
	"gorelay/pkg/packets"
	"time"
)

// GameState represents the current state of the game
type GameState struct {
	ObjectID      int32
	WorldPos      *packets.WorldPosData
	PlayerData    *packets.PlayerData
	GameID        int32
	BuildVer      string
	LastUpdate    time.Time
	LastFrameTime int64
}

// Projectile represents a projectile in the game
type Projectile struct {
	ID         int32
	OwnerID    int32
	OwnerType  int32
	BulletType int32
	Angle      float32
	Damage     int32
	StartTime  int64
	StartPos   *packets.WorldPosData
	Position   *packets.WorldPosData
	Destroyed  bool
}

// Enemy represents an enemy entity in the game
type Enemy struct {
	ObjectID   int32
	ObjectType int32
	Position   *packets.WorldPosData
	HP         int32
	MaxHP      int32
	Defense    int32
	Dead       bool
	LastMove   time.Time
	LastHit    time.Time
	Effects    []int32                // Status effects
	Properties map[string]interface{} // Additional properties from resources
}

// IsDead returns whether the enemy is dead
func (e *Enemy) IsDead() bool {
	return e.Dead || e.HP <= 0
}

// OnGoto updates the enemy's position
func (e *Enemy) OnGoto(x, y float32, timestamp int64) {
	e.Position.X = x
	e.Position.Y = y
	e.LastMove = time.Unix(0, timestamp*int64(time.Millisecond))
}

// OnDamage handles damage taken by the enemy
func (e *Enemy) OnDamage(damage int32, armorPiercing bool) {
	if armorPiercing {
		e.HP -= damage
	} else {
		defense := float32(e.Defense)
		actualDamage := float32(damage) - defense
		if actualDamage < 1 {
			actualDamage = 1
		}
		e.HP -= int32(actualDamage)
	}
	e.LastHit = time.Now()
}

// Player represents another player in the game
type Player struct {
	ObjectID   int32
	Name       string
	Position   *packets.WorldPosData
	Stats      map[string]int32
	Equipment  map[int32]int32 // Slot -> ItemID
	Class      int32
	Level      int32
	Fame       int32
	Guild      string
	LastMove   time.Time
	LastAction time.Time
	Effects    []int32 // Status effects
}

// OnGoto updates the player's position
func (p *Player) OnGoto(x, y float32, timestamp int64) {
	p.Position.X = x
	p.Position.Y = y
	p.LastMove = time.Unix(0, timestamp*int64(time.Millisecond))
}

// OnAction updates the player's last action time
func (p *Player) OnAction() {
	p.LastAction = time.Now()
}

// HasEffect checks if the player has a specific status effect
func (p *Player) HasEffect(effect int32) bool {
	for _, e := range p.Effects {
		if e == effect {
			return true
		}
	}
	return false
}

// Tile represents a map tile
type Tile struct {
	Type     int32
	ObjectID int32
	Position *packets.WorldPosData
}

// Map represents the current game map
type Map struct {
	Name   string
	Width  int32
	Height int32
	Tiles  [][]Tile
	Seed   int32
}
