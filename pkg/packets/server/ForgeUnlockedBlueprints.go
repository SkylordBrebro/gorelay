package server

import (
	"gorelay/pkg/packets/interfaces"
)

// ForgeUnlockedBlueprints represents the server packet for unlocked forge blueprints
type ForgeUnlockedBlueprints struct {
	Success bool
	Unknown []byte
}

// Type returns the packet type for ForgeUnlockedBlueprints
func (p *ForgeUnlockedBlueprints) Type() interfaces.PacketType {
	return interfaces.ForgeUnlockedBlueprints
}

// Read reads the packet data from the provided reader
func (p *ForgeUnlockedBlueprints) Read(r interfaces.Reader) error {
	var err error

	// Read Success
	p.Success, err = r.ReadBool()
	if err != nil {
		return err
	}

	// Read remaining bytes as Unknown
	remainingBytes := r.RemainingBytes()
	if remainingBytes > 0 {
		p.Unknown, err = r.ReadBytes(int(remainingBytes))
		if err != nil {
			return err
		}
	}

	return nil
}

// Write writes the packet data to the provided writer
func (p *ForgeUnlockedBlueprints) Write(w interfaces.Writer) error {
	var err error

	// Write Success
	err = w.WriteBool(p.Success)
	if err != nil {
		return err
	}

	// Write Unknown
	if len(p.Unknown) > 0 {
		err = w.WriteBytes(p.Unknown)
		if err != nil {
			return err
		}
	}

	return nil
}

func (p *ForgeUnlockedBlueprints) ID() int32 {
	return int32(interfaces.ForgeUnlockedBlueprints)
}