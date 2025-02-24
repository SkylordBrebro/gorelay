package server

import (
	"gorelay/pkg/packets/interfaces"
)

// ResultType represents the possible result types for a buy operation
type ResultType int32

const (
	UnknownError       ResultType = -1
	Success            ResultType = 0
	InvalidCharacter   ResultType = 1
	ItemNotFound       ResultType = 2
	NotEnoughGold      ResultType = 3
	InventoryFull      ResultType = 4
	TooLowRank         ResultType = 5
	NotEnoughFame      ResultType = 6
	PetFeedSuccess     ResultType = 7
	TooManyResetsToday ResultType = 10
)

// BuyResult represents a server-side buy result packet
type BuyResult struct {
	Result  int32
	Message string
}

// Type returns the packet type for BuyResult
func (p *BuyResult) Type() interfaces.PacketType {
	return interfaces.BuyResult
}

// ID returns the packet ID
func (p *BuyResult) ID() int32 {
	return int32(interfaces.BuyResult)
}

// Read reads the packet data from the given reader
func (p *BuyResult) Read(r interfaces.Reader) error {
	var err error
	p.Result, err = r.ReadInt32()
	if err != nil {
		return err
	}

	p.Message, err = r.ReadString()
	return err
}

// Write writes the packet data to the given writer
func (p *BuyResult) Write(w interfaces.Writer) error {
	if err := w.WriteInt32(p.Result); err != nil {
		return err
	}

	return w.WriteString(p.Message)
}

// String returns a string representation of the packet
func (p *BuyResult) String() string {
	return "BuyResult"
}

// HasNulls checks if any fields in the packet are null
func (p *BuyResult) HasNulls() bool {
	return false
}

// Structure returns a string representation of the packet structure
func (p *BuyResult) Structure() string {
	return "BuyResult"
}
