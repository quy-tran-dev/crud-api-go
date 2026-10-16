package service_test

import (
	"testing"

	"crud-api/internal/model"
	"crud-api/internal/service"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

// -----------------------------------------------------------------------
// 1. TẠO MOCK PRODUCT REPOSITORY
// -----------------------------------------------------------------------
type MockProductRepository struct {
	mock.Mock
}

func (m *MockProductRepository) Create(product *model.Product) error {
	args := m.Called(product)
	product.ID = 100 // Gán ID giả lập
	return args.Error(0)
}

func (m *MockProductRepository) FindAll(page, limit int, search string) ([]model.Product, int64, error) {
	args := m.Called(page, limit, search)
	return args.Get(0).([]model.Product), args.Get(1).(int64), args.Error(2)
}

func (m *MockProductRepository) FindByID(id uint) (*model.Product, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Product), args.Error(1)
}

func (m *MockProductRepository) Update(product *model.Product) error {
	args := m.Called(product)
	return args.Error(0)
}

func (m *MockProductRepository) Delete(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

// -----------------------------------------------------------------------
// 2. CÁC BÀI TEST CHO PRODUCT SERVICE
// -----------------------------------------------------------------------

func TestCreateProduct_Success(t *testing.T) {
	// ARRANGE
	mockRepo := new(MockProductRepository)
	prodService := service.NewProductService(mockRepo)

	stock := 15
	req := model.CreateProductRequest{
		Name:        "Bàn phím cơ",
		Description: "Bàn phím Bluetooth",
		Price:       1500000,
		Stock:       &stock,
	}

	mockRepo.On("Create", mock.AnythingOfType("*model.Product")).Return(nil)

	// ACT
	product, err := prodService.CreateProduct(req)

	// ASSERT
	assert.NoError(t, err)
	assert.NotNil(t, product)
	assert.Equal(t, uint(100), product.ID)
	assert.Equal(t, req.Name, product.Name)
	assert.Equal(t, req.Price, product.Price)
	mockRepo.AssertExpectations(t)
}

func TestGetProductByID_Success(t *testing.T) {
	// ARRANGE
	mockRepo := new(MockProductRepository)
	prodService := service.NewProductService(mockRepo)

	expectedProduct := &model.Product{
		ID:    5,
		Name:  "Chuột máy tính",
		Price: 500000,
	}

	mockRepo.On("FindByID", uint(5)).Return(expectedProduct, nil)

	// ACT
	product, err := prodService.GetProductByID(5)

	// ASSERT
	assert.NoError(t, err)
	assert.NotNil(t, product)
	assert.Equal(t, uint(5), product.ID)
	assert.Equal(t, "Chuột máy tính", product.Name)
	mockRepo.AssertExpectations(t)
}

func TestGetProductByID_NotFound(t *testing.T) {
	// ARRANGE
	mockRepo := new(MockProductRepository)
	prodService := service.NewProductService(mockRepo)

	// Giả lập GORM không tìm thấy bản ghi
	mockRepo.On("FindByID", uint(999)).Return(nil, gorm.ErrRecordNotFound)

	// ACT
	product, err := prodService.GetProductByID(999)

	// ASSERT
	assert.Nil(t, product)
	assert.ErrorIs(t, err, service.ErrProductNotFound, "Khi DB trả về ErrRecordNotFound thì Service phải ánh xạ sang ErrProductNotFound")
	mockRepo.AssertExpectations(t)
}

func TestDeleteProduct_Success(t *testing.T) {
	// ARRANGE
	mockRepo := new(MockProductRepository)
	prodService := service.NewProductService(mockRepo)

	// Giả lập sản phẩm tồn tại trước khi xóa
	existingProduct := &model.Product{ID: 10, Name: "Màn hình Dell"}
	mockRepo.On("FindByID", uint(10)).Return(existingProduct, nil)
	mockRepo.On("Delete", uint(10)).Return(nil)

	// ACT
	err := prodService.DeleteProduct(10)

	// ASSERT
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}
