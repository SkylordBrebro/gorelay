package events

import (
	"gorelay/pkg/packets"
)

// EventType represents different types of game events
type EventType int

const (
	// Player events
	EventPlayerJoin EventType = iota
	EventPlayerLeave
	EventPlayerMove
	EventPlayerShoot
	EventPlayerHit
	EventPlayerDeath

	// Enemy events
	EventEnemySpawn
	EventEnemyDeath
	EventEnemyShoot

	// Game events
	EventMapChange
	EventTick
	EventChat
)

// Event represents a game event
type Event struct {
	Type   EventType
	Client interface{} // Using interface{} to avoid import cycle
	Packet packets.Packet
	Data   interface{}
}

// EventHandler is a function that handles game events
type EventHandler func(event *Event)

// EventEmitter manages event subscriptions and dispatching
type EventEmitter struct {
	handlers map[EventType][]EventHandler
}

// NewEventEmitter creates a new event emitter
func NewEventEmitter() *EventEmitter {
	return &EventEmitter{
		handlers: make(map[EventType][]EventHandler),
	}
}

// On subscribes to an event type
func (e *EventEmitter) On(eventType EventType, handler EventHandler) {
	e.handlers[eventType] = append(e.handlers[eventType], handler)
}

// Off unsubscribes from an event type
func (e *EventEmitter) Off(eventType EventType, handler EventHandler) {
	if handlers, ok := e.handlers[eventType]; ok {
		for i, h := range handlers {
			if &h == &handler {
				e.handlers[eventType] = append(handlers[:i], handlers[i+1:]...)
				break
			}
		}
	}
}

// Emit dispatches an event to all subscribed handlers
func (e *EventEmitter) Emit(event *Event) {
	if handlers, ok := e.handlers[event.Type]; ok {
		for _, handler := range handlers {
			handler(event)
		}
	}
}

// PlayerEventData contains data specific to player events
type PlayerEventData struct {
	PlayerData *packets.PlayerData
	Position   *packets.WorldPosData
}

// EnemyEventData contains data specific to enemy events
type EnemyEventData struct {
	ObjectID   int32
	ObjectType int32
	Position   *packets.WorldPosData
}

// ChatEventData contains data specific to chat events
type ChatEventData struct {
	Name      string
	Message   string
	Recipient string
}
