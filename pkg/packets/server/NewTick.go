package server

import (
	"gorelay/pkg/packets/dataobjects"
	"gorelay/pkg/packets/interfaces"
)

// NewTick represents the server packet for game tick updates
type NewTick struct {
	TickId           int32
	TickTime         int32
	ServerRealTimeMs int32
	ServerLastRTTMS  uint16
	Statuses         []*dataobjects.Status
}

// Type returns the packet type for NewTick
func (p *NewTick) Type() interfaces.PacketType {
	return interfaces.NewTick
}

// Read reads the packet data from the provided reader
func (p *NewTick) Read(r interfaces.Reader) error {
	var err error

	// Read TickId
	p.TickId, err = r.ReadInt32()
	if err != nil {
		return err
	}

	// Read TickTime
	p.TickTime, err = r.ReadInt32()
	if err != nil {
		return err
	}

	// Read ServerRealTimeMs
	p.ServerRealTimeMs, err = r.ReadInt32()
	if err != nil {
		return err
	}

	// Read ServerLastRTTMS
	p.ServerLastRTTMS, err = r.ReadUInt16()
	if err != nil {
		return err
	}

	// Read Statuses array length
	statusCount, err := r.ReadInt16()
	if err != nil {
		return err
	}

	// Read Statuses array
	p.Statuses = make([]*dataobjects.Status, statusCount)
	for i := 0; i < int(statusCount); i++ {
		p.Statuses[i] = dataobjects.NewStatus()
		err = p.Statuses[i].Read(r)
		if err != nil {
			return err
		}
	}

	return nil
}

// Write writes the packet data to the provided writer
func (p *NewTick) Write(w interfaces.Writer) error {
	var err error

	// Write TickId
	err = w.WriteInt32(p.TickId)
	if err != nil {
		return err
	}

	// Write TickTime
	err = w.WriteInt32(p.TickTime)
	if err != nil {
		return err
	}

	// Write ServerRealTimeMs
	err = w.WriteInt32(p.ServerRealTimeMs)
	if err != nil {
		return err
	}

	// Write ServerLastRTTMS
	err = w.WriteUInt16(p.ServerLastRTTMS)
	if err != nil {
		return err
	}

	// Write Statuses array length
	err = w.WriteInt16(int16(len(p.Statuses)))
	if err != nil {
		return err
	}

	// Write Statuses array
	for _, status := range p.Statuses {
		err = status.Write(w)
		if err != nil {
			return err
		}
	}

	return nil
}
