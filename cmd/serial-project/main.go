package main

import (
	"log"
	"os"
	"os/signal"
	"serial/internal/calculatorAdapter"
	"serial/internal/reciever"
	"serial/internal/sender"
	"serial/internal/serialManager"
	"syscall"

	"github.com/tarm/serial"
)

func calculationStage(packets <-chan reciever.Packet) chan calculatorAdapter.Estimate {
	var adapter calculatorAdapter.CalculatorAdapter = calculatorAdapter.NewFixedCalculatorAdapter(0.5)
	defer adapter.Close()
	estimates := make(chan calculatorAdapter.Estimate, 1000)
	go func() {
		for packet := range packets {
			log.Printf("Received packet: %q", packet.Data)
			adapter.Calculate(estimates, packet.Data)
		}
	}()

	return estimates
}

func serialSendingStage(estimates <-chan calculatorAdapter.Estimate, sm *serialManager.SerialManagement) {
	serialSender := sender.NewSerialSender(sm)
	go func() {
		for estimate := range estimates {
			log.Printf("Received estimate: %f", estimate)
			packet := sender.ComposepPacket(estimate)
			serialSender.Send(packet)
		}
	}()
}

func main() {
	log.SetFlags(log.Lmicroseconds)
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	serialManagement, err := serialManager.NewSerialManagement(&serial.Config{
		Name: "/dev/tty99",
		Baud: 9600,
	})
	if err != nil {
		log.Fatalf("Failed to open serial port: %v", err)
	}
	defer serialManagement.Close()

	packets := make(chan reciever.Packet, 50)
	estimates := calculationStage(packets)
	serialSendingStage(estimates, serialManagement)

	go reciever.Listen(packets, serialManagement)

	func() {
		<-sigChan
		log.Println("Received interrupt, exiting.")
		os.Exit(0)
	}()
}
