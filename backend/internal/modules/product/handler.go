package product

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"go-ecommerce/internal/shared/errors"

	"github.com/gin-gonic/gin"
)

// Handler handles product HTTP requests
type Handler struct {
	service      Service
	categoryRepo CategoryGetter
	brandRepo    BrandGetter
}

// CategoryGetter interface for getting category names
type CategoryGetter interface {
	GetByID(ctx context.Context, id uint) (name string, err error)
}

// BrandGetter interface for getting brand names
type BrandGetter interface {
	GetByID(ctx context.Context, id uint) (name string, err error)
}

// NewHandler creates a new product handler
func NewHandler(service Service, categoryRepo CategoryGetter, brandRepo BrandGetter) *Handler {
	return &Handler{
		service:      service,
		categoryRepo: categoryRepo,
		brandRepo:    brandRepo,
	}
}

// Create handles POST /admin/products
// @Summary Tạo sản phẩm mới
// @Description Tạo sản phẩm với variants và images (multipart/form-data)
// @Tags Products
// @Accept multipart/form-data
// @Produce json
// @Security BearerAuth
// @Success 201 {object} ProductResponse
// @Failure 400 {object} map[string]string
// @Router /admin/products [post]
func (h *Handler) Create(c *gin.Context) {
	var req CreateProductRequest

	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get variants JSON string from form
	variantsJSON := c.PostForm("variants")
	if variantsJSON == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "variants field is required"})
		return
	}

	// Validate JSON format early
	var testVariants []VariantInput
	if err := json.Unmarshal([]byte(variantsJSON), &testVariants); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Invalid variants JSON: %v", err)})
		return
	}

	// Get image files
	form, err := c.MultipartForm()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to parse multipart form"})
		return
	}

	imageFiles := form.File["images"]
	if len(imageFiles) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "At least 1 image is required"})
		return
	}
	if len(imageFiles) > 5 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Maximum 5 images allowed"})
		return
	}

	// Get category and brand names for SKU generation
	categoryName, err := h.categoryRepo.GetByID(c.Request.Context(), req.CategoryID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid category_id"})
		return
	}

	brandName, err := h.brandRepo.GetByID(c.Request.Context(), req.BrandID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid brand_id"})
		return
	}

	res, err := h.service.Create(c.Request.Context(), req, variantsJSON, imageFiles, categoryName, brandName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to create product: %v", err)})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Tạo sản phẩm thành công",
		"data":    res,
	})
}

// GetAll handles GET /admin/products
func (h *Handler) GetAll(c *gin.Context) {
	res, err := h.service.GetAll(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get products"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": res})
}

// GetByID handles GET /admin/products/:id
func (h *Handler) GetByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	res, err := h.service.GetByID(c.Request.Context(), uint(id))
	if err != nil {
		if err == errors.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get product"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": res})
}

// Update handles PUT /admin/products/:id
func (h *Handler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	var req UpdateProductRequest
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	res, err := h.service.Update(c.Request.Context(), uint(id), req)
	if err != nil {
		if err == errors.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update product"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Cập nhật thành công",
		"data":    res,
	})
}

// Delete handles DELETE /admin/products/:id
func (h *Handler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	err = h.service.Delete(c.Request.Context(), uint(id))
	if err != nil {
		if err == errors.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete product"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Xóa sản phẩm thành công"})
}
