package service

import (
	"ad-service-api/internal/models"
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type MongodbAdRepository interface {
	Store(ctx context.Context, ad *models.Advertisement) error
	GetActiveAdCounts(ctx context.Context, now time.Time) (int, error)
	Fetch(ctx context.Context, filter primitive.M, limit, offset int) ([]*models.Advertisement, error)
}

type RedisAdRepository interface {
	IncrAdCountsByDate(ctx context.Context, key string) error
	GetAdCountsByDate(ctx context.Context, key string) (int, error)
}

type AdvertisementService struct {
	adRepo    MongodbAdRepository
	redisRepo RedisAdRepository
}

func NewAdvertisementService(adRepo MongodbAdRepository, redisRepo RedisAdRepository) *AdvertisementService {
	return &AdvertisementService{
		adRepo:    adRepo,
		redisRepo: redisRepo,
	}
}

func (as *AdvertisementService) Store(ctx context.Context, ad *models.Advertisement) error {
	err := as.adRepo.Store(ctx, ad)
	if err != nil {
		return err
	}
	return nil
}

func (as *AdvertisementService) GetActiveAdCounts(ctx context.Context, now time.Time) (int, error) {
	count, err := as.adRepo.GetActiveAdCounts(ctx, now)
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (as *AdvertisementService) Fetch(ctx context.Context, filter primitive.M, limit, offset int) ([]*models.Advertisement, error) {
	ads, err := as.adRepo.Fetch(ctx, filter, limit, offset)
	if err != nil {
		return nil, err
	}
	return ads, nil
}

func (as *AdvertisementService) GetAdCountsByDate(ctx context.Context, today string) (int, error) {
	count, err := as.redisRepo.GetAdCountsByDate(ctx, today)
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (as *AdvertisementService) IncrAdCountsByDate(ctx context.Context, key string) error {
	err := as.redisRepo.IncrAdCountsByDate(ctx, key)
	if err != nil {
		return err
	}
	return nil
}
