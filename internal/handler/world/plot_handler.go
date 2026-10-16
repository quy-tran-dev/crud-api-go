package world

import (
	"errors"
	"net/http"

	"crud-api/internal/handler"
	"crud-api/internal/model"
	"crud-api/internal/service/world"

	"github.com/gin-gonic/gin"
)

type PlotHandler struct {
	plotService world.PlotService
}

func NewPlotHandler(plotService world.PlotService) *PlotHandler {
	return &PlotHandler{plotService: plotService}
}

// CreatePlot: POST /api/v1/stories/:id/plots
func (h *PlotHandler) CreatePlot(c *gin.Context) {
	userID, ok := getUserID(c)
	if !ok {
		return
	}

	storyID, err := parseID(c.Param("id"))
	if err != nil {
		handler.ResponseError(c, http.StatusBadRequest, "ID tác phẩm không hợp lệ", err)
		return
	}

	var req model.CreatePlotRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		handler.ResponseError(c, http.StatusBadRequest, "Dữ liệu tình tiết không hợp lệ", err)
		return
	}

	plot, err := h.plotService.CreatePlot(storyID, userID, req)
	if err != nil {
		if errors.Is(err, world.ErrStoryNotFound) {
			handler.ResponseError(c, http.StatusNotFound, err.Error(), nil)
			return
		}
		handler.ResponseError(c, http.StatusInternalServerError, "Không thể tạo tình tiết", err)
		return
	}

	handler.ResponseSuccess(c, http.StatusCreated, "Tạo tình tiết thành công", plot)
}

// GetPlots: GET /api/v1/stories/:id/plots?type=main_plot
func (h *PlotHandler) GetPlots(c *gin.Context) {
	userID, ok := getUserID(c)
	if !ok {
		return
	}

	storyID, err := parseID(c.Param("id"))
	if err != nil {
		handler.ResponseError(c, http.StatusBadRequest, "ID tác phẩm không hợp lệ", err)
		return
	}

	plotType := c.Query("type")
	plots, err := h.plotService.GetPlots(storyID, userID, plotType)
	if err != nil {
		if errors.Is(err, world.ErrStoryNotFound) {
			handler.ResponseError(c, http.StatusNotFound, err.Error(), nil)
			return
		}
		handler.ResponseError(c, http.StatusInternalServerError, "Lỗi khi lấy danh sách tình tiết", err)
		return
	}

	handler.ResponseSuccess(c, http.StatusOK, "Lấy danh sách tình tiết thành công", plots)
}

// GetPlotByID: GET /api/v1/stories/:id/plots/:plot_id
func (h *PlotHandler) GetPlotByID(c *gin.Context) {
	userID, ok := getUserID(c)
	if !ok {
		return
	}

	storyID, err := parseID(c.Param("id"))
	if err != nil {
		handler.ResponseError(c, http.StatusBadRequest, "ID tác phẩm không hợp lệ", err)
		return
	}

	plotID, err := parseID(c.Param("plot_id"))
	if err != nil {
		handler.ResponseError(c, http.StatusBadRequest, "ID tình tiết không hợp lệ", err)
		return
	}

	plot, err := h.plotService.GetPlotByID(plotID, storyID, userID)
	if err != nil {
		handler.ResponseError(c, http.StatusNotFound, "Không tìm thấy tình tiết", err)
		return
	}

	handler.ResponseSuccess(c, http.StatusOK, "Lấy thông tin tình tiết thành công", plot)
}

// UpdatePlot: PUT /api/v1/stories/:id/plots/:plot_id
func (h *PlotHandler) UpdatePlot(c *gin.Context) {
	userID, ok := getUserID(c)
	if !ok {
		return
	}

	storyID, err := parseID(c.Param("id"))
	if err != nil {
		handler.ResponseError(c, http.StatusBadRequest, "ID tác phẩm không hợp lệ", err)
		return
	}

	plotID, err := parseID(c.Param("plot_id"))
	if err != nil {
		handler.ResponseError(c, http.StatusBadRequest, "ID tình tiết không hợp lệ", err)
		return
	}

	var req model.UpdatePlotRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		handler.ResponseError(c, http.StatusBadRequest, "Dữ liệu cập nhật không hợp lệ", err)
		return
	}

	plot, err := h.plotService.UpdatePlot(plotID, storyID, userID, req)
	if err != nil {
		handler.ResponseError(c, http.StatusInternalServerError, "Lỗi khi cập nhật tình tiết", err)
		return
	}

	handler.ResponseSuccess(c, http.StatusOK, "Cập nhật tình tiết thành công", plot)
}

// DeletePlot: DELETE /api/v1/stories/:id/plots/:plot_id
func (h *PlotHandler) DeletePlot(c *gin.Context) {
	userID, ok := getUserID(c)
	if !ok {
		return
	}

	storyID, err := parseID(c.Param("id"))
	if err != nil {
		handler.ResponseError(c, http.StatusBadRequest, "ID tác phẩm không hợp lệ", err)
		return
	}

	plotID, err := parseID(c.Param("plot_id"))
	if err != nil {
		handler.ResponseError(c, http.StatusBadRequest, "ID tình tiết không hợp lệ", err)
		return
	}

	if err := h.plotService.DeletePlot(plotID, storyID, userID); err != nil {
		handler.ResponseError(c, http.StatusInternalServerError, "Lỗi khi xóa tình tiết", err)
		return
	}

	handler.ResponseSuccess(c, http.StatusOK, "Xóa tình tiết thành công", nil)
}
