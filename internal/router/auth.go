package router

import (
	"example.com/acg-go-demo/internal/config"
	"example.com/acg-go-demo/internal/handler"
	"example.com/acg-go-demo/internal/repository"
	"example.com/acg-go-demo/internal/service"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func AuthRoutes(api *gin.RouterGroup, db *gorm.DB, cfg *config.Config) {
	userRepo := repository.NewUserRepository(db)
	authService := service.NewAuthService(userRepo)

	cookieName := cfg.Auth.CookieName
	if cookieName == "" {
		cookieName = "acggoods_practice_auth"
	}

	authHandler := handler.NewAuthHandler(
		authService,
		cookieName,
		cfg.App.Env == "production",
	)

	api.POST("/auth/login", authHandler.Login)
	api.POST("/auth/logout", authHandler.Logout)
}
