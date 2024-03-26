package service_test

import (
	"ad-service-api/internal/advertisement/service"
	"ad-service-api/internal/models"
	"context"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"ad-service-api/internal/mock"
)

func TestAdvertisementService(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAdRepo := mock.NewMockAdvertisementRepository(ctrl)
	mockAdCountRepo := mock.NewMockAdCountRepository(ctrl)
	adService := service.NewAdvertisementService(mockAdRepo, mockAdCountRepo)
	ctx := context.Background()

	// Test Create method
	ad := &models.Advertisement{}
	mockAdRepo.EXPECT().Create(ctx, ad).Return(nil)
	err := adService.Create(ctx, ad)
	assert.NoError(t, err)

	// Test CountActive method
	now := time.Now()
	mockAdRepo.EXPECT().CountActive(ctx, now).Return(1, nil)
	count, err := adService.CountActive(ctx, now)
	assert.NoError(t, err)
	assert.Equal(t, 1, count)

	// Test Fetch method
	filter := primitive.M{}
	mockAdRepo.EXPECT().Fetch(ctx, filter, 10, 0).Return([]*models.Advertisement{}, nil)
	ads, err := adService.Fetch(ctx, filter, 10, 0)
	assert.NoError(t, err)
	assert.Equal(t, []*models.Advertisement{}, ads)

	// Test IncrByDate and GetByDate methods
	key := "2022-01-01"
	mockAdCountRepo.EXPECT().IncrByDate(ctx, key).Return(nil)
	err = adService.IncrByDate(ctx, key)
	assert.NoError(t, err)

	mockAdCountRepo.EXPECT().GetByDate(ctx, key).Return(1, nil)
	count, err = adService.GetByDate(ctx, key)
	assert.NoError(t, err)
	assert.Equal(t, 1, count)
}
