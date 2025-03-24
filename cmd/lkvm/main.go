package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/tarm/serial"
	ch "github.com/tshelter/ch9329"
)

var (
	port     *serial.Port
	keyboard *ch.KeyboardSender
	mouse    *ch.MouseSender
)

func init() {
	c := &serial.Config{Name: "/dev/ttyUSB0", Baud: 9600, ReadTimeout: time.Millisecond * 50}
	var err error
	port, err = serial.OpenPort(c)
	if err != nil {
		log.Fatal(err)
	}
	keyboard = ch.NewKeyboardSender(port)
	mouse = ch.NewMouseSender(port)
}

func keyDownHandler(w http.ResponseWriter, r *http.Request) {
	key := r.URL.Query().Get("key")
	control := r.URL.Query().Get("control") == "true"
	shift := r.URL.Query().Get("shift") == "true"
	alt := r.URL.Query().Get("alt") == "true"

	var modifiers []string
	if control {
		modifiers = append(modifiers, "ctrl")
	}
	if shift {
		modifiers = append(modifiers, "shift")
	}
	if alt {
		modifiers = append(modifiers, "alt")
	}

	if err := keyboard.Press(key, modifiers); err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)

	if _, err := w.Write([]byte(`{"status":"ok"}`)); err != nil {
		log.Println(err)
		return
	}
}

func keyUpHandler(w http.ResponseWriter, r *http.Request) {
	if err := keyboard.Release(); err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	if _, err := w.Write([]byte(`{"status":"ok"}`)); err != nil {
		log.Fatal(err)
		return
	}
}

func fileServerHandler(w http.ResponseWriter, r *http.Request) {
	http.FileServer(http.Dir("./public")).ServeHTTP(w, r)
}

func main() {
	http.HandleFunc("/api/keydown", keyDownHandler)
	http.HandleFunc("/api/keyup", keyUpHandler)
	http.HandleFunc("/", fileServerHandler)

	fmt.Println("Starting server on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
