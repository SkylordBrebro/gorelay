package server

import (
	"gorelay/pkg/packets/interfaces"
)

// KeyInfoResponse represents the server packet for key information
type KeyInfoResponse struct {
	Name        string
	Description string
	Creator     string
}

// Type returns the packet type for KeyInfoResponse
func (p *KeyInfoResponse) Type() interfaces.PacketType {
	return interfaces.KeyInfoResponse
}

// Read reads the packet data from the provided reader
func (p *KeyInfoResponse) Read(r interfaces.Reader) error {
	var err error

	// Read Name
	p.Name, err = r.ReadUTF32String()
	if err != nil {
		return err
	}

	// Read Description
	p.Description, err = r.ReadUTF32String()
	if err != nil {
		return err
	}

	// Read Creator
	p.Creator, err = r.ReadUTF32String()
	if err != nil {
		return err
	}

	return nil
}

// Write writes the packet data to the provided writer
func (p *KeyInfoResponse) Write(w interfaces.Writer) error {
	var err error

	// Write Name
	err = w.WriteUTF32String(p.Name)
	if err != nil {
		return err
	}

	// Write Description
	err = w.WriteUTF32String(p.Description)
	if err != nil {
		return err
	}

	// Write Creator
	err = w.WriteUTF32String(p.Creator)
	if err != nil {
		return err
	}

	return nil
}

func (p *KeyInfoResponse) ID() int32 {
	return int32(interfaces.KeyInfoResponse)
}