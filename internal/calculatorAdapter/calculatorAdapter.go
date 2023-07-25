package calculatorAdapter

import (
	"context"
	"log"
	pb "serial/internal/grpc"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type CalculatorAdapter interface {
	Calculate(estimates chan<- Estimate, data []byte) error
	Close()
}

type GrpcAdapter struct {
	timeout time.Duration
	client  pb.CalculatorClient
	conn    *grpc.ClientConn
}

type Estimate float32

func NewGrpcAdapter() (*GrpcAdapter, error) {
	conn, err := grpc.Dial("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Printf("did not connect: %v", err)
		return nil, err
	}
	c := pb.NewCalculatorClient(conn)

	return &GrpcAdapter{
		timeout: time.Second,
		client:  c,
		conn:    conn,
	}, nil
}

func (ga *GrpcAdapter) Close() {
	ga.conn.Close()
}

func (ga *GrpcAdapter) Calculate(estimates chan<- Estimate, data []byte) error {
	ctx, cancel := context.WithTimeout(context.Background(), ga.timeout)
	defer cancel()

	r, err := ga.client.Calculate(ctx, &pb.Input{Data: data})
	if err != nil {
		log.Printf("coud not calculate: %v", err)
		return err
	}

	estimates <- Estimate(r.Estimate)
	return nil
}
