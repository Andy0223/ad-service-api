package service_test

import (
	"ad-service-api/internal/advertisement/service"
	"ad-service-api/internal/models"
	"ad-service-api/mocks"
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestAdvertisementService_Create(t *testing.T) {
	mockAdRepo := new(mocks.MockAdvertisementRepository)
	mockAdCountRepo := new(mocks.MockAdCountRepository)

	ad := &models.Advertisement{}
	ctx := context.TODO()

	mockAdRepo.On("Create", ctx, ad).Return(nil)

	s := service.NewAdvertisementService(mockAdRepo, mockAdCountRepo)
	err := s.Create(ctx, ad)

	assert.NoError(t, err)
	mockAdRepo.AssertExpectations(t)
}

func TestAdvertisementService_CountActive(t *testing.T) {
	mockAdRepo := new(mocks.MockAdvertisementRepository)
	mockAdCountRepo := new(mocks.MockAdCountRepository)

	now := time.Now()
	ctx := context.TODO()

	mockAdRepo.On("CountActive", ctx, now).Return(1, nil)

	s := service.NewAdvertisementService(mockAdRepo, mockAdCountRepo)
	count, err := s.CountActive(ctx, now)

	assert.NoError(t, err)
	assert.Equal(t, 1, count)
	mockAdRepo.AssertExpectations(t)
}

func TestAdvertisementService_Fetch(t *testing.T) {
	mockAdRepo := new(mocks.MockAdvertisementRepository)
	mockAdCountRepo := new(mocks.MockAdCountRepository)

	filter := primitive.M{}
	limit := 10
	offset := 0
	ctx := context.TODO()

	mockAdRepo.On("Fetch", ctx, filter, limit, offset).Return([]*models.Advertisement{}, nil)

	s := service.NewAdvertisementService(mockAdRepo, mockAdCountRepo)
	ads, err := s.Fetch(ctx, filter, limit, offset)

	assert.NoError(t, err)
	assert.NotNil(t, ads)
	mockAdRepo.AssertExpectations(t)
}

func TestAdvertisementService_GetByDate(t *testing.T) {
	mockAdRepo := new(mocks.MockAdvertisementRepository)
	mockAdCountRepo := new(mocks.MockAdCountRepository)

	today := "2022-01-01"
	ctx := context.TODO()

	mockAdCountRepo.On("GetByDate", ctx, today).Return(1, nil)

	s := service.NewAdvertisementService(mockAdRepo, mockAdCountRepo)
	count, err := s.GetByDate(ctx, today)

	assert.NoError(t, err)
	assert.Equal(t, 1, count)
	mockAdCountRepo.AssertExpectations(t)
}

func TestAdvertisementService_IncrByDate(t *testing.T) {
	mockAdRepo := new(mocks.MockAdvertisementRepository)
	mockAdCountRepo := new(mocks.MockAdCountRepository)

	key := "2022-01-01"
	ctx := context.TODO()

	mockAdCountRepo.On("IncrByDate", ctx, key).Return(nil)

	s := service.NewAdvertisementService(mockAdRepo, mockAdCountRepo)
	err := s.IncrByDate(ctx, key)

	assert.NoError(t, err)
	mockAdCountRepo.AssertExpectations(t)
}
