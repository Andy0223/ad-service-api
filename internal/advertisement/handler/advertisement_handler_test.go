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
	"go.mongodb.org/mongo-driver/bson/primitive"
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
			AgeStart: 18,
			AgeEnd:   24,
			Gender:   []string{"M", "F"},
			Country:  []string{"US", "JP"},
			Platform: []string{"ios", "web"},
		},
	}

	suite.mockAdService.On("CreateAd", mock.Anything, mock.AnythingOfType("*models.Advertisement")).Return(nil)

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
	// Create a response recorder
	w := httptest.NewRecorder()

	// Create a gin context
	c, _ := gin.CreateTestContext(w)

	c.Request, _ = http.NewRequest(http.MethodGet, "/api/v1/ad", nil)

	// Set up the mock service to expect a call to GetAds and return an empty list and nil error
	suite.mockAdService.On("GetAds", mock.Anything, mock.Anything, mock.Anything).Return([]*models.Advertisement{}, nil)

	suite.h.ListAdHandler(c)

	assert.Equal(suite.T(), http.StatusOK, w.Code)

	// Assert that the GetAds method was called once
	suite.mockAdService.AssertCalled(suite.T(), "GetAds", mock.Anything, mock.Anything, mock.Anything)
}

func (suite *AdvertisementHandlerSuite) TestAdvertisementHandler_DeleteAdHandler() {
	// Create a response recorder
	w := httptest.NewRecorder()

	// Create a gin context
	c, _ := gin.CreateTestContext(w)

	// Convert the string to a primitive.ObjectID
	id, _ := primitive.ObjectIDFromHex("1")
	c.Params = append(c.Params, gin.Param{Key: "id", Value: id.Hex()})

	// Set up the mock service to expect a call to DeleteAdById and return nil error
	suite.mockAdService.On("DeleteAdById", mock.Anything, id).Return(nil)

	suite.h.DeleteAdHandler(c)

	assert.Equal(suite.T(), http.StatusOK, w.Code)

	// Assert that the DeleteAd method was called once
	suite.mockAdService.AssertCalled(suite.T(), "DeleteAdById", mock.Anything, id)
}

func (suite *AdvertisementHandlerSuite) TestAdvertisementHandler_UpdateAdHandler() {
	now := time.Now().Round(time.Second)
	ad := &models.Advertisement{
		Title:   "Test Ad",
		StartAt: now,
		EndAt:   now.Add(24 * time.Hour), // Ends after 24 hours
		Conditions: models.Conditions{
			AgeStart: 18,
			AgeEnd:   24,
			Gender:   []string{"M", "F"},
			Country:  []string{"US", "JP"},
			Platform: []string{"ios", "web"},
		},
	}

	// Create a response recorder
	w := httptest.NewRecorder()

	// Create a gin context
	c, _ := gin.CreateTestContext(w)

	adJson, _ := json.Marshal(ad)
	c.Request, _ = http.NewRequest(http.MethodPut, "/api/v1/ad/1", bytes.NewBuffer(adJson))
	c.Request.Header.Set("Content-Type", "application/json")
	id, _ := primitive.ObjectIDFromHex("1")
	c.Params = append(c.Params, gin.Param{Key: "id", Value: id.Hex()})

	// Set up the mock service to expect a call to UpdateAdById and return nil error
	suite.mockAdService.On("UpdateAdById", mock.Anything, id, mock.AnythingOfType("*models.Advertisement")).Return(nil)

	suite.h.UpdateAdHandler(c)

	// Check the status code
	assert.Equal(suite.T(), http.StatusOK, w.Code)

	// Assert that the UpdateAdById method was called once
	suite.mockAdService.AssertCalled(suite.T(), "UpdateAdById", mock.Anything, id, mock.AnythingOfType("*models.Advertisement"))
}

func TestAdvertisementHandlerSuite(t *testing.T) {
	suite.Run(t, new(AdvertisementHandlerSuite))
}
