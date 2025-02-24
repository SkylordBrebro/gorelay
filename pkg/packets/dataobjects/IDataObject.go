package dataobjects

import (
	"gorelay/pkg/packets/interfaces"
)

// DataObject represents an object that can be read from and written to packets
type DataObject interface {
	Read(r interfaces.Reader) error
	Write(w interfaces.Writer) error
	Clone() DataObject
	String() string
}
