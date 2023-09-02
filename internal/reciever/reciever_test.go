package reciever

import (
	"context"
	"log"
	"math/rand"
	"serial/internal/serialManager"
	"testing"
	"time"
)

func TestListenWith(t *testing.T) {
	// Make random data for n_packets packets
	packetLength := GetPacketLength()
	var n_packets int
	var maxByteToReadInOnce int
	n_packets = 1000
	maxByteToReadInOnce = int(float32(packetLength))
	testListenWithRandomDataByPattern(n_packets, maxByteToReadInOnce, t)
}

func TestListenWithRandomData1(t *testing.T) {
	// Make random data for n_packets packets
	packetLength := GetPacketLength()
	var n_packets int
	var maxByteToReadInOnce int
	n_packets = 1000
	maxByteToReadInOnce = int(float32(packetLength) * 0.5)
	testListenWithRandomDataByPattern(n_packets, maxByteToReadInOnce, t)
}

func TestListenWithRandomData2(t *testing.T) {
	// Make random data for n_packets packets
	packetLength := GetPacketLength()
	var n_packets int
	var maxByteToReadInOnce int
	n_packets = 1000
	maxByteToReadInOnce = int(float32(packetLength) * 5)
	testListenWithRandomDataByPattern(n_packets, maxByteToReadInOnce, t)
}

func testListenWithRandomDataByPattern(n_packets int, maxByteToReadInOnce int, t *testing.T) {
	// n_packets = n_packets

	data := make([]byte, 0, n_packets*GetPacketLength())
	for i := 0; i < n_packets; i++ {
		data = append(data, GetHeader()...)
		data = append(data, randomBytes(GetPacketLength()-len(GetHeader())-len(GetFooter()))...)
		data = append(data, GetFooter()...)
	}

	data = data[3 : len(data)-3]

	slicePoint1 := n_packets * GetPacketLength() / 4
	slicePoint2 := n_packets * GetPacketLength() / 4 * 2
	slicePoint3 := n_packets * GetPacketLength() / 4 * 3
	data = append(data[:slicePoint1], data[slicePoint1+1:]...)
	data = append(data[:slicePoint2], data[slicePoint2+1:]...)
	data = append(data[:slicePoint3], data[slicePoint3+1:]...)

	sm := serialManager.NewMockSerialManager(data, maxByteToReadInOnce)

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	packets := make(chan Packet, int(float32(n_packets)*1.1))

	go func() {
		Listen(ctx, packets, sm)
		close(packets)
	}()

	count := 0
	for range packets {
		count++
		log.Printf("count %d", count)
	}

	n_packets = n_packets - 5 // 5 packets are removed by slicePoint1, slicePoint2, slicePoint3
	if count != n_packets {
		t.Fatalf("expected n_packets %d packets, got %d", n_packets, count)
	}
}

func randomBytes(n int) []byte {
	b := make([]byte, n)
	for i := range b {
		b[i] = byte(rand.Intn(256))
	}
	return b
}
