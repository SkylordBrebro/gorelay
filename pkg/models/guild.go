package models

// GuildRank represents different ranks within a guild
type GuildRank int32

const (
	GuildRankInitiate GuildRank = iota
	GuildRankMember
	GuildRankOfficer
	GuildRankLeader
	GuildRankFounder
)

// String returns the string representation of the guild rank
func (gr GuildRank) String() string {
	switch gr {
	case GuildRankInitiate:
		return "Initiate"
	case GuildRankMember:
		return "Member"
	case GuildRankOfficer:
		return "Officer"
	case GuildRankLeader:
		return "Leader"
	case GuildRankFounder:
		return "Founder"
	default:
		return "Unknown"
	}
}

// Guild represents a guild in the game
type Guild struct {
	ID          int32
	Name        string
	Level       int32
	TotalFame   int32
	CurrentFame int32
	Members     []GuildMember
	Board       []GuildBoardEntry
}

// GuildMember represents a member of a guild
type GuildMember struct {
	Name      string
	Rank      GuildRank
	Fame      int32
	LastSeen  int64
	Character struct {
		Class       int32
		Level       int32
		Fame        int32
		Equipment   []int32
		Stats       map[string]int32
		HasBackpack bool
	}
}

// GuildBoardEntry represents an entry on the guild board
type GuildBoardEntry struct {
	Author    string
	Text      string
	Timestamp int64
}
