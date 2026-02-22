package category

import (
	"context"
	"daymark/internal/models"
	"strings"
)

type Service struct {
	repo *Repository
}

func NewService(r *Repository) *Service {
	return &Service{repo: r}
}

func (s *Service) ListCategories(ctx context.Context) ([]models.Category, error) {
	return s.repo.GetAll(ctx)
}

func (s *Service) GetCategoryByID(ctx context.Context, id uint) (models.Category, error) {
	return s.repo.GetById(ctx, id)
}

func (s *Service) CreateCategory(ctx context.Context, Name string) (models.Category, error) {
	category := models.Category{
		Name: Name,
		Slug: makeSlug(Name),
	}
	if err := s.repo.Create(ctx, category); err != nil {
		return models.Category{}, err
	}
	return category, nil
}

func makeSlug(name string) string {
	slug := strings.ToLower(strings.TrimSpace(name))
	slug = strings.ReplaceAll(slug, " ", "-")
	// Keep only letters, numbers and dashes
	var b strings.Builder
	for _, r := range slug {
		if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') || r == '-' {
			b.WriteRune(r)
		}
	}
	return b.String()
}
