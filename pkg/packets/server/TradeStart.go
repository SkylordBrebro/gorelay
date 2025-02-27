package server

import (
	"gorelay/pkg/packets/dataobjects"
	"gorelay/pkg/packets/interfaces"
)

// TradeStart represents the server packet for starting a trade
type TradeStart struct {
	MyItems         []*dataobjects.Item
	YourName        string
	YourItems       []*dataobjects.Item
	PartnerObjectId int32
}

// Type returns the packet type for TradeStart
func (p *TradeStart) Type() interfaces.PacketType {
	return interfaces.TradeStart
}

// Read reads the packet data from the provided reader
func (p *TradeStart) Read(r interfaces.Reader) error {
	var err error

	// Read MyItems array length
	myItemsCount, err := r.ReadInt16()
	if err != nil {
		return err
	}

	// Read MyItems array
	p.MyItems = make([]*dataobjects.Item, myItemsCount)
	for i := 0; i < int(myItemsCount); i++ {
		p.MyItems[i] = dataobjects.NewItem()
		err = p.MyItems[i].Read(r)
		if err != nil {
			return err
		}
	}

	// Read YourName
	p.YourName, err = r.ReadString()
	if err != nil {
		return err
	}

	// Read YourItems array length
	yourItemsCount, err := r.ReadInt16()
	if err != nil {
		return err
	}

	// Read YourItems array
	p.YourItems = make([]*dataobjects.Item, yourItemsCount)
	for i := 0; i < int(yourItemsCount); i++ {
		p.YourItems[i] = dataobjects.NewItem()
		err = p.YourItems[i].Read(r)
		if err != nil {
			return err
		}
	}

	// Read PartnerObjectId
	p.PartnerObjectId, err = r.ReadInt32()
	if err != nil {
		return err
	}

	return nil
}

// Write writes the packet data to the provided writer
func (p *TradeStart) Write(w interfaces.Writer) error {
	var err error

	// Write MyItems array length
	err = w.WriteInt16(int16(len(p.MyItems)))
	if err != nil {
		return err
	}

	// Write MyItems array
	for _, item := range p.MyItems {
		err = item.Write(w)
		if err != nil {
			return err
		}
	}

	// Write YourName
	err = w.WriteString(p.YourName)
	if err != nil {
		return err
	}

	// Write YourItems array length
	err = w.WriteInt16(int16(len(p.YourItems)))
	if err != nil {
		return err
	}

	// Write YourItems array
	for _, item := range p.YourItems {
		err = item.Write(w)
		if err != nil {
			return err
		}
	}

	// Write PartnerObjectId
	err = w.WriteInt32(p.PartnerObjectId)
	if err != nil {
		return err
	}

	return nil
}

func (p *TradeStart) ID() int32 {
	return int32(interfaces.TradeStart)
}