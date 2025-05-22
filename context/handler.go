package context

import (
	clickhouse2 "discord/pkg/clickhouse"
	"discord/pkg/jetstream"
	"discord/pkg/server/grpc"
	pb "discord/proto"
	"fmt"
)

func AppContextInit() (*AppContextImpl, error) {
	db, err := AppContextDBInit()
	if err != nil {
		fmt.Println("Error initializing DB:", err)
		return nil, err
	}
	stream, err := AppContextStreamInit()
	if err != nil {
		fmt.Println("Error initializing stream:", err)
	}

	server, err := AppContextServerInit()
	if err != nil {
		fmt.Println("Error initializing server:", err)
		return nil, err
	}
	return &AppContextImpl{
		db:     db,
		stream: stream,
		server: server,
	}, nil
}
func AppContextDBInit() (*DB, error) {
	clickhouseConn, err := clickhouse2.InitClickhouse()
	if err != nil {
		return nil, err
	}
	return &DB{Clickhouse: clickhouseConn}, nil
}

func AppContextStreamInit() (*Stream, error) {
	js, err := jetstream.InitJetstream()
	if err != nil {
		return nil, err
	}
	return &Stream{Jetstream: js}, nil
}
func AppContextServerInit() (*Server, error) {
	grpcClient, err := grpc.InitGrpcClient()
	if err != nil {
		return nil, err
	}
	g := &GRPCServer{
		GRPC:        grpcClient,
		GuildClient: pb.NewGuildClient(grpcClient),
	}
	return &Server{GRPC: g}, nil

}
