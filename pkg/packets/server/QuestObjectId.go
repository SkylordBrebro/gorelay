package server

import (
	"gorelay/pkg/packets/interfaces"
)

// QuestObjectId represents the server packet for quest object IDs
type QuestObjectId struct {
	ObjectId    int32
	UnknownInts []int
}

// Type returns the packet type for QuestObjectId
func (p *QuestObjectId) Type() interfaces.PacketType {
	return interfaces.QuestObjectId
}

// Read reads the packet data from the provided reader
func (p *QuestObjectId) Read(r interfaces.Reader) error {
	var err error

	// Read ObjectId
	p.ObjectId, err = r.ReadInt32()
	if err != nil {
		return err
	}

	// Read UnknownInts array length
	length, err := r.ReadCompressedInt()
	if err != nil {
		return err
	}

	// Read UnknownInts array
	p.UnknownInts = make([]int, length)
	for i := 0; i < length; i++ {
		p.UnknownInts[i], err = r.ReadCompressedInt()
		if err != nil {
			return err
		}
	}

	return nil
}

// Write writes the packet data to the provided writer
func (p *QuestObjectId) Write(w interfaces.Writer) error {
	var err error

	// Write ObjectId
	err = w.WriteInt32(p.ObjectId)
	if err != nil {
		return err
	}

	// Write UnknownInts array length
	err = w.WriteCompressedInt(len(p.UnknownInts))
	if err != nil {
		return err
	}

	// Write UnknownInts array
	for _, unknownInt := range p.UnknownInts {
		err = w.WriteCompressedInt(unknownInt)
		if err != nil {
			return err
		}
	}

	return nil
}
