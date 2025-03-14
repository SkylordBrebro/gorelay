﻿package dataobjects

// StatsType represents a byte-based stats type
type StatsType byte

// Stats type constants
var (
	MaximumHP                   = StatsType(0)
	HP                          = StatsType(1)
	Size                        = StatsType(2)
	MaximumMP                   = StatsType(3)
	MP                          = StatsType(4)
	NextLevelExperience         = StatsType(5)
	Experience                  = StatsType(6)
	Level                       = StatsType(7)
	Inventory0                  = StatsType(8)
	Inventory1                  = StatsType(9)
	Inventory2                  = StatsType(10)
	Inventory3                  = StatsType(11)
	Inventory4                  = StatsType(12)
	Inventory5                  = StatsType(13)
	Inventory6                  = StatsType(14)
	Inventory7                  = StatsType(15)
	Inventory8                  = StatsType(16)
	Inventory9                  = StatsType(17)
	Inventory10                 = StatsType(18)
	Inventory11                 = StatsType(19)
	Attack                      = StatsType(20)
	Defense                     = StatsType(21)
	Speed                       = StatsType(22)
	UnknownPBJIGGHPBIH          = StatsType(23)
	Seasonal                    = StatsType(24)
	Skin                        = StatsType(25)
	Vitality                    = StatsType(26)
	Wisdom                      = StatsType(27)
	Dexterity                   = StatsType(28)
	Effects                     = StatsType(29)
	Stars                       = StatsType(30)
	Name                        = StatsType(31)
	Texture1                    = StatsType(32)
	Texture2                    = StatsType(33)
	MerchandiseType             = StatsType(34)
	Credits                     = StatsType(35)
	MerchandisePrice            = StatsType(36)
	PortalUsable                = StatsType(37)
	AccountID                   = StatsType(38)
	AccountFame                 = StatsType(39)
	MerchandiseCurrency         = StatsType(40)
	ObjectConnection            = StatsType(41)
	MerchandiseRemainingCount   = StatsType(42)
	MerchandiseRemainingMinutes = StatsType(43)
	MerchandiseDiscount         = StatsType(44)
	MerchandiseRankRequirement  = StatsType(45)
	HealthBonus                 = StatsType(46)
	ManaBonus                   = StatsType(47)
	AttackBonus                 = StatsType(48)
	DefenseBonus                = StatsType(49)
	SpeedBonus                  = StatsType(50)
	VitalityBonus               = StatsType(51)
	WisdomBonus                 = StatsType(52)
	DexterityBonus              = StatsType(53)
	OwnerAccountID              = StatsType(54)
	RankRequired                = StatsType(55)
	NameChosen                  = StatsType(56)
	CharacterFame               = StatsType(57)
	CharacterFameGoal           = StatsType(58)
	Glowing                     = StatsType(59)
	SinkLevel                   = StatsType(60)
	AltTextureIndex             = StatsType(61)
	GuildName                   = StatsType(62)
	GuildRank                   = StatsType(63)
	OxygenBar                   = StatsType(64)
	XPBoosterActive             = StatsType(65)
	XPBoostTime                 = StatsType(66)
	LootDropBoostTime           = StatsType(67)
	LootTierBoostTime           = StatsType(68)
	HealthPotionCount           = StatsType(69)
	MagicPotionCount            = StatsType(70)
	Unknown71                   = StatsType(71)
	Unknown72                   = StatsType(72)
	Unknown80                   = StatsType(80)
	PetInstanceID               = StatsType(81)
	PetName                     = StatsType(82)
	PetType                     = StatsType(83)
	PetRarity                   = StatsType(84)
	PetMaximumLevel             = StatsType(85)
	PetFamily                   = StatsType(86)
	PetPoints0                  = StatsType(87)
	PetPoints1                  = StatsType(88)
	PetPoints2                  = StatsType(89)
	PetLevel0                   = StatsType(90)
	PetLevel1                   = StatsType(91)
	PetLevel2                   = StatsType(92)
	PetAbilityType0             = StatsType(93)
	PetAbilityType1             = StatsType(94)
	PetAbilityType2             = StatsType(95)
	Effects2                    = StatsType(96)
	FortuneTokens               = StatsType(97)
	SupporterPoints             = StatsType(98)
	SupporterStat               = StatsType(99)
	ChallengerStarBg            = StatsType(100)
	Unknown101                  = StatsType(101)
	ProjectileSpeedMult         = StatsType(102)
	ProjectileLifeMult          = StatsType(103)
	CreatedTimestamp            = StatsType(104)
	PowerupExtraAttack          = StatsType(105)
	PowerupExtraDefense         = StatsType(106)
	PowerupExtraSpeed           = StatsType(107)
	PowerupExtraVitality        = StatsType(108)
	PowerupExtraWisdom          = StatsType(109)
	PowerupExtraDexterity       = StatsType(110)
	PowerupExtraMaxHP           = StatsType(111)
	PowerupExtraMaxMP           = StatsType(112)
	PowerupDamageMult           = StatsType(113)
	PetOwnerObjectID            = StatsType(114)
	GraveAccountID              = StatsType(115)
	Potion1                     = StatsType(116)
	Potion2                     = StatsType(117)
	Potion3                     = StatsType(118)
	PotionBelt                  = StatsType(119)
	ForgeFire                   = StatsType(120)
	Unknown124                  = StatsType(124)
	Unknown127                  = StatsType(127)
	BackpackSlots               = StatsType(130)
)

// IsUTF returns true if the stat type uses UTF-8 encoding
func (s StatsType) IsUTF() bool {
	switch s {
	case Name, AccountID, Experience, GuildName, PetName, GraveAccountID,
		Unknown71, Unknown72, Unknown80, Unknown127:
		return true
	default:
		return false
	}
}

// String returns a string representation of the StatsType
func (s StatsType) String() string {
	return string(s)
}
