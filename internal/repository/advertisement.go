package repository

import (
	"context"
	"fmt"
	"log"
	"time"

	"ad-service-api/internal/config"
	"ad-service-api/internal/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type AdvertisementRepository interface {
	CreateAdvertisement(ad models.Advertisement) error
	CountAdsCreatedToday(today time.Time) (int, error)
	CountActiveAds() (int, error)
	ListAdvertisements(filter bson.M, limit, offset int) ([]models.Advertisement, error)
}

type advertisementRepository struct{}

func NewAdvertisementRepository() AdvertisementRepository {
	return &advertisementRepository{}
}

func (repo *advertisementRepository) CreateAdvertisement(ad models.Advertisement) error {
	collection := config.DB.Collection("ads")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := collection.InsertOne(ctx, ad)
	return err
}

func (repo *advertisementRepository) CountAdsCreatedToday(today time.Time) (int, error) {
	collection := config.DB.Collection("ads")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.M{"createdAt": bson.M{"$gte": today}}

	count, err := collection.CountDocuments(ctx, filter)
	fmt.Println(count)
	if err != nil {
		log.Fatal(err)
	}

	return int(count), nil
}

func (repo *advertisementRepository) CountActiveAds() (int, error) {
	collection := config.DB.Collection("ads")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	now := time.Now()
	filter := bson.M{
		"startAt": bson.M{"$lte": now},
		"endAt":   bson.M{"$gte": now},
	}

	count, err := collection.CountDocuments(ctx, filter)
	if err != nil {
		log.Fatal(err)
	}

	return int(count), nil
}

func (repo *advertisementRepository) ListAdvertisements(filter bson.M, limit, offset int) ([]models.Advertisement, error) {
	collection := config.DB.Collection("ads")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	opts := options.Find().SetSort(bson.D{{Key: "endAt", Value: 1}})
	opts.SetLimit(int64(limit))
	opts.SetSkip(int64(offset))

	cursor, err := collection.Find(ctx, filter, opts)
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
