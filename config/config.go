package config

import (
	"fmt"
	"github.com/joho/godotenv"
	"os"
	"strconv"
)

var (
	Period                      int64
	ClickhouseMessageBatchSize  int
	ClickhouseVoiceLogBatchSize int
	DBHost                      string
	DBPort                      int
	ClickhouseHost              string
	ClickhouseHttpPort          string
	ClickhouseNativePort        string
	ClickhouseUsername          string
	ClickhousePassword          string
	ClickhouseDatabase          string
	RedisHost                   string
	RedisPort                   string
	NatsHost                    string
	NatsPort                    string
	JetstreamName               string
	JetstreamSubjects           string
	JetstreamSubjectsEventMsg   string
	JetstreamSubjectsEventVoice string
)

func LoadConfig() {
	err := godotenv.Load("../.env")
	if err != nil {
		fmt.Println(err)
		panic("Error loading .env file")
	}

	Period = int64(getEnvAsInt("PERIOD", 60))
	ClickhouseMessageBatchSize = getEnvAsInt("CLICKHOUSE_MESSAGE_BATCH_SIZE", 5)
	ClickhouseVoiceLogBatchSize = getEnvAsInt("CLICKHOUSE_VOICE_LOG_BATCH_SIZE", 2)

	DBHost = getEnv("DB_HOST", "localhost")
	DBPort = getEnvAsInt("DB_PORT", 5432)

	ClickhouseMessageBatchSize = getIntEnv("CLICKHOUSE_MESSAGE_BATCH_SIZE", 50)
	ClickhouseVoiceLogBatchSize = getIntEnv("CLICKHOUSE_VOICE_LOG_BATCH_SIZE", 20)
	ClickhouseHost = getEnv("CLICKHOUSE_HOST", "localhost")
	ClickhouseHttpPort = getEnv("CLICKHOUSE_HTTP_PORT", "8125")
	ClickhouseNativePort = getEnv("CLICKHOUSE_NATIVE_PORT", "9005")
	ClickhouseUsername = getEnv("CLICKHOUSE_USERNAME", "root")
	ClickhousePassword = getEnv("CLICKHOUSE_PASSWORD", "root")
	ClickhouseDatabase = getEnv("CLICKHOUSE_DATABASE", "discord")

	RedisHost = getEnv("REDIS_HOST", "localhost")
	RedisPort = getEnv("REDIS_PORT", "6380")

	NatsHost = getEnv("NATS_HOST", "localhost")
	NatsPort = getEnv("NATS_PORT", "4223")
	JetstreamName = getEnv("JETSTREAM_NAME", "discord")
	JetstreamSubjects = getEnv("JETSTREAM_SUBJECTS", "event.*.*")
	JetstreamSubjectsEventMsg = getEnv("JETSTREAM_SUBJECTS_EVENT_MSG", "event.msg")
	JetstreamSubjectsEventVoice = getEnv("JETSTREAM_SUBJECTS_EVENT_VOICE", "event.voice.process")

}

func getEnv(key string, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	if value, exists := os.LookupEnv(key); exists {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

func getIntEnv(key string, defaultValue int) int {
	if value, exists := os.LookupEnv(key); exists {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}
