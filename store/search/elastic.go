package search

import (
	"bytes"
	"context"
	"dev/bluebasooo/video-platform/config"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/elastic/go-elasticsearch/v8"
)

type ElasticDB struct {
	config *config.ElasticConfig
	client *elasticsearch.Client
}

func NewElasticDB(config *config.ElasticConfig) (*ElasticDB, error) {
	client, err := connectToElasticDB(config)
	if err != nil {
		return nil, err
	}

	return &ElasticDB{
		config: config,
		client: client,
	}, nil
}

func connectToElasticDB(config *config.ElasticConfig) (*elasticsearch.Client, error) {
	address := fmt.Sprintf("http://%s", config.Uri())
	cfg := elasticsearch.Config{
		Addresses: []string{address},
		Username:  config.Username,
		Password:  config.Password,
	}

	client, err := elasticsearch.NewClient(cfg)
	if err != nil {
		return nil, err
	}

	return client, nil
}

func (db *ElasticDB) GetClient() *elasticsearch.Client {
	return db.client
}

func (db *ElasticDB) IndexEntity(indexName string, id string, entity interface{}) error {
	es := db.GetClient()

	json, err := json.Marshal(entity)
	if err != nil {
		return err
	}

	resp, err := es.Index(
		indexName,
		strings.NewReader(string(json)),
		es.Index.WithDocumentID(id),
		es.Index.WithRefresh("true"),
		es.Index.WithPretty(),
	)
	defer resp.Body.Close()

	if err != nil {
		return err
	}

	if resp.IsError() {
		return fmt.Errorf("Error indexing entity: %s", resp.Status())
	}

	return nil
}

func (db *ElasticDB) BulkIndexEntities(entities []Instruction) error {
	es := db.GetClient()

	reader, err := prepareEntities(entities)
	if err != nil {
		return err
	}

	resp, err := es.Bulk(reader, es.Bulk.WithContext(context.Background()))
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	if resp.IsError() {
		return fmt.Errorf("Error bulk indexing entities: %s", resp.Status())
	}

	return nil
}

func prepareEntities(entities []Instruction) (*bytes.Reader, error) {
	bulk := BulkRequest{
		Instructions: entities,
	}

	return bulk.toReader(), nil
}
