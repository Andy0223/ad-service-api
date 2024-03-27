package router

import (
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
	col, _ := database.ConnectMongoDB("mongodb://localhost:27017", "2024DcardBackend", "ads")
	rdb, _ := redis.ConnectRedis("localhost:6379", "", 0)

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
