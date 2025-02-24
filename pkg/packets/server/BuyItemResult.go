package server

import (
	"gorelay/pkg/packets/interfaces"
)

// BuyItemResult represents a server-side buy item result packet
type BuyItemResult struct {
	Success bool
	Message string
}

// Type returns the packet type for BuyItemResult
func (p *BuyItemResult) Type() interfaces.PacketType {
	return interfaces.BuyItemResult
}

// ID returns the packet ID
func (p *BuyItemResult) ID() int32 {
	return int32(interfaces.BuyItemResult)
}

// Read reads the packet data from the given reader
func (p *BuyItemResult) Read(r interfaces.Reader) error {
	var err error
	p.Success, err = r.ReadBool()
	if err != nil {
		return err
	}

	p.Message, err = r.ReadString()
	return err
}

// Write writes the packet data to the given writer
func (p *BuyItemResult) Write(w interfaces.Writer) error {
	if err := w.WriteBool(p.Success); err != nil {
		return err
	}

	return w.WriteString(p.Message)
}

// String returns a string representation of the packet
func (p *BuyItemResult) String() string {
	return "BuyItemResult"
}

// HasNulls checks if any fields in the packet are null
func (p *BuyItemResult) HasNulls() bool {
	return false
}

// Structure returns a string representation of the packet structure
func (p *BuyItemResult) Structure() string {
	return "BuyItemResult"
}
