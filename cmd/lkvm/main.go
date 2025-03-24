package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
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

func keyDown(c *gin.Context) {
	key := c.Query("key")
	control := c.Query("control") == "true"
	shift := c.Query("shift") == "true"
	alt := c.Query("alt") == "true"

	modifiers := []string{}
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
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

func keyUp(c *gin.Context) {
	if err := keyboard.Release(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

func main() {
	r := gin.Default()
	r.Static("/", "./public")
	r.POST("api/keydown", keyDown)
	r.POST("api/keyup", keyUp)

	fmt.Println("Starting server on :8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}
