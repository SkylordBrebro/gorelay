package client

import (
	"gorelay/pkg/packets"
	"time"
)

// GameObject represents a game object with its properties
type GameObject struct {
	Type                    int32
	ID                      string
	Enemy                   bool
	Item                    bool
	God                     bool
	Pet                     bool
	SlotType                int32
	BagType                 int32
	Class                   string
	MaxHitPoints            int32
	Defense                 int32
	XPMultiplier            float32
	Projectiles             []ProjectileInfo
	ActivateOnEquip         []StatBonus
	RateOfFire              float32
	NumProjectiles          int32
	ArcGap                  float32
	FameBonus               int32
	FeedPower               int32
	OccupySquare            bool
	FullOccupy              bool
	ProtectFromGroundDamage bool
}

// ProjectileInfo represents projectile properties
type ProjectileInfo struct {
	ID               int32
	ObjectID         string
	Damage           int32
	ArmorPiercing    bool
	MinDamage        int32
	MaxDamage        int32
	Speed            float32
	LifetimeMS       int32
	Parametric       bool
	Wavy             bool
	Boomerang        bool
	MultiHit         bool
	PassesCover      bool
	Amplitude        float32
	Frequency        float32
	Magnitude        float32
	ConditionEffects []ConditionEffect
}

// StatBonus represents stat changes when equipping items
type StatBonus struct {
	StatType int32
	Amount   int32
}

// ConditionEffect represents status effects applied by projectiles
type ConditionEffect struct {
	EffectName string
	Duration   float32
}

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

// Projectile represents an active projectile in the game
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
	Effects    []int32
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
	Effects    []int32
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

// Map represents the current game map
type Map struct {
	Name   string
	Width  int32
	Height int32
	Tiles  [][]int32 // Just store tile types as integers
	Seed   int32
}
