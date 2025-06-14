package main

import (
	"dev/bluebasooo/video-common/db"
	base_server "dev/bluebasooo/video-common/server"
	"dev/bluebasooo/video-recomendator/api/read"
	"dev/bluebasooo/video-recomendator/api/write"
	"dev/bluebasooo/video-recomendator/config"
	db2 "dev/bluebasooo/video-recomendator/db"
	"dev/bluebasooo/video-recomendator/handler"
	"dev/bluebasooo/video-recomendator/repo"
	"dev/bluebasooo/video-recomendator/service"
	"github.com/joho/godotenv"
	"log"
)

var isLocalDevelopment = true

func main() {
	if isLocalDevelopment {
		err := godotenv.Load(".env")
		if err != nil {
			log.Fatal("Error loading .env file")
		}
	}

	appConfig := config.GetApplicationConfig()

	mongoDb, err := db.NewMongoDB(&appConfig.MongoConfig)
	if err != nil {
		log.Fatal(err)
	}
	ch := db2.NewClickhouseDB(&appConfig.ClickhouseConfig)
	service.BucketRepo = repo.NewBucketRepo(mongoDb)
	service.DotsRepo = repo.NewDotsRepo(mongoDb)
	service.MetricsRepo = repo.NewMetricsRepo(ch)
	service.HistoryRepo = repo.NewHistoryRepo(ch)
	service.UpdatesHandler = handler.NewUpdateHandler(100)

	go func() {
		service.Loop()
	}()

	server := base_server.FustestServerUSee{}
	apiInitializer := []base_server.RouteInitializer{
		read.InitPoolApi,
		read.InitHistoryApi,
		write.InitMetricsApi,
	}

	server.Init(apiInitializer)
	server.Start()
}
