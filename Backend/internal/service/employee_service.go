package service

import (
	"context"
	"github.com/Sidarth-Roy/NorthStar-Intelligence/Backend/internal/dto"
	"github.com/Sidarth-Roy/NorthStar-Intelligence/Backend/internal/repository"
	"github.com/Sidarth-Roy/NorthStar-Intelligence/Backend/pkg/model"
)

type EmployeeService interface {
	Create(ctx context.Context, req dto.EmployeeInsertReq) (*dto.EmployeeResponse, error)
	Get(ctx context.Context, id uint) (*dto.EmployeeResponse, error)
	List(ctx context.Context) ([]dto.EmployeeResponse, error)
	Update(ctx context.Context, id uint, req dto.EmployeeUpdateReq) (*dto.EmployeeResponse, error)
	Delete(ctx context.Context, id uint) error
}

type employeeSvc struct{ repo repository.EmployeeRepository }

func NewEmployeeSvc(r repository.EmployeeRepository) EmployeeService { return &employeeSvc{repo: r} }

func (s *employeeSvc) Create(ctx context.Context, req dto.EmployeeInsertReq) (*dto.EmployeeResponse, error) {
	e := &model.Employee{
		EmployeeName: req.EmployeeName,
		Title:        req.Title,
		City:         req.City,
		Country:      req.Country,
		ReportsTo:    req.ReportsTo,
	}
	if err := s.repo.Create(ctx, e); err != nil { return nil, err }
	return mapEmployeeToDTO(e), nil
}

func (s *employeeSvc) Get(ctx context.Context, id uint) (*dto.EmployeeResponse, error) {
	e, err := s.repo.GetByID(ctx, id)
	if err != nil { return nil, err }
	return mapEmployeeToDTO(e), nil
}

func (s *employeeSvc) List(ctx context.Context) ([]dto.EmployeeResponse, error) {
	employees, err := s.repo.GetAll(ctx)
	if err != nil { return nil, err }
	var res []dto.EmployeeResponse
	for _, e := range employees { res = append(res, *mapEmployeeToDTO(&e)) }
	return res, nil
}

func (s *employeeSvc) Update(ctx context.Context, id uint, req dto.EmployeeUpdateReq) (*dto.EmployeeResponse, error) {
	e, err := s.repo.GetByID(ctx, id)
	if err != nil { return nil, err }
	
	e.EmployeeName = req.EmployeeName
	e.Title = req.Title
	e.City = req.City
	e.Country = req.Country
	e.ReportsTo = req.ReportsTo
	
	if err := s.repo.Update(ctx, e); err != nil { return nil, err }
	return mapEmployeeToDTO(e), nil
}

func (s *employeeSvc) Delete(ctx context.Context, id uint) error {
	return s.repo.Delete(ctx, id)
}

func mapEmployeeToDTO(e *model.Employee) *dto.EmployeeResponse {
	mappedOrders := []dto.OrderForEmployeeNestedResponse{} // Initialize as empty slice to avoid null in JSON
	
	for _, o := range e.Orders {
		orderDTO := dto.OrderForEmployeeNestedResponse{
			ID:           o.ID,
			CustomerID:   o.CustomerID,
			OrderDate:    o.OrderDate.Format("2006-01-02"),
			RequiredDate: o.RequiredDate.Format("2006-01-02"),
			ShippedDate:  formatOptionalDate(o.ShippedDate),
			ShipperID:    o.ShipperID,
			Freight:      o.Freight,
			Active:       o.Active,
		}

		// Check the ID or the unique string key instead of the whole struct
		if o.Customer.CustomerID != "" {
			orderDTO.CompanyName = o.Customer.CompanyName
		}

		// Assuming Shipper has a standard 'ID' field
		if o.Shipper.ID != 0 {
			orderDTO.ShipperName = o.Shipper.CompanyName
		}

		mappedOrders = append(mappedOrders, orderDTO)
	}

	res := &dto.EmployeeResponse{
		ID:            e.ID,
		EmployeeName:  e.EmployeeName,
		Title:         e.Title,
		City:          e.City,
		Country:       e.Country,
		ReportsTo:     e.ReportsTo,
		Active:        e.Active,
		Orders:        mappedOrders,
	}

	// FIX: Check if Manager exists before accessing EmployeeName
	if e.Manager != nil {
		res.ReportsToName = e.Manager.EmployeeName
	} else {
		res.ReportsToName = "" // Or "N/A"
	}

	return res
}
