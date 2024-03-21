package handlers

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"ad-service-api/internal/config"
	"ad-service-api/internal/models"
	"ad-service-api/internal/repository"
	"ad-service-api/utils"
	"ad-service-api/validators"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"go.mongodb.org/mongo-driver/bson"
	"golang.org/x/net/context"
)

var (
	adRepo repository.AdvertisementRepository
	ctx    = context.Background()
)

func init() {
	adRepo = repository.NewAdvertisementRepository()
}

func CreateAdvertisement(c *gin.Context) {
	var ad models.Advertisement
	if err := c.ShouldBindJSON(&ad); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// Set the created at time
	ad.CreatedAt = time.Now()
	// Validate the advertisement fields
	if err := validateAdvertisementRequest(ad); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	now := time.Now()
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	todayStr := today.Format("2006-01-02")
	// Ensure daily ad creation limit isn't exceeded
	dailyAdCount, err := config.RDB.Get(ctx, "dailyAdCount:"+todayStr).Int()
	if err == redis.Nil {
		// The result is not in Redis, get it from the database
		dailyAdCount, _ = adRepo.CountAdsCreatedToday(today)
		// Store the result in Redis for 1 hour
		config.RDB.Set(ctx, "dailyAdCount:"+todayStr, dailyAdCount, time.Hour)
	} else if err != nil {
		// An error occurred while getting the result from Redis
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}
	// Ensure daily ad creation limit isn't exceeded
	if dailyAdCount >= 3000 {
		c.JSON(http.StatusForbidden, gin.H{"error": "Cannot create more ads today. Daily limit reached."})
		return
	}
	// Ensure total active ads limit isn't exceeded
	activeAdCount, err := config.RDB.Get(ctx, "activeAdCount").Int()
	if err == redis.Nil {
		// The result is not in Redis, get it from the database
		activeAdCount, _ = adRepo.CountActiveAds()
		// Store the result in Redis for 1 hour
		config.RDB.Set(ctx, "activeAdCount", activeAdCount, time.Hour)
	} else if err != nil {
		// An error occurred while getting the result from Redis
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}
	// Ensure total active ads limit isn't exceeded
	if activeAdCount >= 1000 {
		c.JSON(http.StatusForbidden, gin.H{"error": "Cannot create more ads. Active ads limit reached."})
		return
	}
	// Ad passes all checks; proceed to add
	if err := adRepo.CreateAdvertisement(ad); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create advertisement"})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "Advertisement created successfully"})
}

func ListAdvertisements(c *gin.Context) {
	filter, err := validateQueryParams(c.Request.URL.Query())
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid query parameters"})
		return
	}

	limit, offset := utils.GetPaginationParams(c)
	fmt.Println("Limit:", limit, "Offset:", offset) // 添加日志输出
	// List filtered advertisements
	filteredAds, err := adRepo.ListAdvertisements(filter, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to list advertisements"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"ads": filteredAds})
}

func validateQueryParams(query url.Values) (bson.M, error) {
	filter := bson.M{
		"startAt": bson.M{"$lte": time.Now()},
		"endAt":   bson.M{"$gte": time.Now()},
	}

	// Age condition validation
	if ageStr := query.Get("age"); ageStr != "" {
		age, err := strconv.Atoi(ageStr) // String to int
		if err != nil {
			return nil, fmt.Errorf("invalid age: %v", err)
		}
		if err := validators.ValidateAgeQueryParam(age); err != nil {
			return nil, err
		}

		filter["conditions.ageRange.ageStart"] = bson.M{"$lte": age}
		filter["conditions.ageRange.ageEnd"] = bson.M{"$gte": age}
	}

	// Gender condition validation
	if genders, ok := query["gender"]; ok && len(genders) > 0 {
		if err := validators.ValidateGenders(genders); err != nil {
			return nil, err
		}
		filter["conditions.genders"] = bson.M{"$in": genders}
	}

	// Country condition validation
	if countries, ok := query["country"]; ok && len(countries) > 0 {
		if err := validators.ValidateCountries(countries); err != nil {
			return nil, err
		}
		filter["conditions.countries"] = bson.M{"$in": countries}
	}

	// Platform condition validation
	if platforms, ok := query["platform"]; ok && len(platforms) > 0 {
		if err := validators.ValidatePlatforms(platforms); err != nil {
			return nil, err
		}
		filter["conditions.platforms"] = bson.M{"$in": platforms}
	}

	return filter, nil
}

func validateAdvertisementRequest(ad models.Advertisement) error {
	// Validate startAt and endAt
	if ad.StartAt.After(ad.EndAt) {
		return errors.New("startAt must be before endAt")
	}

	// Validate age range
	if err := validators.ValidateAgeRange(ad.Conditions.AgeRange.AgeStart, ad.Conditions.AgeRange.AgeEnd); err != nil {
		return err
	}

	// Validate genders
	if err := validators.ValidateGenders(ad.Conditions.Genders); err != nil {
		return err
	}

	// Validate countries
	if err := validators.ValidateCountries(ad.Conditions.Countries); err != nil {
		return err
	}

	// Validate platforms
	if err := validators.ValidatePlatforms(ad.Conditions.Platforms); err != nil {
		return err
	}

	return nil
}
