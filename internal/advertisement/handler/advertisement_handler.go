package handler

import (
	"ad-service-api/internal/models"
	"ad-service-api/internal/validators"
	"ad-service-api/utils"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// AdvertisementHandler handles advertisement-related requests
type AdvertisementHandler struct {
	AdvertisementService models.AdvertisementService
}

func NewAdvertisementHandler(e *gin.Engine, adService models.AdvertisementService) {
	handler := &AdvertisementHandler{
		AdvertisementService: adService,
	}

	e.POST("/api/v1/ad", handler.CreateAdHandler)
	e.GET("/api/v1/ads", handler.ListAdHandler)
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
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// Validate the advertisement fields
	if err := validators.CreateAdValueValidation(ad); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	dailyAdCount, _ := h.AdvertisementService.CountAdsCreatedToday(c, today)
	// Ensure daily ad creation limit isn't exceeded
	if dailyAdCount >= 3000 {
		c.JSON(http.StatusForbidden, gin.H{"error": "Cannot create more ads today. Daily limit reached."})
		return
	}
	// Ensure total active ads limit isn't exceeded
	activeAdCount, _ := h.AdvertisementService.CountActiveAds(c, now)
	// Ensure total active ads limit isn't exceeded
	if activeAdCount >= 1000 {
		c.JSON(http.StatusForbidden, gin.H{"error": "Cannot create more ads. Active ads limit reached."})
		return
	}
	// Ad passes all checks; proceed to add
	if err := h.AdvertisementService.CreateAdvertisement(c, &ad); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create advertisement"})
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
	filter, err := validators.ListAdParamsValidation(c.Request.URL.Query())
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid query parameters"})
		return
	}

	limit, offset := utils.GetPaginationParams(c)
	// List filtered advertisements
	filteredAds, err := h.AdvertisementService.ListAdvertisements(c, filter, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to list advertisements"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"ads": filteredAds})
}
