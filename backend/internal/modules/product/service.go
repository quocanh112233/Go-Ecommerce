package product

import (
	"context"
	"encoding/json"
	"fmt"
	"mime/multipart"

	"go-ecommerce/internal/shared/errors"
	"go-ecommerce/pkg/cloudinary"

	"gorm.io/gorm"
)

// Service interface
type Service interface {
	Create(ctx context.Context, req CreateProductRequest, variantsJSON string, imageFiles []*multipart.FileHeader, categoryName, brandName string) (*ProductResponse, error)
	GetByID(ctx context.Context, id uint) (*ProductResponse, error)
	GetAll(ctx context.Context) ([]ProductResponse, error)
	Update(ctx context.Context, id uint, req UpdateProductRequest) (*ProductResponse, error)
	Delete(ctx context.Context, id uint) error
}

type service struct {
	repo       Repository
	cloudinary *cloudinary.Client
}

// NewService creates a new product service
func NewService(repo Repository, cloudinary *cloudinary.Client) Service {
	return &service{repo: repo, cloudinary: cloudinary}
}

func (s *service) Create(ctx context.Context, req CreateProductRequest, variantsJSON string, imageFiles []*multipart.FileHeader, categoryName, brandName string) (*ProductResponse, error) {
	// 1. Validate variants JSON
	if len(variantsJSON) == 0 {
		return nil, fmt.Errorf("variants are required")
	}

	// 2. Parse variants from JSON string
	var variantInputs []VariantInput
	if err := json.Unmarshal([]byte(variantsJSON), &variantInputs); err != nil {
		return nil, fmt.Errorf("invalid variants JSON: %w", err)
	}

	if len(variantInputs) == 0 {
		return nil, fmt.Errorf("at least 1 variant is required")
	}

	// 3. Validate images (min 1, max 5)
	if len(imageFiles) == 0 {
		return nil, fmt.Errorf("at least 1 image is required")
	}
	if len(imageFiles) > 5 {
		return nil, fmt.Errorf("maximum 5 images allowed")
	}

	var createdProduct *Product

	// 4. Start Transaction
	err := s.repo.WithTransaction(func(tx *gorm.DB) error {
		// Step 4.1: Create Product
		product := &Product{
			Name:        req.Name,
			Slug:        generateSlug(req.Name),
			Description: req.Description,
			CategoryID:  req.CategoryID,
			BrandID:     req.BrandID,
			TotalStock:  0, // Will be calculated later
		}

		if err := tx.Create(product).Error; err != nil {
			return fmt.Errorf("failed to create product: %w", err)
		}

		// Step 4.2: Create Variants with 2-step SKU generation
		var variants []ProductVariant
		for _, input := range variantInputs {
			variant := ProductVariant{
				ProductID: product.ID,
				Price:     input.Price,
				Stock:     input.Stock,
				Size:      input.Size,
			}

			// Create variant to get ID
			if err := tx.Create(&variant).Error; err != nil {
				return fmt.Errorf("failed to create variant: %w", err)
			}

			// Generate SKU with the new ID
			sku := generateSKU(categoryName, product.Name, variant.ID)
			if err := tx.Model(&variant).Update("sku", sku).Error; err != nil {
				return fmt.Errorf("failed to update SKU: %w", err)
			}

			variant.SKU = sku
			variants = append(variants, variant)
		}

		// Step 4.3: Upload images to Cloudinary
		var images []ProductImage
		for i, fileHeader := range imageFiles {
			file, err := fileHeader.Open()
			if err != nil {
				return fmt.Errorf("failed to open image file: %w", err)
			}
			defer file.Close()

			result, err := s.cloudinary.Upload(ctx, file, "products")
			if err != nil {
				return fmt.Errorf("failed to upload image: %w", err)
			}

			image := ProductImage{
				ProductID:     product.ID,
				ImageURL:      result.URL,
				ImagePublicID: result.PublicID,
				DisplayOrder:  i + 1, // 1-indexed
			}

			if err := tx.Create(&image).Error; err != nil {
				return fmt.Errorf("failed to save image: %w", err)
			}

			images = append(images, image)
		}

		// Step 4.4: Calculate and update total_stock
		totalStock := calculateTotalStock(variants)
		if err := tx.Model(&product).Update("total_stock", totalStock).Error; err != nil {
			return fmt.Errorf("failed to update total stock: %w", err)
		}

		product.TotalStock = totalStock
		product.Variants = variants
		product.Images = images
		createdProduct = product

		return nil
	})

	if err != nil {
		return nil, err
	}

	return ToProductResponse(createdProduct), nil
}

func (s *service) GetByID(ctx context.Context, id uint) (*ProductResponse, error) {
	product, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, errors.ErrRecordNotFound
	}
	return ToProductResponse(product), nil
}

func (s *service) GetAll(ctx context.Context) ([]ProductResponse, error) {
	products, err := s.repo.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	var responses []ProductResponse
	for _, p := range products {
		responses = append(responses, *ToProductResponse(&p))
	}
	return responses, nil
}

func (s *service) Update(ctx context.Context, id uint, req UpdateProductRequest) (*ProductResponse, error) {
	product, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, errors.ErrRecordNotFound
	}

	if req.Name != "" {
		product.Name = req.Name
		product.Slug = generateSlug(req.Name)
	}
	if req.Description != "" {
		product.Description = req.Description
	}
	if req.CategoryID != 0 {
		product.CategoryID = req.CategoryID
	}
	if req.BrandID != 0 {
		product.BrandID = req.BrandID
	}

	if err := s.repo.Update(ctx, product); err != nil {
		return nil, err
	}

	return ToProductResponse(product), nil
}

func (s *service) Delete(ctx context.Context, id uint) error {
	product, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return errors.ErrRecordNotFound
	}

	// Delete images from Cloudinary
	for _, img := range product.Images {
		if img.ImagePublicID != "" {
			_ = s.cloudinary.Delete(ctx, img.ImagePublicID)
		}
	}

	return s.repo.Delete(ctx, id)
}
