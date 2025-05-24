package main

import (
	"log"
	"net/http"

	"github.com/websocket-chat-service/chat"
	"github.com/websocket-chat-service/websocket/handlers"
)

func main() {

	chatManager := chat.NewChatManager()
	go chatManager.Run()

	// Serve static files
	http.Handle("/", http.FileServer(http.Dir("./frontend")))

	http.HandleFunc("/login", handlers.LoginHandler)

	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		handlers.ServeWebSocket(chatManager, w, r)
	})

	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
