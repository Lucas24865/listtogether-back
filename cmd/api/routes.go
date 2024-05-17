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
	authController, userController, groupController, notifController, listController, err := newControllers()
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

	//Groups
	groups := base.Group("/groups")
	groups.Use(middleware.JwtAuthMiddleware())
	groups.GET("", groupController.GetAll)
	groups.GET(":id", groupController.Get)
	groups.POST("", groupController.Create)
	groups.PUT("", groupController.Edit)

	//Notifications
	notifications := base.Group("/notifications")
	notifications.Use(middleware.JwtAuthMiddleware())
	notifications.GET("", notifController.GetAll)
	notifications.PUT("/accept/:id", notifController.AcceptInvite)
	notifications.PUT("/decline/:id", notifController.DeclineInvite)
	notifications.DELETE("/:id", notifController.Remove)
	notifications.DELETE("", notifController.RemoveAllRead)

	//Lists
	lists := base.Group("/lists")
	lists.GET("", listController.GetAll)
	lists.GET("/:id", listController.Get)
	lists.POST("", listController.Create)
	lists.PUT("", listController.Update)
	lists.DELETE("/:id", listController.Delete)

}

func newControllers() (authController controller.AuthController, userController controller.UserController, groupController controller.GroupController,
	notifController controller.NotificationController, listController controller.ListController, err error) {
	err = repository.FirebaseDB().Connect()
	if err != nil {
		log.Println(err)
		return nil, nil, nil, nil, nil, err
	}

	repo := repository.FirebaseDB()
	userRepo := repository.NewUserRepository(repo)
	notifRepo := repository.NewNotificationRepository(repo)
	groupRepo := repository.NewGroupRepository(repo)
	listRepo := repository.NewListRepository(repo, groupRepo)

	notifService := service.NewNotificationService(notifRepo)
	groupService := service.NewGroupService(groupRepo, notifService)
	userService := service.NewUserService(userRepo, notifService, groupService)
	listService := service.NewListService(listRepo)
	authService := service.NewAuthService(userService)

	authController = controller.NewAuthController(authService)
	userController = controller.NewUserController(userService)
	notifController = controller.NewNotificationController(notifService, userService)
	groupController = controller.NewGroupController(groupService, userService, notifService)
	listController = controller.NewListController(listService)

	return
}
