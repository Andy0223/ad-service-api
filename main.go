package main

import (
	_ "ad-service-api/docs"
	"ad-service-api/internal/advertisement/handler"
	mongodbRepo "ad-service-api/internal/advertisement/repository/mongodb"
	redisRepo "ad-service-api/internal/advertisement/repository/redis"
	adService "ad-service-api/internal/advertisement/service"
	"ad-service-api/internal/config"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func main() {
	db, _ := config.ConnectMongoDB("mongodb://localhost:27017", "2024DcardBackend")
	rdb, _ := config.ConnectRedis("localhost:6379", "", 0)

	r := gin.Default()
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	adRepo := mongodbRepo.NewAdvertisementRepository(db)
	adRedisRepo := redisRepo.NewRedisRepository(rdb)

	adService := adService.NewAdvertisementService(adRepo, adRedisRepo)

	handler.NewAdvertisementHandler(r, adService)

	r.Run(":8080")
}
