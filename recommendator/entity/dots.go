package entity

import (
	"time"
)

type DotHistory struct {
	ID       DotHistoryVersionedId `bson:"_id"`
	History  map[string]float64    `bson:"history"`
	BucketID string                `bson:"bucket_id"`
}

func (d *DotHistory) GetDotID() string {
	return d.ID.DotID
}

func (d *DotHistory) GetValue(id string) float64 {
	val, ok := d.History[id]
	if !ok {
		return 0
	}
	return val
}

type DotHistoryVersionedId struct {
	DotID      string    `bson:"id"`
	DateUpdate time.Time `bson:"date_update"`
}
