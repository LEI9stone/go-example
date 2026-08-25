package main

import (
	"net/http"

	"example.com/acg-go-demo/internal/middleware"
	"example.com/acg-go-demo/internal/response"
	"github.com/gin-gonic/gin"
)

func main() {
	engine := gin.New();
	engine.Use(
		middleware.RequestID(),
		middleware.TraceID(),
		gin.Logger(),
		gin.Recovery(),
	)

	engine.GET("/health", func(c *gin.Context) {
		response.Success(c, gin.H{"status": "ok"})
	})

	engine.GET("/api/hello", func(c *gin.Context) {
		response.Success(c, gin.H{
			"message": "hello from acggoods practice",
		})
	})

	server := &http.Server{
		Addr: ":8080",
		Handler: engine,
	}

	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		panic(err)
	}
}