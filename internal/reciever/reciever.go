package reciever

import (
	"context"
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

func GetPacketLength() int {
	return 50
}

type Packet struct {
	Header []byte
	Data   []byte
	Footer []byte
}

func Listen(ctx context.Context, packets chan<- Packet, sm serialManager.SerialManager) {
	buffer := make([]byte, 0, 10000)

	log.SetFlags(log.Lmicroseconds)

	for {
		select {
		case <-ctx.Done():
			return
		default:
			b := make([]byte, 10000)
			n, err := sm.Read(b)
			if err != nil {
				if err == io.EOF {
					continue
				}
				log.Fatal(err)
				break
			}

			buffer = append(buffer, b[:n]...)
			buffer = TryParsePacket(buffer, packets)
		}
	}
}

func TryParsePacket(buffer []byte, packets chan<- Packet) []byte {
	header := GetHeader()
	footer := GetFooter()
	headerLength := len(header)
	footerLength := len(footer)
	packetLength := GetPacketLength()

	// Maintain a start index
	startIdx := 0
	i_max := len(buffer)

	for startIdx < i_max {
		// Check header
		if i_max-startIdx < headerLength || !equal(buffer[startIdx:startIdx+headerLength], header) {
			startIdx++
			continue
		}

		// Check packet size
		if i_max-startIdx < packetLength {
			break
		}

		// Check footer
		if !equal(buffer[startIdx+packetLength-footerLength:startIdx+packetLength], footer) {
			startIdx++
			continue
		}

		// Create and send packet
		packet := Packet{
			Header: header,
			Data:   buffer[startIdx+headerLength : startIdx+packetLength-footerLength],
			Footer: footer,
		}
		log.Println("Packet received")
		packets <- packet

		// Increment the start index by the packet length
		startIdx += packetLength
	}

	// If the start index has moved, remove the parsed data from the buffer
	if startIdx > 0 {
		copy(buffer, buffer[startIdx:])
		buffer = buffer[:i_max-startIdx]
	}

	return buffer
}

// New equal function to avoid creating subslices
func equal(a, b []byte) bool {
	if len(a) != len(b) {
		return false
	}
	for i, v := range a {
		if v != b[i] {
			return false
		}
	}
	return true
}
