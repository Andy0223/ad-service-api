package service

import (
	"ad-service-api/database"
	"ad-service-api/internal/advertisement/repository"
	"ad-service-api/internal/models"
	"ad-service-api/internal/validators"
	"context"
	"errors"
	"strconv"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type IAdvertisementService interface {
	CreateAd(ctx context.Context, ad *models.Advertisement) error
	GetAds(ctx context.Context, key string, validQueryParams map[string]string) ([]*models.Advertisement, error)
	CountActive(ctx context.Context, now time.Time) (int, error)
	Fetch(ctx context.Context, filter primitive.M, limit, offset int) ([]*models.Advertisement, error)
	GetByDate(ctx context.Context, today string) (int, error)
	IncrByDate(ctx context.Context, key string) error
	GetAdsByKey(ctx context.Context, key string) ([]*models.Advertisement, error)
	SetAdsByKey(ctx context.Context, key string, ads []*models.Advertisement, expiration time.Duration) error
	DeleteAdsCacheByPattern(ctx context.Context, pattern string) error
	IsAdExpired(ad []*models.Advertisement, now time.Time) bool
	DeleteAdById(ctx context.Context, id primitive.ObjectID) error
	GetAdById(ctx context.Context, id primitive.ObjectID) (*models.Advertisement, error)
	UpdateAdById(ctx context.Context, id primitive.ObjectID, ad *models.Advertisement) error
}

type AdvertisementService struct {
	adRepo      repository.IAdvertisementRepository
	adRedisRepo repository.IAdRedisRepository
}

func NewAdvertisementService(adRepo repository.IAdvertisementRepository, adRedisRepo repository.IAdRedisRepository) IAdvertisementService {
	return &AdvertisementService{
		adRepo:      adRepo,
		adRedisRepo: adRedisRepo,
	}
}

func (as *AdvertisementService) CreateAd(c context.Context, ad *models.Advertisement) error {
	now := time.Now()
	today := now.Format("2006-01-02")

	// Validate the advertisement fields
	if err := validators.CreateAndUpdateAdValueValidation(*ad); err != nil {
		return err
	}

	dailyAdCount, err := as.GetByDate(c, today)
	if err != nil {
		return err
	}

	// Ensure daily ad creation limit isn't exceeded
	if dailyAdCount >= 3000 {
		return errors.New("cannot create more ads today. Daily limit reached")
	}

	// Ensure total active ads limit isn't exceeded
	activeAdCount, err := as.CountActive(c, now)
	if err != nil {
		return err
	}

	if activeAdCount >= 1000 {
		return errors.New("cannot create more ads. Active ads limit reached")
	}

	// Ad passes all checks; proceed to add
	if err := as.adRepo.InsertAd(c, ad); err != nil {
		return err
	}

	// Increment the Redis counter for today's ads
	if err := as.IncrByDate(c, today); err != nil {
		return err
	}

	// Invalidate the cache for the list of ads
	if err := as.DeleteAdsCacheByPattern(c, "ads:*"); err != nil {
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
	count, err := as.adRedisRepo.GetByDate(ctx, today)
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (as *AdvertisementService) IncrByDate(ctx context.Context, key string) error {
	err := as.adRedisRepo.IncrByDate(ctx, key)
	if err != nil {
		return err
	}
	return nil
}

func (as *AdvertisementService) GetAdsByKey(ctx context.Context, key string) ([]*models.Advertisement, error) {
	ads, err := as.adRedisRepo.GetAdsByKey(ctx, key)
	if err != nil {
		return nil, err
	}
	return ads, nil
}

// SetAdsByKey sets the advertisements associated with the specified key in Redis.
func (as *AdvertisementService) SetAdsByKey(ctx context.Context, key string, ads []*models.Advertisement, expiration time.Duration) error {
	err := as.adRedisRepo.SetAdsByKey(ctx, key, ads, expiration)
	if err != nil {
		return err
	}
	return nil
}

// DeleteAdsByPattern deletes the advertisements associated with the specified pattern in Redis.
func (as *AdvertisementService) DeleteAdsCacheByPattern(ctx context.Context, pattern string) error {
	err := as.adRedisRepo.DeleteAdsCacheByPattern(ctx, pattern)
	if err != nil {
		return err
	}
	return nil
}

// FilterExpiredAds filters out expired ads from a slice of ads
func (as *AdvertisementService) IsAdExpired(ads []*models.Advertisement, now time.Time) bool {
	for _, ad := range ads {
		if now.After(ad.EndAt) {
			return true
		}
	}
	return false
}

// delete ad
func (as *AdvertisementService) DeleteAdById(c context.Context, oid primitive.ObjectID) error {
	ad, err := as.GetAdById(c, oid)
	if err != nil {
		return err
	}
	if ad == nil {
		return errors.New("advertisement not found")
	}
	if err := as.adRepo.DeleteAdById(c, oid); err != nil {
		return err
	}
	if err := as.DeleteAdsCacheByPattern(c, "ads:*"); err != nil {
		return err
	}
	return nil
}

// get ad by id
func (as *AdvertisementService) GetAdById(ctx context.Context, id primitive.ObjectID) (*models.Advertisement, error) {
	ad, err := as.adRepo.GetAdById(ctx, id)
	if err != nil {
		return nil, err
	}
	return ad, nil
}

// updtae ad
func (as *AdvertisementService) UpdateAdById(c context.Context, oid primitive.ObjectID, ad *models.Advertisement) error {
	advertisement, err := as.GetAdById(c, oid)
	if err != nil {
		return err
	}
	if advertisement == nil {
		return errors.New("advertisement not found")
	}
	if err := validators.CreateAndUpdateAdValueValidation(*ad); err != nil {
		return err
	}
	if err := as.adRepo.UpdateAdById(c, oid, ad); err != nil {
		return err
	}
	if err := as.DeleteAdsCacheByPattern(c, "ads:*"); err != nil {
		return err
	}
	return nil
}

func (as *AdvertisementService) GetAds(c context.Context, key string, validQueryParams map[string]string) ([]*models.Advertisement, error) {
	now := time.Now()

	// Try to get the result from Redis first
	result, _ := as.GetAdsByKey(c, key)

	// Check if the ad from redis is expired
	isAdexpired := as.IsAdExpired(result, now)

	if result == nil || isAdexpired {
		// Create a filter, limit, and offset based on the query parameters
		filter := database.CreateFilter(validQueryParams)
		limit, _ := strconv.Atoi(validQueryParams["limit"])
		offset, _ := strconv.Atoi(validQueryParams["offset"])

		// If the result is not in Redis, get it from the database
		filteredAds, err := as.Fetch(c, filter, limit, offset)
		if err != nil {
			return nil, err
		}

		// Store the result in Redis for future use
		err = as.SetAdsByKey(c, key, filteredAds, time.Hour)
		if err != nil {
			return nil, err
		}

		return filteredAds, nil
	} else {
		return result, nil
	}
}
