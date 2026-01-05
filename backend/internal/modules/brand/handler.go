package brand

import (
	"net/http"
	"strconv"

	"go-ecommerce/internal/shared/errors"

	"github.com/gin-gonic/gin"
)

// Handler handles brand HTTP requests
type Handler struct {
	service Service
}

// NewHandler creates a new brand handler
func NewHandler(service Service) *Handler {
	return &Handler{service: service}
}

// Create handles POST /admin/brands
// @Summary Tạo thương hiệu mới
// @Description Tạo thương hiệu với logo (multipart/form-data)
// @Tags Brands
// @Accept multipart/form-data
// @Produce json
// @Security BearerAuth
// @Param name formData string true "Tên thương hiệu"
// @Param description formData string false "Mô tả"
// @Param logo formData file false "Logo"
// @Success 201 {object} BrandResponse
// @Failure 400 {object} map[string]string
// @Router /admin/brands [post]
func (h *Handler) Create(c *gin.Context) {
	var req CreateBrandRequest

	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get logo file (optional)
	file, _, err := c.Request.FormFile("logo")
	if err != nil && err != http.ErrMissingFile {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Lỗi đọc file logo"})
		return
	}
	if file != nil {
		defer file.Close()
	}

	res, err := h.service.Create(c.Request.Context(), req, file)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Lỗi tạo thương hiệu"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Tạo thương hiệu thành công",
		"data":    res,
	})
}

// GetAll handles GET /admin/brands
func (h *Handler) GetAll(c *gin.Context) {
	res, err := h.service.GetAll(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Lỗi lấy danh sách"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": res})
}

// GetByID handles GET /admin/brands/:id
func (h *Handler) GetByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID không hợp lệ"})
		return
	}

	res, err := h.service.GetByID(c.Request.Context(), uint(id))
	if err != nil {
		if err == errors.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Không tìm thấy thương hiệu"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Lỗi hệ thống"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": res})
}

// Update handles PUT /admin/brands/:id
func (h *Handler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID không hợp lệ"})
		return
	}

	var req UpdateBrandRequest
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get logo file (optional)
	file, _, err := c.Request.FormFile("logo")
	if err != nil && err != http.ErrMissingFile {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Lỗi đọc file logo"})
		return
	}
	if file != nil {
		defer file.Close()
	}

	res, err := h.service.Update(c.Request.Context(), uint(id), req, file)
	if err != nil {
		if err == errors.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Không tìm thấy thương hiệu"})
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

// Delete handles DELETE /admin/brands/:id
func (h *Handler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID không hợp lệ"})
		return
	}

	err = h.service.Delete(c.Request.Context(), uint(id))
	if err != nil {
		if err == errors.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Không tìm thấy thương hiệu"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Lỗi hệ thống"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Xóa thương hiệu thành công"})
}
