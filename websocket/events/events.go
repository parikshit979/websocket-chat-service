package events

import (
	"encoding/json"
	"time"
)

type EventType string

const (
	EventSendMessage    = "send_message"
	EventReceiveMessage = "receive_message"
	EventChangeRoom     = "change_room"
)

// Event is the Messages sent over the websocket
// Used to differ between different actions
type Event struct {
	// Type is the message type sent
	Type string `json:"type"`
	// Payload is the data Based on the Type
	Payload json.RawMessage `json:"payload"`
}

// ReceiveMessageEvent is the payload sent in the
// send_message event from the WebSocket.
type ReceiveMessageEvent struct {
	Message string `json:"message"`
	From    string `json:"from"`
	Room    string `json:"room"`
}

// SendMessageEvent is the payload sent in the
// receive_message event to the WebSocket.
type SendMessageEvent struct {
	ReceiveMessageEvent
	Sent time.Time `json:"sent"`
}

type ChangeRoomEvent struct {
	Room string `json:"room"`
}
