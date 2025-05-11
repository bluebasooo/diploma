package service

import (
	"dev/bluebasooo/video-platform/utils"
	"sync"
)

type UploadStatus int

const (
	Ok UploadStatus = iota
	Processing
	Error
	NonStarted
)

var (
	cacheTasks = TaskCache{
		tasks: make(map[string]*Task),
		mu:    sync.Mutex{},
	}
)

type Task struct {
	ID  string
	Ops map[string]*Operation
}

type Operation struct {
	Hash   string
	Status UploadStatus
}

type TaskCache struct { // better to use redis
	tasks map[string]*Task
	mu    sync.Mutex
}

func (cache *TaskCache) enqueue(taskId string, opsSz int64) ([]string, error) {
	cache.mu.Lock()
	defer cache.mu.Unlock()

	task := Task{
		ID:  taskId,
		Ops: make(map[string]*Operation),
	}

	result := make([]string, 0, opsSz)
	for i := int64(0); i < opsSz; i++ {
		hash := utils.GetRandomHashId()
		result = append(result, hash)
		task.Ops[hash] = &Operation{
			Hash:   hash,
			Status: NonStarted,
		}
	}
	cache.tasks[taskId] = &task

	return result, nil
}

func (cache *TaskCache) doneOp(TaskID string, hashOp string) {
	cache.mu.Lock()
	defer cache.mu.Unlock()

	task, ok := cache.tasks[TaskID]
	if !ok {
		return
	}

	op, ok := task.Ops[hashOp]
	if !ok {
		return
	}
	op.Status = Ok
}

func (cache *TaskCache) doneTask(taskID string) {
	cache.mu.Lock()
	defer cache.mu.Unlock()

	task, ok := cache.tasks[taskID]
	if !ok {
		return
	}

	for _, op := range task.Ops {
		if op.Status != Ok {
			return
		}
	}

	delete(cache.tasks, taskID)
}

func (cache *TaskCache) status(taskID string, hash string) UploadStatus {
	cache.mu.Lock()
	defer cache.mu.Unlock()

	task, ok := cache.tasks[hash]
	if !ok {
		return NonStarted
	}

	return task.Ops[hash].Status
}

func (cache *TaskCache) checkOfDone(taskID string) bool {
	cache.mu.Lock()
	defer cache.mu.Unlock()

	task, ok := cache.tasks[taskID]
	if !ok {
		return true // task is done
	}

	for _, op := range task.Ops {
		if op.Status != Ok {
			return false
		}
	}

	return true
}

func (cache *TaskCache) err(taskID string, hash string) {
	cache.mu.Lock()
	defer cache.mu.Unlock()

	task, ok := cache.tasks[taskID]
	if !ok {
		return
	}

	op, ok := task.Ops[hash]
	if !ok {
		return
	}
	op.Status = Error
}
