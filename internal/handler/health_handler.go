package handler

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type HealthHandler struct {
	db *gorm.DB
}

func NewHealthHandler(db *gorm.DB) *HealthHandler {
	return &HealthHandler{db: db}
}

func (h *HealthHandler) RegisterRoutes(r *gin.Engine) {
	// Endpoint kiểm tra tổng quan (Server + Database)
	r.GET("/health", h.CheckHealth)
	
	// Liveness Probe: Kiểm tra tiến trình server còn sống không
	r.GET("/health/live", h.CheckLiveness)
	
	// Readiness Probe: Kiểm tra ứng dụng và DB đã sẵn sàng phục vụ request chưa
	r.GET("/health/ready", h.CheckReadiness)
}

func (h *HealthHandler) CheckHealth(c *gin.Context) {
	dbStatus := "healthy"
	httpStatus := http.StatusOK

	// Ping kiểm tra kết nối thực tế tới PostgreSQL trong tối đa 2 giây
	ctx, cancel := context.WithTimeout(c.Request.Context(), 2*time.Second)
	defer cancel()

	sqlDB, err := h.db.DB()
	if err != nil || sqlDB.PingContext(ctx) != nil {
		dbStatus = "unhealthy"
		httpStatus = http.StatusServiceUnavailable
	}

	c.JSON(httpStatus, gin.H{
		"status":    ternary(httpStatus == http.StatusOK, "healthy", "unhealthy"),
		"timestamp": time.Now().UTC().Format(time.RFC3339),
		"services": gin.H{
			"database": dbStatus,
			"server":   "healthy",
		},
	})
}

func (h *HealthHandler) CheckLiveness(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status": "alive",
		"time":   time.Now().UTC().Format(time.RFC3339),
	})
}

func (h *HealthHandler) CheckReadiness(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 2*time.Second)
	defer cancel()

	sqlDB, err := h.db.DB()
	if err != nil || sqlDB.PingContext(ctx) != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"status":  "not_ready",
			"message": "Không thể kết nối đến cơ sở dữ liệu",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "ready",
		"message": "Hệ thống sẵn sàng phục vụ!",
	})
}

func ternary(cond bool, a, b string) string {
	if cond {
		return a
	}
	return b
}
