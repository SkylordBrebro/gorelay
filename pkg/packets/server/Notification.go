package server

import (
	"fmt"
	"gorelay/pkg/packets/interfaces"
)

// Notification represents the server packet for notifications
type Notification struct {
	// Notification types
	NotificationType    byte
	NotificationUnknown byte
	Message             string
	ObjectId            int32
	Color               int32
}

// Notification type constants
const (
	NotificationTypePlayer          = 0
	NotificationTypeSystem          = 1
	NotificationTypeError           = 2
	NotificationTypeSticky          = 3
	NotificationTypeGlobal          = 4
	NotificationTypeRealmQueue      = 5
	NotificationTypeObject          = 6
	NotificationTypePlayerDeath     = 7
	NotificationTypePortalOpened    = 8
	NotificationTypeBlueprintUnlock = 20
	NotificationTypeWithIcon        = 21
	NotificationTypeFameBonus       = 22
	NotificationTypeForgeFire       = 23
)

// Type returns the packet type for Notification
func (p *Notification) Type() interfaces.PacketType {
	return interfaces.Notification
}

// Read reads the packet data from the provided reader
func (p *Notification) Read(r interfaces.Reader) error {
	var err error

	// Read NotificationType
	p.NotificationType, err = r.ReadByte()
	if err != nil {
		return err
	}

	// Read NotificationUnknown
	p.NotificationUnknown, err = r.ReadByte()
	if err != nil {
		return err
	}

	// Only read additional fields if NotificationType is Object (6)
	if p.NotificationType == NotificationTypeObject {
		// Read Message
		p.Message, err = r.ReadString()
		if err != nil {
			return err
		}

		// Read ObjectId
		p.ObjectId, err = r.ReadInt32()
		if err != nil {
			return err
		}

		// Read Color
		p.Color, err = r.ReadInt32()
		if err != nil {
			return err
		}
	}

	return nil
}

// Write writes the packet data to the provided writer
func (p *Notification) Write(w interfaces.Writer) error {
	var err error

	// Write NotificationType
	err = w.WriteByte(p.NotificationType)
	if err != nil {
		return err
	}

	// Write NotificationUnknown
	err = w.WriteByte(p.NotificationUnknown)
	if err != nil {
		return err
	}

	// Only write additional fields if NotificationType is Object (6)
	if p.NotificationType == NotificationTypeObject {
		// Write Message
		err = w.WriteString(p.Message)
		if err != nil {
			return err
		}

		// Write ObjectId
		err = w.WriteInt32(p.ObjectId)
		if err != nil {
			return err
		}

		// Write Color
		err = w.WriteInt32(p.Color)
		if err != nil {
			return err
		}
	}

	return nil
}

// CreateNotification creates a new notification with the given object ID and message
func CreateNotification(objectId int32, message string) *Notification {
	return CreateNotificationWithColor(objectId, 65535, message)
}

// CreateNotificationWithColor creates a new notification with the given object ID, color, and message
func CreateNotificationWithColor(objectId int32, color int32, message string) *Notification {
	notification := &Notification{
		NotificationType: NotificationTypeObject,
		ObjectId:         objectId,
		Message:          fmt.Sprintf("{\"key\":\"server.plus_symbol\",\"tokens\":{\"amount\":\"%s\"}}", message),
		Color:            color,
	}
	return notification
}

func (p *Notification) ID() int32 {
	return int32(interfaces.Notification)
}