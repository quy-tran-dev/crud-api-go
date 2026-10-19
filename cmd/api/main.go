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
	worldHandler "crud-api/internal/handler/world"
	"crud-api/internal/middleware"
	"crud-api/internal/model"
	"crud-api/internal/repository"
	"crud-api/internal/service"
	worldService "crud-api/internal/service/world"

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
	// Quản lý Sản phẩm
	productRepo := repository.NewProductRepository(db)
	productService := service.NewProductService(productRepo)
	productHandler := handler.NewProductHandler(productService)

	// Quản lý Xác thực & Người dùng
	userRepo := repository.NewUserRepository(db)
	authService := service.NewAuthService(userRepo, cfg)
	authHandler := handler.NewAuthHandler(authService)

	// Quản lý Kịch bản & World-Building (Script Management)
	storyRepo := repository.NewStoryRepository(db)
	categoryRepo := repository.NewCategoryRepository(db)
	locationRepo := repository.NewLocationRepository(db)
	factionRepo := repository.NewFactionRepository(db)
	characterRepo := repository.NewCharacterRepository(db)
	plotRepo := repository.NewPlotRepository(db)
	itemRepo := repository.NewItemRepository(db)
	frameworkRepo := repository.NewFrameworkRepository(db)

	storyService := worldService.NewStoryService(storyRepo, categoryRepo, locationRepo, factionRepo, characterRepo, plotRepo, itemRepo)
	categoryService := worldService.NewCategoryService(categoryRepo, storyRepo)
	locationService := worldService.NewLocationService(locationRepo, storyRepo)
	factionService := worldService.NewFactionService(factionRepo, storyRepo)
	characterService := worldService.NewCharacterService(characterRepo, categoryRepo, storyRepo)
	plotService := worldService.NewPlotService(plotRepo, storyRepo)
	itemService := worldService.NewItemService(itemRepo, storyRepo)
	frameworkService := worldService.NewFrameworkService(frameworkRepo, storyRepo, categoryRepo)

	worldHandlers := &worldHandler.WorldModuleHandlers{
		StoryHandler:     worldHandler.NewStoryHandler(storyService),
		CategoryHandler:  worldHandler.NewCategoryHandler(categoryService),
		LocationHandler:  worldHandler.NewLocationHandler(locationService),
		FactionHandler:   worldHandler.NewFactionHandler(factionService),
		CharacterHandler: worldHandler.NewCharacterHandler(characterService),
		PlotHandler:      worldHandler.NewPlotHandler(plotService),
		ItemHandler:      worldHandler.NewItemHandler(itemService),
		FrameworkHandler: worldHandler.NewFrameworkHandler(frameworkService),
	}

	// Middlewares
	authMiddleware := middleware.AuthMiddleware(authService)
	adminMiddleware := middleware.RequireRole(model.RoleAdmin)

	// 4. Khởi tạo Gin Router sạch (loại bỏ cảnh báo GIN-debug WARNING)
	r := gin.New()
	r.Use(gin.Logger(), gin.Recovery())

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
		authHandler.RegisterRoutes(apiV1, authMiddleware)
		productHandler.RegisterRoutes(apiV1, authMiddleware, adminMiddleware)
		worldHandler.RegisterWorldRoutes(apiV1, authMiddleware, worldHandlers)
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
