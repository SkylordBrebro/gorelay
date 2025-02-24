package models

// EventType represents different types of events that can occur
type EventType string

const (
	// Client events
	EventClientConnect      EventType = "client_connect"
	EventClientDisconnect   EventType = "client_disconnect"
	EventClientReady        EventType = "client_ready"
	EventClientArrived      EventType = "client_arrived"
	EventClientConnectError EventType = "client_connect_error"

	// Game events
	EventMapChange    EventType = "map_change"
	EventTick         EventType = "tick"
	EventChat         EventType = "chat"
	EventDeath        EventType = "death"
	EventNewTick      EventType = "new_tick"
	EventShowEffect   EventType = "show_effect"
	EventGoto         EventType = "goto"
	EventUpdate       EventType = "update"
	EventNotification EventType = "notification"

	// Player events
	EventPlayerShoot EventType = "player_shoot"
	EventEnemyHit    EventType = "enemy_hit"
	EventAoeHit      EventType = "aoe_hit"
	EventDamage      EventType = "damage"
	EventStatChange  EventType = "stat_change"
)

// Event represents an event in the game
type Event struct {
	Type    EventType   `json:"type"`
	Time    int64       `json:"time"`
	Source  string      `json:"source"`
	Data    interface{} `json:"data"`
	Context interface{} `json:"context,omitempty"`
}

// EventHandler represents a function that handles events
type EventHandler func(*Event)

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

// Clear removes all event handlers
func (e *EventEmitter) Clear() {
	e.handlers = make(map[EventType][]EventHandler)
}
