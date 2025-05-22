package repo

import (
	mongo "dev/bluebasooo/video-common/db"
)

type DotsRepo struct {
	db *mongo.MongoDB
}

func NewDotsRepo(db *mongo.MongoDB) *DotsRepo {
	return &DotsRepo{db: db}
}

func 