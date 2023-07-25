package reciever

import (
	"bytes"
	"io"
	"log"
	"serial/internal/serialManager"
)

// This function style make it possible to const byte[]
func GetHeader() []byte {
	return []byte{0xBB}
}

// This function style make it possible to const byte[]
func GetFooter() []byte {
	return []byte{0xAA}
}

type Packet struct {
	Header []byte
	Data   []byte
	Footer []byte
}

func Listen(packets chan<- Packet, sm serialManager.SerialManager) {
	defer sm.Close()
	buffer := make([]byte, 0)
	func() {
		for {
			b := make([]byte, 10)
			n, err := sm.Read(b)
			if err != nil {
				if err == io.EOF {
					continue
				}
				log.Fatal(err)
				break
			}

			buffer = append(buffer, b[:n]...)

			if packet, rest, found := TryParsePacket(buffer); found {
				packets <- packet
				buffer = rest
			}
		}
	}()
}

func TryParsePacket(buffer []byte) (Packet, []byte, bool) {
	headerIndex := bytes.Index(buffer, GetHeader())
	if headerIndex == -1 {
		return Packet{}, buffer, false
	}

	footerIndex := bytes.Index(buffer[headerIndex:], GetFooter())
	if footerIndex == -1 {
		return Packet{}, buffer, false
	}

	packetEndIndex := headerIndex + footerIndex + len(GetFooter())
	if packetEndIndex > len(buffer) {
		return Packet{}, buffer, false
	}

	packet := Packet{
		Header: GetHeader(),
		Data:   buffer[headerIndex+len(GetHeader()) : footerIndex],
		Footer: GetFooter(),
	}
	rest := buffer[packetEndIndex:]

	return packet, rest, true
}
