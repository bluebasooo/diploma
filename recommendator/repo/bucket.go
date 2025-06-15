package repo

import (
	"context"
	"dev/bluebasooo/video-common/db"
	"dev/bluebasooo/video-recomendator/entity"
	"errors"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
	"maps"
	"sync"

	"go.mongodb.org/mongo-driver/bson"
)

type BucketRepo struct {
	db           *db.MongoDB
	bucketCache  map[string]entity.Bucket
	bucketLocked map[string]*sync.Mutex
	mu           sync.RWMutex
}

func (r *BucketRepo) fromCache(id string) (*entity.Bucket, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	val, ok := r.bucketCache[id]
	if !ok {
		return nil, false
	}

	return &val, true
}

func (r *BucketRepo) LockOnBucket(id string) bool {
	_, ok := r.bucketLocked[id]
	if !ok {
		r.bucketLocked[id] = &sync.Mutex{}
	}
	r.bucketLocked[id].Lock()
	return true
}

func (r *BucketRepo) UnlockOnBucket(id string) {
	r.bucketLocked[id].Unlock()
}

// if not found - no fail, but not return
func (r *BucketRepo) fromCacheBatch(ids ...string) []entity.Bucket {
	r.mu.RLock()
	defer r.mu.RUnlock()

	buckets := make([]entity.Bucket, 0, len(ids))
	for _, id := range ids {
		val, ok := r.bucketCache[id]
		if !ok {
			continue
		}
		buckets = append(buckets, val)
	}

	return buckets
}

func (r *BucketRepo) allToCache(buckets map[string]entity.Bucket) {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.bucketCache = maps.Clone(buckets)
}

func (r *BucketRepo) toCache(buckets ...entity.Bucket) {
	r.mu.Lock()
	defer r.mu.Unlock()

	for _, bucket := range buckets {
		r.bucketCache[bucket.ID] = bucket
	}
}

// maybe need deep copy
func (r *BucketRepo) copyCache() map[string]entity.Bucket {
	r.mu.RLock()
	defer r.mu.RUnlock()

	return maps.Clone(r.bucketCache)
}

func NewBucketRepo(db *db.MongoDB) *BucketRepo {
	return &BucketRepo{
		db:           db,
		bucketCache:  make(map[string]entity.Bucket),
		bucketLocked: make(map[string]*sync.Mutex),
		mu:           sync.RWMutex{},
	}
}

type BucketToAdd struct {
	CanAdd       bool
	DistToCenter float64
	DotId        string
}

// smth bad i think
func (r *BucketRepo) GetAllBuckets(ctx context.Context) (map[string]entity.Bucket, error) {
	if len(r.bucketCache) != 0 {
		return r.copyCache(), nil
	}

	collection := r.db.GetCollection("buckets")

	var buckets []entity.Bucket
	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}

	err = cursor.All(ctx, &buckets)
	if err != nil {
		return nil, err
	}

	if len(buckets) == 0 {
		bucket := entity.Bucket{
			ID:                       primitive.NewObjectID().Hex(),
			BucketDotsToDistToCenter: make(map[entity.VideoDotId]float64),
			BucketCenter:             make(map[entity.VideoId]float64),
			IsSeparated:              false,
		}
		err = r.UpsertBuckets(ctx, bucket)
		if err != nil {
			return nil, err
		}

		return r.copyCache(), nil
	}

	bucketById := make(map[string]entity.Bucket)
	for _, val := range buckets {
		bucketById[val.ID] = val
	}

	r.allToCache(bucketById)

	return bucketById, nil
}

func (r *BucketRepo) GetBucket(ctx context.Context, id string) (*entity.Bucket, error) {
	val, err := r.GetBuckets(ctx, id)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	if len(val) == 0 {
		log.Fatal("bucket not found")
		return nil, errors.New("bucket not found")
	}

	return &(val[0]), nil
}

func (r *BucketRepo) GetBuckets(ctx context.Context, ids ...string) ([]entity.Bucket, error) {
	if len(r.bucketCache) == 0 {
		_, err := r.GetAllBuckets(ctx)
		if err != nil {
			return nil, err
		}
	}

	buckets := r.fromCacheBatch(ids...)

	return buckets, nil
}

// all fields should be full
func (r *BucketRepo) UpsertBuckets(ctx context.Context, buckets ...entity.Bucket) error {
	collection := r.db.GetCollection("buckets")

	updates := make([]mongo.WriteModel, 0, len(buckets))

	for _, bucket := range buckets {
		filter := bson.M{"_id": bucket.ID}
		update := bson.M{"$set": bson.M{
			"dots_to_dist_to_center": bucket.BucketDotsToDistToCenter,
			"center":                 bucket.BucketCenter,
			"is_separated":           bucket.IsSeparated,
		},
		}

		model := mongo.NewUpdateOneModel().
			SetFilter(filter).
			SetUpdate(update).
			SetUpsert(true)

		updates = append(updates, model)
	}

	_, err := collection.BulkWrite(ctx, updates)
	if err != nil {
		return err
	}

	r.toCache(buckets...)

	return nil
}
