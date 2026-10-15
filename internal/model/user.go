package model

import "time"

const (
	RoleAdmin = "admin"
	RoleUser  = "user"
)

// User đại diện cho bảng users trong PostgreSQL
type User struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Username  string    `gorm:"size:100;uniqueIndex;not null" json:"username"`
	Email     string    `gorm:"size:255;uniqueIndex;not null" json:"email"`
	Password  string    `gorm:"size:255;not null" json:"-"` // Dấu "-" đảm bảo không bao giờ serialize password ra JSON
	Role      string    `gorm:"size:20;not null;default:'user'" json:"role"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// RegisterRequest nhận dữ liệu khi đăng ký
type RegisterRequest struct {
	Username string `json:"username" binding:"required,min=3,max=50"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
	Role     string `json:"role" binding:"omitempty,oneof=admin user"`
}

// LoginRequest nhận dữ liệu khi đăng nhập
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// UserResponse là DTO trả về thông tin user an toàn (không có password)
type UserResponse struct {
	ID        uint      `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	Role      string    `json:"role"`
	CreatedAt time.Time `json:"created_at"`
}

// AuthResponse trả về thông tin sau khi đăng nhập / đăng ký thành công
type AuthResponse struct {
	User  UserResponse `json:"user"`
	Token string       `json:"token"`
}

// ToResponse chuyển đổi struct User sang UserResponse
func (u *User) ToResponse() UserResponse {
	return UserResponse{
		ID:        u.ID,
		Username:  u.Username,
		Email:     u.Email,
		Role:      u.Role,
		CreatedAt: u.CreatedAt,
	}
}
