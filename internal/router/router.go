package router

import (
	"ad-service-api/internal/handlers"
	"ad-service-api/internal/middleware"

	"github.com/gin-gonic/gin"
)

// RegisterRoutes registers advertisement-related routes.
func RegisterRoutes() *gin.Engine {
	router := gin.Default()

	// 注册中间件
	router.Use(middleware.Logger())

	group := router.Group("/api/v1")
	{
		group.POST("/ad", handlers.CreateAdvertisement)
		group.GET("/ads", handlers.ListAdvertisements)
	}

	return router
}
