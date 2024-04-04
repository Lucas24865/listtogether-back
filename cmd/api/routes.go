package main

import (
	"ListTogetherAPI/internal/controller"
	"ListTogetherAPI/internal/middleware"
	"ListTogetherAPI/internal/repository"
	"ListTogetherAPI/internal/service"
	"github.com/gin-gonic/gin"
	"log"
)

func mapRoutes(router *gin.Engine) {
	authController, userController, err := newControllers()
	if err != nil {
		log.Panicf(err.Error())
	}

	base := router.Group("/api")

	//auth
	authBase := base.Group("/auth")
	authBase.POST("/register", authController.Register)
	authBase.POST("/login", authController.Login)

	protected := base.Group("/user")
	protected.Use(middleware.JwtAuthMiddleware())
	protected.GET("", userController.Get)

}

func newControllers() (controller.AuthController, controller.UserController, error) {
	err := repository.FirebaseDB().Connect()
	if err != nil {
		log.Println(err)
		return nil, nil, err
	}

	repo := repository.FirebaseDB()
	userRepo := repository.NewUserRepository(repo)

	authService := service.NewAuthService(userRepo)
	authController := controller.NewAuthController(authService)

	userService := service.NewUserService(userRepo)
	userController := controller.NewUserController(userService)

	return authController, userController, nil
}
