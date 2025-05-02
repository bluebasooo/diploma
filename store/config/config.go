package config

import "os"

func GetApplicationConfig() *ApplicationConfig {
	return &ApplicationConfig{
		MongoConfig: MongoConfig{ // TODO right env names
			URI:              os.Getenv("MONGO_URI"),
			User:             os.Getenv("MONGO_USER"),
			Password:         os.Getenv("MONGO_PASSWORD"),
			DatabaseName:     os.Getenv("MONGO_DATABASE"),
			CollectionsNames: []string{},
		},
		MinioConfig: MinioConfig{
			URI:        os.Getenv("MINIO_URI"),
			AccessKey:  os.Getenv("MINIO_ACCESS_KEY"),
			SecretKey:  os.Getenv("MINIO_SECRET_KEY"),
			BucketName: os.Getenv("MINIO_BUCKET_NAME"),
		},
		ElasticConfig: ElasticConfig{
			URI: os.Getenv("ELASTIC_URI"),
		},
	}
}

type ApplicationConfig struct {
	MongoConfig   MongoConfig
	MinioConfig   MinioConfig
	ElasticConfig ElasticConfig
}

type MongoConfig struct {
	URI              string
	User             string
	Password         string
	DatabaseName     string
	CollectionsNames []string
}

type MinioConfig struct {
	URI        string
	AccessKey  string
	SecretKey  string
	BucketName string
}

type ElasticConfig struct {
	URI      string
	Username string
	Password string
}
