package broker

import (
	"dev/bluebasooo/video-recomendator/service/broker/executor"
	"dev/bluebasooo/video-recomendator/service/broker/planner"
	"github.com/segmentio/kafka-go"
)

type Broker[F any] interface {
	AsyncExecution(payload *F) error
	EventLoop(errs chan<- error)
}

type BrokerKafka[F any] struct {
	plan    *planner.Planner[F]
	execute *executor.Executor[F]
}

func NewBrokerKafka[F any](writer *kafka.Writer, reader *kafka.Reader, execute executor.Execute[F]) Broker[F] {
	return &BrokerKafka[F]{
		plan:    planner.NewPlanner[F](writer),
		execute: executor.NewExecutor[F](reader, execute),
	}
}

func (b *BrokerKafka[F]) AsyncExecution(payload *F) error {
	err := b.plan.Plan(payload)
	if err != nil {
		return err
	}

	return nil
}

func (b *BrokerKafka[F]) EventLoop(errs chan<- error) {
	for {
		errs <- b.execute.Execute()
	}
}

type DummyBroker[F any] struct {
	queued  chan *F
	execute executor.Execute[F]
}

func NewDummyBroker[F any](exec executor.Execute[F]) Broker[F] {
	return &DummyBroker[F]{
		queued:  make(chan *F),
		execute: exec,
	}
}

func (b *DummyBroker[F]) AsyncExecution(payload *F) error {
	b.queued <- payload
	return nil
}

func (b *DummyBroker[F]) EventLoop(errs chan<- error) {
	for {
		payload := <-b.queued
		errs <- b.execute(payload)
	}
}
