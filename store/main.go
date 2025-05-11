package main

import (
	read_api "dev/bluebasooo/video-platform/api/read"
	write_api "dev/bluebasooo/video-platform/api/write"
	base_server "dev/bluebasooo/video-platform/base/server"
	config2 "dev/bluebasooo/video-platform/config"
	"dev/bluebasooo/video-platform/db"
	obj_storage "dev/bluebasooo/video-platform/obj-storage"
	"dev/bluebasooo/video-platform/search"
	"dev/bluebasooo/video-platform/service"
	"github.com/joho/godotenv"
	"log"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	config := config2.GetApplicationConfig()

	mongo, _ := db.NewMongoDB(&config.MongoConfig)
	elastic, _ := search.NewElasticDB(&config.ElasticConfig)
	service.InitReindexer(elastic)
	minio := obj_storage.NewObjectStorage(&config.MinioConfig)

	service.InitRepos(mongo, elastic, minio)

	server := base_server.FustestServerUSee{}
	apiInitializer := []base_server.RouteInitializer{
		read_api.InitFileMetaApi,
		read_api.InitFileReadApi,
		read_api.InitCommentsReadApi,
		read_api.InitPreviewApi,
		read_api.InitSearchApi,
		read_api.InitReadAuthorApi,
		write_api.InitWriteFileApi,
		write_api.InitWriteCommentsApi,
		write_api.InitWriteAuthorsApi,
		write_api.InitWritePreviewApi,
	}

	server.Init(apiInitializer)
	server.Start()
}
