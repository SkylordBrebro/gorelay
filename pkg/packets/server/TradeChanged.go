package server

import (
	"gorelay/pkg/packets/interfaces"
)

// TradeChanged represents the server packet for trade changes
type TradeChanged struct {
	Offers []bool
}

// Type returns the packet type for TradeChanged
func (p *TradeChanged) Type() interfaces.PacketType {
	return interfaces.TradeChanged
}

// Read reads the packet data from the provided reader
func (p *TradeChanged) Read(r interfaces.Reader) error {
	var err error

	// Read Offers array length
	offersLength, err := r.ReadInt16()
	if err != nil {
		return err
	}

	// Read Offers array
	p.Offers = make([]bool, offersLength)
	for i := 0; i < int(offersLength); i++ {
		p.Offers[i], err = r.ReadBool()
		if err != nil {
			return err
		}
	}

	return nil
}

// Write writes the packet data to the provided writer
func (p *TradeChanged) Write(w interfaces.Writer) error {
	var err error

	// Write Offers array length
	err = w.WriteInt16(int16(len(p.Offers)))
	if err != nil {
		return err
	}

	// Write Offers array
	for _, offer := range p.Offers {
		err = w.WriteBool(offer)
		if err != nil {
			return err
		}
	}

	return nil
}

func (p *TradeChanged) ID() int32 {
	return int32(interfaces.TradeChanged)
}