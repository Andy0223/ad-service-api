package mongodb

import (
	"ad-service-api/internal/models"
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type AdvertisementRepository struct {
	mongoDB *mongo.Database
}

func NewAdvertisementRepository(mongoDB *mongo.Database) *AdvertisementRepository {
	return &AdvertisementRepository{
		mongoDB: mongoDB,
	}
}

func (repo *AdvertisementRepository) Store(ctx context.Context, ad *models.Advertisement) error {
	collection := repo.mongoDB.Collection("ads")
	_, err := collection.InsertOne(ctx, ad)
	if err != nil {
		return err
	}

	return nil
}

func (repo *AdvertisementRepository) GetActiveAdCounts(ctx context.Context, now time.Time) (int, error) {
	collection := repo.mongoDB.Collection("ads")
	filter := bson.M{
		"startAt": bson.M{"$lte": now},
		"endAt":   bson.M{"$gte": now},
	}

	count, err := collection.CountDocuments(ctx, filter)
	if err != nil {
		return 0, err
	}

	return int(count), nil
}

func (repo *AdvertisementRepository) Fetch(ctx context.Context, filter bson.M, limit, offset int) ([]*models.Advertisement, error) {
	collection := repo.mongoDB.Collection("ads")
	findOptions := options.Find().SetLimit(int64(limit)).SetSkip(int64(offset))
	cursor, err := collection.Find(ctx, filter, findOptions)
	if err != nil {
		return nil, err
	}

	var ads []*models.Advertisement
	if err = cursor.All(ctx, &ads); err != nil {
		return nil, err
	}

	return ads, nil
}
