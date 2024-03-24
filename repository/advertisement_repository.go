package repository

import (
	"context"
	"errors"
	"log"
	"time"

	"ad-service-api/internal/config"
	"ad-service-api/internal/models"

	"github.com/go-redis/redis"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type AdvertisementRepository interface {
	CreateAdvertisement(ad models.Advertisement) error
	CountAdsCreatedToday() (int, error)
	CountActiveAds() (int, error)
	ListAdvertisements(filter bson.M, limit, offset int) ([]models.Advertisement, error)
}

type advertisementRepository struct {
	collection *mongo.Collection
}

func NewAdvertisementRepository() AdvertisementRepository {
	return &advertisementRepository{
		collection: config.DB.Collection("ads"),
	}
}

func (repo *advertisementRepository) CreateAdvertisement(ad models.Advertisement) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

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

func (repo *advertisementRepository) CountAdsCreatedToday() (int, error) {
	ctx := context.Background()
	// Get the count from Redis
	today := time.Now().Format("2006-01-02")
	count, err := config.RDB.Get(ctx, "ads:"+today).Int()
	if err == redis.Nil {
		return 0, errors.New("key not found")
	} else if err != nil {
		return 0, err
	}

	return count, nil
}

func (repo *advertisementRepository) CountActiveAds() (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	now := time.Now()
	filter := bson.M{
		"startAt": bson.M{"$lte": now},
		"endAt":   bson.M{"$gte": now},
	}

	count, err := repo.collection.CountDocuments(ctx, filter)
	if err != nil {
		log.Fatal(err)
	}

	return int(count), nil
}

func (repo *advertisementRepository) ListAdvertisements(filter bson.M, limit, offset int) ([]models.Advertisement, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	opts := options.Find().SetSort(bson.D{{Key: "endAt", Value: 1}})
	opts.SetLimit(int64(limit))
	opts.SetSkip(int64(offset))

	cursor, err := repo.collection.Find(ctx, filter, opts)
	if err != nil {
		log.Fatal(err)
	}
	defer cursor.Close(ctx)

	var ads []models.Advertisement
	if err = cursor.All(ctx, &ads); err != nil {
		log.Fatal(err)
	}

	return ads, nil
}
