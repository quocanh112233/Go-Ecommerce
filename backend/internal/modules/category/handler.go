package category

import (
	"net/http"
	"strconv"

	"go-ecommerce/internal/shared/errors"

	"github.com/gin-gonic/gin"
)

// Handler handles category HTTP requests
type Handler struct {
	service Service
}

// NewHandler creates a new category handler
func NewHandler(service Service) *Handler {
	return &Handler{service: service}
}

// Create handles POST /admin/categories
// @Summary Tạo danh mục mới
// @Description Tạo một danh mục sản phẩm mới (Admin only)
// @Tags Categories
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body CreateCategoryRequest true "Category data"
// @Success 201 {object} CategoryResponse
// @Failure 400 {object} map[string]string
// @Router /admin/categories [post]
func (h *Handler) Create(c *gin.Context) {
	var req CreateCategoryRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	res, err := h.service.Create(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Lỗi tạo danh mục"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Tạo danh mục thành công",
		"data":    res,
	})
}

// GetAll handles GET /admin/categories
// @Summary Lấy tất cả danh mục
// @Description Lấy danh sách tất cả danh mục (Admin only)
// @Tags Categories
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {array} CategoryResponse
// @Router /admin/categories [get]
func (h *Handler) GetAll(c *gin.Context) {
	res, err := h.service.GetAll(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Lỗi lấy danh sách"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": res})
}

// GetByID handles GET /admin/categories/:id
// @Summary Lấy danh mục theo ID
// @Description Lấy chi tiết một danh mục (Admin only)
// @Tags Categories
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Category ID"
// @Success 200 {object} CategoryResponse
// @Failure 404 {object} map[string]string
// @Router /admin/categories/{id} [get]
func (h *Handler) GetByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID không hợp lệ"})
		return
	}

	res, err := h.service.GetByID(c.Request.Context(), uint(id))
	if err != nil {
		if err == errors.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Không tìm thấy danh mục"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Lỗi hệ thống"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": res})
}

// Update handles PUT /admin/categories/:id
// @Summary Cập nhật danh mục
// @Description Cập nhật thông tin danh mục (Admin only)
// @Tags Categories
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Category ID"
// @Param request body UpdateCategoryRequest true "Category data"
// @Success 200 {object} CategoryResponse
// @Failure 404 {object} map[string]string
// @Router /admin/categories/{id} [put]
func (h *Handler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID không hợp lệ"})
		return
	}

	var req UpdateCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	res, err := h.service.Update(c.Request.Context(), uint(id), req)
	if err != nil {
		if err == errors.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Không tìm thấy danh mục"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Lỗi hệ thống"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Cập nhật thành công",
		"data":    res,
	})
}

// Delete handles DELETE /admin/categories/:id
// @Summary Xóa danh mục
// @Description Xóa một danh mục (Admin only)
// @Tags Categories
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Category ID"
// @Success 200 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /admin/categories/{id} [delete]
func (h *Handler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID không hợp lệ"})
		return
	}

	err = h.service.Delete(c.Request.Context(), uint(id))
	if err != nil {
		if err == errors.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Không tìm thấy danh mục"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Lỗi hệ thống"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Xóa danh mục thành công"})
}
