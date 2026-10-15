package handler

import (
	"errors"
	"net/http"

	"crud-api/internal/model"
	"crud-api/internal/service"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	authService service.AuthService
}

func NewAuthHandler(authService service.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

func (h *AuthHandler) RegisterRoutes(rg *gin.RouterGroup, authMiddleware gin.HandlerFunc) {
	auth := rg.Group("/auth")
	{
		auth.POST("/register", h.Register)
		auth.POST("/login", h.Login)
		
		// Tuyến đường yêu cầu đăng nhập
		auth.GET("/me", authMiddleware, h.GetMe)
	}
}

// Register godoc
// POST /api/v1/auth/register
func (h *AuthHandler) Register(c *gin.Context) {
	var req model.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		ResponseError(c, http.StatusBadRequest, "Dữ liệu đăng ký không hợp lệ", err)
		return
	}

	res, err := h.authService.Register(req)
	if err != nil {
		if errors.Is(err, service.ErrEmailExists) || errors.Is(err, service.ErrUsernameExists) {
			ResponseError(c, http.StatusConflict, err.Error(), nil)
			return
		}
		ResponseError(c, http.StatusInternalServerError, "Đăng ký thất bại", err)
		return
	}

	ResponseSuccess(c, http.StatusCreated, "Đăng ký tài khoản thành công", res)
}

// Login godoc
// POST /api/v1/auth/login
func (h *AuthHandler) Login(c *gin.Context) {
	var req model.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		ResponseError(c, http.StatusBadRequest, "Dữ liệu đăng nhập không hợp lệ", err)
		return
	}

	res, err := h.authService.Login(req)
	if err != nil {
		if errors.Is(err, service.ErrInvalidCreds) {
			ResponseError(c, http.StatusUnauthorized, "Email hoặc mật khẩu không chính xác", nil)
			return
		}
		ResponseError(c, http.StatusInternalServerError, "Đăng nhập thất bại", err)
		return
	}

	ResponseSuccess(c, http.StatusOK, "Đăng nhập thành công", res)
}

// GetMe godoc
// GET /api/v1/auth/me
func (h *AuthHandler) GetMe(c *gin.Context) {
	userIDVal, exists := c.Get("userID")
	if !exists {
		ResponseError(c, http.StatusUnauthorized, "Chưa xác thực người dùng", nil)
		return
	}

	userID := userIDVal.(uint)
	profile, err := h.authService.GetProfile(userID)
	if err != nil {
		ResponseError(c, http.StatusNotFound, "Không tìm thấy thông tin người dùng", err)
		return
	}

	ResponseSuccess(c, http.StatusOK, "Lấy thông tin người dùng thành công", profile)
}
