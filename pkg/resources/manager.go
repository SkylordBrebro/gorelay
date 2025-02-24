package resources

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

// ResourceManager handles game resource definitions
type ResourceManager struct {
	Objects map[int32]ObjectDef
	Tiles   map[int32]TileDef
	Pets    map[int32]PetDef
}

// ObjectDef defines properties of game objects
type ObjectDef struct {
	ID          int32           `json:"id"`
	ObjectType  string          `json:"type"`
	DisplayName string          `json:"displayName"`
	Class       string          `json:"class"`
	MaxHP       int32           `json:"maxHP"`
	Defense     int32           `json:"defense"`
	SlotTypes   []int32         `json:"slotTypes"`
	Equipment   []int32         `json:"equipment"`
	Tex1        int32           `json:"tex1"`
	Tex2        int32           `json:"tex2"`
	Size        int32           `json:"size"`
	Enemy       bool            `json:"enemy"`
	God         bool            `json:"god"`
	Projectiles []ProjectileDef `json:"projectiles"`
}

// TileDef defines properties of map tiles
type TileDef struct {
	ID        int32   `json:"id"`
	Name      string  `json:"name"`
	NoWalk    bool    `json:"noWalk"`
	MinDamage int32   `json:"minDamage"`
	MaxDamage int32   `json:"maxDamage"`
	Speed     float32 `json:"speed"`
	Sink      bool    `json:"sink"`
	Push      bool    `json:"push"`
	SinkLevel int32   `json:"sinkLevel"`
	Tex       int32   `json:"tex"`
}

// PetDef defines properties of pets
type PetDef struct {
	ID            int32  `json:"id"`
	ObjectType    string `json:"type"`
	DisplayName   string `json:"displayName"`
	Rarity        string `json:"rarity"`
	Family        string `json:"family"`
	FirstAbility  string `json:"firstAbility"`
	SecondAbility string `json:"secondAbility"`
	ThirdAbility  string `json:"thirdAbility"`
}

// ProjectileDef defines properties of projectiles
type ProjectileDef struct {
	ID          int32   `json:"id"`
	ObjectType  string  `json:"type"`
	Speed       float32 `json:"speed"`
	MinDamage   int32   `json:"minDamage"`
	MaxDamage   int32   `json:"maxDamage"`
	Size        float32 `json:"size"`
	Amplitude   float32 `json:"amplitude"`
	Frequency   float32 `json:"frequency"`
	Wavy        bool    `json:"wavy"`
	Parametric  bool    `json:"parametric"`
	Boomerang   bool    `json:"boomerang"`
	MultiHit    bool    `json:"multiHit"`
	PassesCover bool    `json:"passesCover"`
}

// NewResourceManager creates a new resource manager instance
func NewResourceManager() *ResourceManager {
	return &ResourceManager{
		Objects: make(map[int32]ObjectDef),
		Tiles:   make(map[int32]TileDef),
		Pets:    make(map[int32]PetDef),
	}
}

// LoadResources loads game resources from JSON files
func (rm *ResourceManager) LoadResources(objectsPath, tilesPath string) error {
	// Load Objects.json
	objData, err := ioutil.ReadFile(objectsPath)
	if err != nil {
		return fmt.Errorf("failed to read objects file: %v", err)
	}

	var objects []ObjectDef
	if err := json.Unmarshal(objData, &objects); err != nil {
		return fmt.Errorf("failed to parse objects file: %v", err)
	}

	for _, obj := range objects {
		rm.Objects[obj.ID] = obj
	}

	// Load Tiles.json
	tileData, err := ioutil.ReadFile(tilesPath)
	if err != nil {
		return fmt.Errorf("failed to read tiles file: %v", err)
	}

	var tiles []TileDef
	if err := json.Unmarshal(tileData, &tiles); err != nil {
		return fmt.Errorf("failed to parse tiles file: %v", err)
	}

	for _, tile := range tiles {
		rm.Tiles[tile.ID] = tile
	}

	return nil
}

// GetObject returns an object definition by ID
func (rm *ResourceManager) GetObject(id int32) (ObjectDef, bool) {
	obj, ok := rm.Objects[id]
	return obj, ok
}

// GetTile returns a tile definition by ID
func (rm *ResourceManager) GetTile(id int32) (TileDef, bool) {
	tile, ok := rm.Tiles[id]
	return tile, ok
}

// GetPet returns a pet definition by ID
func (rm *ResourceManager) GetPet(id int32) (PetDef, bool) {
	pet, ok := rm.Pets[id]
	return pet, ok
}
