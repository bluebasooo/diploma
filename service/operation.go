package service

import (
	"crypto/rand"
	"encoding/hex"
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
		tasks: make(map[string]Task),
		mu:    sync.Mutex{},
	}
)

type Task struct {
	ID  string
	Ops map[string]Operation
}

type Operation struct {
	Hash   string
	Status UploadStatus
}

type TaskCache struct { // better to use redis
	tasks map[string]Task
	mu    sync.Mutex
}

func (cache *TaskCache) enqueue(hash string, opsSz int64) ([]string, error) {
	cache.mu.Lock()
	defer cache.mu.Unlock()

	task := Task{
		ID:  hash,
		Ops: make(map[string]Operation),
	}

	result := make([]string, 0, opsSz)
	for i := int64(0); i < opsSz; i++ {
		hash, err := GenerateRandom64CharHash()
		if err != nil {
			return nil, err
		}
		result = append(result, hash)
		task.Ops[hash] = Operation{
			Hash:   hash,
			Status: NonStarted,
		}
	}
	cache.tasks[hash] = task

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

func GenerateRandom64CharHash() (string, error) {
	// Создаем байтовый массив длиной 32 (так как каждый байт превратится в 2 hex символа)
	bytes := make([]byte, 32)

	// Заполняем массив криптографически стойкими случайными числами
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}

	// Конвертируем байты в 64-символьную hex строку
	return hex.EncodeToString(bytes), nil
}
