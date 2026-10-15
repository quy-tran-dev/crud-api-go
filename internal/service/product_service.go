package service

import (
	"errors"

	"crud-api/internal/model"
	"crud-api/internal/repository"

	"gorm.io/gorm"
)

var (
	ErrProductNotFound = errors.New("không tìm thấy sản phẩm")
)

type ProductService interface {
	CreateProduct(req model.CreateProductRequest) (*model.Product, error)
	GetProducts(page, limit int, search string) ([]model.Product, int64, error)
	GetProductByID(id uint) (*model.Product, error)
	UpdateProduct(id uint, req model.UpdateProductRequest) (*model.Product, error)
	DeleteProduct(id uint) error
}

type productService struct {
	repo repository.ProductRepository
}

func NewProductService(repo repository.ProductRepository) ProductService {
	return &productService{repo: repo}
}

func (s *productService) CreateProduct(req model.CreateProductRequest) (*model.Product, error) {
	product := &model.Product{
		Name:        req.Name,
		Description: req.Description,
		Price:       req.Price,
		Stock:       *req.Stock,
	}

	if err := s.repo.Create(product); err != nil {
		return nil, err
	}

	return product, nil
}

func (s *productService) GetProducts(page, limit int, search string) ([]model.Product, int64, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	return s.repo.FindAll(page, limit, search)
}

func (s *productService) GetProductByID(id uint) (*model.Product, error) {
	product, err := s.repo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrProductNotFound
		}
		return nil, err
	}
	return product, nil
}

func (s *productService) UpdateProduct(id uint, req model.UpdateProductRequest) (*model.Product, error) {
	product, err := s.GetProductByID(id)
	if err != nil {
		return nil, err
	}

	if req.Name != "" {
		product.Name = req.Name
	}
	if req.Description != nil {
		product.Description = *req.Description
	}
	if req.Price != nil {
		product.Price = *req.Price
	}
	if req.Stock != nil {
		product.Stock = *req.Stock
	}

	if err := s.repo.Update(product); err != nil {
		return nil, err
	}

	return product, nil
}

func (s *productService) DeleteProduct(id uint) error {
	_, err := s.GetProductByID(id)
	if err != nil {
		return err
	}

	return s.repo.Delete(id)
}
