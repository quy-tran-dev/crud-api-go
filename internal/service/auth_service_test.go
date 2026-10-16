package service_test

import (
	"testing"

	"crud-api/internal/config"
	"crud-api/internal/model"
	"crud-api/internal/service"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// -----------------------------------------------------------------------
// 1. TẠO MOCK REPOSITORY (Giả lập tầng Database không cần kết nối thật)
// -----------------------------------------------------------------------
type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) Create(user *model.User) error {
	args := m.Called(user)
	// Gán ID giả lập khi tạo mới thành công
	user.ID = 1
	return args.Error(0)
}

func (m *MockUserRepository) FindByEmail(email string) (*model.User, error) {
	args := m.Called(email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.User), args.Error(1)
}

func (m *MockUserRepository) FindByUsername(username string) (*model.User, error) {
	args := m.Called(username)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.User), args.Error(1)
}

func (m *MockUserRepository) FindByID(id uint) (*model.User, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.User), args.Error(1)
}

// Helper hàm tạo cấu hình mẫu cho bài test
func testConfig() *config.Config {
	return &config.Config{
		JWTSecret:          "test_secret_key_123456",
		JWTExpirationHours: "24",
	}
}

// -----------------------------------------------------------------------
// 2. CÁC BÀI TEST CHO CHỨC NĂNG ĐĂNG KÝ (REGISTER)
// -----------------------------------------------------------------------

func TestRegister_Success(t *testing.T) {
	// [A] - ARRANGE (Chuẩn bị)
	mockRepo := new(MockUserRepository)
	cfg := testConfig()
	authService := service.NewAuthService(mockRepo, cfg)

	req := model.RegisterRequest{
		Username: "vana",
		Email:    "vana@gmail.com",
		Password: "password123",
		Role:     "user",
	}

	// Thiết lập kỳ vọng cho Mock:
	// Khi gọi FindByEmail -> trả về nil (chưa tồn tại)
	mockRepo.On("FindByEmail", req.Email).Return(nil, gorm.ErrRecordNotFound)
	// Khi gọi FindByUsername -> trả về nil (chưa tồn tại)
	mockRepo.On("FindByUsername", req.Username).Return(nil, gorm.ErrRecordNotFound)
	// Khi gọi Create -> trả về nil (không có lỗi)
	mockRepo.On("Create", mock.AnythingOfType("*model.User")).Return(nil)

	// [A] - ACT (Thực thi)
	resp, err := authService.Register(req)

	// [A] - ASSERT (Kiểm chứng kết quả)
	assert.NoError(t, err, "Đăng ký thành công thì err phải là nil")
	assert.NotNil(t, resp, "Kết quả trả về không được rỗng")
	assert.Equal(t, req.Email, resp.User.Email)
	assert.Equal(t, req.Username, resp.User.Username)
	assert.NotEmpty(t, resp.Token, "Sau khi đăng ký phải tự sinh ra Token")

	// Đảm bảo các hàm mock đã được gọi đúng như kỳ vọng
	mockRepo.AssertExpectations(t)
}

func TestRegister_EmailAlreadyExists(t *testing.T) {
	// ARRANGE
	mockRepo := new(MockUserRepository)
	cfg := testConfig()
	authService := service.NewAuthService(mockRepo, cfg)

	req := model.RegisterRequest{
		Username: "vana",
		Email:    "exists@gmail.com",
		Password: "password123",
	}

	// Giả lập đã có user trùng email
	existingUser := &model.User{ID: 99, Email: req.Email}
	mockRepo.On("FindByEmail", req.Email).Return(existingUser, nil)

	// ACT
	resp, err := authService.Register(req)

	// ASSERT
	assert.Nil(t, resp)
	assert.ErrorIs(t, err, service.ErrEmailExists, "Phải báo lỗi trùng Email")
	mockRepo.AssertExpectations(t)
}

func TestRegister_UsernameAlreadyExists(t *testing.T) {
	// ARRANGE
	mockRepo := new(MockUserRepository)
	cfg := testConfig()
	authService := service.NewAuthService(mockRepo, cfg)

	req := model.RegisterRequest{
		Username: "existing_user",
		Email:    "new@gmail.com",
		Password: "password123",
	}

	// Email chưa tồn tại nhưng username đã tồn tại
	mockRepo.On("FindByEmail", req.Email).Return(nil, gorm.ErrRecordNotFound)
	mockRepo.On("FindByUsername", req.Username).Return(&model.User{ID: 10, Username: req.Username}, nil)

	// ACT
	resp, err := authService.Register(req)

	// ASSERT
	assert.Nil(t, resp)
	assert.ErrorIs(t, err, service.ErrUsernameExists, "Phải báo lỗi trùng Username")
	mockRepo.AssertExpectations(t)
}

