package service

import (
	"ad-service-api/internal/advertisement/repository"
	"ad-service-api/internal/models"
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type IAdvertisementService interface {
	Create(ctx context.Context, ad *models.Advertisement) error
	CountActive(ctx context.Context, now time.Time) (int, error)
	Fetch(ctx context.Context, filter primitive.M, limit, offset int) ([]*models.Advertisement, error)
	GetByDate(ctx context.Context, today string) (int, error)
	IncrByDate(ctx context.Context, key string) error
	GetAdsByKey(ctx context.Context, key string) ([]*models.Advertisement, error)
	SetAdsByKey(ctx context.Context, key string, ads []*models.Advertisement, expiration time.Duration) error
	DeleteAdsByPattern(ctx context.Context, pattern string) error
}

type AdvertisementService struct {
	adRepo      repository.IAdvertisementRepository
	adCountRepo repository.IAdCountRepository
}

func NewAdvertisementService(adRepo repository.IAdvertisementRepository, adCountRepo repository.IAdCountRepository) IAdvertisementService {
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

func (as *AdvertisementService) GetAdsByKey(ctx context.Context, key string) ([]*models.Advertisement, error) {
	ads, err := as.adCountRepo.GetAdsByKey(ctx, key)
	if err != nil {
		return nil, err
	}
	return ads, nil
}

// SetAdsByKey sets the advertisements associated with the specified key in Redis.
func (as *AdvertisementService) SetAdsByKey(ctx context.Context, key string, ads []*models.Advertisement, expiration time.Duration) error {
	err := as.adCountRepo.SetAdsByKey(ctx, key, ads, expiration)
	if err != nil {
		return err
	}
	return nil
}

// DeleteAdsByPattern deletes the advertisements associated with the specified pattern in Redis.
func (as *AdvertisementService) DeleteAdsByPattern(ctx context.Context, pattern string) error {
	err := as.adCountRepo.DeleteAdsByPattern(ctx, pattern)
	if err != nil {
		return err
	}
	return nil
}
