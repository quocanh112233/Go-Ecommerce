package user

import (
	"net/http"

	"go-ecommerce/internal/shared/errors"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// Handler xử lý các request liên quan đến User
type Handler struct {
	service Service
}

// NewHandler khởi tạo Handler
func NewHandler(service Service) *Handler {
	return &Handler{service: service}
}

// Register xử lý request đăng ký tài khoản
// @Summary Đăng ký tài khoản mới
// @Description Tạo tài khoản user với email và password
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body RegisterRequest true "Thông tin đăng ký"
// @Success 201 {object} UserResponse
// @Failure 400 {object} map[string]string
// @Failure 409 {object} map[string]string
// @Router /auth/register [post]
func (h *Handler) Register(c *gin.Context) {
	var req RegisterRequest

	// 1. Parse & Validate JSON
	// ShouldBindJSON sẽ tự động kiểm tra các tag binding trong DTO (required, email...)
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 2. Gọi Service
	res, err := h.service.Register(c.Request.Context(), req)
	if err != nil {
		// Map các lỗi từ Service sang HTTP Status Code
		if err == errors.ErrEmailAlreadyExists {
			c.JSON(http.StatusConflict, gin.H{"error": "Email đã tồn tại"})
			return
		}

		// Các lỗi khác (DB lỗi, Server lỗi...)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Lỗi hệ thống, vui lòng thử lại sau"})
		return
	}

	// 3. Trả về kết quả thành công (201 Created)
	c.JSON(http.StatusCreated, gin.H{
		"message": "Đăng ký thành công",
		"data":    res,
	})
}

// Login xử lý request đăng nhập
// @Summary Đăng nhập
// @Description Đăng nhập bằng email và password
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body LoginRequest true "Thông tin đăng nhập"
// @Success 200 {object} LoginResponse
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Router /auth/login [post]
func (h *Handler) Login(c *gin.Context) {
	var req LoginRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	res, err := h.service.Login(c.Request.Context(), req)
	if err != nil {
		if err == errors.ErrInvalidCredentials {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Email hoặc mật khẩu không đúng"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Lỗi hệ thống"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Đăng nhập thành công",
		"data":    res,
	})
}

// RefreshToken xử lý request làm mới token
// @Summary Làm mới Access Token
// @Description Dùng Refresh Token để lấy Access Token mới
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body map[string]string true "Refresh Token"
// @Success 200 {object} LoginResponse
// @Failure 401 {object} map[string]string
// @Router /auth/refresh-token [post]
func (h *Handler) RefreshToken(c *gin.Context) {
	var req struct {
		RefreshToken string `json:"refresh_token" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	res, err := h.service.RefreshToken(c.Request.Context(), req.RefreshToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Token không hợp lệ hoặc đã hết hạn"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Làm mới token thành công",
		"data":    res,
	})
}

// Logout xử lý request đăng xuất
// @Summary Đăng xuất
// @Description Xóa session, vô hiệu hóa Refresh Token
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body map[string]string true "Refresh Token"
// @Success 200 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Router /auth/logout [post]
func (h *Handler) Logout(c *gin.Context) {
	var req struct {
		RefreshToken string `json:"refresh_token" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.service.Logout(c.Request.Context(), req.RefreshToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Token không hợp lệ"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Đăng xuất thành công",
	})
}

// GetProfile lấy thông tin user hiện tại
// @Summary Lấy thông tin cá nhân
// @Description Lấy thông tin user đang đăng nhập
// @Tags User
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} UserResponse
// @Failure 401 {object} map[string]string
// @Router /me [get]
func (h *Handler) GetProfile(c *gin.Context) {
	// Lấy userID từ context (đã được set bởi AuthMiddleware)
	userIDStr, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	// Parse UUID
	userID, err := uuid.Parse(userIDStr.(string))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid user ID"})
		return
	}

	res, err := h.service.GetProfile(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": res,
	})
}
