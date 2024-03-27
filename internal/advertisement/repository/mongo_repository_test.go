package repository_test

// import (
// 	"ad-service-api/internal/advertisement/repository"
// 	"ad-service-api/internal/models"
// 	"context"
// 	"testing"

// 	"github.com/stretchr/testify/mock"
// 	"go.mongodb.org/mongo-driver/mongo"
// 	"go.mongodb.org/mongo-driver/mongo/options"
// )

// // CollectionAPI is an interface that includes the methods you need to mock.
// type CollectionAPI interface {
// 	InsertOne(ctx context.Context, document interface{}, opts ...*options.InsertOneOptions) (*mongo.InsertOneResult, error)
// 	Find(ctx context.Context, filter interface{}, opts ...*options.FindOptions) (*mongo.Cursor, error)
// 	CountDocuments(ctx context.Context, filter interface{}, opts ...*options.CountOptions) (int64, error)
// }

// // Collection is a wrapper around *mongo.Collection that implements CollectionAPI.
// type Collection struct {
// 	collection *mongo.Collection
// }

// func (c *Collection) InsertOne(ctx context.Context, document interface{}, opts ...*options.InsertOneOptions) (*mongo.InsertOneResult, error) {
// 	return c.collection.InsertOne(ctx, document, opts...)
// }

// func (c *Collection) Find(ctx context.Context, filter interface{}, opts ...*options.FindOptions) (*mongo.Cursor, error) {
// 	return c.collection.Find(ctx, filter, opts...)
// }

// func (c *Collection) CountDocuments(ctx context.Context, filter interface{}, opts ...*options.CountOptions) (int64, error) {
// 	return c.collection.CountDocuments(ctx, filter, opts...)
// }

// // MockCollection is a mock implementation of CollectionAPI.
// type MockCollection struct {
// 	mock.Mock
// }

// func (m *MockCollection) InsertOne(ctx context.Context, document interface{}, opts ...*options.InsertOneOptions) (*mongo.InsertOneResult, error) {
// 	args := m.Called(ctx, document)
// 	return args.Get(0).(*mongo.InsertOneResult), args.Error(1)
// }

// func (m *MockCollection) Find(ctx context.Context, filter interface{}, opts ...*options.FindOptions) (*mongo.Cursor, error) {
// 	args := m.Called(ctx, filter)
// 	return args.Get(0).(*mongo.Cursor), args.Error(1)
// }

// // Now you can use MockCollection in your tests.
// func TestAdvertisementRepository(t *testing.T) {
// 	mockCollection := new(MockCollection)
// 	adRepo := repository.NewAdvertisementRepository(mockCollection)

// 	ctx := context.Background()
// 	ad := &models.Advertisement{Title: "Test Ad"}

// 	// Test Create
// 	mockCollection.On("InsertOne", ctx, ad).Return(&mongo.InsertOneResult{}, nil)
// 	err := adRepo.Create(ctx, ad)
// 	if err != nil {
// 		t.Errorf("Error should be nil: %v", err)
// 	}

// 	// Test CountActive
// 	mockCollection.On("CountDocuments", ctx, mock.Anything).Return(int64(1), nil)
// 	count, err := adRepo.CountActive(ctx, mock.Anything)
// 	if err != nil {
// 		t.Errorf("Error should be nil: %v", err)
// 	}

// 	// Test Fetch
// 	mockCollection.On("Find", ctx, mock.Anything).Return(&mongo.Cursor{}, nil)
// 	_, err = adRepo.Fetch(ctx, mock.Anything, 10, 0)
// 	if err != nil {
// 		t.Errorf("Error should be nil: %v", err)
// 	}

// 	mockCollection.AssertExpectations(t)
// }
