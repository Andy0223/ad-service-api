package service

import (
	"ad-service-api/internal/models"
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type AdvertisementService struct {
	adRepo models.AdvertisementRepository
}

func NewAdvertisementService(adRepo models.AdvertisementRepository) models.AdvertisementService {
	return &AdvertisementService{
		adRepo: adRepo,
	}
}

func (as *AdvertisementService) CreateAdvertisement(ctx context.Context, ad *models.Advertisement) error {
	err := as.adRepo.CreateAdvertisement(ctx, ad)
	if err != nil {
		return err
	}
	return nil
}

func (as *AdvertisementService) CountAdsCreatedToday(ctx context.Context, today string) (int, error) {
	count, err := as.adRepo.CountAdsCreatedToday(ctx, today)
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (as *AdvertisementService) CountActiveAds(ctx context.Context, now time.Time) (int, error) {
	count, err := as.adRepo.CountActiveAds(ctx, now)
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (as *AdvertisementService) ListAdvertisements(ctx context.Context, filter primitive.M, limit, offset int) ([]*models.Advertisement, error) {
	ads, err := as.adRepo.ListAdvertisements(ctx, filter, limit, offset)
	if err != nil {
		return nil, err
	}
	return ads, nil
}
