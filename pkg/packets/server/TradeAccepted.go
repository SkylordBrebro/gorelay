package server

import (
	"gorelay/pkg/packets/interfaces"
)

// TradeAccepted represents the server packet for trade accepted
type TradeAccepted struct {
	MyOffers   []bool
	YourOffers []bool
}

// Type returns the packet type for TradeAccepted
func (p *TradeAccepted) Type() interfaces.PacketType {
	return interfaces.TradeAccepted
}

// Read reads the packet data from the provided reader
func (p *TradeAccepted) Read(r interfaces.Reader) error {
	var err error

	// Read MyOffers array length
	myOffersLength, err := r.ReadInt16()
	if err != nil {
		return err
	}

	// Read MyOffers array
	p.MyOffers = make([]bool, myOffersLength)
	for i := 0; i < int(myOffersLength); i++ {
		p.MyOffers[i], err = r.ReadBool()
		if err != nil {
			return err
		}
	}

	// Read YourOffers array length
	yourOffersLength, err := r.ReadInt16()
	if err != nil {
		return err
	}

	// Read YourOffers array
	p.YourOffers = make([]bool, yourOffersLength)
	for i := 0; i < int(yourOffersLength); i++ {
		p.YourOffers[i], err = r.ReadBool()
		if err != nil {
			return err
		}
	}

	return nil
}

// Write writes the packet data to the provided writer
func (p *TradeAccepted) Write(w interfaces.Writer) error {
	var err error

	// Write MyOffers array length
	err = w.WriteInt16(int16(len(p.MyOffers)))
	if err != nil {
		return err
	}

	// Write MyOffers array
	for _, offer := range p.MyOffers {
		err = w.WriteBool(offer)
		if err != nil {
			return err
		}
	}

	// Write YourOffers array length
	err = w.WriteInt16(int16(len(p.YourOffers)))
	if err != nil {
		return err
	}

	// Write YourOffers array
	for _, offer := range p.YourOffers {
		err = w.WriteBool(offer)
		if err != nil {
			return err
		}
	}

	return nil
}
