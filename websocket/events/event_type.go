package events

type EventType string

const (
	EventSendMessage    = "send_message"
	EventReceiveMessage = "receive_message"
	EventChangeRoom     = "change_room"
)

type EventHandler func() error
