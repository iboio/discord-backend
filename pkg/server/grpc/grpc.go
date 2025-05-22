package grpc

import (
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
)

func InitGrpcClient() (*grpc.ClientConn, error) {
	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}

	var err error
	conn, err := grpc.NewClient("localhost:50050", opts...)
	if err != nil {
		log.Printf("Failed to connect to gRPC server: %v", err)
		return nil, err
	}

	log.Println("Connected to gRPC server")
	return conn, nil
}
