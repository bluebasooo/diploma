package db

import (
	"dev/bluebasooo/video-recomendator/config"

	"github.com/ClickHouse/clickhouse-go/v2"
)

type ClickhouseDB struct {
	db clickhouse.Conn
}

func NewClickhouseDB(config config.ApplicationConfig) *ClickhouseDB {
	opts := clickhouse.Options{
		Addr: []string{config.ClickhouseConfig.URI},
		Auth: clickhouse.Auth{
			Database: config.ClickhouseConfig.Database,
			Username: config.ClickhouseConfig.User,
			Password: config.ClickhouseConfig.Password,
		},
	}

	db, err := clickhouse.Open(&opts)

	if err != nil {
		panic(err)
	}
	return &ClickhouseDB{db: db}
}

func (c *ClickhouseDB) Close() error {
	return c.db.Close()
}

func (c *ClickhouseDB) Db() clickhouse.Conn {
	return c.db
}
