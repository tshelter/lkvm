package main

import (
	"log"
	"net/http"

	"github.com/tshelter/lkvm/internal/handlers"
)

func main() {
	fs := http.FileServer(http.Dir("public"))

	http.Handle("/", fs)
	http.HandleFunc("/ws", handlers.NewWebSocketEventHandler())

	log.Println("Listening on :8000...")
	log.Fatal(http.ListenAndServe(":8000", nil))
}
