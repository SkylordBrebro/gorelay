package models

// StatType represents different types of stats that can be applied to entities
type StatType int32

const (
	// Basic stats
	MAXHPSTAT        StatType = 0
	HPSTAT           StatType = 1
	SIZESTAT         StatType = 2
	MAXMPSTAT        StatType = 3
	MPSTAT           StatType = 4
	NEXTLEVELEXPSTAT StatType = 5
	EXPSTAT          StatType = 6
	LEVELSTAT        StatType = 7

	// Inventory slots 0-11
	INVENTORY0STAT  StatType = 8
	INVENTORY1STAT  StatType = 9
	INVENTORY2STAT  StatType = 10
	INVENTORY3STAT  StatType = 11
	INVENTORY4STAT  StatType = 12
	INVENTORY5STAT  StatType = 13
	INVENTORY6STAT  StatType = 14
	INVENTORY7STAT  StatType = 15
	INVENTORY8STAT  StatType = 16
	INVENTORY9STAT  StatType = 17
	INVENTORY10STAT StatType = 18
	INVENTORY11STAT StatType = 19

	// Base stats
	ATTACKSTAT              StatType = 20
	DEFENSESTAT             StatType = 21
	SPEEDSTAT               StatType = 22
	TEXTURESTAT             StatType = 25
	VITALITYSTAT            StatType = 26
	WISDOMSTAT              StatType = 27
	DEXTERITYSTAT           StatType = 28
	CONDITIONSTAT           StatType = 29
	NUMSTARSSTAT            StatType = 30
	NAMESTAT                StatType = 31
	TEX1STAT                StatType = 32
	TEX2STAT                StatType = 33
	MERCHANDISETYPESTAT     StatType = 34
	CREDITSSTAT             StatType = 35
	MERCHANDISEPRICESTAT    StatType = 36
	ACTIVESTAT              StatType = 37
	ACCOUNTIDSTAT           StatType = 38
	FAMESTAT                StatType = 39
	MERCHANDISECURRENCYSTAT StatType = 40
	CONNECTSTAT             StatType = 41
	MERCHANDISECOUNTSTAT    StatType = 42
	MERCHANDISEMINSLEFTSTAT StatType = 43
	MERCHANDISEDISCOUNTSTAT StatType = 44
	MERCHANDISERANKREQSTAT  StatType = 45

	// Stat boosts
	MAXHPBOOSTSTAT     StatType = 46
	MAXMPBOOSTSTAT     StatType = 47
	ATTACKBOOSTSTAT    StatType = 48
	DEFENSEBOOSTSTAT   StatType = 49
	SPEEDBOOSTSTAT     StatType = 50
	VITALITYBOOSTSTAT  StatType = 51
	WISDOMBOOSTSTAT    StatType = 52
	DEXTERITYBOOSTSTAT StatType = 53

	// Account and character stats
	OWNERACCOUNTIDSTAT     StatType = 54
	RANKREQUIREDSTAT       StatType = 55
	NAMECHOSENSTAT         StatType = 56
	CURRFAMESTAT           StatType = 57
	NEXTCLASSQUESTFAMESTAT StatType = 58
	LEGENDARYRANKSTAT      StatType = 59
	SINKLEVELSTAT          StatType = 60
	ALTTEXTURESTAT         StatType = 61
	GUILDNAMESTAT          StatType = 62
	GUILDRANKSTAT          StatType = 63
	BREATHSTAT             StatType = 64
	XPBOOSTEDSTAT          StatType = 65
	XPTIMERSTAT            StatType = 66
	LDTIMERSTAT            StatType = 67
	LTTIMERSTAT            StatType = 68

	// Potion and backpack stats
	HEALTHPOTIONSTACKSTAT StatType = 69
	MAGICPOTIONSTACKSTAT  StatType = 70
	BACKPACK0STAT         StatType = 71
	BACKPACK1STAT         StatType = 72
	BACKPACK2STAT         StatType = 73
	BACKPACK3STAT         StatType = 74
	BACKPACK4STAT         StatType = 75
	BACKPACK5STAT         StatType = 76
	BACKPACK6STAT         StatType = 77
	BACKPACK7STAT         StatType = 78
	HASBACKPACKSTAT       StatType = 79

	// Unknown stat
	UNKNOWN80 StatType = 80

	// Pet stats
	PETINSTANCEIDSTAT         StatType = 81
	PETNAMESTAT               StatType = 82
	PETTYPESTAT               StatType = 83
	PETRARITYSTAT             StatType = 84
	PETMAXABILITYPOWERSTAT    StatType = 85
	PETFAMILYSTAT             StatType = 86
	PETFIRSTABILITYPOINTSTAT  StatType = 87
	PETSECONDABILITYPOINTSTAT StatType = 88
	PETTHIRDABILITYPOINTSTAT  StatType = 89
	PETFIRSTABILITYPOWERSTAT  StatType = 90
	PETSECONDABILITYPOWERSTAT StatType = 91
	PETTHIRDABILITYPOWERSTAT  StatType = 92
	PETFIRSTABILITYTYPESTAT   StatType = 93
	PETSECONDABILITYTYPESTAT  StatType = 94
	PETTHIRDABILITYTYPESTAT   StatType = 95

	// Additional stats
	NEWCONSTAT           StatType = 96
	FORTUNETOKENSTAT     StatType = 97
	SUPPORTERPOINTSSTAT  StatType = 98
	SUPPORTERSTAT        StatType = 99
	CHALLENGERSTARBGSTAT StatType = 100

	// Projectile modifiers
	PROJECTILESPEEDMULT StatType = 102
	PROJECTILELIFEMULT  StatType = 103

	// Additional stats
	OPENEDATTIMESTAMP     StatType = 104
	EXALTEDATK            StatType = 105
	EXALTEDDEFENSE        StatType = 106
	EXALTEDSPD            StatType = 107
	EXALTEDVIT            StatType = 108
	EXALTEDWIS            StatType = 109
	EXALTEDDEX            StatType = 110
	EXALTEDHP             StatType = 111
	EXALTEDMP             StatType = 112
	EXALTATIONBONUSDMG    StatType = 113
	EXALTATIONICREDUCTION StatType = 114
	GRAVEACCOUNTID        StatType = 115
	POTIONONETYPE         StatType = 116
	POTIONTWOTYPE         StatType = 117
	POTIONTHREETYPE       StatType = 118
	POTIONBELT            StatType = 119
	FORGEFIRE             StatType = 120
	UNKNOWN121            StatType = 121
	UNKNOWN123            StatType = 123
)
