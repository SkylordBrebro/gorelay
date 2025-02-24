package server

import (
	"gorelay/pkg/packets/interfaces"
)

// TradeResult represents the result of a trade
type TradeResult int32

// Trade result constants
const (
	TradeResultSuccess TradeResult = iota
	TradeResultCancelled
	TradeResultUnknown
)

// TradeDone represents the server packet for completed trades
type TradeDone struct {
	Code        int32
	Description string
}

// Result returns the trade result as a TradeResult enum
func (p *TradeDone) Result() TradeResult {
	return TradeResult(p.Code)
}

// Type returns the packet type for TradeDone
func (p *TradeDone) Type() interfaces.PacketType {
	return interfaces.TradeDone
}

// Read reads the packet data from the provided reader
func (p *TradeDone) Read(r interfaces.Reader) error {
	var err error

	// Read Code
	p.Code, err = r.ReadInt32()
	if err != nil {
		return err
	}

	// Read Description
	p.Description, err = r.ReadString()
	if err != nil {
		return err
	}

	return nil
}

// Write writes the packet data to the provided writer
func (p *TradeDone) Write(w interfaces.Writer) error {
	var err error

	// Write Code
	err = w.WriteInt32(p.Code)
	if err != nil {
		return err
	}

	// Write Description
	err = w.WriteString(p.Description)
	if err != nil {
		return err
	}

	return nil
}
