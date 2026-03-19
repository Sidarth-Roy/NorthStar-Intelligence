package service

import (
	"github.com/Sidarth-Roy/NorthStar-Intelligence/Backend/internal/dto"
	"github.com/Sidarth-Roy/NorthStar-Intelligence/Backend/internal/repository"
	"github.com/Sidarth-Roy/NorthStar-Intelligence/Backend/pkg/model"
)

type ProductService struct {
	repo repository.ProductRepository
}

func NewProductService(repo repository.ProductRepository) *ProductService {
	return &ProductService{repo: repo}
}

func (s *ProductService) Create(req dto.ProductCreateRequest) (*dto.ProductResponse, error) {
	p := &model.Product{
		ProductName:     req.ProductName,
		QuantityPerUnit: req.QuantityPerUnit,
		UnitPrice:       req.UnitPrice,
		CategoryID:      req.CategoryID,
		Discontinued:    req.Discontinued,
	}

	if err := s.repo.Create(p); err != nil {
		return nil, err
	}

	return s.mapToResponse(p), nil
}

func (s *ProductService) GetByID(id uint) (*dto.ProductResponse, error) {
	p, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}
	return s.mapToResponse(p), nil
}

// Helper to map Model -> DTO
func (s *ProductService) mapToResponse(p *model.Product) *dto.ProductResponse {
	return &dto.ProductResponse{
		ID:              p.ID,
		ProductName:     p.ProductName,
		UnitPrice:       p.UnitPrice,
		QuantityPerUnit: p.QuantityPerUnit,
		CategoryID:      p.CategoryID,
		Active:          p.Active,
		CreatedAt:       p.CreatedAt,
		ModifiedAt:      p.UpdatedAt,
	}
}