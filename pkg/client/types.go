package client

import (
	"gorelay/pkg/packets"
	"time"
)

// Projectile represents a projectile in the game
type Projectile struct {
	OwnerType  int32
	OwnerID    int32
	BulletID   int32
	BulletType int32
	Angle      float32
	StartTime  int64
	StartPos   *packets.WorldPosData
	Damage     int32
}

// Enemy represents an enemy entity in the game
type Enemy struct {
	ObjectID   int32
	ObjectType int32
	Position   *packets.WorldPosData
	HP         int32
	MaxHP      int32
	Dead       bool
	LastMove   time.Time
}

// IsDead returns whether the enemy is dead
func (e *Enemy) IsDead() bool {
	return e.Dead
}

// OnGoto updates the enemy's position
func (e *Enemy) OnGoto(x, y float32, timestamp int64) {
	e.Position.X = x
	e.Position.Y = y
	e.LastMove = time.Unix(0, timestamp*int64(time.Millisecond))
}

// Player represents another player in the game
type Player struct {
	ObjectID int32
	Name     string
	Position *packets.WorldPosData
	Stats    map[string]int32
	LastMove time.Time
}

// OnGoto updates the player's position
func (p *Player) OnGoto(x, y float32, timestamp int64) {
	p.Position.X = x
	p.Position.Y = y
	p.LastMove = time.Unix(0, timestamp*int64(time.Millisecond))
}
