package xmldata

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"time"
)

// Global access points
var Objects *GameData
var Tiles *GroundData

// GameObjects holds the collection of all game objects
type GameObjects struct {
	Objects []GameObject `xml:"Object"`
}

// GameObject represents a single game object with its properties
type GameObject struct {
	Type      string `xml:"type,attr"`
	ID        string `xml:"id,attr"`
	DisplayID string `xml:"DisplayId"`
	Class     string `xml:"Class"`
	
	// Combat stats
	Size         int    `xml:"Size"`
	MaxHitPoints int    `xml:"MaxHitPoints"`
	Defense      int    `xml:"Defense"`
	
	// Additional tags as boolean flags
	Enemy  *struct{} `xml:"Enemy"`
	God    *struct{} `xml:"God"`
	Quest  *struct{} `xml:"Quest"`
	Oryx   *struct{} `xml:"Oryx"`
	
	// Labels string for quick checking
	Labels string `xml:"Labels"`
	
	// Projectiles
	Projectiles []Projectile `xml:"Projectile"`
}

// Projectile represents a projectile definition
type Projectile struct {
	ID         string `xml:"id,attr"`
	ObjectID   string `xml:"ObjectId"`
	Damage     int    `xml:"Damage"`
	Speed      float32 `xml:"Speed"`
	LifetimeMS float32    `xml:"LifetimeMS"`
	Size       int    `xml:"Size"`
	MultiHit   *struct{} `xml:"MultiHit"`
}

// GroundTypes holds all ground type definitions
type GroundTypes struct {
	Grounds []GroundType `xml:"GroundType"`
}

// GroundType represents a single ground type
type GroundType struct {
	Type        string `xml:"type,attr"`
	ID          string `xml:"id,attr"`
	DisplayID   string `xml:"DisplayId"`
	NoWalk      *struct{} `xml:"NoWalk"`
	Speed       float32 `xml:"Speed"`
	SinkLevel   int `xml:"SinkLevel"`
	Texture     *TextureData `xml:"Texture"`
}

// TextureData holds texture file and index information
type TextureData struct {
	File  string `xml:"File"`
	Index int    `xml:"Index"`
}

// Optimized lookup structures for game objects
type GameData struct {
	ObjectsByID   map[string]*GameObject
	ObjectsByType map[string]*GameObject
	ObjectsByTypeID map[int]*GameObject // For lookup by numeric type ID
}

// Optimized lookup structures for ground types
type GroundData struct {
	GroundsByID   map[string]*GroundType
	GroundsByType map[string]*GroundType  
	GroundsByTypeID map[int]*GroundType // For lookup by numeric type ID
}

// LoadGameObjects loads and parses the XML data file
func LoadGameObjects(filename string) (*GameObjects, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("error opening file: %w", err)
	}
	defer file.Close()

	data, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf("error reading file: %w", err)
	}

	var gameObjects GameObjects
	err = xml.Unmarshal(data, &gameObjects)
	if err != nil {
		return nil, fmt.Errorf("error unmarshalling XML: %w", err)
	}

	return &gameObjects, nil
}

// LoadGroundTypes loads and parses the ground types XML file
func LoadGroundTypes(filename string) (*GroundTypes, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("error opening file: %w", err)
	}
	defer file.Close()

	data, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf("error reading file: %w", err)
	}

	var groundTypes GroundTypes
	err = xml.Unmarshal(data, &groundTypes)
	if err != nil {
		return nil, fmt.Errorf("error unmarshalling XML: %w", err)
	}

	return &groundTypes, nil
}

// CreateGameData processes the game objects into optimized lookup structures
func CreateGameData(gameObjects *GameObjects) *GameData {
	data := &GameData{
		ObjectsByID:   make(map[string]*GameObject),
		ObjectsByType: make(map[string]*GameObject),
		ObjectsByTypeID: make(map[int]*GameObject),
	}

	for i := range gameObjects.Objects {
		obj := &gameObjects.Objects[i]
		
		// Store by ID for quick lookup
		if obj.ID != "" {
			data.ObjectsByID[obj.ID] = obj
		}
		
		// Store by Type for quick lookup (string)
		if obj.Type != "" {
			data.ObjectsByType[obj.Type] = obj
			
			// Also store by numeric type ID (convert hex or decimal)
			typeID, err := parseTypeID(obj.Type)
			if err == nil {
				data.ObjectsByTypeID[typeID] = obj
			}
		}
	}

	return data
}

// CreateGroundData processes the ground types into optimized lookup structures
func CreateGroundData(groundTypes *GroundTypes) *GroundData {
	data := &GroundData{
		GroundsByID:   make(map[string]*GroundType),
		GroundsByType: make(map[string]*GroundType),
		GroundsByTypeID: make(map[int]*GroundType),
	}

	for i := range groundTypes.Grounds {
		ground := &groundTypes.Grounds[i]
		
		// Store by ID for quick lookup
		if ground.ID != "" {
			data.GroundsByID[ground.ID] = ground
		}
		
		// Store by Type for quick lookup (string)
		if ground.Type != "" {
			data.GroundsByType[ground.Type] = ground
			
			// Also store by numeric type ID (convert hex or decimal)
			typeID, err := parseTypeID(ground.Type)
			if err == nil {
				data.GroundsByTypeID[typeID] = ground
			}
		}
	}

	return data
}

// Helper function to parse type IDs which could be hex (0xNNNN) or decimal
func parseTypeID(typeStr string) (int, error) {
	// Check if it's a hex value (0xNNNN format)
	if len(typeStr) > 2 && typeStr[0:2] == "0x" {
		parsed, err := strconv.ParseInt(typeStr[2:], 16, 32)
		if err != nil {
			return -1, fmt.Errorf("could not parse type ID: %s", typeStr)
		}
		return int(parsed), nil
	}
	
	// Try parsing as decimal
	val, err := strconv.Atoi(typeStr)
	if err != nil {
		return 0, fmt.Errorf("could not parse type ID: %s", typeStr)
	}
	return val, nil
}

// LoadAssets loads all game assets (objects and ground types)
func LoadAssets() {
	start := time.Now()
	
	// Load game objects
	parsed, err := LoadGameObjects("Xml/Objects.xml")
	if err != nil {
		fmt.Printf("Error loading game objects: %v\n", err)
		return
	}
	
	elapsed := time.Since(start)
	fmt.Printf("Objects.xml loaded in %s\n", elapsed)
	start = time.Now()

	Objects = CreateGameData(parsed)
	
	fmt.Printf("%d Objects parsed in %s\n", len(parsed.Objects), time.Since(start))
	start = time.Now()
	
	// Load ground types
	parsedGrounds, err := LoadGroundTypes("Xml/GroundTypes.xml")
	if err != nil {
		fmt.Printf("Error loading ground types: %v\n", err)
		return
	}
	
	fmt.Printf("GroundTypes.xml loaded in %s\n", time.Since(start))
	start = time.Now()
	
	Tiles = CreateGroundData(parsedGrounds)
	
	fmt.Printf("%d GroundTypes parsed in %s\n", len(parsedGrounds.Grounds), time.Since(start))
}

// GetObjectByTypeID returns a game object by its numeric type ID
func GetObjectByTypeID(typeID int) *GameObject {
	if Objects == nil {
		return nil
	}
	return Objects.ObjectsByTypeID[typeID]
}

// GetGroundByTypeID returns a ground type by its numeric type ID
func GetGroundByTypeID(typeID int) *GroundType {
	if Tiles == nil {
		return nil
	}
	return Tiles.GroundsByTypeID[typeID]
}