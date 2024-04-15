package router

import (
	"os"

	"ad-service-api/database"
	"ad-service-api/internal/advertisement/handler"
	"ad-service-api/internal/advertisement/repository"
	"ad-service-api/internal/advertisement/service"
	"ad-service-api/internal/middleware"
	"ad-service-api/redis"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// NewRouter creates a new router
func NewRouter() *gin.Engine {
	mongoUsername := os.Getenv("MONGO_USERNAME")
	mongoPassword := os.Getenv("MONGO_PASSWORD")
	mongoHost := os.Getenv("MONGO_HOST")
	mongoDb := os.Getenv("MONGO_DB")
	mongoCollection := os.Getenv("MONGO_COLLECTION")
	redisHost := os.Getenv("REDIS_HOST")
	redisPassword := os.Getenv("REDIS_PASSWORD")
	redisDb := 0
	col, _ := database.ConnectMongoDB(mongoUsername, mongoPassword, mongoHost, mongoDb, mongoCollection)
	rdb, _ := redis.ConnectRedis(redisHost, redisPassword, redisDb)

	adRepo := repository.NewAdvertisementRepository(col)
	adRedisRepo := repository.NewAdRedisRepository(rdb)
	adService := service.NewAdvertisementService(adRepo, adRedisRepo)
	adHandler := handler.NewAdvertisementHandler(adService)

	r := gin.Default()
	r.Use(middleware.Logger())

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	adRoutes := r.Group("/api/v1")
	{
		adRoutes.POST("/ad", adHandler.CreateAdHandler)
		adRoutes.GET("/ad", adHandler.ListAdHandler)
		adRoutes.DELETE("/ad/:id", adHandler.DeleteAdHandler)
		adRoutes.PUT("/ad/:id", adHandler.UpdateAdHandler)
	}

	return r
}
