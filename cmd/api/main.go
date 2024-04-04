package main

import (
	"github.com/gin-gonic/gin"
)

func main() {
	if err := SetupRouter().Run(":8080"); err != nil {
		panic(err.Error())
	}
}

func SetupRouter() *gin.Engine {
	router := gin.Default()

	mapRoutes(router)
	return router
}
