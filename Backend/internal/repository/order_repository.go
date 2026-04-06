package repository

import (
	"context"
	"github.com/Sidarth-Roy/NorthStar-Intelligence/Backend/pkg/model"
	"gorm.io/gorm"
)

type OrderRepository interface {
	Create(ctx context.Context, o *model.Order) error
	GetByID(ctx context.Context, id uint) (*model.Order, error)
	GetAll(ctx context.Context) ([]model.Order, error)
	Update(ctx context.Context, o *model.Order) error
	Delete(ctx context.Context, id uint) error
	DeleteOrderDetail(ctx context.Context, id uint) error
}

type orderRepo struct{ db *gorm.DB }

func NewOrderRepo(db *gorm.DB) OrderRepository { return &orderRepo{db: db} }

func (r *orderRepo) Create(ctx context.Context, o *model.Order) error {
	// GORM automatically runs creates with associations in a transaction
	return r.db.WithContext(ctx).Create(o).Error
}

func (r *orderRepo) GetByID(ctx context.Context, id uint) (*model.Order, error) {
	var o model.Order
	// Preload fetches the associated OrderDetails automatically
	err := r.db.WithContext(ctx).
        Preload("Customer").      // Loads CustomerName
        Preload("Employee").      // Loads EmployeeName
        Preload("Shipper").       // Loads ShipperName
        Preload("OrderDetails").
        Preload("OrderDetails.Product").
		First(&o, id).Error
	return &o, err
}

// func (r *orderRepo) GetAll(ctx context.Context) ([]model.Order, error) {
// 	var orders []model.Order
// 	err := r.db.WithContext(ctx).Preload("OrderDetails").Find(&orders).Error
// 	return orders, err
// }
func (r *orderRepo) GetAll(ctx context.Context) ([]model.Order, error) {
    var orders []model.Order
    err := r.db.WithContext(ctx).
        Preload("Customer").      // Loads CustomerName
        Preload("Employee").      // Loads EmployeeName
        Preload("Shipper").       // Loads ShipperName
        Preload("OrderDetails").
        Preload("OrderDetails.Product"). // Nested Preload for ProductName
        Find(&orders).Error
    return orders, err
}

func (r *orderRepo) Update(ctx context.Context, o *model.Order) error {
	// FullSaveAssociations replaces existing details with the new slice
	return r.db.WithContext(ctx).Session(&gorm.Session{FullSaveAssociations: true}).Save(o).Error
}

func (r *orderRepo) Delete(ctx context.Context, id uint) error {
	// Select ensures we cascade soft-delete to the associated OrderDetails
	return r.db.WithContext(ctx).Select("OrderDetails").Delete(&model.Order{Base: model.Base{ID: id}}).Error
}

func (r *orderRepo) DeleteOrderDetail(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&model.OrderDetail{Base: model.Base{ID: id}}).Error
}