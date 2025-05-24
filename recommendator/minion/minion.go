package minion

import (
	"context"
	mongo "dev/bluebasooo/video-common/db"
	"dev/bluebasooo/video-recomendator/db"
	"dev/bluebasooo/video-recomendator/queue"
	"dev/bluebasooo/video-recomendator/repo"
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
		indexRepo:      repo.NewIndexRepo(ch),
		historyRepo:    repo.NewHistoryRepo(ch),
	}
}

func NewMiniMinion(kafkaConsumer *queue.KafkaConsumer, mongo *mongo.MongoDB) *MiniMinion {
	return &MiniMinion{
		kafkaConsumer: kafkaConsumer,
		mongo:         mongo,
	}
}

func (g *GrandMinion) PlanUpdate(ctx context.Context) error {
	updates, err := g.Pull(ctx)
	if err != nil {
		return err
	}


}

func (g *MiniMinion) ExecuteUpdate(ctx context.Context) error {

}
