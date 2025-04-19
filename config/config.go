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
			URI:      os.Getenv("MINIO_URI"),
			User:     os.Getenv("MINIO_USER"),
			Password: os.Getenv("MINIO_PASSWORD"),
		},
	}
}

type ApplicationConfig struct {
	MongoConfig MongoConfig
	MinioConfig MinioConfig
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
