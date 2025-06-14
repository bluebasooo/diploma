package handler

import "sync"

type UpdateHandler struct {
	Producer chan bool
	cnt      int
	max      int
	mu       sync.Mutex
}

func NewUpdateHandler(max int) *UpdateHandler {
	return &UpdateHandler{
		make(chan bool, 1),
		0,
		max,
		sync.Mutex{},
	}
}

func (handler *UpdateHandler) Increment(plus int) {
	handler.mu.Lock()
	defer handler.mu.Unlock()

	handler.cnt += plus

	if handler.max < handler.cnt {
		handler.cnt = handler.cnt % handler.max
		handler.Producer <- true
	}
}
