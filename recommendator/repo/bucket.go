package repo

import (
	"context"
	"dev/bluebasooo/video-common/db"
	"dev/bluebasooo/video-recomendator/entity"
	"go.mongodb.org/mongo-driver/mongo"
	"maps"
	"sync"

	"go.mongodb.org/mongo-driver/bson"
)

type BucketRepo struct {
	db          *db.MongoDB
	bucketCache map[string]entity.Bucket
	mu          sync.RWMutex
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

// if not found - no fail, but not return
func (r *BucketRepo) fromCacheBatch(ids ...string) []entity.Bucket {
	r.mu.RLock()
	defer r.mu.RUnlock()

	buckets := make([]entity.Bucket, len(ids))
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
		db:          db,
		bucketCache: make(map[string]entity.Bucket),
		mu:          sync.RWMutex{},
	}
}

func (r *BucketRepo) GetBucket(ctx context.Context, id string) (*entity.Bucket, error) {
	val, ok := r.fromCache(id)
	if ok {
		return val, nil
	}

	collection := r.db.GetCollection("buckets")

	var bucket entity.Bucket
	err := collection.FindOne(ctx, bson.M{"_id": id}).Decode(&bucket)
	if err != nil {
		return nil, err
	}
	r.toCache(bucket)

	return &bucket, nil
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

	cursor.All(ctx, &buckets)

	bucketById := make(map[string]entity.Bucket)
	for _, val := range buckets {
		bucketById[val.ID] = val
	}

	r.allToCache(bucketById)

	return bucketById, nil
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
