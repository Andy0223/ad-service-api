package handler_test

import (
	"ad-service-api/internal/advertisement/handler"
	"ad-service-api/internal/models"
	"ad-service-api/mocks"
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type AdvertisementHandlerSuite struct {
	suite.Suite
	mockAdService *mocks.MockAdvertisementService
	h             *handler.AdvertisementHandler
}

func (suite *AdvertisementHandlerSuite) SetupTest() {
	suite.mockAdService = new(mocks.MockAdvertisementService)
	suite.h = handler.NewAdvertisementHandler(suite.mockAdService)
}

func (suite *AdvertisementHandlerSuite) TestAdvertisementHandler_CreateAdHandler() {
	now := time.Now().Round(time.Second)
	ad := &models.Advertisement{
		Title:   "Test Ad",
		StartAt: now,
		EndAt:   now.Add(24 * time.Hour), // Ends after 24 hours
		Conditions: models.Conditions{
			AgeRange: models.AgeRange{
				AgeStart: 18,
				AgeEnd:   24,
			},
			Genders:   []string{"M", "F"},
			Countries: []string{"US", "JP"},
			Platforms: []string{"ios", "web"},
		},
	}

	suite.mockAdService.On("GetByDate", mock.Anything, now.Format("2006-01-02")).Return(1, nil)
	suite.mockAdService.On("CountActive", mock.AnythingOfType("*gin.Context"), mock.AnythingOfType("time.Time")).Return(1, nil)
	suite.mockAdService.On("Create", mock.Anything, ad).Return(nil)
	suite.mockAdService.On("IncrByDate", mock.Anything, "ads:"+now.Format("2006-01-02")).Return(nil)

	// Create a response recorder
	w := httptest.NewRecorder()

	// Create a gin context
	c, _ := gin.CreateTestContext(w)

	adJson, _ := json.Marshal(ad)
	c.Request, _ = http.NewRequest(http.MethodPost, "/api/v1/ad", bytes.NewBuffer(adJson))
	c.Request.Header.Set("Content-Type", "application/json")

	suite.h.CreateAdHandler(c)

	assert.Equal(suite.T(), http.StatusCreated, w.Code)
	suite.mockAdService.AssertExpectations(suite.T())
}

func (suite *AdvertisementHandlerSuite) TestAdvertisementHandler_ListAdHandler() {
	expectedLimit := 5
	expectedOffset := 0

	// Provide some specific test ad data
	expectedAds := []*models.Advertisement{
		{Title: "Test Ad 1"},
		{Title: "Test Ad 2"},
	}

	// Set the call expectation for the mock method, return specific test data
	suite.mockAdService.On("Fetch", mock.AnythingOfType("*gin.Context"), mock.AnythingOfType("primitive.M"), expectedLimit, expectedOffset).Return(expectedAds, nil)

	// Create response recorder and gin context
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest(http.MethodGet, "/api/v1/ads", nil)

	// Call the handler
	suite.h.ListAdHandler(c)

	// Verify the response status code
	assert.Equal(suite.T(), http.StatusOK, w.Code)

	// Parse the response body, verify if the returned ad data meets expectations
	var adsResponse gin.H
	err := json.NewDecoder(w.Body).Decode(&adsResponse)
	assert.NoError(suite.T(), err)

	adsInterface, ok := adsResponse["ads"].([]interface{})
	assert.True(suite.T(), ok, "The type of the returned ad list should match the expected value")

	adsValue := make([]*models.Advertisement, len(adsInterface))
	for i, v := range adsInterface {
		bytes, err := json.Marshal(v)
		assert.NoError(suite.T(), err)

		ad := &models.Advertisement{}
		err = json.Unmarshal(bytes, ad)
		assert.NoError(suite.T(), err)

		adsValue[i] = ad
	}
	assert.Equal(suite.T(), expectedAds, adsValue, "The returned ad list should match the expected value")

	// Verify if the call expectation of the mock service is met
	suite.mockAdService.AssertExpectations(suite.T())
}

func TestAdvertisementHandlerSuite(t *testing.T) {
	suite.Run(t, new(AdvertisementHandlerSuite))
}
