package repository_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"ad-service-api/internal/advertisement/repository"
	"ad-service-api/internal/models"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/integration/mtest"
)

func TestAdvertisementRepository_Create(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))

	mt.Run("Create", func(mt *mtest.T) {
		// Set up the mock responses
		mt.AddMockResponses(mtest.CreateSuccessResponse())
		mt.AddMockResponses(mtest.CreateCursorResponse(1, "foo.bar", mtest.FirstBatch, bson.D{{Key: "n", Value: int32(1)}}))

		repo := repository.NewAdvertisementRepository(mt.Coll)
		ctx := context.Background()
		ad := &models.Advertisement{
			Title: "Test Ad",
			// Populate the Advertisement struct
		}

		err := repo.Create(ctx, ad)
		assert.Nil(t, err)

		count, err := mt.Coll.CountDocuments(ctx, bson.M{})
		fmt.Println(count)
		assert.Nil(t, err)
		assert.Equal(t, int64(1), count, "expected document count to increase after insert")
	})
}

func TestAdvertisementRepository_CountActive(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))

	now := time.Now()

	mt.Run("CountActive", func(mt *mtest.T) {
		// Set up the mock responses
		mt.AddMockResponses(mtest.CreateCursorResponse(1, "foo.bar", mtest.FirstBatch, bson.D{{Key: "n", Value: int32(1)}}))

		repo := repository.NewAdvertisementRepository(mt.Coll)
		ctx := context.Background()

		// Insert sample advertisements
		mt.Coll.InsertOne(ctx, bson.M{"startAt": now.Add(-time.Hour), "endAt": now.Add(time.Hour)}) // Should be counted as active

		// Add another mock response for the CountDocuments operation
		mt.AddMockResponses(mtest.CreateCursorResponse(1, "foo.bar", mtest.FirstBatch, bson.D{{Key: "n", Value: int32(1)}}))

		count, err := repo.CountActive(ctx, now)
		fmt.Println(count)
		assert.Nil(t, err)
		assert.Equal(t, 1, count, "expected count of active advertisements to be correct")
	})
}

func TestAdvertisementRepository_Fetch(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))

	mt.Run("Fetch", func(mt *mtest.T) {
		repo := repository.NewAdvertisementRepository(mt.Coll)
		ctx := context.Background()

		// Insert sample advertisements
		mt.Coll.InsertOne(ctx, bson.M{"title": "Ad 1"})
		mt.Coll.InsertOne(ctx, bson.M{"title": "Ad 2"})

		// Set up the mock responses
		mt.AddMockResponses(mtest.CreateCursorResponse(1, "foo.bar", mtest.FirstBatch, bson.D{{Key: "_id", Value: primitive.NewObjectID()}, {Key: "title", Value: "Ad 1"}}))
		mt.AddMockResponses(mtest.CreateCursorResponse(1, "foo.bar", mtest.NextBatch, bson.D{{Key: "_id", Value: primitive.NewObjectID()}, {Key: "title", Value: "Ad 2"}}))
		mt.AddMockResponses(mtest.CreateCursorResponse(1, "foo.bar", mtest.NextBatch))

		ads, err := repo.Fetch(ctx, bson.M{}, 10, 0)
		fmt.Println(ads)
		assert.Nil(t, err)
		assert.Len(t, ads, 2, "expected number of advertisements to match")
	})
}
