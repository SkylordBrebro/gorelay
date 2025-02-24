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

// String returns the string representation of the condition effect
func (ce ConditionEffect) String() string {
	switch ce {
	case ConditionEffectDead:
		return "Dead"
	case ConditionEffectQuiet:
		return "Quiet"
	case ConditionEffectWeak:
		return "Weak"
	case ConditionEffectSlowed:
		return "Slowed"
	case ConditionEffectSick:
		return "Sick"
	case ConditionEffectDazed:
		return "Dazed"
	case ConditionEffectStunned:
		return "Stunned"
	case ConditionEffectBlind:
		return "Blind"
	case ConditionEffectHallucinating:
		return "Hallucinating"
	case ConditionEffectDrunk:
		return "Drunk"
	case ConditionEffectConfused:
		return "Confused"
	case ConditionEffectStunImmune:
		return "Stun Immune"
	case ConditionEffectInvisible:
		return "Invisible"
	case ConditionEffectParalyzed:
		return "Paralyzed"
	case ConditionEffectSpeedy:
		return "Speedy"
	case ConditionEffectBleeding:
		return "Bleeding"
	case ConditionEffectArmorBroken:
		return "Armor Broken"
	case ConditionEffectBerserk:
		return "Berserk"
	case ConditionEffectHealing:
		return "Healing"
	case ConditionEffectDamaging:
		return "Damaging"
	case ConditionEffectBrave:
		return "Brave"
	case ConditionEffectStasis:
		return "Stasis"
	case ConditionEffectInvincible:
		return "Invincible"
	case ConditionEffectInvulnerable:
		return "Invulnerable"
	case ConditionEffectArmored:
		return "Armored"
	case ConditionEffectCursed:
		return "Cursed"
	default:
		return "None"
	}
}
