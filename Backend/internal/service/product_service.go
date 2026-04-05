package service

import (
	"context"
	"github.com/Sidarth-Roy/NorthStar-Intelligence/Backend/internal/dto"
	"github.com/Sidarth-Roy/NorthStar-Intelligence/Backend/internal/repository"
	"github.com/Sidarth-Roy/NorthStar-Intelligence/Backend/pkg/model"
)

type ProductService interface {
	Create(ctx context.Context, req dto.ProductUpsertReq) (*dto.ProductResponse, error)
	Get(ctx context.Context, id uint) (*dto.ProductResponse, error)
	List(ctx context.Context) ([]dto.ProductResponse, error)
	Update(ctx context.Context, id uint, req dto.ProductUpsertReq) (*dto.ProductResponse, error)
	Delete(ctx context.Context, id uint) error
}

type productSvc struct{ repo repository.ProductRepository }

func NewProductSvc(r repository.ProductRepository) ProductService { return &productSvc{repo: r} }

func (s *productSvc) Create(ctx context.Context, req dto.ProductUpsertReq) (*dto.ProductResponse, error) {
	p := &model.Product{
		ProductName: req.ProductName, 
		UnitPrice: req.UnitPrice, 
		CategoryID: req.CategoryID,
		QuantityPerUnit: req.QuantityPerUnit,
	}
	if err := s.repo.Create(ctx, p); err != nil { return nil, err }
	return mapToDTO(p), nil
}

func (s *productSvc) Get(ctx context.Context, id uint) (*dto.ProductResponse, error) {
	p, err := s.repo.GetByID(ctx, id)
	if err != nil { return nil, err }
	return mapToDTO(p), nil
}

func (s *productSvc) List(ctx context.Context) ([]dto.ProductResponse, error) {
	products, err := s.repo.GetAll(ctx)
	if err != nil { return nil, err }
	var res []dto.ProductResponse
	for _, p := range products { res = append(res, *mapToDTO(&p)) }
	return res, nil
}

func (s *productSvc) Update(ctx context.Context, id uint, req dto.ProductUpsertReq) (*dto.ProductResponse, error) {
	p, err := s.repo.GetByID(ctx, id)
	if err != nil { return nil, err }
	
	p.ProductName = req.ProductName
	p.UnitPrice = req.UnitPrice
	p.CategoryID = req.CategoryID
	
	if err := s.repo.Update(ctx, p); err != nil { return nil, err }
	return mapToDTO(p), nil
}

func (s *productSvc) Delete(ctx context.Context, id uint) error {
	return s.repo.Delete(ctx, id)
}

func mapToDTO(p *model.Product) *dto.ProductResponse {
	return &dto.ProductResponse{
		ID: p.ID, 
		ProductName: p.ProductName, 
		UnitPrice: p.UnitPrice,
		CategoryID: p.CategoryID, 
		Active: p.Active, 
		// ModifiedAt: p.UpdatedAt.String(),
	}
}