package sender

import (
	"serial/internal/calculatorAdapter"
)

// This function style makes header const byte[]
func GetHeader() []byte {
	return []byte{0xBB}
}

func GetFooter() []byte {
	return []byte{0xAA}
}

func GetInstructionA() []byte {
	return []byte{0o33, 0o33, 0o33, 0o33}
}

func GetInstructionB() []byte {
	return []byte{0o66, 0o66, 0o66, 0o66}
}

func calculateCheckSum(b []byte) byte {
	var sum byte
	for _, v := range b {
		sum += v
	}
	return sum
}

func ComposepPacket(estimate calculatorAdapter.Estimate) []byte {
	packet := make([]byte, 7)

	if estimate > 0.5 {
		copy(packet[0:2], GetHeader())
		copy(packet[2:6], GetInstructionA())
		checkSum := calculateCheckSum(packet[0:6])
		packet[6] = checkSum
	} else {
		copy(packet[0:2], GetHeader())
		copy(packet[2:6], GetInstructionB())
		checkSum := calculateCheckSum(packet[0:6])
		packet[6] = checkSum
	}

	return packet
}
