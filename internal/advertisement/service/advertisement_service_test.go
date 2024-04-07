package service_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"ad-service-api/internal/advertisement/service"
	"ad-service-api/internal/models"
	"ad-service-api/mocks"
)

type AdvertisementServiceSuite struct {
	suite.Suite
	mockAdRepo      *mocks.MockAdvertisementRepository
	mockAdRedisRepo *mocks.MockAdRedisRepository
	s               service.IAdvertisementService
	ctx             context.Context
}

func (suite *AdvertisementServiceSuite) SetupTest() {
	suite.mockAdRepo = new(mocks.MockAdvertisementRepository)
	suite.mockAdRedisRepo = new(mocks.MockAdRedisRepository)
	suite.s = service.NewAdvertisementService(suite.mockAdRepo, suite.mockAdRedisRepo)
	suite.ctx = context.TODO()
}

func (suite *AdvertisementServiceSuite) TestAdvertisementService_Create() {
	ad := &models.Advertisement{}

	suite.mockAdRepo.On("Create", suite.ctx, ad).Return(nil)

	err := suite.s.Create(suite.ctx, ad)

	assert.NoError(suite.T(), err)
	suite.mockAdRepo.AssertExpectations(suite.T())
}

func (suite *AdvertisementServiceSuite) TestAdvertisementService_CountActive() {
	now := time.Now()

	suite.mockAdRepo.On("CountActive", suite.ctx, now).Return(1, nil)

	count, err := suite.s.CountActive(suite.ctx, now)

	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), 1, count)
	suite.mockAdRepo.AssertExpectations(suite.T())
}

func (suite *AdvertisementServiceSuite) TestAdvertisementService_Fetch() {
	filter := primitive.M{}
	limit := 10
	offset := 0

	suite.mockAdRepo.On("Fetch", suite.ctx, filter, limit, offset).Return([]*models.Advertisement{}, nil)

	ads, err := suite.s.Fetch(suite.ctx, filter, limit, offset)

	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), ads)
	suite.mockAdRepo.AssertExpectations(suite.T())
}

func (suite *AdvertisementServiceSuite) TestAdvertisementService_GetByDate() {
	today := "2022-01-01"

	suite.mockAdRedisRepo.On("GetByDate", suite.ctx, today).Return(1, nil)

	count, err := suite.s.GetByDate(suite.ctx, today)

	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), 1, count)
	suite.mockAdRedisRepo.AssertExpectations(suite.T())
}

func (suite *AdvertisementServiceSuite) TestAdvertisementService_IncrByDate() {
	key := "2022-01-01"

	suite.mockAdRedisRepo.On("IncrByDate", suite.ctx, key).Return(nil)

	err := suite.s.IncrByDate(suite.ctx, key)

	assert.NoError(suite.T(), err)
	suite.mockAdRedisRepo.AssertExpectations(suite.T())
}

func (suite *AdvertisementServiceSuite) TestAdvertisementService_GetAdsByKey() {
	key := "test"

	suite.mockAdRedisRepo.On("GetAdsByKey", suite.ctx, key).Return([]*models.Advertisement{}, nil)

	ads, err := suite.s.GetAdsByKey(suite.ctx, key)

	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), ads)
	suite.mockAdRedisRepo.AssertExpectations(suite.T())
}

func (suite *AdvertisementServiceSuite) TestAdvertisementService_SetAdsByKey() {
	key := "test"
	ads := []*models.Advertisement{}
	expiration := time.Second

	suite.mockAdRedisRepo.On("SetAdsByKey", suite.ctx, key, ads, expiration).Return(nil)

	err := suite.s.SetAdsByKey(suite.ctx, key, ads, expiration)

	assert.NoError(suite.T(), err)
	suite.mockAdRedisRepo.AssertExpectations(suite.T())
}

func (suite *AdvertisementServiceSuite) TestAdvertisementService_DeleteAdsByPattern() {
	pattern := "test"

	suite.mockAdRedisRepo.On("DeleteAdsByPattern", suite.ctx, pattern).Return(nil)

	err := suite.s.DeleteAdsByPattern(suite.ctx, pattern)

	assert.NoError(suite.T(), err)
	suite.mockAdRedisRepo.AssertExpectations(suite.T())
}

func (suite *AdvertisementServiceSuite) TestAdvertisementService_IsAdExpired() {
	ads := []*models.Advertisement{
		{Title: "Test Ad 1", EndAt: time.Now().Add(time.Hour)},
		{Title: "Test Ad 2", EndAt: time.Now().Add(-time.Hour)},
	}

	isExpired := suite.s.IsAdExpired(ads)

	assert.True(suite.T(), isExpired)
}

func TestAdvertisementServiceSuite(t *testing.T) {
	suite.Run(t, new(AdvertisementServiceSuite))
}
