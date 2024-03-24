package repository

import (
	"ad-service-api/internal/config"
	"ad-service-api/internal/models"
	"context"
	"errors"
	"time"

	"github.com/go-redis/redis/v8"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type AdvertisementRepository struct {
	collection  *mongo.Collection
	redisClient *redis.Client
}

func NewAdvertisementRepository(collection *mongo.Collection, redisClient *redis.Client) *AdvertisementRepository {
	return &AdvertisementRepository{
		collection:  collection,
		redisClient: redisClient,
	}
}

func (repo *AdvertisementRepository) CreateAdvertisement(ctx context.Context, ad *models.Advertisement) error {
	_, err := repo.collection.InsertOne(ctx, ad)
	if err != nil {
		return err
	}

	// Increment the Redis counter for today's ads
	today := time.Now().Format("2006-01-02")
	err = config.RDB.Incr(ctx, "ads:"+today).Err()
	if err != nil {
		return err
	}

	return nil
}

func (repo *AdvertisementRepository) CountAdsCreatedToday(ctx context.Context, today string) (int, error) {
	count, err := config.RDB.Get(ctx, "ads:"+today).Int()
	if err == redis.Nil {
		return 0, errors.New("key not found")
	}
	return count, nil
}

func (repo *AdvertisementRepository) CountActiveAds(ctx context.Context, now time.Time) (int, error) {
	filter := bson.M{
		"startAt": bson.M{"$lte": now},
		"endAt":   bson.M{"$gte": now},
	}

	count, err := repo.collection.CountDocuments(ctx, filter)
	if err != nil {
		return 0, err
	}

	return int(count), nil
}

func (repo *AdvertisementRepository) ListAdvertisements(ctx context.Context, filter bson.M, limit, offset int) ([]*models.Advertisement, error) {
	findOptions := options.Find().SetLimit(int64(limit)).SetSkip(int64(offset))
	cursor, err := repo.collection.Find(ctx, filter, findOptions)
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
