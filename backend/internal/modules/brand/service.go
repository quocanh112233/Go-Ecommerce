package brand

import (
	"context"
	"fmt"
	"mime/multipart"
	"regexp"
	"strings"

	"go-ecommerce/internal/shared/errors"
	"go-ecommerce/pkg/cloudinary"
)

// ProductChecker interface for checking product existence
type ProductChecker interface {
	HasProductsWithBrand(ctx context.Context, brandID uint) (bool, error)
}

// generateSlug creates a URL-friendly slug from name
func generateSlug(name string) string {
	slug := strings.ToLower(strings.TrimSpace(name))
	slug = strings.ReplaceAll(slug, " ", "-")
	reg := regexp.MustCompile(`[^a-z0-9-]+`)
	slug = reg.ReplaceAllString(slug, "")
	reg = regexp.MustCompile(`-+`)
	slug = reg.ReplaceAllString(slug, "-")
	return strings.Trim(slug, "-")
}

// Service interface
type Service interface {
	Create(ctx context.Context, req CreateBrandRequest, logo multipart.File) (*BrandResponse, error)
	GetByID(ctx context.Context, id uint) (*BrandResponse, error)
	GetAll(ctx context.Context) ([]BrandResponse, error)
	Update(ctx context.Context, id uint, req UpdateBrandRequest, logo multipart.File) (*BrandResponse, error)
	Delete(ctx context.Context, id uint) error
}

type service struct {
	repo           Repository
	cloudinary     *cloudinary.Client
	productChecker ProductChecker
}

// NewService creates a new brand service
func NewService(repo Repository, cloudinary *cloudinary.Client, productChecker ProductChecker) Service {
	return &service{repo: repo, cloudinary: cloudinary, productChecker: productChecker}
}

func (s *service) Create(ctx context.Context, req CreateBrandRequest, logo multipart.File) (*BrandResponse, error) {
	brand := &Brand{
		Name:        req.Name,
		Slug:        generateSlug(req.Name),
		Description: req.Description,
	}

	// Upload logo if provided
	if logo != nil {
		result, err := s.cloudinary.Upload(ctx, logo, "brands")
		if err != nil {
			return nil, err
		}
		brand.LogoURL = result.URL
		brand.LogoPublicID = result.PublicID
	}

	if err := s.repo.Create(ctx, brand); err != nil {
		return nil, err
	}

	return ToBrandResponse(brand), nil
}

func (s *service) GetByID(ctx context.Context, id uint) (*BrandResponse, error) {
	brand, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, errors.ErrRecordNotFound
	}
	return ToBrandResponse(brand), nil
}

func (s *service) GetAll(ctx context.Context) ([]BrandResponse, error) {
	brands, err := s.repo.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	var responses []BrandResponse
	for _, b := range brands {
		responses = append(responses, *ToBrandResponse(&b))
	}
	return responses, nil
}

func (s *service) Update(ctx context.Context, id uint, req UpdateBrandRequest, logo multipart.File) (*BrandResponse, error) {
	brand, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, errors.ErrRecordNotFound
	}

	if req.Name != "" {
		brand.Name = req.Name
		brand.Slug = generateSlug(req.Name)
	}
	if req.Description != "" {
		brand.Description = req.Description
	}

	// Upload new logo if provided
	if logo != nil {
		// Delete old logo first
		if brand.LogoPublicID != "" {
			_ = s.cloudinary.Delete(ctx, brand.LogoPublicID)
		}
		result, err := s.cloudinary.Upload(ctx, logo, "brands")
		if err != nil {
			return nil, err
		}
		brand.LogoURL = result.URL
		brand.LogoPublicID = result.PublicID
	}

	if err := s.repo.Update(ctx, brand); err != nil {
		return nil, err
	}

	return ToBrandResponse(brand), nil
}

func (s *service) Delete(ctx context.Context, id uint) error {
	brand, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return errors.ErrRecordNotFound
	}

	// Check if any products use this brand
	if s.productChecker != nil {
		hasProducts, err := s.productChecker.HasProductsWithBrand(ctx, id)
		if err != nil {
			return err
		}
		if hasProducts {
			return fmt.Errorf("cannot delete brand: products are using this brand")
		}
	}

	// Delete logo from Cloudinary
	if brand.LogoPublicID != "" {
		_ = s.cloudinary.Delete(ctx, brand.LogoPublicID)
	}

	return s.repo.Delete(ctx, id)
}
