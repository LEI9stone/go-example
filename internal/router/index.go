package router

import (
	"example.com/acg-go-demo/internal/response"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func RegisterRoutes(engine *gin.Engine, db *gorm.DB) {
	api := engine.Group("/api")

	api.GET("/health", func(c *gin.Context) {
		response.Success(c, gin.H{"status": "ok"})
	})

	api.GET("/hello", func(c *gin.Context) {
		response.Success(c, gin.H{
			"message": "hello from acggoods practice",
		})
	})

	UsersRoutes(api, db)
}
