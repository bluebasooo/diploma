package config

import "os"

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
			Topic:      os.Getenv("KAFKA_TOPIC"),

		}
	}
}

type ApplicationConfig struct {
	ClickhouseConfig ClickhouseConfig
	KafkaReaderConfig KafkaReaderConfig
	KafkaWriterConfig KafkaWriterConfig
}

type KafkaReaderConfig struct {
	BrokerHost string
	BrokerPort string
	Topic      string
	GroupId    string
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
