package main

import (
	_ "ad-service-api/docs"
	_adHandler "ad-service-api/internal/advertisement/handler"
	_adRepo "ad-service-api/internal/advertisement/repository"
	_adService "ad-service-api/internal/advertisement/service"
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
	collection := db.Collection("ads")

	adRepo := _adRepo.NewAdvertisementRepository(collection, rdb)
	adService := _adService.NewAdvertisementService(adRepo)
	_adHandler.NewAdvertisementHandler(r, adService)

	r.Run(":8080") // 在 0.0.0.0:8080 上启动服务器
}
