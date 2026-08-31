package router

import (
	"example.com/acg-go-demo/internal/handler"
	"example.com/acg-go-demo/internal/repository"
	"example.com/acg-go-demo/internal/service"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func UsersRoutes(api *gin.RouterGroup, db *gorm.DB) {

	userRepo := repository.NewUserRepository(db)
	userService := service.NewUserService(userRepo)
	userHandler := handler.NewUserHandler(userService)

	api.POST("/users/register", userHandler.Register)
	api.GET("/users/:id", userHandler.Get)
	api.PATCH("/users/:id/nickname", userHandler.UpdateNickname)
	api.DELETE("/users/:id", userHandler.Delete)
}
