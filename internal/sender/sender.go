package sender

import "serial/internal/serialManager"

type Sender interface {
	Send(packet []byte)
}

type SerialSender struct {
	serialManager serialManager.SerialManager
}

func NewSerialSender(serialManager serialManager.SerialManager) *SerialSender {
	return &SerialSender{
		serialManager: serialManager,
	}
}

func (ss *SerialSender) Send(packet []byte) {
	ss.serialManager.Write(packet)
}
