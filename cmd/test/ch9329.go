package main

import (
	"errors"
	_ "io"
	"sync"
	"time"

	"github.com/tarm/serial"
)

type CH9329 struct {
	PortName      string
	BaudRate      int
	LeftStatus    int
	serialPort    *serial.Port
	mediaKeyTable map[MediaKey][]byte
	mutex         sync.Mutex
}

type MediaKey int
type KeyGroup byte
type CommandCode byte
type MouseButtonCode byte

const (
	CharKey   KeyGroup = 0x02
	MediaKeyG KeyGroup = 0x03
)

const (
	GET_INFO             CommandCode = 0x01
	SEND_KB_GENERAL_DATA CommandCode = 0x02
	SEND_KB_MEDIA_DATA   CommandCode = 0x03
	SEND_MS_ABS_DATA     CommandCode = 0x04
	SEND_MS_REL_DATA     CommandCode = 0x05
)

const (
	EJECT MediaKey = iota
	CDSTOP
	PREVTRACK
	NEXTTRACK
	PLAYPAUSE
	MUTE
	VOLUMEDOWN
	VOLUMEUP
)

var mediaKeyTable = map[MediaKey][]byte{
	EJECT:      {0x02, 0x80, 0x00, 0x00},
	CDSTOP:     {0x02, 0x40, 0x00, 0x00},
	PREVTRACK:  {0x02, 0x20, 0x00, 0x00},
	NEXTTRACK:  {0x02, 0x10, 0x00, 0x00},
	PLAYPAUSE:  {0x02, 0x08, 0x00, 0x00},
	MUTE:       {0x02, 0x04, 0x00, 0x00},
	VOLUMEDOWN: {0x02, 0x02, 0x00, 0x00},
	VOLUMEUP:   {0x02, 0x01, 0x00, 0x00},
}

func NewCH9329(portName string, baudRate int) (*CH9329, error) {
	config := &serial.Config{Name: portName, Baud: baudRate, ReadTimeout: time.Second}
	port, err := serial.OpenPort(config)
	if err != nil {
		return nil, err
	}

	ch := &CH9329{
		PortName:      portName,
		BaudRate:      baudRate,
		serialPort:    port,
		mediaKeyTable: mediaKeyTable,
	}
	return ch, nil
}

func (ch *CH9329) sendPacket(data []byte) error {
	ch.mutex.Lock()
	defer ch.mutex.Unlock()

	_, err := ch.serialPort.Write(data)
	if err != nil {
		return err
	}

	time.Sleep(time.Millisecond)
	return nil
}

func (ch *CH9329) createPacketArray(arrList []int, addCheckSum bool) []byte {
	bytePacket := make([]byte, len(arrList))
	for i, v := range arrList {
		bytePacket[i] = byte(v)
	}

	if addCheckSum {
		sum := 0
		for _, v := range arrList {
			sum += v
		}
		bytePacket = append(bytePacket, byte(sum&0xFF))
	}
	return bytePacket
}

func (ch *CH9329) keyDown(keyGroup KeyGroup, keys ...byte) error {
	if len(keys) > 6 {
		return errors.New("too many keys")
	}

	data := []int{0x57, 0xAB, 0x00, int(keyGroup), 0x08}
	for _, key := range keys {
		data = append(data, int(key))
	}

	for len(data) < 13 {
		data = append(data, 0)
	}

	packet := ch.createPacketArray(data, true)
	return ch.sendPacket(packet)
}

func (ch *CH9329) keyUpAll(keyGroup KeyGroup) error {
	var packet []byte
	if keyGroup == CharKey {
		packet = []byte{0x57, 0xAB, 0x00, 0x02, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x0C}
	} else {
		packet = []byte{0x57, 0xAB, 0x00, 0x03, 0x04, 0x02, 0x00, 0x00, 0x00, 0x0B}
	}
	return ch.sendPacket(packet)
}

func (ch *CH9329) keyPressMedia(mediaKey MediaKey) error {
	data, exists := ch.mediaKeyTable[mediaKey]
	if !exists {
		return errors.New("invalid media key")
	}
	return ch.sendPacket(data)
}

func (ch *CH9329) mouseMoveRel(x, y int) error {
	if x > 127 {
		x = 127
	} else if x < -128 {
		x = -128
	}
	if y > 127 {
		y = 127
	} else if y < -128 {
		y = -128
	}

	if x < 0 {
		x += 256
	}
	if y < 0 {
		y += 256
	}

	data := []int{0x57, 0xAB, 0x00, 0x05, 0x05, 0x01, ch.LeftStatus, x, y, 0x00}
	packet := ch.createPacketArray(data, true)
	return ch.sendPacket(packet)
}

func (ch *CH9329) mouseButtonDown(buttonCode MouseButtonCode) error {
	data := []int{0x57, 0xAB, 0x00, 0x05, 0x05, 0x01, int(buttonCode), 0x00, 0x00, 0x00}
	if buttonCode == 0x01 {
		ch.LeftStatus = 1
	}

	packet := ch.createPacketArray(data, true)
	return ch.sendPacket(packet)
}

func (ch *CH9329) mouseButtonUpAll() error {
	ch.LeftStatus = 0
	packet := []byte{0x57, 0xAB, 0x00, 0x05, 0x05, 0x01, 0x00, 0x00, 0x00, 0x00, 0x0D}
	return ch.sendPacket(packet)
}
