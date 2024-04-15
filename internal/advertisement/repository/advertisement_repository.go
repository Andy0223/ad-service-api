package repository

import (
	"ad-service-api/internal/models"
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

//go:generate mockery --name=IAdvertisementRepository --structname=MockAdvertisementRepository --output=mocks --dir=./internal/advertisement/repository --inpackage --with-expecter --testonly
type IAdvertisementRepository interface {
	InsertAd(ctx context.Context, ad *models.Advertisement) error
	CountActive(ctx context.Context, now time.Time) (int, error)
	Fetch(ctx context.Context, filter bson.M, limit, offset int) ([]*models.Advertisement, error)
	DeleteAdById(ctx context.Context, id primitive.ObjectID) error
	GetAdById(ctx context.Context, id primitive.ObjectID) (*models.Advertisement, error)
	UpdateAdById(ctx context.Context, id primitive.ObjectID, ad *models.Advertisement) error
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
func (r *AdvertisementRepository) InsertAd(ctx context.Context, ad *models.Advertisement) error {
	_, err := r.collection.InsertOne(ctx, ad)
	if err != nil {
		return fmt.Errorf("failed to insert advertisement: %w", err)
	}
	return nil
}

// CountActive returns the count of active advertisements based on the provided timestamp.
func (r *AdvertisementRepository) CountActive(ctx context.Context, now time.Time) (int, error) {
	filter := bson.M{
		"startAt": bson.M{"$lte": now},
		"endAt":   bson.M{"$gte": now},
	}

	count, err := r.collection.CountDocuments(ctx, filter)
	if err != nil {
		return 0, fmt.Errorf("failed to count active advertisements: %w", err)
	}

	return int(count), nil
}

// Fetch retrieves advertisements from the MongoDB collection based on the provided filter, limit, and offset.
func (r *AdvertisementRepository) Fetch(ctx context.Context, filter bson.M, limit, offset int) ([]*models.Advertisement, error) {
	findOptions := options.Find().SetLimit(int64(limit)).SetSkip(int64(offset)).SetSort(bson.D{{Key: "endAt", Value: 1}})
	cursor, err := r.collection.Find(ctx, filter, findOptions)
	if err != nil {
		return nil, fmt.Errorf("failed to find advertisements: %w", err)
	}
	defer cursor.Close(ctx)

	var ads []*models.Advertisement
	if err := cursor.All(ctx, &ads); err != nil {
		return nil, fmt.Errorf("failed to decode advertisements: %w", err)
	}

	return ads, nil
}

// delete ads from the MongoDB collection based on the provided filter
func (r *AdvertisementRepository) DeleteAdById(ctx context.Context, id primitive.ObjectID) error {
	filter := bson.M{"_id": id}
	_, err := r.collection.DeleteOne(ctx, filter)
	if err != nil {
		return fmt.Errorf("failed to delete advertisements: %w", err)
	}
	return nil
}

// GetAdById retrieves an advertisement from the MongoDB collection based on the provided ID.
func (r *AdvertisementRepository) GetAdById(ctx context.Context, id primitive.ObjectID) (*models.Advertisement, error) {
	filter := bson.M{"_id": id}

	var ad models.Advertisement
	if err := r.collection.FindOne(ctx, filter).Decode(&ad); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to find advertisement: %w", err)
	}

	return &ad, nil
}

// UpdateAdById updates an advertisement in the MongoDB collection based on the provided ID.
func (r *AdvertisementRepository) UpdateAdById(ctx context.Context, id primitive.ObjectID, ad *models.Advertisement) error {
	filter := bson.M{"_id": id}
	update := bson.M{"$set": ad}

	_, err := r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return fmt.Errorf("failed to update advertisement: %w", err)
	}

	return nil
}
