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

// String returns the string representation of the tile terrain
func (tt TileTerrain) String() string {
	switch tt {
	case TerrainGrass:
		return "Grass"
	case TerrainSand:
		return "Sand"
	case TerrainStone:
		return "Stone"
	case TerrainIce:
		return "Ice"
	case TerrainSnow:
		return "Snow"
	case TerrainLava:
		return "Lava"
	case TerrainWater:
		return "Water"
	case TerrainSpace:
		return "Space"
	default:
		return "None"
	}
}

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
