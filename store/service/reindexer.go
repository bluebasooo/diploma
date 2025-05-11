package service

import (
	"bytes"
	"context"
	"dev/bluebasooo/video-platform/search"
	"encoding/json"
	"fmt"
	"github.com/elastic/go-elasticsearch/v8/esutil"
	"log"
	"time"
)

var bulkIndexer esutil.BulkIndexer

type opType string

const (
	Index  opType = "index"
	Create        = "create"
	Update        = "update"
	Delete        = "delete"
)

func InitReindexer(es *search.ElasticDB) esutil.BulkIndexer {
	bi, _ := esutil.NewBulkIndexer(esutil.BulkIndexerConfig{
		Client:        es.GetClient(),
		NumWorkers:    25,
		FlushInterval: 1 * time.Second,
	})

	bulkIndexer = bi

	return bi
}

func Schedule[T any](opType opType, indexName string, id string, bodies []T) {
	items := make([]esutil.BulkIndexerItem, 0, len(bodies))

	for _, body := range bodies {
		items = append(items, *buildBulkItem(opType, indexName, id, body))
	}

	for _, item := range items {
		err := bulkIndexer.Add(context.Background(), item)
		if err != nil {
			log.Println(err)
		}
	}

	fmt.Print("hello")
}

func buildBulkItem(oper opType, indexName string, id string, body any) *esutil.BulkIndexerItem {
	data, _ := json.Marshal(body)
	return &esutil.BulkIndexerItem{
		Index:      indexName,
		Action:     string(oper),
		DocumentID: id,
		Body:       bytes.NewReader(data),
	}
}
