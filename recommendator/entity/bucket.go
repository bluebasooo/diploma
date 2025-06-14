package entity

type VideoDotId = string
type VideoId = string

type Bucket struct {
	ID                       string                 `bson:"_id"`
	BucketDotsToDistToCenter map[VideoDotId]float64 `bson:"dots_to_dist_to_center"`
	BucketCenter             map[VideoId]float64    `bson:"center"`
}
