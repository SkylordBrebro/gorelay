// Decompiled with JetBrains decompiler
// Type: prankster.Proxy.Networking.Packets.Client.QuestRedeem
// Assembly: prankster, Version=1.0.0.1, Culture=neutral, PublicKeyToken=null
// MVID: 674C3C29-3FFB-46FB-A4BE-03322F13731C
// Assembly location: \\hv\e$\rotmg\multisource\pranksterREAL-cleaned.exe

package client

import (
	"gorelay/pkg/packets"
	"gorelay/pkg/packets/interfaces"
)

// QuestRedeem represents a packet for redeeming a quest
type QuestRedeem struct {
	*packets.BasePacket
	QuestID string
	Slots   []int32
	ItemIDs []int32
}

// NewQuestRedeem creates a new QuestRedeem packet
func NewQuestRedeem() *QuestRedeem {
	return &QuestRedeem{
		BasePacket: packets.NewPacket(interfaces.QuestRedeem, byte(interfaces.QuestRedeem)),
	}
}

// Type returns the packet type
func (q *QuestRedeem) Type() interfaces.PacketType {
	return interfaces.QuestRedeem
}

// Read reads the packet data from a PacketReader
func (q *QuestRedeem) Read(r *packets.PacketReader) error {
	var err error
	q.QuestID, err = r.ReadString()
	if err != nil {
		return err
	}

	slotCount, err := r.ReadInt16()
	if err != nil {
		return err
	}

	q.Slots = make([]int32, slotCount)
	q.ItemIDs = make([]int32, slotCount)

	for i := 0; i < int(slotCount); i++ {
		q.Slots[i], err = r.ReadInt32()
		if err != nil {
			return err
		}
		q.ItemIDs[i], err = r.ReadInt32()
		if err != nil {
			return err
		}
	}

	return nil
}

// Write writes the packet data to a PacketWriter
func (q *QuestRedeem) Write(w *packets.PacketWriter) error {
	if err := w.WriteString(q.QuestID); err != nil {
		return err
	}

	if err := w.WriteInt16(int16(len(q.Slots))); err != nil {
		return err
	}

	for i := 0; i < len(q.Slots); i++ {
		if err := w.WriteInt32(q.Slots[i]); err != nil {
			return err
		}
		if err := w.WriteInt32(q.ItemIDs[i]); err != nil {
			return err
		}
	}

	return nil
}
