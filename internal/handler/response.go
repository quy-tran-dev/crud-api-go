package handler

import (
	"math"

	"github.com/gin-gonic/gin"
)

type APIResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
	Error   any    `json:"error,omitempty"`
}

type PaginationMeta struct {
	CurrentPage int   `json:"current_page"`
	Limit       int   `json:"limit"`
	TotalItems  int64 `json:"total_items"`
	TotalPages  int   `json:"total_pages"`
}

type PaginatedResponse struct {
	Success    bool           `json:"success"`
	Message    string         `json:"message"`
	Data       any            `json:"data"`
	Pagination PaginationMeta `json:"pagination"`
}

func ResponseSuccess(c *gin.Context, statusCode int, message string, data any) {
	c.JSON(statusCode, APIResponse{
		Success: true,
		Message: message,
		Data:    data,
	})
}

func ResponsePaginated(c *gin.Context, statusCode int, message string, data any, page, limit int, total int64) {
	totalPages := int(math.Ceil(float64(total) / float64(limit)))
	c.JSON(statusCode, PaginatedResponse{
		Success: true,
		Message: message,
		Data:    data,
		Pagination: PaginationMeta{
			CurrentPage: page,
			Limit:       limit,
			TotalItems:  total,
			TotalPages:  totalPages,
		},
	})
}

func ResponseError(c *gin.Context, statusCode int, message string, err any) {
	errMessage := ""
	if e, ok := err.(error); ok {
		errMessage = e.Error()
	} else if s, ok := err.(string); ok {
		errMessage = s
	}

	c.JSON(statusCode, APIResponse{
		Success: false,
		Message: message,
		Error:   errMessage,
	})
}
