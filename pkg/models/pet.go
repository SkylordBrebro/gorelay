package models

// PetAbility represents different abilities that pets can have
type PetAbility int32

const (
	PetAbilityAttackClose PetAbility = iota
	PetAbilityAttackMid
	PetAbilityAttackFar
	PetAbilityElectric
	PetAbilityHeal
	PetAbilityMagicHeal
	PetAbilitySavage
	PetAbilityDecoy
	PetAbilityRising
	PetAbilityMayhemClose
	PetAbilityMayhemMid
	PetAbilityMayhemFar
)

// String returns the string representation of the pet ability
func (pa PetAbility) String() string {
	switch pa {
	case PetAbilityAttackClose:
		return "Attack Close"
	case PetAbilityAttackMid:
		return "Attack Mid"
	case PetAbilityAttackFar:
		return "Attack Far"
	case PetAbilityElectric:
		return "Electric"
	case PetAbilityHeal:
		return "Heal"
	case PetAbilityMagicHeal:
		return "Magic Heal"
	case PetAbilitySavage:
		return "Savage"
	case PetAbilityDecoy:
		return "Decoy"
	case PetAbilityRising:
		return "Rising Fury"
	case PetAbilityMayhemClose:
		return "Mayhem Close"
	case PetAbilityMayhemMid:
		return "Mayhem Mid"
	case PetAbilityMayhemFar:
		return "Mayhem Far"
	default:
		return "Unknown"
	}
}

// Pet represents a pet entity
type Pet struct {
	ID            int32
	Name          string
	Type          int32
	Rarity        string
	Family        string
	FirstAbility  PetAbility
	SecondAbility PetAbility
	ThirdAbility  PetAbility
	Level         int32
	MaxLevel      int32
	Abilities     map[PetAbility]struct {
		Level     int32
		Points    int32
		MaxPoints int32
	}
}

// PetYardType represents different types of pet yards
type PetYardType int32

const (
	PetYardCommon PetYardType = iota
	PetYardUncommon
	PetYardRare
	PetYardLegendary
	PetYardDivine
)

// String returns the string representation of the pet yard type
func (pyt PetYardType) String() string {
	switch pyt {
	case PetYardCommon:
		return "Common"
	case PetYardUncommon:
		return "Uncommon"
	case PetYardRare:
		return "Rare"
	case PetYardLegendary:
		return "Legendary"
	case PetYardDivine:
		return "Divine"
	default:
		return "Unknown"
	}
}

// MaxLevel returns the maximum level for pets in this yard type
func (pyt PetYardType) MaxLevel() int32 {
	switch pyt {
	case PetYardCommon:
		return 30
	case PetYardUncommon:
		return 50
	case PetYardRare:
		return 70
	case PetYardLegendary:
		return 90
	case PetYardDivine:
		return 100
	default:
		return 0
	}
}
