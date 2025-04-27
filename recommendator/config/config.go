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
	}
}

type ApplicationConfig struct {
	ClickhouseConfig ClickhouseConfig
}

type ClickhouseConfig struct {
	URI      string
	User     string
	Password string
	Database string
}
