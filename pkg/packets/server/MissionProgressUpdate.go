package server

import (
	"gorelay/pkg/packets/interfaces"
)

// MissionProgressUpdate represents the server packet for mission progress updates
type MissionProgressUpdate struct {
	MissionId  int32
	Progress   int32
	Objectives []int32
}

// Type returns the packet type for MissionProgressUpdate
func (p *MissionProgressUpdate) Type() interfaces.PacketType {
	return interfaces.MissionProgressUpdate
}

// Read reads the packet data from the provided reader
func (p *MissionProgressUpdate) Read(r interfaces.Reader) error {
	var err error

	// Read MissionId
	p.MissionId, err = r.ReadInt32()
	if err != nil {
		return err
	}

	// Read Progress
	p.Progress, err = r.ReadInt32()
	if err != nil {
		return err
	}

	// Read Objectives array length
	objectivesCount, err := r.ReadInt16()
	if err != nil {
		return err
	}

	// Read Objectives array
	p.Objectives = make([]int32, objectivesCount)
	for i := 0; i < int(objectivesCount); i++ {
		p.Objectives[i], err = r.ReadInt32()
		if err != nil {
			return err
		}
	}

	return nil
}

// Write writes the packet data to the provided writer
func (p *MissionProgressUpdate) Write(w interfaces.Writer) error {
	var err error

	// Write MissionId
	err = w.WriteInt32(p.MissionId)
	if err != nil {
		return err
	}

	// Write Progress
	err = w.WriteInt32(p.Progress)
	if err != nil {
		return err
	}

	// Write Objectives array length
	err = w.WriteInt16(int16(len(p.Objectives)))
	if err != nil {
		return err
	}

	// Write Objectives array
	for _, objective := range p.Objectives {
		err = w.WriteInt32(objective)
		if err != nil {
			return err
		}
	}

	return nil
}
