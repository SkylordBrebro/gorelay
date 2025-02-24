package models

// MapInfo contains information about a game map
type MapInfo struct {
	Width               int32  `json:"width"`
	Height              int32  `json:"height"`
	Name                string `json:"name"`
	DisplayName         string `json:"displayName"`
	Difficulty          int32  `json:"difficulty"`
	Background          int32  `json:"background"`
	AllowPlayerTeleport bool   `json:"allowPlayerTeleport"`
	ShowDisplays        bool   `json:"showDisplays"`
	MaxPlayers          int32  `json:"maxPlayers"`
	ClientXML           []byte `json:"clientXML,omitempty"`
	ExtraXML            []byte `json:"extraXML,omitempty"`
	Music               string `json:"music,omitempty"`
	Seed                int32  `json:"seed"`
}

// MapStats contains statistics about a map
type MapStats struct {
	NumPlayers    int32
	NumEnemies    int32
	NumQuests     int32
	NumPortals    int32
	UpdateTime    int64
	ElapsedTime   int64
	MaxPlayerSeen int32
}
