package handlers

import (
	"log"
	"net/http"

	"github.com/websocket-chat-service/chat"
	"github.com/websocket-chat-service/websocket/wsinterface"
)

// ServeWebSocket handles websocket requests from the peer.
func ServeWebSocket(chatManager *chat.ChatManager, w http.ResponseWriter, r *http.Request) {
	log.Println("Serving websocket")
	webSocketConn, err := wsinterface.NewWebsocketInterface(wsinterface.GorillaWebsocket, w, r)
	if err != nil {
		log.Printf("Failed to create websocket socket interface: %s", err.Error())
	}
	client := chat.NewChatClient(chatManager, webSocketConn)

	// Read messages from websocket and write messages to websocket with below goroutines.
	go client.WriteMessagesToWebSocket()
	go client.ReadMessagesFromWebSocket()
}
