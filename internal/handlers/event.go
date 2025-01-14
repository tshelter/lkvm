package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/tshelter/lkvm/internal/dto"
)

type WebsocketTransport struct {
	upgrade websocket.Upgrader
}

func NewWebSocketEventHandler() func(w http.ResponseWriter, r *http.Request) {
	return WebsocketTransport{
		upgrade: websocket.Upgrader{ReadBufferSize: 1024, WriteBufferSize: 1024},
	}.WebSocketEventHandler
}

func (t WebsocketTransport) WebSocketEventHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Incoming request to upgrade to WebSocket")

	conn, err := t.upgrade.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("WebSocket upgrade error: %v", err)
		writeJSON(w, http.StatusInternalServerError, dto.ErrorResponse{ErrorMsg: "Unable to upgrade WebSocket connection"})
		return
	}
	log.Println("WebSocket connection established")
	defer func(conn *websocket.Conn) {
		err := conn.Close()
		if err != nil {
			log.Printf("Error closing WebSocket connection: %v", err)
		}
	}(conn)

	for {
		var event dto.Event
		if err := conn.ReadJSON(&event); err != nil {
			log.Printf("Error reading JSON from WebSocket: %v", err)
			if writeErr := conn.WriteJSON(dto.ErrorResponse{ErrorMsg: "Invalid event format"}); writeErr != nil {
				log.Printf("Error sending error response: %v", writeErr)
			}
			continue
		}

		log.Printf("Received event: %+v", event)
	}
}

func writeJSON(w http.ResponseWriter, status int, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(v); err != nil {
		log.Printf("Error writing JSON response: %v", err)
	}
}
