package chat

import "github.com/websocket-chat-service/websocket/events"

// ChatManager maintains the set of active clients and broadcasts messages to the
// clients.
type ChatManager struct {
	// Registered clients.
	clients map[*ChatClient]bool

	// Inbound messages from the clients.
	broadcast chan events.Event

	// Register requests from the clients.
	register chan *ChatClient

	// Unregister requests from clients.
	unregister chan *ChatClient
}

func NewChatManager() *ChatManager {
	return &ChatManager{
		broadcast:  make(chan events.Event),
		register:   make(chan *ChatClient),
		unregister: make(chan *ChatClient),
		clients:    make(map[*ChatClient]bool),
	}
}

func (cm *ChatManager) Run() {
	for {
		select {
		case client := <-cm.register:
			cm.clients[client] = true
		case client := <-cm.unregister:
			if _, ok := cm.clients[client]; ok {
				delete(cm.clients, client)
				close(client.send)
			}
		case message := <-cm.broadcast:
			for client := range cm.clients {
				select {
				case client.send <- message:
				default:
					close(client.send)
					delete(cm.clients, client)
				}
			}
		}
	}
}

func (cm *ChatManager) GetClients() map[*ChatClient]bool {
	return cm.clients
}
