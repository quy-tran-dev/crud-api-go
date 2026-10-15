package handler

import (
	"errors"
	"net/http"
	"strconv"

	"crud-api/internal/model"
	"crud-api/internal/service"

	"github.com/gin-gonic/gin"
)

type ProductHandler struct {
	service service.ProductService
}

func NewProductHandler(s service.ProductService) *ProductHandler {
	return &ProductHandler{service: s}
}

func (h *ProductHandler) RegisterRoutes(rg *gin.RouterGroup) {
	products := rg.Group("/products")
	{
		products.POST("", h.CreateProduct)
		products.GET("", h.GetProducts)
		products.GET("/:id", h.GetProductByID)
		products.PUT("/:id", h.UpdateProduct)
		products.DELETE("/:id", h.DeleteProduct)
	}
}

// CreateProduct godoc
// POST /api/v1/products
func (h *ProductHandler) CreateProduct(c *gin.Context) {
	var req model.CreateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		ResponseError(c, http.StatusBadRequest, "Dữ liệu đầu vào không hợp lệ", err)
		return
	}

	product, err := h.service.CreateProduct(req)
	if err != nil {
		ResponseError(c, http.StatusInternalServerError, "Không thể tạo sản phẩm", err)
		return
	}

	ResponseSuccess(c, http.StatusCreated, "Tạo sản phẩm thành công", product)
}

// GetProducts godoc
// GET /api/v1/products?page=1&limit=10&search=keyword
func (h *ProductHandler) GetProducts(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	search := c.Query("search")

	products, total, err := h.service.GetProducts(page, limit, search)
	if err != nil {
		ResponseError(c, http.StatusInternalServerError, "Không thể lấy danh sách sản phẩm", err)
		return
	}

	ResponsePaginated(c, http.StatusOK, "Lấy danh sách sản phẩm thành công", products, page, limit, total)
}

// GetProductByID godoc
// GET /api/v1/products/:id
func (h *ProductHandler) GetProductByID(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		ResponseError(c, http.StatusBadRequest, "ID sản phẩm phải là số nguyên dương", err)
		return
	}

	product, err := h.service.GetProductByID(uint(id))
	if err != nil {
		if errors.Is(err, service.ErrProductNotFound) {
			ResponseError(c, http.StatusNotFound, "Không tìm thấy sản phẩm", err)
			return
		}
		ResponseError(c, http.StatusInternalServerError, "Lỗi khi lấy thông tin sản phẩm", err)
		return
	}

	ResponseSuccess(c, http.StatusOK, "Lấy thông tin sản phẩm thành công", product)
}

// UpdateProduct godoc
// PUT /api/v1/products/:id
func (h *ProductHandler) UpdateProduct(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		ResponseError(c, http.StatusBadRequest, "ID sản phẩm phải là số nguyên dương", err)
		return
	}

	var req model.UpdateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		ResponseError(c, http.StatusBadRequest, "Dữ liệu cập nhật không hợp lệ", err)
		return
	}

	product, err := h.service.UpdateProduct(uint(id), req)
	if err != nil {
		if errors.Is(err, service.ErrProductNotFound) {
			ResponseError(c, http.StatusNotFound, "Không tìm thấy sản phẩm để cập nhật", err)
			return
		}
		ResponseError(c, http.StatusInternalServerError, "Lỗi khi cập nhật sản phẩm", err)
		return
	}

	ResponseSuccess(c, http.StatusOK, "Cập nhật sản phẩm thành công", product)
}

// DeleteProduct godoc
// DELETE /api/v1/products/:id
func (h *ProductHandler) DeleteProduct(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		ResponseError(c, http.StatusBadRequest, "ID sản phẩm phải là số nguyên dương", err)
		return
	}

	if err := h.service.DeleteProduct(uint(id)); err != nil {
		if errors.Is(err, service.ErrProductNotFound) {
			ResponseError(c, http.StatusNotFound, "Không tìm thấy sản phẩm để xóa", err)
			return
		}
		ResponseError(c, http.StatusInternalServerError, "Lỗi khi xóa sản phẩm", err)
		return
	}

	ResponseSuccess(c, http.StatusOK, "Xóa sản phẩm thành công", nil)
}
