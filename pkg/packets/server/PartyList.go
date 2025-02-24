package server

import (
	"gorelay/pkg/packets/dataobjects"
	"gorelay/pkg/packets/interfaces"
)

// PartyList represents the server packet for party listings
type PartyList struct {
	PacketNumber byte
	Parties      []*dataobjects.PartyInfo
}

// Type returns the packet type for PartyList
func (p *PartyList) Type() interfaces.PacketType {
	return interfaces.PartyList
}

// Read reads the packet data from the provided reader
func (p *PartyList) Read(r interfaces.Reader) error {
	var err error

	// Read PacketNumber
	p.PacketNumber, err = r.ReadByte()
	if err != nil {
		return err
	}

	// Read Parties array length
	partiesCount, err := r.ReadInt16()
	if err != nil {
		return err
	}

	// Read Parties array
	p.Parties = make([]*dataobjects.PartyInfo, partiesCount)
	for i := 0; i < int(partiesCount); i++ {
		p.Parties[i] = dataobjects.NewPartyInfo()
		err = p.Parties[i].Read(r)
		if err != nil {
			return err
		}
	}

	return nil
}

// Write writes the packet data to the provided writer
func (p *PartyList) Write(w interfaces.Writer) error {
	var err error

	// Write PacketNumber
	err = w.WriteByte(p.PacketNumber)
	if err != nil {
		return err
	}

	// Write Parties array length
	err = w.WriteInt16(int16(len(p.Parties)))
	if err != nil {
		return err
	}

	// Write Parties array
	for _, party := range p.Parties {
		err = party.Write(w)
		if err != nil {
			return err
		}
	}

	return nil
}
