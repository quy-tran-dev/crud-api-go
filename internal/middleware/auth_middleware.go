package middleware

import (
	"net/http"
	"strings"

	"crud-api/internal/handler"
	"crud-api/internal/service"

	"github.com/gin-gonic/gin"
)

// AuthMiddleware kiểm tra tính hợp lệ của JWT Bearer token
func AuthMiddleware(authService service.AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			handler.ResponseError(c, http.StatusUnauthorized, "Yêu cầu đăng nhập (Thiếu Authorization Header)", nil)
			c.Abort()
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			handler.ResponseError(c, http.StatusUnauthorized, "Định dạng Token phải là 'Bearer <token>'", nil)
			c.Abort()
			return
		}

		tokenString := parts[1]
		claims, err := authService.ValidateToken(tokenString)
		if err != nil {
			handler.ResponseError(c, http.StatusUnauthorized, "Token không hợp lệ hoặc đã hết hạn", err)
			c.Abort()
			return
		}

		// Lưu thông tin người dùng vào Gin Context để các handler tiếp theo sử dụng
		c.Set("userID", claims.UserID)
		c.Set("username", claims.Username)
		c.Set("userRole", claims.Role)

		c.Next()
	}
}

// RequireRole kiểm tra quyền của người dùng (Role-Based Authorization)
func RequireRole(allowedRoles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		roleVal, exists := c.Get("userRole")
		if !exists {
			handler.ResponseError(c, http.StatusForbidden, "Không xác định được quyền hạn người dùng", nil)
			c.Abort()
			return
		}

		userRole, ok := roleVal.(string)
		if !ok {
			handler.ResponseError(c, http.StatusForbidden, "Quyền hạn không hợp lệ", nil)
			c.Abort()
			return
		}

		for _, r := range allowedRoles {
			if r == userRole {
				c.Next()
				return
			}
		}

		handler.ResponseError(c, http.StatusForbidden, "Bạn không có quyền thực hiện thao tác này (Yêu cầu quyền: "+strings.Join(allowedRoles, ", ")+")", nil)
		c.Abort()
	}
}
