package category

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	"go-ecommerce/internal/shared/errors"
)

// ProductChecker interface for checking product existence
type ProductChecker interface {
	HasProductsWithCategory(ctx context.Context, categoryID uint) (bool, error)
}

// generateSlug creates a URL-friendly slug from name
func generateSlug(name string) string {
	// Convert to lowercase and replace spaces with hyphens
	slug := strings.ToLower(strings.TrimSpace(name))
	slug = strings.ReplaceAll(slug, " ", "-")
	// Remove special characters
	reg := regexp.MustCompile(`[^a-z0-9-]+`)
	slug = reg.ReplaceAllString(slug, "")
	// Remove multiple hyphens
	reg = regexp.MustCompile(`-+`)
	slug = reg.ReplaceAllString(slug, "-")
	return strings.Trim(slug, "-")
}

// Service interface
type Service interface {
	Create(ctx context.Context, req CreateCategoryRequest) (*CategoryResponse, error)
	GetByID(ctx context.Context, id uint) (*CategoryResponse, error)
	GetAll(ctx context.Context) ([]CategoryResponse, error)
	Update(ctx context.Context, id uint, req UpdateCategoryRequest) (*CategoryResponse, error)
	Delete(ctx context.Context, id uint) error
}

type service struct {
	repo           Repository
	productChecker ProductChecker
}

// NewService creates a new category service
func NewService(repo Repository, productChecker ProductChecker) Service {
	return &service{repo: repo, productChecker: productChecker}
}

func (s *service) Create(ctx context.Context, req CreateCategoryRequest) (*CategoryResponse, error) {
	category := &Category{
		Name:        req.Name,
		Slug:        generateSlug(req.Name),
		Description: req.Description,
	}

	if err := s.repo.Create(ctx, category); err != nil {
		return nil, err
	}

	return ToCategoryResponse(category), nil
}

func (s *service) GetByID(ctx context.Context, id uint) (*CategoryResponse, error) {
	category, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, errors.ErrRecordNotFound
	}

	return ToCategoryResponse(category), nil
}

func (s *service) GetAll(ctx context.Context) ([]CategoryResponse, error) {
	categories, err := s.repo.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	var responses []CategoryResponse
	for _, c := range categories {
		responses = append(responses, *ToCategoryResponse(&c))
	}

	return responses, nil
}

func (s *service) Update(ctx context.Context, id uint, req UpdateCategoryRequest) (*CategoryResponse, error) {
	category, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, errors.ErrRecordNotFound
	}

	if req.Name != "" {
		category.Name = req.Name
		category.Slug = generateSlug(req.Name)
	}
	if req.Description != "" {
		category.Description = req.Description
	}

	if err := s.repo.Update(ctx, category); err != nil {
		return nil, err
	}

	return ToCategoryResponse(category), nil
}

func (s *service) Delete(ctx context.Context, id uint) error {
	_, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return errors.ErrRecordNotFound
	}

	// Check if any products use this category
	if s.productChecker != nil {
		hasProducts, err := s.productChecker.HasProductsWithCategory(ctx, id)
		if err != nil {
			return err
		}
		if hasProducts {
			return fmt.Errorf("cannot delete category: products are using this category")
		}
	}

	return s.repo.Delete(ctx, id)
}
