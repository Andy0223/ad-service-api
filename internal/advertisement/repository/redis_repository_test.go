package repository_test

import (
	"context"
	"testing"

	"github.com/go-redis/redismock/v8"
	"github.com/stretchr/testify/assert"

	"ad-service-api/internal/advertisement/repository"
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

	mock.ExpectGet("ads:testKey").SetVal("1")

	count, err := repo.GetByDate(context.Background(), "testKey")
	assert.NoError(t, err)
	assert.Equal(t, 1, count)
	assert.NoError(t, mock.ExpectationsWereMet())
}
