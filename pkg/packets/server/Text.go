package server

import (
	"gorelay/pkg/packets/interfaces"
)

// Text represents the server packet for text messages
type Text struct {
	Name        string
	ObjectId    int32
	NumStars    int16
	BubbleTime  byte
	Recipient   string
	RawText     string
	CleanText   string
	IsSupporter bool
	StarBg      int32
}

// Type returns the packet type for Text
func (p *Text) Type() interfaces.PacketType {
	return interfaces.Text
}

// Read reads the packet data from the provided reader
func (p *Text) Read(r interfaces.Reader) error {
	var err error

	// Read Name
	p.Name, err = r.ReadString()
	if err != nil {
		return err
	}

	// Read ObjectId
	p.ObjectId, err = r.ReadInt32()
	if err != nil {
		return err
	}

	// Read NumStars
	p.NumStars, err = r.ReadInt16()
	if err != nil {
		return err
	}

	// Read BubbleTime
	p.BubbleTime, err = r.ReadByte()
	if err != nil {
		return err
	}

	// Read Recipient
	p.Recipient, err = r.ReadString()
	if err != nil {
		return err
	}

	// Read RawText
	p.RawText, err = r.ReadString()
	if err != nil {
		return err
	}

	// Read CleanText
	p.CleanText, err = r.ReadString()
	if err != nil {
		return err
	}

	// Read IsSupporter
	p.IsSupporter, err = r.ReadBool()
	if err != nil {
		return err
	}

	// Read StarBg
	p.StarBg, err = r.ReadInt32()
	if err != nil {
		return err
	}

	return nil
}

// Write writes the packet data to the provided writer
func (p *Text) Write(w interfaces.Writer) error {
	var err error

	// Write Name
	err = w.WriteString(p.Name)
	if err != nil {
		return err
	}

	// Write ObjectId
	err = w.WriteInt32(p.ObjectId)
	if err != nil {
		return err
	}

	// Write NumStars
	err = w.WriteInt16(p.NumStars)
	if err != nil {
		return err
	}

	// Write BubbleTime
	err = w.WriteByte(p.BubbleTime)
	if err != nil {
		return err
	}

	// Write Recipient
	err = w.WriteString(p.Recipient)
	if err != nil {
		return err
	}

	// Write RawText
	err = w.WriteString(p.RawText)
	if err != nil {
		return err
	}

	// Write CleanText
	err = w.WriteString(p.CleanText)
	if err != nil {
		return err
	}

	// Write IsSupporter
	err = w.WriteBool(p.IsSupporter)
	if err != nil {
		return err
	}

	// Write StarBg
	err = w.WriteInt32(p.StarBg)
	if err != nil {
		return err
	}

	return nil
}

// CreateOryxNotification creates a Text packet for an Oryx notification
func CreateOryxNotification(sender, message string) *Text {
	return &Text{
		BubbleTime: 0,
		CleanText:  message,
		Name:       "#" + sender,
		NumStars:   -1,
		ObjectId:   -1,
		Recipient:  "",
		RawText:    message,
	}
}

// CreateAnnouncement creates a Text packet for an announcement
func CreateAnnouncement(message string) *Text {
	return &Text{
		BubbleTime: 0,
		CleanText:  message,
		Name:       "",
		NumStars:   -1,
		ObjectId:   -1,
		Recipient:  "",
		RawText:    message,
	}
}

// CreateCustomText creates a Text packet with a custom name and message
func CreateCustomText(name, message string) *Text {
	return &Text{
		BubbleTime: 0,
		CleanText:  message,
		Name:       name,
		NumStars:   -1,
		ObjectId:   -1,
		Recipient:  "",
		RawText:    message,
	}
}
