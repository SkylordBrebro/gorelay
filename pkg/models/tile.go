package models

// Tile represents a map tile
type Tile struct {
	Type      int32
	X         int32
	Y         int32
	Terrain   TileTerrain
	Region    string
	ObjectID  int32
	UpdatedAt int64
}

// TileTerrain represents different types of tile terrain
type TileTerrain int32

const (
	TerrainNone TileTerrain = iota
	TerrainGrass
	TerrainSand
	TerrainStone
	TerrainIce
	TerrainSnow
	TerrainLava
	TerrainWater
	TerrainSpace
)

// TileProperties represents properties of a tile type
type TileProperties struct {
	ID         int32
	NoWalk     bool
	MinDamage  int32
	MaxDamage  int32
	Speed      float32
	Push       bool
	Sink       bool
	SinkLevel  int32
	Texture    int32
	ObjectName string
	EdgeTex    int32
	CornerTex  int32
	Terrain    TileTerrain
	Region     string
}

// MapTile represents a tile in the game map
type MapTile struct {
	X          int32
	Y          int32
	Type       int32
	Terrain    TileTerrain
	Properties TileProperties
	Objects    []int32
}
