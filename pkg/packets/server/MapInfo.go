package server

import (
	"gorelay/pkg/packets/interfaces"
)

// MapInfo represents the server packet for map information
type MapInfo struct {
	Width               int32
	Height              int32
	Name                string
	DisplayName         string
	RealmName           string
	Seed                int32
	Background          int32
	Difficulty          float32
	AllowPlayerTeleport bool
	NoSave              bool
	ShowDisplays        bool
	MaxPlayers          int16
	GameOpenedTime      int32
	ServerVersion       string
	BGColor             int32
	ViewRadius          byte
	DungeonModifiers    string
	DungeonModifiers2   string
	DungeonModifiers3   string
	Unknown             int16
	MaxRealmScore       int32
	CurrentRealmScore   int32
}

// Type returns the packet type for MapInfo
func (p *MapInfo) Type() interfaces.PacketType {
	return interfaces.MapInfo
}

// Read reads the packet data from the provided reader
func (p *MapInfo) Read(r interfaces.Reader) error {
	var err error

	// Read Width
	p.Width, err = r.ReadInt32()
	if err != nil {
		return err
	}

	// Read Height
	p.Height, err = r.ReadInt32()
	if err != nil {
		return err
	}

	// Read Name
	p.Name, err = r.ReadString()
	if err != nil {
		return err
	}

	// Read DisplayName
	p.DisplayName, err = r.ReadString()
	if err != nil {
		return err
	}

	// Read RealmName
	p.RealmName, err = r.ReadString()
	if err != nil {
		return err
	}

	// Read Seed
	p.Seed, err = r.ReadInt32()
	if err != nil {
		return err
	}

	// Read Background
	p.Background, err = r.ReadInt32()
	if err != nil {
		return err
	}

	// Read Difficulty
	p.Difficulty, err = r.ReadFloat32()
	if err != nil {
		return err
	}

	// Read AllowPlayerTeleport
	p.AllowPlayerTeleport, err = r.ReadBool()
	if err != nil {
		return err
	}

	// Read NoSave
	p.NoSave, err = r.ReadBool()
	if err != nil {
		return err
	}

	// Read ShowDisplays
	p.ShowDisplays, err = r.ReadBool()
	if err != nil {
		return err
	}

	// Read MaxPlayers
	p.MaxPlayers, err = r.ReadInt16()
	if err != nil {
		return err
	}

	// Read GameOpenedTime
	p.GameOpenedTime, err = r.ReadInt32()
	if err != nil {
		return err
	}

	// Read ServerVersion
	p.ServerVersion, err = r.ReadString()
	if err != nil {
		return err
	}

	// Check if we have more data to read
	if r.RemainingBytes() > 3 {
		// Read BGColor
		p.BGColor, err = r.ReadInt32()
		if err != nil {
			return err
		}
	}

	if r.RemainingBytes() > 0 {
		// Read ViewRadius
		p.ViewRadius, err = r.ReadByte()
		if err != nil {
			return err
		}
	}

	if r.RemainingBytes() > 0 {
		// Read DungeonModifiers
		p.DungeonModifiers, err = r.ReadString()
		if err != nil {
			return err
		}

		// Read DungeonModifiers2
		p.DungeonModifiers2, err = r.ReadString()
		if err != nil {
			return err
		}

		// Read DungeonModifiers3
		p.DungeonModifiers3, err = r.ReadString()
		if err != nil {
			return err
		}
	}

	if r.RemainingBytes() > 0 {
		// Read Unknown
		p.Unknown, err = r.ReadInt16()
		if err != nil {
			return err
		}
	}

	if r.RemainingBytes() > 7 {
		// Read MaxRealmScore
		p.MaxRealmScore, err = r.ReadInt32()
		if err != nil {
			return err
		}

		// Read CurrentRealmScore
		p.CurrentRealmScore, err = r.ReadInt32()
		if err != nil {
			return err
		}
	}

	return nil
}

// Write writes the packet data to the provided writer
func (p *MapInfo) Write(w interfaces.Writer) error {
	var err error

	// Write Width
	err = w.WriteInt32(p.Width)
	if err != nil {
		return err
	}

	// Write Height
	err = w.WriteInt32(p.Height)
	if err != nil {
		return err
	}

	// Write Name
	err = w.WriteString(p.Name)
	if err != nil {
		return err
	}

	// Write DisplayName
	err = w.WriteString(p.DisplayName)
	if err != nil {
		return err
	}

	// Write RealmName
	err = w.WriteString(p.RealmName)
	if err != nil {
		return err
	}

	// Write Seed
	err = w.WriteInt32(p.Seed)
	if err != nil {
		return err
	}

	// Write Background
	err = w.WriteInt32(p.Background)
	if err != nil {
		return err
	}

	// Write Difficulty
	err = w.WriteFloat32(p.Difficulty)
	if err != nil {
		return err
	}

	// Write AllowPlayerTeleport
	err = w.WriteBool(p.AllowPlayerTeleport)
	if err != nil {
		return err
	}

	// Write NoSave
	err = w.WriteBool(p.NoSave)
	if err != nil {
		return err
	}

	// Write ShowDisplays
	err = w.WriteBool(p.ShowDisplays)
	if err != nil {
		return err
	}

	// Write MaxPlayers
	err = w.WriteInt16(p.MaxPlayers)
	if err != nil {
		return err
	}

	// Write GameOpenedTime
	err = w.WriteInt32(p.GameOpenedTime)
	if err != nil {
		return err
	}

	// Write ServerVersion
	err = w.WriteString(p.ServerVersion)
	if err != nil {
		return err
	}

	// Write BGColor
	err = w.WriteInt32(p.BGColor)
	if err != nil {
		return err
	}

	// Write ViewRadius
	err = w.WriteByte(p.ViewRadius)
	if err != nil {
		return err
	}

	// Write DungeonModifiers
	err = w.WriteString(p.DungeonModifiers)
	if err != nil {
		return err
	}

	// Write DungeonModifiers2
	err = w.WriteString(p.DungeonModifiers2)
	if err != nil {
		return err
	}

	// Write DungeonModifiers3
	err = w.WriteString(p.DungeonModifiers3)
	if err != nil {
		return err
	}

	// Write Unknown
	err = w.WriteInt16(p.Unknown)
	if err != nil {
		return err
	}

	// Write MaxRealmScore
	err = w.WriteInt32(p.MaxRealmScore)
	if err != nil {
		return err
	}

	// Write CurrentRealmScore
	err = w.WriteInt32(p.CurrentRealmScore)
	if err != nil {
		return err
	}

	return nil
}
