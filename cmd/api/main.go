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

	"crud-api/internal/config"
	"crud-api/internal/database"
	"crud-api/internal/handler"
	"crud-api/internal/repository"
	"crud-api/internal/service"

	"github.com/gin-gonic/gin"
)

func main() {
	// 1. Tải cấu hình ứng dụng
	cfg := config.LoadConfig()

	// 2. Kết nối cơ sở dữ liệu PostgreSQL & AutoMigrate
	db, err := database.ConnectDB(cfg)
	if err != nil {
		log.Fatalf("Không thể khởi tạo cơ sở dữ liệu: %v", err)
	}

	// 3. Khởi tạo tầng Repository, Service, Handler (Dependency Injection)
	productRepo := repository.NewProductRepository(db)
	productService := service.NewProductService(productRepo)
	productHandler := handler.NewProductHandler(productService)

	// 4. Khởi tạo Gin Router
	r := gin.Default()

	// CORS Middleware đơn giản
	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	})

	healthHandler := handler.NewHealthHandler(db)
	healthHandler.RegisterRoutes(r)

	// API V1 routes
	apiV1 := r.Group("/api/v1")
	{
		productHandler.RegisterRoutes(apiV1)
	}

	// 5. Khởi động HTTP Server với cơ chế Graceful Shutdown
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%s", cfg.AppPort),
		Handler:      r,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		log.Printf("Server đang chạy tại cổng http://localhost:%s", cfg.AppPort)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Lỗi khởi động server: %v\n", err)
		}
	}()

	// Chờ tín hiệu ngắt (SIGINT, SIGTERM) để tắt server an toàn
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Đang nhận tín hiệu tắt máy...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server bị buộc tắt: %v", err)
	}

	log.Println("Server đã tắt an toàn!")
}
