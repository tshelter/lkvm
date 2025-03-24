package main

import (
	"github.com/tarm/serial"
	"log"
	"time"
)

func main() {
	ch, err := NewCH9329("/dev/ttyUSB0", 115200)
	if err != nil {
		log.Fatalf("Failed to initialize CH9329: %v", err)
	}
	defer func(serialPort *serial.Port) {
		err := serialPort.Close()
		if err != nil {
			log.Fatalf("Failed to close serial port: %v", err)
		}
	}(ch.serialPort)

	// Define the ASCII key codes for "Hello World" (including a space and uppercase letters)
	message := []byte{
		0x0B, // H
		0x08, // e
		0x0F, // l
		0x0F, // l
		0x12, // o
		0x2C, // Space
		0x1A, // W
		0x12, // o
		0x13, // r
		0x0F, // l
		0x07, // d
	}

	time.Sleep(2 * time.Second) // Add delay before sending a message
	log.Println("Sending message 'Hello World'...")

	// Press and release each key
	for _, key := range message {
		err := ch.keyDown(CharKey, key)
		if err != nil {
			log.Printf("Failed to press key %X: %v", key, err)
		}
		time.Sleep(50 * time.Millisecond) // Add delay between key presses
		err = ch.keyUpAll(CharKey)
		if err != nil {
			log.Printf("Failed to release key %X: %v", key, err)
		}
	}

	log.Println("Message 'Hello World' sent successfully!")
}
