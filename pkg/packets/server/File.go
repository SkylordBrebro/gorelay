package server

import (
	"gorelay/pkg/packets/interfaces"
)

// File represents the server packet for file data
type File struct {
	Name     string
	Contents []byte
}

// Type returns the packet type for File
func (p *File) Type() interfaces.PacketType {
	return interfaces.File
}

// Read reads the packet data from the provided reader
func (p *File) Read(r interfaces.Reader) error {
	var err error

	// Read Name
	p.Name, err = r.ReadString()
	if err != nil {
		return err
	}

	// Read Contents length
	contentsLength, err := r.ReadInt32()
	if err != nil {
		return err
	}

	// Read Contents
	p.Contents, err = r.ReadBytes(int(contentsLength))
	if err != nil {
		return err
	}

	return nil
}

// Write writes the packet data to the provided writer
func (p *File) Write(w interfaces.Writer) error {
	var err error

	// Write Name
	err = w.WriteString(p.Name)
	if err != nil {
		return err
	}

	// Write Contents length
	err = w.WriteInt32(int32(len(p.Contents)))
	if err != nil {
		return err
	}

	// Write Contents
	err = w.WriteBytes(p.Contents)
	if err != nil {
		return err
	}

	return nil
}

func (p *File) ID() int32 {
	return int32(interfaces.File)
}