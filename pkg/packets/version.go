package packets

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

// VersionManager handles packet versioning and updates
type VersionManager struct {
	buildVersion string
	packetMap    map[string]int
	idToName     map[int]string
}

// NewVersionManager creates a new version manager
func NewVersionManager() *VersionManager {
	return &VersionManager{
		packetMap: make(map[string]int),
		idToName:  make(map[int]string),
	}
}

// LoadVersions loads packet versions from a JSON file
func (vm *VersionManager) LoadVersions(path string) error {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return fmt.Errorf("failed to read versions file: %v", err)
	}

	var versions struct {
		BuildVersion string         `json:"buildVersion"`
		PacketIds    map[string]int `json:"packetIds"`
	}

	if err := json.Unmarshal(data, &versions); err != nil {
		return fmt.Errorf("failed to parse versions file: %v", err)
	}

	vm.buildVersion = versions.BuildVersion
	vm.packetMap = versions.PacketIds

	// Build reverse mapping
	for name, id := range vm.packetMap {
		vm.idToName[id] = name
	}

	return nil
}

// GetPacketID returns the ID for a given packet name
func (vm *VersionManager) GetPacketID(name string) (int, error) {
	if id, ok := vm.packetMap[name]; ok {
		return id, nil
	}
	return 0, fmt.Errorf("no ID found for packet %s", name)
}

// GetPacketName returns the name for a given packet ID
func (vm *VersionManager) GetPacketName(id int) (string, error) {
	if name, ok := vm.idToName[id]; ok {
		return name, nil
	}
	return "", fmt.Errorf("no name found for packet ID %d", id)
}

// GetBuildVersion returns the current build version
func (vm *VersionManager) GetBuildVersion() string {
	return vm.buildVersion
}

// UpdateVersions updates the version information and saves it to disk
func (vm *VersionManager) UpdateVersions(buildVersion string, packetMap map[string]int, path string) error {
	versions := struct {
		BuildVersion string         `json:"buildVersion"`
		PacketIds    map[string]int `json:"packetIds"`
	}{
		BuildVersion: buildVersion,
		PacketIds:    packetMap,
	}

	data, err := json.MarshalIndent(versions, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal versions: %v", err)
	}

	if err := ioutil.WriteFile(path, data, 0644); err != nil {
		return fmt.Errorf("failed to write versions file: %v", err)
	}

	vm.buildVersion = buildVersion
	vm.packetMap = packetMap

	// Update reverse mapping
	vm.idToName = make(map[int]string)
	for name, id := range vm.packetMap {
		vm.idToName[id] = name
	}

	return nil
}
