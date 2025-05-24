package wsinterface

import (
	"fmt"
	"net/http"
	"time"
)

type WebsocketPackageType string

const (
	GorillaWebsocket WebsocketPackageType = "gorilla"
)

type MessageType int

const (
	TextMessageType MessageType = iota
	CloseMessageType
	PingMessageType
)

type WebsocketInterface interface {
	ReadMessage() (messageType int, p []byte, err error)
	WriteMessage(messageType MessageType, data []byte) error
	SetReadLimit(limit int64)
	SetReadDeadline(t time.Time) error
	SetWriteDeadline(t time.Time) error
	SetPingHandler(h func(appData string) error)
	SetPongHandler(h func(appData string) error)
	Close() error
}

func NewWebsocketInterface(packageType WebsocketPackageType, w http.ResponseWriter, r *http.Request) (WebsocketInterface, error) {
	switch packageType {
	case GorillaWebsocket:
		return newGorillaWebsocketConnection(w, r)
	default:
		return nil, fmt.Errorf("unsupported websocket package")
	}
}
