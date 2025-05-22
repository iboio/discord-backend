package context

import (
	pb "discord/proto"
	"github.com/ClickHouse/clickhouse-go/v2"
	"github.com/nats-io/nats.go"
	"google.golang.org/grpc"
)

type AppContext interface {
	DB() *DB
	Stream() *Stream
	Server() *Server
}

type DB struct {
	Clickhouse clickhouse.Conn
}
type Stream struct {
	Jetstream nats.JetStreamContext
}
type Server struct {
	GRPC *GRPCServer
}
type GRPCServer struct {
	GRPC        *grpc.ClientConn
	GuildClient pb.GuildClient
}
type AppContextImpl struct {
	db     *DB
	stream *Stream
	server *Server
}

func (c *AppContextImpl) DB() *DB {
	return c.db
}

func (c *AppContextImpl) Stream() *Stream {
	return c.stream
}

func (c *AppContextImpl) Server() *Server {
	return c.server
}
