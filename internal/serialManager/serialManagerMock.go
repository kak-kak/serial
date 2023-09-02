package serialManager

import (
	"bytes"
	"log"
	"math/rand"
	"time"
)

type MockSerialManager struct {
	dataBuffer          *bytes.Buffer
	maxByteToReadInOnce int
	amountN             int
}

// NewMockSerialManager creates a new mock SerialManager with the provided data
func NewMockSerialManager(data []byte, maxByteToReadInOnce int) *MockSerialManager {
	return &MockSerialManager{
		dataBuffer:          bytes.NewBuffer(data),
		maxByteToReadInOnce: maxByteToReadInOnce,
		amountN:             0,
	}
}

// Close does nothing and returns no error
func (m *MockSerialManager) Close() error {
	return nil
}

func (m *MockSerialManager) Read(b []byte) (int, error) {
	// Determine how many bytes to read
	n := rand.Intn(m.maxByteToReadInOnce)

	m.amountN = m.amountN + n

	// Simulate latency
	time.Sleep(time.Duration(rand.Intn(10)) * time.Millisecond)

	log.Printf("%d / %d", m.amountN, m.dataBuffer.Len())

	// Read up to n bytes
	temp := make([]byte, n)
	n, err := m.dataBuffer.Read(temp)
	if err != nil {
		return 0, err
	}

	// Copy the read bytes to b
	copy(b, temp[:n])

	return n, nil
}

// Write does nothing and returns a success
func (m *MockSerialManager) Write(b []byte) (int, error) {
	return len(b), nil
}
