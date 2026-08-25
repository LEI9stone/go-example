package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

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

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			panic(err)
		}
	}()

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)

	<-signals

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Printf("server shutdonw failed: %v", err)
	}
}
