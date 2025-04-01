package clickhouse

import (
	"context"
	"discord/config"
	"fmt"
	"log"

	"github.com/ClickHouse/clickhouse-go/v2"
)

type ConnectionClickhouse struct {
	Connection     clickhouse.Conn
	MessageRecords []string
	VoiceRecords   []string
}

func InitClickhouse() (*ConnectionClickhouse, error) {
	fmt.Println(config.ClickhouseHost, config.ClickhouseNativePort)
	conn, err := clickhouse.Open(
		&clickhouse.Options{
			Addr: []string{fmt.Sprintf("%s:%s", config.ClickhouseHost, config.ClickhouseNativePort)},
			Auth: clickhouse.Auth{
				Database: config.ClickhouseDatabase,
				Username: config.ClickhouseUsername,
				Password: config.ClickhousePassword,
			},
			Protocol: clickhouse.Native,
		})

	if err != nil {
		log.Printf("Error connecting to ClickHouse: %v\n", err)
		return nil, err
	}

	ctx := context.TODO()
	if err := conn.Ping(ctx); err != nil {
		log.Printf("Failed to ping ClickHouse: %v\n", err)
		return nil, err
	}

	log.Println("Connected to ClickHouse successfully!")

	return &ConnectionClickhouse{
		MessageRecords: []string{},
		VoiceRecords:   []string{},
		Connection:     conn,
	}, nil
}
