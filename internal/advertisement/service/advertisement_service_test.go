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
	mockAdCountRepo *mocks.MockAdCountRepository
	s               service.IAdvertisementService
	ctx             context.Context
}

func (suite *AdvertisementServiceSuite) SetupTest() {
	suite.mockAdRepo = new(mocks.MockAdvertisementRepository)
	suite.mockAdCountRepo = new(mocks.MockAdCountRepository)
	suite.s = service.NewAdvertisementService(suite.mockAdRepo, suite.mockAdCountRepo)
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

	suite.mockAdCountRepo.On("GetByDate", suite.ctx, today).Return(1, nil)

	count, err := suite.s.GetByDate(suite.ctx, today)

	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), 1, count)
	suite.mockAdCountRepo.AssertExpectations(suite.T())
}

func (suite *AdvertisementServiceSuite) TestAdvertisementService_IncrByDate() {
	key := "2022-01-01"

	suite.mockAdCountRepo.On("IncrByDate", suite.ctx, key).Return(nil)

	err := suite.s.IncrByDate(suite.ctx, key)

	assert.NoError(suite.T(), err)
	suite.mockAdCountRepo.AssertExpectations(suite.T())
}

func TestAdvertisementServiceSuite(t *testing.T) {
	suite.Run(t, new(AdvertisementServiceSuite))
}
