package service

import (
	"ad-service-api/internal/advertisement/repository/mongodb"
	"ad-service-api/internal/advertisement/repository/redis"
	"ad-service-api/internal/models"
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Advertisement interface {
	Create(ctx context.Context, ad *models.Advertisement) error
	CountActive(ctx context.Context, now time.Time) (int, error)
	Fetch(ctx context.Context, filter primitive.M, limit, offset int) ([]*models.Advertisement, error)
	IncrByDate(ctx context.Context, key string) error
	GetByDate(ctx context.Context, key string) (int, error)
}

type AdvertisementService struct {
	adRepo      mongodb.AdvertisementRepository
	adCountRepo redis.AdCountRepository
}

func NewAdvertisementService(adRepo mongodb.AdvertisementRepository, adCountRepo redis.AdCountRepository) *AdvertisementService {
	return &AdvertisementService{
		adRepo:      adRepo,
		adCountRepo: adCountRepo,
	}
}

func (as *AdvertisementService) Create(ctx context.Context, ad *models.Advertisement) error {
	err := as.adRepo.Create(ctx, ad)
	if err != nil {
		return err
	}
	return nil
}

func (as *AdvertisementService) CountActive(ctx context.Context, now time.Time) (int, error) {
	count, err := as.adRepo.CountActive(ctx, now)
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

func (as *AdvertisementService) GetByDate(ctx context.Context, today string) (int, error) {
	count, err := as.adCountRepo.GetByDate(ctx, today)
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (as *AdvertisementService) IncrByDate(ctx context.Context, key string) error {
	err := as.adCountRepo.IncrByDate(ctx, key)
	if err != nil {
		return err
	}
	return nil
}
