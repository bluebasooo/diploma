package db

import (
	"context"
	"dev/bluebasooo/video-platform/config"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoDB struct {
	config      *config.MongoConfig
	client      *mongo.Client
	db          *mongo.Database
	collections map[string]*mongo.Collection
}

func NewMongoDB(config *config.MongoConfig) (*MongoDB, error) {
	client, err := connectToMongoDB(config)
	if err != nil {
		return nil, err
	}

	db := MongoDB{
		config: config,
		client: client,
	}

	db.initDb(config.DatabaseName)
	db.initAllCollections(config.CollectionsNames)

	return &db, nil
}

func connectToMongoDB(config *config.MongoConfig) (*mongo.Client, error) {
	opts := options.Client()
	opts.SetAuth(options.Credential{
		AuthSource: "admin",
		Username:   config.User,
		Password:   config.Password,
	})
	opts.ApplyURI(config.URI)

	client, err := mongo.Connect(context.TODO(), opts)
	if err != nil {
		return nil, err
	}

	return client, nil
}

// use this func to get the db - maybe in this we will ping
func (db *MongoDB) initDb(databaseName string) {
	if db.db == nil {
		db.db = db.client.Database(databaseName)
	}
}

func (db *MongoDB) GetCollection(name string) *mongo.Collection {
	_, ok := db.collections[name]
	if !ok && db.isCollectionSupported(name) {
		db.initAllCollections(db.config.CollectionsNames)
	}
	return db.collections[name]
}

func (db *MongoDB) isCollectionSupported(name string) bool {
	for _, collection := range db.config.CollectionsNames {
		if collection == name {
			return true
		}
	}
	return false
}

func (db *MongoDB) initAllCollections(collectionNames []string) {
	for _, collectionName := range collectionNames {
		db.collections[collectionName] = db.db.Collection(collectionName)
	}
}

func (db *MongoDB) Close() error {
	return db.client.Disconnect(context.TODO())
}
