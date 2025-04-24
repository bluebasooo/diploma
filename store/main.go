package main

import (
	read_api "dev/bluebasooo/video-platform/api/read"
	write_api "dev/bluebasooo/video-platform/api/write"
	base_server "dev/bluebasooo/video-platform/base/server"
)

func main() {
	server := base_server.FustestServerUSee{}
	apiInitializer := []base_server.RouteInitializer{
		read_api.InitFileMetaApi,
		read_api.InitFileReadApi,
		write_api.InitFileWriteApi,
	}

	server.Init(apiInitializer)
	server.Start()
}
