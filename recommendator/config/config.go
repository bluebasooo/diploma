package config

import (
	"os"

	"dev/bluebasooo/video-common/config"
)

func GetApplicationConfig() *ApplicationConfig {
	return &ApplicationConfig{
		ClickhouseConfig: ClickhouseConfig{
			URI:      os.Getenv("CLICKHOUSE_URI"),
			User:     os.Getenv("CLICKHOUSE_USER"),
			Password: os.Getenv("CLICKHOUSE_PASSWORD"),
			Database: os.Getenv("CLICKHOUSE_DATABASE"),
		},
		KafkaWriterConfig: KafkaWriterConfig{
			BrokerPort: os.Getenv("KAFKA_BROKER_URL"),
			BrokerHost: os.Getenv("KAFKA_BROKER_URL"),
			Topic:      os.Getenv("KAFKA_TOPIC"),
		},
		KafkaReaderConfig: KafkaReaderConfig{
			BrokerHost: os.Getenv("KAFKA_BROKER_URL"),
			BrokerPort: os.Getenv("KAFKA_BROKER_PORT"),
			Topic:      os.Getenv("KAFKA_TOPIC"),
			GroupID:    os.Getenv("KAFKA_GROUP_ID"),
		},
		MongoConfig: config.MongoConfig{
			Host:         os.Getenv("DB_MONGO_HOST"),
			Port:         os.Getenv("DB_MONGO_PORT"),
			User:         os.Getenv("DB_MONGO_USER"),
			Password:     os.Getenv("DB_MONGO_PASSWORD"),
			DatabaseName: os.Getenv("DB_MONGO_DATABASE_NAME"),
			CollectionsNames: []string{
				"dots",
				"buckets",
			},
		},
		Debug: true,
	}
}

type ApplicationConfig struct {
	ClickhouseConfig  ClickhouseConfig
	KafkaReaderConfig KafkaReaderConfig
	KafkaWriterConfig KafkaWriterConfig
	MongoConfig       config.MongoConfig
	Debug             bool
}

type KafkaReaderConfig struct {
	BrokerHost string
	BrokerPort string
	Topic      string
	GroupID    string
}

type KafkaWriterConfig struct {
	BrokerHost string
	BrokerPort string
	Topic      string
}

type ClickhouseConfig struct {
	URI      string
	User     string
	Password string
	Database string
}
