package service

import (
	"context"
	"github.com/Sidarth-Roy/NorthStar-Intelligence/Backend/internal/dto"
	"github.com/Sidarth-Roy/NorthStar-Intelligence/Backend/internal/repository"
	"github.com/Sidarth-Roy/NorthStar-Intelligence/Backend/pkg/model"
)

type OrderService interface {
	Create(ctx context.Context, req dto.OrderUpsertReq) (*dto.OrderResponse, error)
	Get(ctx context.Context, id uint) (*dto.OrderResponse, error)
	List(ctx context.Context) ([]dto.OrderResponse, error)
	Update(ctx context.Context, id uint, req dto.OrderUpsertReq) (*dto.OrderResponse, error)
	Delete(ctx context.Context, id uint) error
}

type orderSvc struct{ repo repository.OrderRepository }

func NewOrderSvc(r repository.OrderRepository) OrderService { return &orderSvc{repo: r} }

func (s *orderSvc) Create(ctx context.Context, req dto.OrderUpsertReq) (*dto.OrderResponse, error) {
	order := mapDTOToOrderModel(req)
	
	if err := s.repo.Create(ctx, order); err != nil { 
		return nil, err 
	}
	return mapOrderToDTO(order), nil
}

func (s *orderSvc) Get(ctx context.Context, id uint) (*dto.OrderResponse, error) {
	o, err := s.repo.GetByID(ctx, id)
	if err != nil { return nil, err }
	return mapOrderToDTO(o), nil
}

func (s *orderSvc) List(ctx context.Context) ([]dto.OrderResponse, error) {
	orders, err := s.repo.GetAll(ctx)
	if err != nil { return nil, err }
	var res []dto.OrderResponse
	for _, o := range orders { res = append(res, *mapOrderToDTO(&o)) }
	return res, nil
}

func (s *orderSvc) Update(ctx context.Context, id uint, req dto.OrderUpsertReq) (*dto.OrderResponse, error) {
	existingOrder, err := s.repo.GetByID(ctx, id)
	if err != nil { return nil, err }
	
	// Update header fields
	existingOrder.CustomerID = req.CustomerID
	existingOrder.EmployeeID = req.EmployeeID
	existingOrder.OrderDate = req.OrderDate
	existingOrder.RequiredDate = req.RequiredDate
	existingOrder.ShippedDate = req.ShippedDate
	existingOrder.ShipperID = req.ShipperID
	existingOrder.Freight = req.Freight

	// Replace existing details
	var updatedDetails []model.OrderDetail
	for _, detail := range req.OrderDetails {
		updatedDetails = append(updatedDetails, model.OrderDetail{
			ProductID: detail.ProductID,
			UnitPrice: detail.UnitPrice,
			Quantity:  detail.Quantity,
			Discount:  detail.Discount,
		})
	}
	existingOrder.OrderDetails = updatedDetails

	if err := s.repo.Update(ctx, existingOrder); err != nil { return nil, err }
	return mapOrderToDTO(existingOrder), nil
}

func (s *orderSvc) Delete(ctx context.Context, id uint) error {
	return s.repo.Delete(ctx, id)
}

// Helpers
func mapDTOToOrderModel(req dto.OrderUpsertReq) *model.Order {
	var details []model.OrderDetail
	for _, d := range req.OrderDetails {
		details = append(details, model.OrderDetail{
			ProductID: d.ProductID,
			UnitPrice: d.UnitPrice,
			Quantity:  d.Quantity,
			Discount:  d.Discount,
		})
	}

	return &model.Order{
		CustomerID:   req.CustomerID,
		EmployeeID:   req.EmployeeID,
		OrderDate:    req.OrderDate,
		RequiredDate: req.RequiredDate,
		ShippedDate:  req.ShippedDate,
		ShipperID:    req.ShipperID,
		Freight:      req.Freight,
		OrderDetails: details,
	}
}

func mapOrderToDTO(o *model.Order) *dto.OrderResponse {
	var detailResponses []dto.OrderDetailResponse
	for _, d := range o.OrderDetails {
		detailResponses = append(detailResponses, dto.OrderDetailResponse{
			ID:        d.ID,
			ProductID: d.ProductID,
			UnitPrice: d.UnitPrice,
			Quantity:  d.Quantity,
			Discount:  d.Discount,
		})
	}

	return &dto.OrderResponse{
		ID:           o.ID,
		CustomerID:   o.CustomerID,
		EmployeeID:   o.EmployeeID,
		OrderDate:    o.OrderDate,
		RequiredDate: o.RequiredDate,
		ShippedDate:  o.ShippedDate,
		ShipperID:    o.ShipperID,
		Freight:      o.Freight,
		Active:       o.Active,
		ModifiedAt:   o.UpdatedAt.String(),
		OrderDetails: detailResponses,
	}
}