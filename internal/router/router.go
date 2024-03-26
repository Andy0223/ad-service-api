package router

import (
	"ad-service-api/internal/advertisement/handler"
	"ad-service-api/internal/advertisement/repository/mongodb"
	"ad-service-api/internal/advertisement/repository/redis"
	"ad-service-api/internal/advertisement/service"
	"ad-service-api/internal/config"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// NewRouter creates a new router
func NewRouter() *gin.Engine {
	db, _ := config.ConnectMongoDB("mongodb://localhost:27017", "2024DcardBackend")
	collection := db.Collection("ads")
	rdb, _ := config.ConnectRedis("localhost:6379", "", 0)

	adRepo := mongodb.NewAdvertisementRepository(collection)
	adCountRepo := redis.NewAdCountRepository(rdb)
	adService := service.NewAdvertisementService(adRepo, adCountRepo)
	adHandler := handler.NewAdvertisementHandler(adService)

	r := gin.Default()
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	r.POST("/api/v1/ad", adHandler.CreateAdHandler)
	r.GET("/api/v1/ads", adHandler.ListAdHandler)

	return r
}
