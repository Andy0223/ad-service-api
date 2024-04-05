package handler

import (
	"ad-service-api/internal/advertisement/service"
	"ad-service-api/internal/models"
	"ad-service-api/internal/validators"
	"ad-service-api/utils"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type AdvertisementHandler struct {
	AdvertisementService service.IAdvertisementService
}

func NewAdvertisementHandler(adService service.IAdvertisementService) *AdvertisementHandler {
	return &AdvertisementHandler{
		AdvertisementService: adService,
	}
}

// CreateAdHandler creates a new advertisement
// @Summary Create new advertisement
// @Description Create new advertisement with the input payload
// @ID create-ad
// @Accept  json
// @Produce  json
// @Param ad body models.Advertisement true "Create ad"
// @Success 201 {object} models.Advertisement
// @Router /api/v1/ad [post]
func (h *AdvertisementHandler) CreateAdHandler(c *gin.Context) {
	var ad models.Advertisement
	now := time.Now()
	today := now.Format("2006-01-02")
	if err := c.ShouldBindJSON(&ad); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON format: " + err.Error()})
		return
	}
	// Validate the advertisement fields
	if err := validators.CreateAdValueValidation(ad); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid advertisement data: " + err.Error()})
		return
	}

	dailyAdCount, err := h.AdvertisementService.GetByDate(c, today)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get daily ad count: " + err.Error()})
		return
	}
	// Ensure daily ad creation limit isn't exceeded
	if dailyAdCount >= 3000 {
		c.JSON(http.StatusForbidden, gin.H{"error": "Cannot create more ads today. Daily limit reached."})
		return
	}
	// Ensure total active ads limit isn't exceeded
	activeAdCount, err := h.AdvertisementService.CountActive(c, now)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get active ad count: " + err.Error()})
		return
	}
	// Ensure total active ads limit isn't exceeded
	if activeAdCount >= 1000 {
		c.JSON(http.StatusForbidden, gin.H{"error": "Cannot create more ads. Active ads limit reached."})
		return
	}
	// Ad passes all checks; proceed to add
	if err := h.AdvertisementService.Create(c, &ad); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create advertisement: " + err.Error()})
		return
	}
	// Increment the Redis counter for today's ads
	if err := h.AdvertisementService.IncrByDate(c, today); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to increment ad count: " + err.Error()})
		return
	}

	// Invalidate the cache for the list of ads
	if err := h.AdvertisementService.DeleteAdsByPattern(c, "ads:*"); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to invalidate cache: " + err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Advertisement created successfully"})
}

// ListAdHandler lists all advertisements
// @Summary List all advertisements
// @Description Get a list of all advertisements
// @ID get-ads
// @Produce  json
// @Success 200 {array} models.Advertisement
// @Router /api/v1/ads [get]
func (h *AdvertisementHandler) ListAdHandler(c *gin.Context) {
	queryParams, err := validators.ListAdParamsValidation(c.Request.URL.Query())
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid query parameters: " + err.Error()})
		return
	}

	// Generate a unique key for this set of query parameters
	key := utils.GenerateRedisKey(queryParams)

	// Try to get the result from Redis first
	result, _ := h.AdvertisementService.GetAdsByKey(c, key)

	if result == nil {
		// Create a filter, limit, and offset based on the query parameters
		filter := utils.CreateFilter(queryParams)
		limit, _ := strconv.Atoi(queryParams["limit"])
		offset, _ := strconv.Atoi(queryParams["offset"])

		// If the result is not in Redis, get it from the database
		filteredAds, err := h.AdvertisementService.Fetch(c, filter, limit, offset)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to list advertisements: " + err.Error()})
			return
		}

		// Store the result in Redis for future use
		err = h.AdvertisementService.SetAdsByKey(c, key, filteredAds, 1000*time.Millisecond)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to cache advertisements: " + err.Error()})
			return
		}

		// Return the result
		c.JSON(http.StatusOK, gin.H{"ads": filteredAds})
	} else {
		// Return the result from Redis
		fmt.Println("result", result)
		c.JSON(http.StatusOK, gin.H{"ads": result})
	}
}
