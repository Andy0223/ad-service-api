package repository_test

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/go-redis/redismock/v8"
	"github.com/stretchr/testify/assert"

	"ad-service-api/internal/advertisement/repository"
	"ad-service-api/internal/models"
)

func TestAdCountRepository_IncrByDate(t *testing.T) {
	db, mock := redismock.NewClientMock()
	repo := repository.NewAdCountRepository(db)

	mock.ExpectIncr("testKey").SetVal(1)

	err := repo.IncrByDate(context.Background(), "testKey")
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestAdCountRepository_GetByDate(t *testing.T) {
	db, mock := redismock.NewClientMock()
	repo := repository.NewAdCountRepository(db)

	mock.ExpectGet("testKey").SetVal("1")

	count, err := repo.GetByDate(context.Background(), "testKey")
	assert.NoError(t, err)
	assert.Equal(t, 1, count)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestAdCountRepository_GetAdByKey(t *testing.T) {
	db, mock := redismock.NewClientMock()
	repo := repository.NewAdCountRepository(db)

	ads := []*models.Advertisement{
		{Title: "test1"},
	}
	adsJson, _ := json.Marshal(ads)

	mock.ExpectGet("ads:testKey").SetVal(string(adsJson))

	returnedAds, err := repo.GetAdsByKey(context.Background(), "ads:testKey")
	assert.NoError(t, err)
	assert.Equal(t, ads, returnedAds)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestAdCountRepository_SetAdsByKey(t *testing.T) {
	db, mock := redismock.NewClientMock()
	repo := repository.NewAdCountRepository(db)

	ads := []*models.Advertisement{
		{Title: "test1"},
	}
	adsJson, _ := json.Marshal(ads)

	mock.ExpectSet("testKey", adsJson, 0).SetVal("OK")

	err := repo.SetAdsByKey(context.Background(), "testKey", ads, 0)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestAdCountRepository_DeleteAdsByPattern(t *testing.T) {
	db, mock := redismock.NewClientMock()
	repo := repository.NewAdCountRepository(db)

	mock.ExpectKeys("testKey").SetVal([]string{"ads:testKey"})
	mock.ExpectDel("ads:testKey").SetVal(1)

	err := repo.DeleteAdsByPattern(context.Background(), "testKey")
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}
