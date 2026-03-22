package service

import (
	"context"
	"github.com/Sidarth-Roy/NorthStar-Intelligence/Backend/internal/dto"
	"github.com/Sidarth-Roy/NorthStar-Intelligence/Backend/internal/repository"
	"github.com/Sidarth-Roy/NorthStar-Intelligence/Backend/pkg/model"
)

type CategoryService interface {
	Create(ctx context.Context, req dto.CategoryUpsertReq) (*dto.CategoryResponse, error)
	Get(ctx context.Context, id uint) (*dto.CategoryResponse, error)
	List(ctx context.Context) ([]dto.CategoryResponse, error)
	Update(ctx context.Context, id uint, req dto.CategoryUpsertReq) (*dto.CategoryResponse, error)
	Delete(ctx context.Context, id uint) error
}

type categorySvc struct{ repo repository.CategoryRepository }

func NewCategorySvc(r repository.CategoryRepository) CategoryService { return &categorySvc{repo: r} }

func (s *categorySvc) Create(ctx context.Context, req dto.CategoryUpsertReq) (*dto.CategoryResponse, error) {
	c := &model.Category{
		CategoryName: req.CategoryName,
		Description:  req.Description,
	}
	if err := s.repo.Create(ctx, c); err != nil { return nil, err }
	return mapCategoryToDTO(c), nil
}

func (s *categorySvc) Get(ctx context.Context, id uint) (*dto.CategoryResponse, error) {
	c, err := s.repo.GetByID(ctx, id)
	if err != nil { return nil, err }
	return mapCategoryToDTO(c), nil
}

func (s *categorySvc) List(ctx context.Context) ([]dto.CategoryResponse, error) {
	categories, err := s.repo.GetAll(ctx)
	if err != nil { return nil, err }
	var res []dto.CategoryResponse
	for _, c := range categories { res = append(res, *mapCategoryToDTO(&c)) }
	return res, nil
}

func (s *categorySvc) Update(ctx context.Context, id uint, req dto.CategoryUpsertReq) (*dto.CategoryResponse, error) {
	c, err := s.repo.GetByID(ctx, id)
	if err != nil { return nil, err }
	
	c.CategoryName = req.CategoryName
	c.Description = req.Description
	
	if err := s.repo.Update(ctx, c); err != nil { return nil, err }
	return mapCategoryToDTO(c), nil
}

func (s *categorySvc) Delete(ctx context.Context, id uint) error {
	return s.repo.Delete(ctx, id)
}

func mapCategoryToDTO(c *model.Category) *dto.CategoryResponse {
	return &dto.CategoryResponse{
		ID:           c.ID,
		CategoryName: c.CategoryName,
		Description:  c.Description,
		Active:       c.Active,
		ModifiedAt:   c.UpdatedAt.String(),
	}
}