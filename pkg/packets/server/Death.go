package server

import (
	"gorelay/pkg/packets/interfaces"
)

// Death represents the server packet for player death information
type Death struct {
	AccountId string
	CharId    int
	KilledBy  string
	Unknown1  int32
	Unknown2  int16
	Stats     string
}

// Type returns the packet type for Death
func (p *Death) Type() interfaces.PacketType {
	return interfaces.Death
}

// Read reads the packet data from the provided reader
func (p *Death) Read(r interfaces.Reader) error {
	var err error

	// Read AccountId
	p.AccountId, err = r.ReadString()
	if err != nil {
		return err
	}

	// Read CharId
	p.CharId, err = r.ReadCompressedInt()
	if err != nil {
		return err
	}

	// Read KilledBy
	p.KilledBy, err = r.ReadString()
	if err != nil {
		return err
	}

	// Read Unknown1
	p.Unknown1, err = r.ReadInt32()
	if err != nil {
		return err
	}

	// Read Unknown2
	p.Unknown2, err = r.ReadInt16()
	if err != nil {
		return err
	}

	// Read Stats
	p.Stats, err = r.ReadString()
	if err != nil {
		return err
	}

	return nil
}

// Write writes the packet data to the provided writer
func (p *Death) Write(w interfaces.Writer) error {
	var err error

	// Write AccountId
	err = w.WriteString(p.AccountId)
	if err != nil {
		return err
	}

	// Write CharId
	err = w.WriteCompressedInt(p.CharId)
	if err != nil {
		return err
	}

	// Write KilledBy
	err = w.WriteString(p.KilledBy)
	if err != nil {
		return err
	}

	// Write Unknown1
	err = w.WriteInt32(p.Unknown1)
	if err != nil {
		return err
	}

	// Write Unknown2
	err = w.WriteInt16(p.Unknown2)
	if err != nil {
		return err
	}

	// Write Stats
	err = w.WriteString(p.Stats)
	if err != nil {
		return err
	}

	return nil
}
