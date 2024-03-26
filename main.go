package main

import (
	_ "ad-service-api/docs"
	"ad-service-api/internal/router"
)

func main() {
	router := router.NewRouter()
	router.Run(":8080")
}
