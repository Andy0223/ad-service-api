package handler

import (
	"ad-service-api/internal/advertisement/service"
	"ad-service-api/internal/models"
	"ad-service-api/internal/validators"
	"ad-service-api/redis"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
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
	if err := c.ShouldBindJSON(&ad); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON format: " + err.Error()})
		return
	}
	err := h.AdvertisementService.CreateAd(c, &ad)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "Advertisement created successfully"})
}

// ListAdHandler lists all advertisements with optional query parameters
// @Summary List all advertisements with optional query parameters
// @Description Get a list of all advertisements with optional query parameters
// @ID get-ads
// @Produce  json
// @Success 200 {array} models.Advertisement
// @Router /api/v1/ad [get]
func (h *AdvertisementHandler) ListAdHandler(c *gin.Context) {
	// Validate the query parameters
	validQueryParams, err := validators.ListAdParamsValidation(c.Request.URL.Query())
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid query parameters: " + err.Error()})
		return
	}

	// Generate a unique key for this set of query parameters
	key := redis.GenerateRedisKey(validQueryParams)

	// Try to get the result from the service
	result, err := h.AdvertisementService.GetAds(c, key, validQueryParams)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to list advertisements: " + err.Error()})
		return
	}

	// Return the result
	c.JSON(http.StatusOK, gin.H{"ads": result})
}

// delete the advertisement with given id
// @Summary Delete the advertisement with given id
// @Description Delete the advertisement with given id
// @ID delete-ad
// @Produce  json
// @Param id path string true "Advertisement ID"
// @Success 200 {string} string
// @Router /api/v1/ad/{id} [delete]
func (h *AdvertisementHandler) DeleteAdHandler(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid advertisement ID"})
		return
	}
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid advertisement ID"})
		return
	}
	err = h.AdvertisementService.DeleteAdById(c, oid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete advertisement: " + err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Advertisement deleted successfully"})
}

// update the advertisement with given id
// @Summary Update the advertisement with given id
// @Description Update the advertisement with given id
// @ID update-ad
// @Accept  json
// @Produce  json
// @Param id path string true "Advertisement ID"
// @Param ad body models.Advertisement true "Update ad"
// @Success 200 {object} models.Advertisement
// @Router /api/v1/ad/{id} [put]
func (h *AdvertisementHandler) UpdateAdHandler(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid advertisement ID"})
		return
	}
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid advertisement ID"})
		return
	}
	var ad models.Advertisement
	if err := c.ShouldBindJSON(&ad); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON format: " + err.Error()})
		return
	}
	err = h.AdvertisementService.UpdateAdById(c, oid, &ad)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update advertisement: " + err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Advertisement updated successfully"})
}
