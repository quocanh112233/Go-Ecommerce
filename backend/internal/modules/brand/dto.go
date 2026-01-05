package brand

import "time"

// CreateBrandRequest - Request body for creating brand (form-data)
type CreateBrandRequest struct {
	Name        string `form:"name" binding:"required,min=2,max=100"`
	Description string `form:"description"`
}

// UpdateBrandRequest - Request body for updating brand
type UpdateBrandRequest struct {
	Name        string `form:"name" binding:"omitempty,min=2,max=100"`
	Description string `form:"description"`
}

// BrandResponse - Response DTO
type BrandResponse struct {
	ID          uint      `json:"id"`
	Name        string    `json:"name"`
	Slug        string    `json:"slug"`
	Description string    `json:"description"`
	LogoURL     string    `json:"logo_url"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// ToBrandResponse converts entity to response DTO
func ToBrandResponse(b *Brand) *BrandResponse {
	return &BrandResponse{
		ID:          b.ID,
		Name:        b.Name,
		Slug:        b.Slug,
		Description: b.Description,
		LogoURL:     b.LogoURL,
		CreatedAt:   b.CreatedAt,
		UpdatedAt:   b.UpdatedAt,
	}
}
