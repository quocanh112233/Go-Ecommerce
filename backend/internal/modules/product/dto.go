package product

import "time"

// VariantInput represents variant data from form
type VariantInput struct {
	Price float64 `json:"price" binding:"required,min=0"`
	Stock int     `json:"stock" binding:"required,min=0"`
	Size  string  `json:"size" binding:"required,min=1,max=50"`
}

// CreateProductRequest - Request body (multipart/form-data)
type CreateProductRequest struct {
	Name        string `form:"name" binding:"required,min=2,max=255"`
	Description string `form:"description"`
	CategoryID  uint   `form:"category_id" binding:"required"`
	BrandID     uint   `form:"brand_id" binding:"required"`
	// Variants will be parsed from JSON string in form-data
	// Images will be uploaded files
}

// UpdateProductRequest - Request body for update
type UpdateProductRequest struct {
	Name        string `form:"name" binding:"omitempty,min=2,max=255"`
	Description string `form:"description"`
	CategoryID  uint   `form:"category_id"`
	BrandID     uint   `form:"brand_id"`
	// Allow updating variants and images
}

// ProductResponse - Full response with nested data
type ProductResponse struct {
	ID          uint              `json:"id"`
	Name        string            `json:"name"`
	Slug        string            `json:"slug"`
	Description string            `json:"description"`
	CategoryID  uint              `json:"category_id"`
	BrandID     uint              `json:"brand_id"`
	TotalStock  int               `json:"total_stock"`
	RatingAvg   float64           `json:"rating_avg"`
	ReviewCount int               `json:"review_count"`
	Variants    []VariantResponse `json:"variants"`
	Images      []ImageResponse   `json:"images"`
	CreatedAt   time.Time         `json:"created_at"`
	UpdatedAt   time.Time         `json:"updated_at"`
}

// VariantResponse - Variant DTO
type VariantResponse struct {
	ID    uint    `json:"id"`
	Price float64 `json:"price"`
	Stock int     `json:"stock"`
	Size  string  `json:"size"`
	SKU   string  `json:"sku"`
}

// ImageResponse - Image DTO
type ImageResponse struct {
	ID           uint   `json:"id"`
	ImageURL     string `json:"image_url"`
	DisplayOrder int    `json:"display_order"`
}

// ToProductResponse converts entity to response DTO
func ToProductResponse(p *Product) *ProductResponse {
	var variants []VariantResponse
	for _, v := range p.Variants {
		variants = append(variants, VariantResponse{
			ID:    v.ID,
			Price: v.Price,
			Stock: v.Stock,
			Size:  v.Size,
			SKU:   v.SKU,
		})
	}

	var images []ImageResponse
	for _, img := range p.Images {
		images = append(images, ImageResponse{
			ID:           img.ID,
			ImageURL:     img.ImageURL,
			DisplayOrder: img.DisplayOrder,
		})
	}

	return &ProductResponse{
		ID:          p.ID,
		Name:        p.Name,
		Slug:        p.Slug,
		Description: p.Description,
		CategoryID:  p.CategoryID,
		BrandID:     p.BrandID,
		TotalStock:  p.TotalStock,
		RatingAvg:   p.RatingAvg,
		ReviewCount: p.ReviewCount,
		Variants:    variants,
		Images:      images,
		CreatedAt:   p.CreatedAt,
		UpdatedAt:   p.UpdatedAt,
	}
}
