package service

import (
	"errors"
	"fmt"
	"strconv"
	"time"

	"crud-api/internal/config"
	"crud-api/internal/model"
	"crud-api/internal/repository"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

var (
	ErrEmailExists    = errors.New("email đã được sử dụng")
	ErrUsernameExists = errors.New("tên đăng nhập đã được sử dụng")
	ErrInvalidCreds   = errors.New("email hoặc mật khẩu không chính xác")
	ErrUserNotFound   = errors.New("không tìm thấy người dùng")
	ErrInvalidToken   = errors.New("token không hợp lệ hoặc đã hết hạn")
)

type JWTClaims struct {
	UserID   uint   `json:"user_id"`
	Username string `json:"username"`
	Role     string `json:"role"`
	jwt.RegisteredClaims
}

type AuthService interface {
	Register(req model.RegisterRequest) (*model.AuthResponse, error)
	Login(req model.LoginRequest) (*model.AuthResponse, error)
	GetProfile(userID uint) (*model.UserResponse, error)
	GenerateToken(user *model.User) (string, error)
	ValidateToken(tokenString string) (*JWTClaims, error)
}

type authService struct {
	repo repository.UserRepository
	cfg  *config.Config
}

func NewAuthService(repo repository.UserRepository, cfg *config.Config) AuthService {
	return &authService{
		repo: repo,
		cfg:  cfg,
	}
}

func (s *authService) Register(req model.RegisterRequest) (*model.AuthResponse, error) {
	// 1. Kiểm tra email đã tồn tại chưa
	if existing, err := s.repo.FindByEmail(req.Email); err == nil && existing != nil {
		return nil, ErrEmailExists
	}

	// 2. Kiểm tra username đã tồn tại chưa
	if existing, err := s.repo.FindByUsername(req.Username); err == nil && existing != nil {
		return nil, ErrUsernameExists
	}

	// 3. Băm mật khẩu (Hash password bằng Bcrypt)
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("lỗi khi mã hóa mật khẩu: %w", err)
	}

	// 4. Thiết lập role mặc định nếu không truyền
	role := model.RoleUser
	if req.Role != "" {
		role = req.Role
	}

	user := &model.User{
		Username: req.Username,
		Email:    req.Email,
		Password: string(hashedPassword),
		Role:     role,
	}

	// 5. Lưu vào Database
	if err := s.repo.Create(user); err != nil {
		return nil, fmt.Errorf("lỗi khi lưu người dùng: %w", err)
	}

	// 6. Tự động sinh JWT token
	token, err := s.GenerateToken(user)
	if err != nil {
		return nil, err
	}

	return &model.AuthResponse{
		User:  user.ToResponse(),
		Token: token,
	}, nil
}

func (s *authService) Login(req model.LoginRequest) (*model.AuthResponse, error) {
	// 1. Tìm user theo Email
	user, err := s.repo.FindByEmail(req.Email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrInvalidCreds
		}
		return nil, err
	}

	// 2. So sánh mật khẩu đã hash
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return nil, ErrInvalidCreds
	}

	// 3. Sinh token JWT
	token, err := s.GenerateToken(user)
	if err != nil {
		return nil, err
	}

	return &model.AuthResponse{
		User:  user.ToResponse(),
		Token: token,
	}, nil
}

func (s *authService) GetProfile(userID uint) (*model.UserResponse, error) {
	user, err := s.repo.FindByID(userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}
	res := user.ToResponse()
	return &res, nil
}

func (s *authService) GenerateToken(user *model.User) (string, error) {
	hours, err := strconv.Atoi(s.cfg.JWTExpirationHours)
	if err != nil || hours <= 0 {
		hours = 24
	}

	expirationTime := time.Now().Add(time.Duration(hours) * time.Hour)

	claims := &JWTClaims{
		UserID:   user.ID,
		Username: user.Username,
		Role:     user.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(s.cfg.JWTSecret))
	if err != nil {
		return "", fmt.Errorf("không thể ký token: %w", err)
	}

	return tokenString, nil
}

func (s *authService) ValidateToken(tokenString string) (*JWTClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("phương thức mã hóa không hợp lệ: %v", token.Header["alg"])
		}
		return []byte(s.cfg.JWTSecret), nil
	})

	if err != nil {
		return nil, ErrInvalidToken
	}

	if claims, ok := token.Claims.(*JWTClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, ErrInvalidToken
}
