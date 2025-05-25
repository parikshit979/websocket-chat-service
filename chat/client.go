package chat

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/websocket-chat-service/websocket/events"
	"github.com/websocket-chat-service/websocket/wsinterface"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 512
)

// Client is a middleman between the websocket connection and the chat manager.
type ChatClient struct {
	chatManager *ChatManager

	// The websocket connection.
	conn wsinterface.WebsocketInterface

	// Unbuffered channel of outbound messages.
	send chan events.Event

	// Client current chat room.
	chatRoom string
}

func NewChatClient(chatManager *ChatManager, conn wsinterface.WebsocketInterface) *ChatClient {
	client := &ChatClient{
		chatManager: chatManager,
		conn:        conn,
		send:        make(chan events.Event),
	}
	client.chatManager.register <- client
	return client
}

func (c *ChatClient) GetChatManger() *ChatManager {
	return c.chatManager
}

func (c *ChatClient) GetChatRoom() string {
	return c.chatRoom
}

func (c *ChatClient) createBroadcastEventMessage(event events.Event) (*events.Event, error) {
	// Marshal payload into ReceiveMessageEvent.
	var receiveMessageEvent events.ReceiveMessageEvent
	if err := json.Unmarshal(event.Payload, &receiveMessageEvent); err != nil {
		return nil, fmt.Errorf("bad payload in request: %s", err.Error())
	}

	// Prepare SendMessageEvent for broadcast.
	var sendMessageMessage events.SendMessageEvent

	sendMessageMessage.Sent = time.Now()
	sendMessageMessage.Message = receiveMessageEvent.Message
	sendMessageMessage.From = receiveMessageEvent.From
	sendMessageMessage.Room = receiveMessageEvent.Room

	data, err := json.Marshal(sendMessageMessage)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal broadcast message: %v", err)
	}

	// Place payload into otugoing event.
	var outgoingEvent events.Event
	outgoingEvent.Payload = data
	outgoingEvent.Type = events.EventReceiveMessage

	return &outgoingEvent, nil
}

// ReadMessagseFromWebSocket reads message from the websocket connection to the chat manager.
//
// The application runs ReadMessagesFromWebSocket in a per-connection goroutine. The application
// ensures that there is at most one reader on a connection by executing all
// reads from this goroutine.
func (c *ChatClient) ReadMessagesFromWebSocket() {
	defer func() {
		c.chatManager.unregister <- c
		c.conn.Close()
	}()
	c.conn.SetReadLimit(maxMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error { c.conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			log.Printf("Failed to to read message: %s", err.Error())
			break
		}

		// Marshal incoming data into a Event struct
		var event events.Event
		if err := json.Unmarshal(message, &event); err != nil {
			log.Printf("error marshalling message: %s", err.Error())
			continue
		}

		if event.Type == events.EventChangeRoom {
			// Marshal Payload into wanted format
			var changeRoomEvent events.ChangeRoomEvent
			if err := json.Unmarshal(event.Payload, &changeRoomEvent); err != nil {
				log.Printf("bad payload in request for change room event: %s", err.Error())
			}
			c.chatRoom = changeRoomEvent.Room
		} else {
			var outgoingEvent *events.Event
			outgoingEvent, err = c.createBroadcastEventMessage(event)
			if err != nil {
				log.Printf("error while creating brodcast message: %s", err.Error())
			}
			c.chatManager.broadcast <- *outgoingEvent
		}
	}
}

// WriteMessageToWebSocket writes messagess from the chat manager to the websocket connection.
//
// A goroutine running WriteMessagesToWebSocket is started for each connection. The
// application ensures that there is at most one writer to a connection by
// executing all writes from this goroutine.
func (c *ChatClient) WriteMessagesToWebSocket() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()
	for {
		select {
		case event, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// The chat manager closed the channel.
				c.conn.WriteMessage(wsinterface.CloseMessageType, []byte{})
				return
			}

			message, err := json.Marshal(event)
			if err != nil {
				log.Printf("failed to marshal event message: %s", err.Error())
				return
			}

			err = c.conn.WriteMessage(wsinterface.TextMessageType, message)
			if err != nil {
				log.Printf("Failed to write message: %s", err.Error())
				return
			}
		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(wsinterface.PingMessageType, []byte{}); err != nil {
				return
			}
		}
	}
}
