package client

import (
	"gorelay/pkg/packets"
	"gorelay/pkg/packets/interfaces"
)

// AcceptTrade represents a trade acceptance packet
type AcceptTrade struct {
	*packets.BasePacket
	MyOffers   []bool
	YourOffers []bool
}

// NewAcceptTrade creates a new AcceptTrade packet
func NewAcceptTrade() *AcceptTrade {
	return &AcceptTrade{
		BasePacket: packets.NewPacket(interfaces.AcceptTrade, byte(interfaces.AcceptTrade)),
	}
}

// Type returns the packet type
func (p *AcceptTrade) Type() interfaces.PacketType {
	return interfaces.AcceptTrade
}

// Read reads the packet data from the reader
func (p *AcceptTrade) Read(r interfaces.Reader) error {
	// Read MyOffers
	myOffersLen, err := r.ReadInt16()
	if err != nil {
		return err
	}

	p.MyOffers = make([]bool, myOffersLen)
	for i := 0; i < int(myOffersLen); i++ {
		b, err := r.ReadByte()
		if err != nil {
			return err
		}
		p.MyOffers[i] = b != 0
	}

	// Read YourOffers
	yourOffersLen, err := r.ReadInt16()
	if err != nil {
		return err
	}

	p.YourOffers = make([]bool, yourOffersLen)
	for i := 0; i < int(yourOffersLen); i++ {
		b, err := r.ReadByte()
		if err != nil {
			return err
		}
		p.YourOffers[i] = b != 0
	}

	return nil
}

// Write writes the packet data to the writer
func (p *AcceptTrade) Write(w interfaces.Writer) error {
	// Write MyOffers
	if err := w.WriteInt16(int16(len(p.MyOffers))); err != nil {
		return err
	}
	for _, offer := range p.MyOffers {
		if err := w.WriteByte(boolToByte(offer)); err != nil {
			return err
		}
	}

	// Write YourOffers
	if err := w.WriteInt16(int16(len(p.YourOffers))); err != nil {
		return err
	}
	for _, offer := range p.YourOffers {
		if err := w.WriteByte(boolToByte(offer)); err != nil {
			return err
		}
	}

	return nil
}

func boolToByte(b bool) byte {
	if b {
		return 1
	}
	return 0
}
