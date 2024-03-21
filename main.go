package main

import (
	"ad-service-api/internal/config"
	"ad-service-api/internal/router"
)

func main() {
	config.ConnectMongoDB("mongodb://localhost:27017", "2024DcardBackend")
	config.ConnectRedis()
	router := router.RegisterRoutes() // 创建路由

	router.Run(":8080") // 在 0.0.0.0:8080 上启动服务器
}
