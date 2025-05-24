package wsinterface

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

type Gorilla struct {
	conn *websocket.Conn
}

func newGorillaWebsocketConnection(w http.ResponseWriter, r *http.Request) (*Gorilla, error) {
	log.Println("Creating wesocket connection with Gorilla package")
	upgrader := websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return nil, fmt.Errorf("error while creating gorilla websocket instance: %s", err.Error())
	}

	g := &Gorilla{
		conn: conn,
	}

	return g, nil
}

func (g *Gorilla) ReadMessage() (messageType int, p []byte, err error) {
	messageType, message, err := g.conn.ReadMessage()
	if err != nil {
		if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
			log.Printf("unexpected error while reading message: %s", err.Error())
		}
		return messageType, message, fmt.Errorf("error while reading message: %s", err.Error())
	}

	newline := []byte{'\n'}
	space := []byte{' '}
	message = bytes.TrimSpace(bytes.Replace(message, newline, space, -1))
	return messageType, message, nil
}

func (g *Gorilla) WriteMessage(messageType MessageType, data []byte) error {
	switch messageType {
	case TextMessageType:
		w, err := g.conn.NextWriter(websocket.TextMessage)
		if err != nil {
			return fmt.Errorf("error while create writer instance: %s", err.Error())
		}
		w.Write(data)

		if err := w.Close(); err != nil {
			return fmt.Errorf("error while closing writer: %s", err.Error())
		}
	case CloseMessageType:
		if err := g.conn.WriteMessage(websocket.CloseMessage, []byte{}); err != nil {
			return fmt.Errorf("error while writing CloseMessage: %s", err.Error())
		}
	case PingMessageType:
		if err := g.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
			return fmt.Errorf("error while writing PingMessage: %s", err.Error())
		}
	default:
		return fmt.Errorf("unsupported message type")

	}
	return nil
}

func (g *Gorilla) SetReadLimit(limit int64) {
	g.conn.SetReadLimit(limit)
}

func (g *Gorilla) SetReadDeadline(pongWait time.Time) error {
	err := g.conn.SetReadDeadline(pongWait)
	if err != nil {
		return fmt.Errorf("error while setting read deadline: %s", err.Error())
	}
	return nil
}

func (g *Gorilla) SetWriteDeadline(writeWait time.Time) error {
	err := g.conn.SetWriteDeadline(writeWait)
	if err != nil {
		return fmt.Errorf("error while setting write deadline: %s", err.Error())
	}
	return nil
}

func (g *Gorilla) SetPingHandler(h func(appData string) error) {
}

func (g *Gorilla) SetPongHandler(h func(appData string) error) {
	g.conn.SetPongHandler(h)
}

func (g *Gorilla) Close() error {
	if err := g.conn.Close(); err != nil {
		return fmt.Errorf("error while closing socket connection: %s", err.Error())
	}
	return nil
}
