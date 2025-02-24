package server

import (
	"gorelay/pkg/packets/interfaces"
)

// AccountList represents a server-side account list packet
type AccountList struct {
	AccountListId int32
	AccountIds    []string
	LockAction    int32
}

// Type returns the packet type for AccountList
func (p *AccountList) Type() interfaces.PacketType {
	return interfaces.AccountList
}

// ID returns the packet ID
func (p *AccountList) ID() int32 {
	return int32(interfaces.AccountList)
}

// Read reads the packet data from the given reader
func (p *AccountList) Read(r interfaces.Reader) error {
	var err error
	p.AccountListId, err = r.ReadInt32()
	if err != nil {
		return err
	}

	length, err := r.ReadInt16()
	if err != nil {
		return err
	}

	p.AccountIds = make([]string, length)
	for i := 0; i < int(length); i++ {
		p.AccountIds[i], err = r.ReadString()
		if err != nil {
			return err
		}
	}

	p.LockAction, err = r.ReadInt32()
	return err
}

// Write writes the packet data to the given writer
func (p *AccountList) Write(w interfaces.Writer) error {
	if err := w.WriteInt32(p.AccountListId); err != nil {
		return err
	}

	if err := w.WriteInt16(int16(len(p.AccountIds))); err != nil {
		return err
	}

	for _, accountId := range p.AccountIds {
		if err := w.WriteString(accountId); err != nil {
			return err
		}
	}

	return w.WriteInt32(p.LockAction)
}

// String returns a string representation of the packet
func (p *AccountList) String() string {
	return "AccountList"
}

// HasNulls checks if any fields in the packet are null
func (p *AccountList) HasNulls() bool {
	return false
}

// Structure returns a string representation of the packet structure
func (p *AccountList) Structure() string {
	return "AccountList"
}
