package handler

import (
	"discord/config"
	"discord/internal/cron"
	"discord/internal/processor"
	"discord/pkg/clickhouse"
	"discord/pkg/jetstream"
	redisdb "discord/pkg/redis"
	"fmt"
)

func Handler() {
	config.LoadConfig()
	redisdb.InitRedis(0)
	redisdb.InitRedis(1)
	fmt.Println("redis connection success")
	ch, err := clickhouse.InitClickhouse()
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("clickhouse connection success")
	js, err := jetstream.InitJetstream()
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("jetstream connection success")
	p := processor.HandlerProcessor(js, ch)
	p.Process()
	fmt.Println("processor started")
	c := cron.HandlerProcessor(ch)
	c.StartCron()
	fmt.Println("cron started")
	select {}
}
