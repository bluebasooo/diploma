package minion

import (
	"context"
	mongo "dev/bluebasooo/video-common/db"
	"dev/bluebasooo/video-recomendator/db"
	"dev/bluebasooo/video-recomendator/queue"
)

type GrandMinion struct {
	kafkaProducer *queue.KafkaProducer
	indexRepo     *repo.IndexRepo
	historyRepo   *repo.HistoryRepo
}

type MiniMinion struct {
	kafkaConsumer *queue.KafkaConsumer
	repo          *repo.DotsRepo
}

func NewGrandMinion(kafkaProducer *queue.KafkaProducer, ch *db.ClickhouseDB) *GrandMinion {
	return &GrandMinion{
		kafkaProducer: kafkaProducer,
		ch:            ch,
	}
}

func NewMiniMinion(kafkaConsumer *queue.KafkaConsumer, mongo *mongo.MongoDB) *MiniMinion {
	return &MiniMinion{
		kafkaConsumer: kafkaConsumer,
		mongo:         mongo,
	}
}

func (g *GrandMinion) PlanUpdate(ctx context.Context) error {
	g.ch.Db().Query(ctx context.Context, query string, args ...any)
}

func (g *MiniMinion) ExecuteUpdate(ctx context.Context) error {

}
