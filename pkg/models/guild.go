package models

// GuildRank represents different ranks within a guild
type GuildRank int32

const (
	GuildRankNoRank   GuildRank = -1
	GuildRankInitiate GuildRank = 0
	GuildRankMember   GuildRank = 10
	GuildRankOfficer  GuildRank = 20
	GuildRankLeader   GuildRank = 30
	GuildRankFounder  GuildRank = 40
)

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
