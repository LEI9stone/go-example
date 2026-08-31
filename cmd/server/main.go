package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"example.com/acg-go-demo/internal/app"
	"example.com/acg-go-demo/internal/middleware"
	"example.com/acg-go-demo/internal/router"
	"github.com/gin-gonic/gin"

	appconfig "example.com/acg-go-demo/internal/config"
)

func main() {

	engine := gin.New()
	engine.Use(
		middleware.RequestID(),
		middleware.TraceID(),
		gin.Logger(),
		gin.Recovery(),
	)

	cfg, err := appconfig.Load()

	if err != nil {
		log.Fatalf("load config failed: %v", err)
	}

	db, err := app.OpenDatabase(cfg)

	if err != nil {
		log.Fatalf("open database failed: %v", err)
	}

	router.RegisterRoutes(engine, db, cfg)

	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", cfg.App.Port),
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
