package server

import (
	"gorelay/pkg/packets/interfaces"
)

// MultipleMissionsProgressUpdate represents the server packet for multiple mission progress updates
type MultipleMissionsProgressUpdate struct {
	UnknownString string
}

// Type returns the packet type for MultipleMissionsProgressUpdate
func (p *MultipleMissionsProgressUpdate) Type() interfaces.PacketType {
	return interfaces.MultipleMissionsProgressUpdate
}

// Read reads the packet data from the provided reader
func (p *MultipleMissionsProgressUpdate) Read(r interfaces.Reader) error {
	var err error

	// Read UnknownString
	p.UnknownString, err = r.ReadString()
	if err != nil {
		return err
	}

	return nil
}

// Write writes the packet data to the provided writer
func (p *MultipleMissionsProgressUpdate) Write(w interfaces.Writer) error {
	// Write UnknownString
	return w.WriteString(p.UnknownString)
}

func (p *MultipleMissionsProgressUpdate) ID() int32 {
	return int32(interfaces.MultipleMissionsProgressUpdate)
}