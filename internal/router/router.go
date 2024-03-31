package router

import (
	"os"

	"ad-service-api/database"
	"ad-service-api/internal/advertisement/handler"
	"ad-service-api/internal/advertisement/repository"
	"ad-service-api/internal/advertisement/service"
	"ad-service-api/redis"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// NewRouter creates a new router
func NewRouter() *gin.Engine {
	mongoUri := os.Getenv("MONGO_URI")
	mongoDb := os.Getenv("MONGO_DB")
	mongoCollection := os.Getenv("MONGO_COLLECTION")
	redisHost := os.Getenv("REDIS_HOST")
	redisPassword := os.Getenv("REDIS_PASSWORD")
	redisDb := 0
	col, _ := database.ConnectMongoDB(mongoUri, mongoDb, mongoCollection)
	rdb, _ := redis.ConnectRedis(redisHost, redisPassword, redisDb)

	adRepo := repository.NewAdvertisementRepository(col)
	adCountRepo := repository.NewAdCountRepository(rdb)
	adService := service.NewAdvertisementService(adRepo, adCountRepo)
	adHandler := handler.NewAdvertisementHandler(adService)

	r := gin.Default()
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	r.POST("/api/v1/ad", adHandler.CreateAdHandler)
	r.GET("/api/v1/ads", adHandler.ListAdHandler)

	return r
}