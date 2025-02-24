package models

// ConditionEffect represents different status effects that can be applied to entities
type ConditionEffect int32

const (
	ConditionEffectNone ConditionEffect = iota
	ConditionEffectDead
	ConditionEffectQuiet
	ConditionEffectWeak
	ConditionEffectSlowed
	ConditionEffectSick
	ConditionEffectDazed
	ConditionEffectStunned
	ConditionEffectBlind
	ConditionEffectHallucinating
	ConditionEffectDrunk
	ConditionEffectConfused
	ConditionEffectStunImmune
	ConditionEffectInvisible
	ConditionEffectParalyzed
	ConditionEffectSpeedy
	ConditionEffectBleeding
	ConditionEffectArmorBroken
	ConditionEffectBerserk
	ConditionEffectHealing
	ConditionEffectDamaging
	ConditionEffectBrave
	ConditionEffectStasis
	ConditionEffectInvincible
	ConditionEffectInvulnerable
	ConditionEffectArmored
	ConditionEffectCursed
)
