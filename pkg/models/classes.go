package models

// CharacterClass represents the object types of all classes in the game
type CharacterClass int32

const (
	// Class object type IDs
	ClassRogue       CharacterClass = 768
	ClassArcher      CharacterClass = 775
	ClassWizard      CharacterClass = 782
	ClassPriest      CharacterClass = 784
	ClassWarrior     CharacterClass = 797
	ClassKnight      CharacterClass = 798
	ClassPaladin     CharacterClass = 799
	ClassAssassin    CharacterClass = 800
	ClassNecromancer CharacterClass = 801
	ClassHuntress    CharacterClass = 802
	ClassMystic      CharacterClass = 803
	ClassTrickster   CharacterClass = 804
	ClassSorcerer    CharacterClass = 805
	ClassNinja       CharacterClass = 806
	ClassSamurai     CharacterClass = 785
)
