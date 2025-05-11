package service

import (
	"dev/bluebasooo/video-platform/db"
	obj_storage "dev/bluebasooo/video-platform/obj-storage"
	"dev/bluebasooo/video-platform/repo"
	"dev/bluebasooo/video-platform/search"
)

var authorRepository *repo.AuthorRepository
var commRepo *repo.CommentRepo
var fileMetaRepository *repo.FileMetaRepository
var previewRepo *repo.PreviewRepository
var searchRepo *repo.SearchRepo
var fileRepository *repo.FileRepository

func InitRepos(db *db.MongoDB, indexer *search.ElasticDB, storage *obj_storage.ObjectStorage) {
	authorRepository = repo.NewAuthorRepository(db)
	commRepo = repo.NewCommentRepo(db)
	fileMetaRepository = repo.NewFileMetaRepository(db)
	previewRepo = repo.NewPreviewRepository(db)

	searchRepo = repo.NewSearchRepo(indexer)

	fileRepository = repo.NewFileRepository(*storage, "video-parts")
}
