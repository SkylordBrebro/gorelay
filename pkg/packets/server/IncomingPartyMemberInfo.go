package server

import (
	"gorelay/pkg/packets/interfaces"
)

// PartyPlayer represents a player in a party
type PartyPlayer struct {
	// Add fields as needed based on the implementation
	// Since the original doesn't show the fields, I'm making an assumption
	Name      string
	ObjectId  int32
	ClassType int32
}

// NewPartyPlayer creates a new PartyPlayer by reading from a packet reader
func NewPartyPlayer(r interfaces.Reader) (PartyPlayer, error) {
	var player PartyPlayer
	var err error

	player.Name, err = r.ReadString()
	if err != nil {
		return player, err
	}

	player.ObjectId, err = r.ReadInt32()
	if err != nil {
		return player, err
	}

	player.ClassType, err = r.ReadInt32()
	return player, err
}

// Write writes the party player data to a packet writer
func (p *PartyPlayer) Write(w interfaces.Writer) error {
	var err error

	err = w.WriteString(p.Name)
	if err != nil {
		return err
	}

	err = w.WriteInt32(p.ObjectId)
	if err != nil {
		return err
	}

	err = w.WriteInt32(p.ClassType)
	return err
}

// IncomingPartyMemberInfo represents the server packet for party member information
type IncomingPartyMemberInfo struct {
	PartyId     uint32
	Unknown2    uint16
	MaxSize     byte
	PartyPlayer []PartyPlayer
	Description string
}

// Type returns the packet type for IncomingPartyMemberInfo
func (p *IncomingPartyMemberInfo) Type() interfaces.PacketType {
	return interfaces.IncomingPartyMemberInfo
}

// Read reads the packet data from the provided reader
func (p *IncomingPartyMemberInfo) Read(r interfaces.Reader) error {
	var err error

	// Read PartyId
	p.PartyId, err = r.ReadUInt32()
	if err != nil {
		return err
	}

	// Read Unknown2
	p.Unknown2, err = r.ReadUInt16()
	if err != nil {
		return err
	}

	// Read MaxSize
	p.MaxSize, err = r.ReadByte()
	if err != nil {
		return err
	}

	// Read PartyPlayer array length
	playerCount, err := r.ReadInt16()
	if err != nil {
		return err
	}

	// Read PartyPlayer array
	p.PartyPlayer = make([]PartyPlayer, playerCount)
	for i := 0; i < int(playerCount); i++ {
		p.PartyPlayer[i], err = NewPartyPlayer(r)
		if err != nil {
			return err
		}
	}

	// Read Description
	p.Description, err = r.ReadString()
	if err != nil {
		return err
	}

	return nil
}

// Write writes the packet data to the provided writer
func (p *IncomingPartyMemberInfo) Write(w interfaces.Writer) error {
	var err error

	// Write PartyId
	err = w.WriteUInt32(p.PartyId)
	if err != nil {
		return err
	}

	// Write Unknown2
	err = w.WriteUInt16(p.Unknown2)
	if err != nil {
		return err
	}

	// Write MaxSize
	err = w.WriteByte(p.MaxSize)
	if err != nil {
		return err
	}

	// Write PartyPlayer array length
	err = w.WriteInt16(int16(len(p.PartyPlayer)))
	if err != nil {
		return err
	}

	// Write PartyPlayer array
	for _, player := range p.PartyPlayer {
		err = player.Write(w)
		if err != nil {
			return err
		}
	}

	// Write Description
	err = w.WriteString(p.Description)
	if err != nil {
		return err
	}

	return nil
}
