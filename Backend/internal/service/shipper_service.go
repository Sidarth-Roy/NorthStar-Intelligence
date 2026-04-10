package service

import (
	"context"

	"github.com/Sidarth-Roy/NorthStar-Intelligence/Backend/internal/dto"
	"github.com/Sidarth-Roy/NorthStar-Intelligence/Backend/internal/repository"
	"github.com/Sidarth-Roy/NorthStar-Intelligence/Backend/pkg/model"
)

type ShipperService interface {
	Create(ctx context.Context, req dto.ShipperInsertReq) (*dto.ShipperResponse, error)
	Get(ctx context.Context, id uint) (*dto.ShipperResponse, error)
	List(ctx context.Context) ([]dto.ShipperResponse, error)
	Update(ctx context.Context, id uint, req dto.ShipperUpdateReq) (*dto.ShipperResponse, error)
	Delete(ctx context.Context, id uint) error
}

type shipperSvc struct{ repo repository.ShipperRepository }

func NewShipperSvc(r repository.ShipperRepository) ShipperService { return &shipperSvc{repo: r} }

func (s *shipperSvc) Create(ctx context.Context, req dto.ShipperInsertReq) (*dto.ShipperResponse, error) {
	ship := &model.Shipper{
		CompanyName: req.CompanyName,
	}
	if err := s.repo.Create(ctx, ship); err != nil { return nil, err }
	return mapShipperToDTO(ship), nil
}

func (s *shipperSvc) Get(ctx context.Context, id uint) (*dto.ShipperResponse, error) {
	ship, err := s.repo.GetByID(ctx, id)
	if err != nil { return nil, err }
	return mapShipperToDTO(ship), nil
}

func (s *shipperSvc) List(ctx context.Context) ([]dto.ShipperResponse, error) {
	shippers, err := s.repo.GetAll(ctx)
	if err != nil { return nil, err }
	var res []dto.ShipperResponse
	for _, ship := range shippers { res = append(res, *mapShipperToDTO(&ship)) }
	return res, nil
}

func (s *shipperSvc) Update(ctx context.Context, id uint, req dto.ShipperUpdateReq) (*dto.ShipperResponse, error) {
	ship, err := s.repo.GetByID(ctx, id)
	if err != nil { return nil, err }
	
	ship.CompanyName = req.CompanyName
	
	if err := s.repo.Update(ctx, ship); err != nil { return nil, err }
	return mapShipperToDTO(ship), nil
}

func (s *shipperSvc) Delete(ctx context.Context, id uint) error {
	return s.repo.Delete(ctx, id)
}

func mapShipperToDTO(s *model.Shipper) *dto.ShipperResponse {
	var mappedOrders []dto.OrderForShipperNestedResponse
	for _, o := range s.Orders {
		mappedOrders = append(mappedOrders, dto.OrderForShipperNestedResponse{
			ID:           o.ID,
			CustomerID:   o.Customer.ID,
			CompanyName:  o.Customer.CompanyName,
			OrderDate:    o.OrderDate.Format("2006-01-02"),
			RequiredDate: o.RequiredDate.Format("2006-01-02"),
			ShippedDate:  formatOptionalDate(o.ShippedDate),
			ShipperID:    o.ShipperID,
			ShipperName:  o.Shipper.CompanyName,
			Freight:      o.Freight,
			Active:       o.Active,
		})
	}
	if mappedOrders == nil {
		mappedOrders = []dto.OrderForShipperNestedResponse{} // Return null instead of empty array if no orders
	}
	return &dto.ShipperResponse{
		ID:          s.ID,
		CompanyName: s.CompanyName,
		Active:      s.Active,
		Orders:      mappedOrders,
		// ModifiedAt:  s.UpdatedAt.String(),
	}
}