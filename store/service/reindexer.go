package service

import (
	"dev/bluebasooo/video-platform/search"
)

var limit = 1000
var queue = NewQueue[search.Instruction]()

type opType int

const (
	Index opType = iota
	Create
	Update
	Delete
)

func Shedule(opType opType, indexName string, id string, body search.Body) {
	instruction := search.Instruction{
		Meta: *resolveMeta(opType, indexName, id),
		Body: body,
	}
	queue.Push(instruction)

	go Execute()
}

func Execute() error {
	if queue.Size() < limit {
		return nil
	}

	instructions, err := queue.BatchPop(limit)
	if err != nil {
		return err
	}

	return searchRepo.BulkIndexEntities(instructions)
}

func resolveMeta(opType opType, indexName string, id string) *search.Meta {
	switch opType {
	case Index:
		return &search.Meta{Index: &search.IndexPreamble{Index: indexName, ID: id}}
	case Create:
		return &search.Meta{Create: &search.IndexPreamble{Index: indexName, ID: id}}
	case Update:
		return &search.Meta{Update: &search.DocIndexPreamble{Doc: search.IndexPreamble{Index: indexName, ID: id}}}
	case Delete:
		return &search.Meta{Delete: &search.IndexPreamble{Index: indexName, ID: id}}
	}
	return nil
}
