package service

import (
	"errors"
	"sync"
)

type Queue[T any] struct {
	items []T
	mu    sync.Mutex
}

func NewQueue[T any]() *Queue[T] {
	return &Queue[T]{
		items: make([]T, 0),
		mu:    sync.Mutex{},
	}
}

func (q *Queue[T]) Push(item T) {
	q.mu.Lock()
	defer q.mu.Unlock()

	q.items = append(q.items, item)
}

func (q *Queue[T]) BatchPop(limit int) ([]T, error) {
	q.mu.Lock()
	defer q.mu.Unlock()

	if len(q.items) == 0 {
		return nil, errors.New("stack is empty")
	}

	items := q.items[:limit]
	q.items = q.items[limit:]

	return items, nil
}

func (q *Queue[T]) Size() int {
	q.mu.Lock()
	defer q.mu.Unlock()
	return len(q.items)
}
