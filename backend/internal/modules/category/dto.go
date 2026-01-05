package category

import "time"

// CreateCategoryRequest - Request body for creating category
type CreateCategoryRequest struct {
	Name        string `json:"name" binding:"required,min=2,max=100"`
	Description string `json:"description"`
}

// UpdateCategoryRequest - Request body for updating category
type UpdateCategoryRequest struct {
	Name        string `json:"name" binding:"omitempty,min=2,max=100"`
	Description string `json:"description"`
}

// CategoryResponse - Response DTO
type CategoryResponse struct {
	ID          uint      `json:"id"`
	Name        string    `json:"name"`
	Slug        string    `json:"slug"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// ToCategoryResponse converts entity to response DTO
func ToCategoryResponse(c *Category) *CategoryResponse {
	return &CategoryResponse{
		ID:          c.ID,
		Name:        c.Name,
		Slug:        c.Slug,
		Description: c.Description,
		CreatedAt:   c.CreatedAt,
		UpdatedAt:   c.UpdatedAt,
	}
}
