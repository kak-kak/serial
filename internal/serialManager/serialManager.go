package serialManager

import (
	"sync"

	"github.com/tarm/serial"
)

type SerialManager interface {
	Close() error
	Read(b []byte) (int, error)
	Write(b []byte) (int, error)
}

type SerialManagement struct {
	port *serial.Port
	lock sync.Mutex
}

func NewSerialManagement(config *serial.Config) (*SerialManagement, error) {
	port, err := serial.OpenPort(config)
	if err != nil {
		return nil, err
	}
	return &SerialManagement{port: port}, nil
}

func (sm *SerialManagement) Close() error {
	sm.lock.Lock()
	defer sm.lock.Unlock()
	return sm.port.Close()
}

func (sm *SerialManagement) Read(b []byte) (int, error) {
	sm.lock.Lock()
	defer sm.lock.Unlock()
	return sm.port.Read(b)
}

func (sm *SerialManagement) Write(b []byte) (int, error) {
	sm.lock.Lock()
	defer sm.lock.Unlock()
	return sm.port.Write(b)
}
