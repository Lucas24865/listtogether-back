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
	authController, userController, groupController, notifController, err := newControllers()
	if err != nil {
		log.Panicf(err.Error())
	}

	base := router.Group("/api")

	//auth
	authBase := base.Group("/auth")
	authBase.POST("/register", authController.Register)
	authBase.POST("/login", authController.Login)

	//User
	user := base.Group("/user")
	user.Use(middleware.JwtAuthMiddleware())
	user.GET("", userController.Get)
	user.POST("", userController.Get)

	//Groups
	groups := base.Group("/groups")
	groups.Use(middleware.JwtAuthMiddleware())
	groups.GET("", groupController.GetAll)
	groups.POST("", groupController.Create)

	//Notifications
	notifications := base.Group("/notifications")
	notifications.Use(middleware.JwtAuthMiddleware())
	notifications.GET("", notifController.GetAll)
	notifications.POST("/accept/:id", notifController.Accept)
	notifications.GET("/decline/:id", notifController.Decline)
	notifications.DELETE("/:id", notifController.Remove)

}

func newControllers() (controller.AuthController, controller.UserController, controller.GroupController, controller.NotificationController, error) {
	err := repository.FirebaseDB().Connect()
	if err != nil {
		log.Println(err)
		return nil, nil, nil, nil, err
	}

	repo := repository.FirebaseDB()
	userRepo := repository.NewUserRepository(repo)
	notifRepo := repository.NewNotificationRepository(repo)
	groupRepo := repository.NewGroupRepository(repo)

	authService := service.NewAuthService(userRepo)
	notifService := service.NewNotificationService(notifRepo)
	userService := service.NewUserService(userRepo, notifService)
	groupService := service.NewGroupService(groupRepo)

	authController := controller.NewAuthController(authService)
	userController := controller.NewUserController(userService)
	notifController := controller.NewNotificationController(notifService)
	groupController := controller.NewGroupController(groupService, userService, notifService)

	return authController, userController, groupController, notifController, nil
}
