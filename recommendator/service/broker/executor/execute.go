package executor

import (
	"context"
	"encoding/json"
	"github.com/segmentio/kafka-go"
)

type Execute[F any] = func(payload *F) error

type Executor[F any] struct {
	dataConsumer *kafka.Reader
	exec         Execute[F]
}

func NewExecutor[F any](
	dataConsumer *kafka.Reader,
	exec Execute[F],
) *Executor[F] {
	return &Executor[F]{
		dataConsumer: dataConsumer,
		exec:         exec,
	}
}

func (ex *Executor[F]) Execute() error {
	message, err := ex.dataConsumer.FetchMessage(context.Background())
	if err != nil {
		return err
	}

	var value F
	err = json.Unmarshal(message.Value, &value)
	if err != nil {
		return err
	}

	err = ex.exec(&value)
	if err != nil {
		return err
	}

	err = ex.dataConsumer.CommitMessages(context.Background(), message)
	if err != nil {
		return err
	}

	return nil
}
