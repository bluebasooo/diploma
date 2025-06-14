package db

import (
	"dev/bluebasooo/video-recomendator/config"

	"github.com/ClickHouse/clickhouse-go/v2"
)

type ClickhouseDB struct {
	db clickhouse.Conn
}

func NewClickhouseDB(config *config.ClickhouseConfig) *ClickhouseDB {
	opts := clickhouse.Options{
		Addr: []string{config.URI},
		Auth: clickhouse.Auth{
			Database: config.Database,
			Username: config.User,
			Password: config.Password,
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
