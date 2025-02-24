package server

import (
	"gorelay/pkg/packets/interfaces"
)

// NewCharacterInformation represents the server packet for new character information
type NewCharacterInformation struct {
	CharacterXml string
}

// Type returns the packet type for NewCharacterInformation
func (p *NewCharacterInformation) Type() interfaces.PacketType {
	return interfaces.NewCharacterInformation
}

// Read reads the packet data from the provided reader
func (p *NewCharacterInformation) Read(r interfaces.Reader) error {
	var err error
	p.CharacterXml, err = r.ReadString()
	return err
}

// Write writes the packet data to the provided writer
func (p *NewCharacterInformation) Write(w interfaces.Writer) error {
	return w.WriteString(p.CharacterXml)
}
