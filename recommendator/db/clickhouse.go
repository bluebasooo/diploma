package db

import (
	"context"
	"dev/bluebasooo/video-recomendator/config"
	"log"

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
		log.Fatal(err)
	}

	err = db.Ping(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	ch := &ClickhouseDB{db: db}

	ch.migrations()

	return ch
}

var migrateMetrics = `
CREATE TABLE IF NOT EXISTS metrics (
    id         UUID       DEFAULT generateUUIDv4(),
    user_id    String,
    video_id   String,
    view_id    String,
    type       String,
    value      Float64,
    created_at DateTime64(3)
)
ENGINE = MergeTree
PARTITION BY toYYYYMM(created_at)
ORDER BY (user_id, video_id, view_id, type, created_at);
`

var migrateHistory = `
CREATE TABLE IF NOT EXISTS history (
    id         UUID       DEFAULT generateUUIDv4(),
    user_id    String,
    video_id   String,
    value      Float64,
    created_at DateTime64(3)
)
ENGINE = MergeTree
PARTITION BY toYYYYMM(created_at)
ORDER BY (user_id, video_id, created_at);
`

func (c *ClickhouseDB) migrations() {
	err := c.db.Exec(context.Background(), migrateMetrics)
	if err != nil {
		log.Fatal(err)
	}

	err = c.db.Exec(context.Background(), migrateHistory)
	if err != nil {
		log.Fatal(err)
	}
}

func (c *ClickhouseDB) Close() error {
	return c.db.Close()
}

func (c *ClickhouseDB) Db() clickhouse.Conn {
	return c.db
}
