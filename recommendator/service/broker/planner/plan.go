package planner

import (
	"bytes"
	"context"
	"dev/bluebasooo/video-recomendator/service/broker/task "
	"encoding/json"
	"github.com/segmentio/kafka-go"
)

type Planner[F any] struct {
	dataProducer *kafka.Writer
	taskType     task.TaskType
}

func NewPlanner[F any](dataProducer *kafka.Writer) *Planner[F] {
	return &Planner[F]{
		dataProducer: dataProducer,
	}
}

func (p *Planner[F]) Plan(payload *F) error {
	writer := bytes.NewBuffer(make([]byte, 0))
	err := json.NewEncoder(writer).Encode(payload)
	message := kafka.Message{
		Topic: p.taskType,
		Value: writer.Bytes(),
	}

	err = p.dataProducer.WriteMessages(context.Background(), message)
	if err != nil {
		return err
	}

	return nil
}
