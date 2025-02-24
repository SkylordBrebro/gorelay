package events

import (
	"gorelay/pkg/packets"
)

// EventType represents different types of events
type EventType int

const (
	// Connection events
	EventConnect EventType = iota
	EventDisconnect
	EventReconnect
	EventServerSwitch

	// Game state events
	EventTick
	EventUpdate
	EventNewTick
	EventShowEffect
	EventDeath
	EventReconnectRequest
	EventGotoAck
	EventAoeAck
	EventShootAck
	EventPlayerHit
	EventEnemyHit
	EventOtherHit

	// Player events
	EventPlayerShoot
	EventPlayerMove
	EventPlayerTeleport
	EventPlayerDamage
	EventPlayerHeal
	EventPlayerChat
	EventPlayerText

	// Enemy events
	EventEnemyShoot
	EventEnemyMove
	EventEnemyDeath
	EventNewEnemy
	EventEnemyUpdate

	// Projectile events
	EventProjectileSpawn
	EventProjectileDestroy
	EventProjectileHit

	// Item events
	EventInventoryUpdate
	EventItemDrop
	EventItemPickup

	// Map events
	EventMapInfo
	EventTileUpdate
)

// Event represents an event in the game
type Event struct {
	Type   EventType
	Client interface{}
	Packet packets.Packet
	Data   interface{}
}

// EventEmitter handles event dispatching
type EventEmitter struct {
	handlers map[EventType][]func(*Event)
}

// NewEventEmitter creates a new event emitter
func NewEventEmitter() *EventEmitter {
	return &EventEmitter{
		handlers: make(map[EventType][]func(*Event)),
	}
}

// On registers a handler for an event type
func (e *EventEmitter) On(eventType EventType, handler func(*Event)) {
	if _, exists := e.handlers[eventType]; !exists {
		e.handlers[eventType] = make([]func(*Event), 0)
	}
	e.handlers[eventType] = append(e.handlers[eventType], handler)
}

// Off removes a handler for an event type
func (e *EventEmitter) Off(eventType EventType, handler func(*Event)) {
	if handlers, exists := e.handlers[eventType]; exists {
		for i, h := range handlers {
			if &h == &handler {
				e.handlers[eventType] = append(handlers[:i], handlers[i+1:]...)
				break
			}
		}
	}
}

// Emit dispatches an event to all registered handlers
func (e *EventEmitter) Emit(event *Event) {
	if handlers, exists := e.handlers[event.Type]; exists {
		for _, handler := range handlers {
			handler(event)
		}
	}
}

// Event data structures
type PlayerEventData struct {
	PlayerData interface{} // Will be replaced with proper PlayerData type
	Position   interface{} // Will be replaced with proper Position type
}

type EnemyEventData struct {
	Enemy    interface{} // Will be replaced with proper Enemy type
	Position interface{} // Will be replaced with proper Position type
}

type ProjectileEventData struct {
	OwnerID      int32
	ProjectileID int32
	Position     interface{} // Will be replaced with proper Position type
	Damage       int32
}

type MapEventData struct {
	Width  int32
	Height int32
	Name   string
	Seed   int32
}

type ItemEventData struct {
	ItemID   int32
	SlotID   int32
	Position interface{} // Will be replaced with proper Position type
}
