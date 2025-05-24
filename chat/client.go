package chat

import (
	"log"
	"time"

	"github.com/gorilla/websocket"

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

// Client is a middleman between the websocket connection and the hub.
type ChatClient struct {
	chatManager *ChatManager

	// The websocket connection.
	conn wsinterface.WebsocketInterface

	// Buffered channel of outbound messages.
	send chan []byte
}

func NewChatClient(chatManager *ChatManager, conn wsinterface.WebsocketInterface) *ChatClient {
	client := &ChatClient{
		chatManager: chatManager,
		conn:        conn,
		send:        make(chan []byte),
	}
	client.chatManager.register <- client
	return client
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
		c.chatManager.broadcast <- message
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
		case message, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// The chat manager closed the channel.
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			err := c.conn.WriteMessage(wsinterface.TextMessageType, message)
			if err != nil {
				log.Printf("Failed to write message: %s", err.Error())
				return
			}
		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, []byte{}); err != nil {
				return
			}
		}
	}
}
