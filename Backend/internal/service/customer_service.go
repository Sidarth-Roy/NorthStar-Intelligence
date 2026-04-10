package service

import (
	"context"
	"github.com/Sidarth-Roy/NorthStar-Intelligence/Backend/internal/dto"
	"github.com/Sidarth-Roy/NorthStar-Intelligence/Backend/internal/repository"
	"github.com/Sidarth-Roy/NorthStar-Intelligence/Backend/pkg/model"
)

type CustomerService interface {
	Create(ctx context.Context, req dto.CustomerInsertReq) (*dto.CustomerResponse, error)
	Get(ctx context.Context, id uint) (*dto.CustomerResponse, error)
	List(ctx context.Context) ([]dto.CustomerResponse, error)
	Update(ctx context.Context, id uint, req dto.CustomerUpdateReq) (*dto.CustomerResponse, error)
	Delete(ctx context.Context, id uint) error
}

type customerSvc struct{ repo repository.CustomerRepository }

func NewCustomerSvc(r repository.CustomerRepository) CustomerService { return &customerSvc{repo: r} }

func (s *customerSvc) Create(ctx context.Context, req dto.CustomerInsertReq) (*dto.CustomerResponse, error) {
	c := &model.Customer{
		CustomerID:   req.CustomerID,
		CompanyName:  req.CompanyName,
		ContactName:  req.ContactName,
		ContactTitle: req.ContactTitle,
		City:         req.City,
		Country:      req.Country,
	}
	if err := s.repo.Create(ctx, c); err != nil { return nil, err }
	return mapCustomerToDTO(c), nil
}

func (s *customerSvc) Get(ctx context.Context, id uint) (*dto.CustomerResponse, error) {
	c, err := s.repo.GetByID(ctx, id)
	if err != nil { return nil, err }
	return mapCustomerToDTO(c), nil
}

func (s *customerSvc) List(ctx context.Context) ([]dto.CustomerResponse, error) {
	customers, err := s.repo.GetAll(ctx)
	if err != nil { return nil, err }
	var res []dto.CustomerResponse
	for _, c := range customers { res = append(res, *mapCustomerToDTO(&c)) }
	return res, nil
}

func (s *customerSvc) Update(ctx context.Context, id uint, req dto.CustomerUpdateReq) (*dto.CustomerResponse, error) {
	c, err := s.repo.GetByID(ctx, id)
	if err != nil { return nil, err }
	
	c.CustomerID = req.CustomerID
	c.CompanyName = req.CompanyName
	c.ContactName = req.ContactName
	c.ContactTitle = req.ContactTitle
	c.City = req.City
	c.Country = req.Country
	
	if err := s.repo.Update(ctx, c); err != nil { return nil, err }
	return mapCustomerToDTO(c), nil
}

func (s *customerSvc) Delete(ctx context.Context, id uint) error {
	return s.repo.Delete(ctx, id)
}

func mapCustomerToDTO(c *model.Customer) *dto.CustomerResponse {
	return &dto.CustomerResponse{
		ID:           c.ID,
		CustomerID:   c.CustomerID,
		CompanyName:  c.CompanyName,
		ContactName:  c.ContactName,
		ContactTitle: c.ContactTitle,
		City:         c.City,
		Country:      c.Country,
		Active:       c.Active,
		// ModifiedAt:   c.UpdatedAt.String(),
	}
}