package server

import (
	"gorelay/pkg/packets/interfaces"
)

// QuestRedeemResponse represents the server packet for quest redeem responses
type QuestRedeemResponse struct {
	Success bool
	Message string
}

// Type returns the packet type for QuestRedeemResponse
func (p *QuestRedeemResponse) Type() interfaces.PacketType {
	return interfaces.QuestRedeemResponse
}

// Read reads the packet data from the provided reader
func (p *QuestRedeemResponse) Read(r interfaces.Reader) error {
	var err error

	// Read Success
	p.Success, err = r.ReadBool()
	if err != nil {
		return err
	}

	// Read Message
	p.Message, err = r.ReadString()
	if err != nil {
		return err
	}

	return nil
}

// Write writes the packet data to the provided writer
func (p *QuestRedeemResponse) Write(w interfaces.Writer) error {
	var err error

	// Write Success
	err = w.WriteBool(p.Success)
	if err != nil {
		return err
	}

	// Write Message
	err = w.WriteString(p.Message)
	if err != nil {
		return err
	}

	return nil
}

func (p *QuestRedeemResponse) ID() int32 {
	return int32(interfaces.QuestRedeemResponse)
}