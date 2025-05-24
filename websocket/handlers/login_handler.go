package handlers

import (
	"encoding/json"
	"log"
	"net/http"
)

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req LoginRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	// Simple authentication logic (replace with real auth)
	if req.Username == "admin" && req.Password == "password" {
		log.Printf("Successful logged-in user: %s", req.Username)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"message":"Login successful"}`))
	} else {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
	}
}
