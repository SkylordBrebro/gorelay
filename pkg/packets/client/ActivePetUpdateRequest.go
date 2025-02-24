package client

import (
	"gorelay/pkg/packets"
	"gorelay/pkg/packets/interfaces"
)

// Pet command constants
const (
	FollowPet   = 1
	UnfollowPet = 2
	ReleasePet  = 3
)

// ActivePetUpdateRequest represents a pet update request packet
type ActivePetUpdateRequest struct {
	*packets.BasePacket
	CommandID byte
	PetID     uint32
}

// NewActivePetUpdateRequest creates a new ActivePetUpdateRequest packet
func NewActivePetUpdateRequest() *ActivePetUpdateRequest {
	return &ActivePetUpdateRequest{
		BasePacket: packets.NewPacket(interfaces.ActivePetUpdateRequest, byte(interfaces.ActivePetUpdateRequest)),
	}
}

// Type returns the packet type
func (p *ActivePetUpdateRequest) Type() interfaces.PacketType {
	return interfaces.ActivePetUpdateRequest
}

// Read reads the packet data from the reader
func (p *ActivePetUpdateRequest) Read(r *packets.PacketReader) error {
	var err error
	p.CommandID, err = r.ReadByte()
	if err != nil {
		return err
	}

	p.PetID, err = r.ReadUInt32()
	return err
}

// Write writes the packet data to the writer
func (p *ActivePetUpdateRequest) Write(w *packets.PacketWriter) error {
	if err := w.WriteByte(p.CommandID); err != nil {
		return err
	}
	return w.WriteUInt32(p.PetID)
}
