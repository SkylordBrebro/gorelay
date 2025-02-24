package client

import (
	"gorelay/pkg/packets/dataobjects"
	"gorelay/pkg/packets/interfaces"
)

// StartUse represents a client-side start use packet
type StartUse struct {
	Time  int32
	Start *dataobjects.Location
	End   *dataobjects.Location
}

// Type returns the packet type for StartUse
func (p *StartUse) Type() interfaces.PacketType {
	return interfaces.StartUse
}

// Read reads the packet data from the given reader
func (p *StartUse) Read(r interfaces.Reader) error {
	var err error
	p.Time, err = r.ReadInt32()
	if err != nil {
		return err
	}

	p.Start = dataobjects.NewLocation()
	if err := p.Start.Read(r); err != nil {
		return err
	}

	p.End = dataobjects.NewLocation()
	return p.End.Read(r)
}

// Write writes the packet data to the given writer
func (p *StartUse) Write(w interfaces.Writer) error {
	if err := w.WriteInt32(p.Time); err != nil {
		return err
	}
	if err := p.Start.Write(w); err != nil {
		return err
	}
	return p.End.Write(w)
}
