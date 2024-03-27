package repository

import (
	"ad-service-api/internal/models"
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

//go:generate mockery --name=IAdvertisementRepository --structname=MockAdvertisementRepository --output=mocks --dir=./internal/advertisement/repository --inpackage --with-expecter --testonly
type IAdvertisementRepository interface {
	Create(ctx context.Context, ad *models.Advertisement) error
	CountActive(ctx context.Context, now time.Time) (int, error)
	Fetch(ctx context.Context, filter bson.M, limit, offset int) ([]*models.Advertisement, error)
}

// AdvertisementRepositoryImpl implements the AdvertisementRepository interface.
type AdvertisementRepository struct {
	collection *mongo.Collection
}

// NewAdvertisementRepository creates a new instance of AdvertisementRepositoryImpl.
func NewAdvertisementRepository(collection *mongo.Collection) IAdvertisementRepository {
	return &AdvertisementRepository{
		collection: collection,
	}
}

// Create inserts a new advertisement document into the MongoDB collection.
func (r *AdvertisementRepository) Create(ctx context.Context, ad *models.Advertisement) error {
	_, err := r.collection.InsertOne(ctx, ad)
	return err
}

// CountActive returns the count of active advertisements based on the provided timestamp.
func (r *AdvertisementRepository) CountActive(ctx context.Context, now time.Time) (int, error) {
	filter := bson.M{
		"startAt": bson.M{"$lte": now},
		"endAt":   bson.M{"$gte": now},
	}

	count, err := r.collection.CountDocuments(ctx, filter)

	return int(count), err
}

// Fetch retrieves advertisements from the MongoDB collection based on the provided filter, limit, and offset.
func (r *AdvertisementRepository) Fetch(ctx context.Context, filter bson.M, limit, offset int) ([]*models.Advertisement, error) {
	findOptions := options.Find().SetLimit(int64(limit)).SetSkip(int64(offset))
	cursor, err := r.collection.Find(ctx, filter, findOptions)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var ads []*models.Advertisement
	if err := cursor.All(ctx, &ads); err != nil {
		return nil, err
	}

	return ads, nil
}
