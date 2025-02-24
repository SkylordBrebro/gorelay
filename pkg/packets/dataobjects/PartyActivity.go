package dataobjects

// PartyActivity represents the type of activity a party is engaged in
type PartyActivity byte

const (
	// Dungeons represents dungeon activity
	Dungeons PartyActivity = iota
	// Realms represents realm activity
	Realms
	// Other represents other activity
	Other
)
