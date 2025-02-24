package models

// PetAbility represents different abilities that pets can have
type PetAbility int32

const (
	PetAbilityAttackClose PetAbility = 0x192
	PetAbilityAttackMid   PetAbility = 0x194
	PetAbilityAttackFar   PetAbility = 0x195
	PetAbilityElectric    PetAbility = 0x196
	PetAbilityHeal        PetAbility = 0x197
	PetAbilityMagicHeal   PetAbility = 0x198
	PetAbilitySavage      PetAbility = 0x199
	PetAbilityDecoy       PetAbility = 0x19a
	PetAbilityRisingFury  PetAbility = 0x19b
)

// Pet represents a pet entity
type Pet struct {
	ID            int32
	Name          string
	Type          int32
	Rarity        int32
	Family        int32
	FirstAbility  PetAbility
	SecondAbility PetAbility
	ThirdAbility  PetAbility
	MaxLevel      int32
	Abilities     map[PetAbility]struct {
		Level     int32
		Points    int32
		MaxPoints int32
	}
	OwnerID    int32
	ObjectID   int32
	HP         int32
	Size       int32
	Condition  int32
	Texture    int32
	InstanceID int32
}