// -----------------------------------------------------------------------
// 3. CÁC BÀI TEST CHO CHỨC NĂNG ĐĂNG NHẬP (LOGIN)
// -----------------------------------------------------------------------

func TestLogin_Success(t *testing.T) {
	// ARRANGE
	mockRepo := new(MockUserRepository)
	cfg := testConfig()
	authService := service.NewAuthService(mockRepo, cfg)

	rawPassword := "mypassword123"
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(rawPassword), bcrypt.DefaultCost)

	userInDB := &model.User{
		ID:       1,
		Username: "vana",
		Email:    "vana@gmail.com",
		Password: string(hashedPassword),
		Role:     "user",
	}

	// Mock tìm thấy user theo Email
	mockRepo.On("FindByEmail", userInDB.Email).Return(userInDB, nil)

	// ACT
	loginReq := model.LoginRequest{
		Email:    "vana@gmail.com",
		Password: rawPassword,
	}
	resp, err := authService.Login(loginReq)

	// ASSERT
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.NotEmpty(t, resp.Token, "Đăng nhập thành công phải trả về JWT token")
	assert.Equal(t, userInDB.Email, resp.User.Email)
	mockRepo.AssertExpectations(t)
}

func TestLogin_WrongPassword(t *testing.T) {
	// ARRANGE
	mockRepo := new(MockUserRepository)
	cfg := testConfig()
	authService := service.NewAuthService(mockRepo, cfg)

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("correct_password"), bcrypt.DefaultCost)
	userInDB := &model.User{
		ID:       1,
		Email:    "vana@gmail.com",
		Password: string(hashedPassword),
	}

	mockRepo.On("FindByEmail", userInDB.Email).Return(userInDB, nil)

	// ACT: Nhập sai mật khẩu
	loginReq := model.LoginRequest{
		Email:    "vana@gmail.com",
		Password: "wrong_password_123",
	}
	resp, err := authService.Login(loginReq)

	// ASSERT
	assert.Nil(t, resp)
	assert.ErrorIs(t, err, service.ErrInvalidCreds, "Mật khẩu sai phải báo lỗi không hợp lệ")
	mockRepo.AssertExpectations(t)
}

func TestLogin_UserNotFound(t *testing.T) {
	// ARRANGE
	mockRepo := new(MockUserRepository)
	cfg := testConfig()
	authService := service.NewAuthService(mockRepo, cfg)

	mockRepo.On("FindByEmail", "notfound@gmail.com").Return(nil, gorm.ErrRecordNotFound)

	// ACT
	loginReq := model.LoginRequest{
		Email:    "notfound@gmail.com",
		Password: "any_password",
	}
	resp, err := authService.Login(loginReq)

	// ASSERT
	assert.Nil(t, resp)
	assert.ErrorIs(t, err, service.ErrInvalidCreds, "Không tìm thấy email phải báo lỗi không hợp lệ")
	mockRepo.AssertExpectations(t)
}

// -----------------------------------------------------------------------
// 4. TEST TOKEN VALIDATION (XÁC THỰC JWT)
// -----------------------------------------------------------------------

func TestValidateToken_SuccessAndFail(t *testing.T) {
	mockRepo := new(MockUserRepository)
	cfg := testConfig()
	authService := service.NewAuthService(mockRepo, cfg)

	user := &model.User{
		ID:       42,
		Username: "boss",
		Role:     "admin",
	}

	// Sinh token hợp lệ
	tokenStr, err := authService.GenerateToken(user)
	assert.NoError(t, err)
	assert.NotEmpty(t, tokenStr)

	// Test case 1: Token hợp lệ
	claims, err := authService.ValidateToken(tokenStr)
	assert.NoError(t, err)
	assert.Equal(t, uint(42), claims.UserID)
	assert.Equal(t, "admin", claims.Role)

	// Test case 2: Token rác / token giả mạo
	invalidClaims, err := authService.ValidateToken("fake.token.string")
	assert.Error(t, err)
	assert.Nil(t, invalidClaims)
}
