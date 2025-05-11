package config

import "os"

func GetApplicationConfig() *ApplicationConfig {
	return &ApplicationConfig{
		MongoConfig: MongoConfig{ // TODO right env names
			Host:         os.Getenv("DB_MONGO_HOST"),
			Port:         os.Getenv("DB_MONGO_PORT"),
			User:         os.Getenv("DB_MONGO_USER"),
			Password:     os.Getenv("DB_MONGO_PASSWORD"),
			DatabaseName: os.Getenv("DB_MONGO_DATABASE_NAME"),
			CollectionsNames: []string{
				"authors",
				"comments",
				"file_meta",
				"video_previews",
			},
		},
		MinioConfig: MinioConfig{
			Host:       os.Getenv("MINIO_HOST"),
			Port:       os.Getenv("MINIO_PORT"),
			AccessKey:  os.Getenv("MINIO_ROOT_USER"),
			SecretKey:  os.Getenv("MINIO_ROOT_PASSWORD"),
			BucketName: os.Getenv("MINIO_BUCKET_NAME"),
		},
		ElasticConfig: ElasticConfig{
			Host: os.Getenv("ELASTIC_HOST"),
			Port: os.Getenv("ELASTIC_PORT"),
		},
	}
}

type ApplicationConfig struct {
	MongoConfig   MongoConfig
	MinioConfig   MinioConfig
	ElasticConfig ElasticConfig
}

type MongoConfig struct {
	Host             string
	Port             string
	User             string
	Password         string
	DatabaseName     string
	CollectionsNames []string
}

func (mongo *MongoConfig) Uri() string {
	return mongo.Host + ":" + mongo.Port
}

type MinioConfig struct {
	Host       string
	Port       string
	AccessKey  string
	SecretKey  string
	BucketName string
}

func (minio *MinioConfig) Uri() string {
	return minio.Host + ":" + minio.Port
}

type ElasticConfig struct {
	Host     string
	Port     string
	Username string
	Password string
}

func (elastic *ElasticConfig) Uri() string {
	return elastic.Host + ":" + elastic.Port
}

type ExternalResource interface {
	Uri() string
}
