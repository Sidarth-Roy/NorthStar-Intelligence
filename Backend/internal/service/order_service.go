package service

import (
	"context"
	"time"

	"github.com/Sidarth-Roy/NorthStar-Intelligence/Backend/internal/dto"
	"github.com/Sidarth-Roy/NorthStar-Intelligence/Backend/internal/repository"
	"github.com/Sidarth-Roy/NorthStar-Intelligence/Backend/pkg/model"
	"github.com/Sidarth-Roy/NorthStar-Intelligence/Backend/pkg/logger"
	"go.uber.org/zap"
)

type OrderService interface {
	Create(ctx context.Context, req dto.OrderInsertReq) (*dto.OrderResponse, error)
	Get(ctx context.Context, id uint) (*dto.OrderResponse, error)
	List(ctx context.Context) ([]dto.OrderResponse, error)
	Update(ctx context.Context, id uint, req dto.OrderUpdateReq) (*dto.OrderResponse, error)
	Delete(ctx context.Context, id uint) error
	DeleteOrderDetail(ctx context.Context, id uint) error
}

type orderSvc struct{ repo repository.OrderRepository }

func NewOrderSvc(r repository.OrderRepository) OrderService { return &orderSvc{repo: r} }

func (s *orderSvc) Create(ctx context.Context, req dto.OrderInsertReq) (*dto.OrderResponse, error) {
	order := mapDTOToOrderModel(req)
	
	if err := s.repo.Create(ctx, order); err != nil { 
		return nil, err 
	}

	fullOrder, err := s.repo.GetByID(ctx, order.ID)
	if err != nil {
		return nil, err
	}
	return mapOrderToDTO(fullOrder), nil
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
	for i := range orders { 
        // Use the index to get the actual pointer to the slice element
        res = append(res, *mapOrderToDTO(&orders[i])) 
    }
	return res, nil
}

func (s *orderSvc) Update(ctx context.Context, id uint, req dto.OrderUpdateReq) (*dto.OrderResponse, error) {
	existingOrder, err := s.repo.GetByID(ctx, id)
	if err != nil { return nil, err }
	
	// Update header fields
	existingOrder.CustomerID = req.CustomerID
	existingOrder.EmployeeID = req.EmployeeID
	existingOrder.OrderDate = parseDate(req.OrderDate)
	existingOrder.RequiredDate = parseDate(req.RequiredDate)
	existingOrder.ShippedDate = parseOptionalDate(req.ShippedDate)
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

	// 🔄 NEW: Reload the order to get the Names for the newly attached OrderDetails
	updatedOrder, err := s.repo.GetByID(ctx, existingOrder.ID)
	if err != nil {
		return nil, err
	}
	return mapOrderToDTO(updatedOrder), nil
}

func (s *orderSvc) Delete(ctx context.Context, id uint) error {
	return s.repo.Delete(ctx, id)
}

func (s *orderSvc) DeleteOrderDetail(ctx context.Context, id uint) error {
	return s.repo.DeleteOrderDetail(ctx, id)
}


// --- HELPERS ---

// parseDate converts string "YYYY-MM-DD" to time.Time
func parseDate(dateStr string) time.Time {
	t, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		logger.Get().Error("error parsing date", zap.Error(err), zap.String("input", dateStr))		
		return time.Time{} // Return zero time if invalid
	}
	return t			
}

// parseOptionalDate handles nullable strings for *time.Time
func parseOptionalDate(dateStr string) *time.Time {
	if dateStr == "" || dateStr == "NULL" {
		return nil
	}
	t, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		return nil
	}
	return &t
}

// Helpers
func mapDTOToOrderModel(req dto.OrderInsertReq) *model.Order {
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
		Base: 		  model.Base{ID: req.OrderID},
		CustomerID:   req.CustomerID,
		EmployeeID:   req.EmployeeID,
		OrderDate:    parseDate(req.OrderDate),
		RequiredDate: parseDate(req.RequiredDate),
		ShippedDate:  parseOptionalDate(req.ShippedDate),
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
			ProductName: d.Product.ProductName,
			UnitPrice: d.UnitPrice,
			Quantity:  d.Quantity,
			Discount:  d.Discount,
		})
	}

	return &dto.OrderResponse{
		ID:           o.ID,
		CustomerID:   o.CustomerID,
		CustomerName: o.Customer.CompanyName,
		EmployeeID:   o.EmployeeID,
		EmployeeName: o.Employee.EmployeeName,
		OrderDate:    o.OrderDate.Format("2006-01-02"),
		RequiredDate: o.RequiredDate.Format("2006-01-02"),
		ShippedDate:  formatOptionalDate(o.ShippedDate),
		ShipperID:    o.ShipperID,
		ShipperName:  o.Shipper.CompanyName,
		Freight:      o.Freight,
		Active:       o.Active,
		// ModifiedAt:   o.UpdatedAt.String(),
		OrderDetails: detailResponses,
	}
}

// Helper function to handle *time.Time pointers safely
func formatOptionalDate(t *time.Time) string {
    if t == nil {
        return "" // or "N/A"
    }
    return t.Format("2006-01-02")
}