package repository

import (
	"context"
	"github.com/Sidarth-Roy/NorthStar-Intelligence/Backend/pkg/model"
	"gorm.io/gorm"
)

type ShipperRepository interface {
	Create(ctx context.Context, s *model.Shipper) error
	GetByID(ctx context.Context, id uint) (*model.Shipper, error)
	GetAll(ctx context.Context) ([]model.Shipper, error)
	Update(ctx context.Context, s *model.Shipper) error
	Delete(ctx context.Context, id uint) error
}

type shipperRepo struct{ db *gorm.DB }

func NewShipperRepo(db *gorm.DB) ShipperRepository { return &shipperRepo{db: db} }

func (r *shipperRepo) Create(ctx context.Context, s *model.Shipper) error {
	return r.db.WithContext(ctx).Create(s).Error
}

func (r *shipperRepo) GetByID(ctx context.Context, id uint) (*model.Shipper, error) {
	var s model.Shipper
	err := r.db.WithContext(ctx).First(&s, id).Error
	return &s, err
}

func (r *shipperRepo) GetAll(ctx context.Context) ([]model.Shipper, error) {
	var shippers []model.Shipper
	err := r.db.WithContext(ctx).Find(&shippers).Error
	return shippers, err
}

func (r *shipperRepo) Update(ctx context.Context, s *model.Shipper) error {
	return r.db.WithContext(ctx).Save(s).Error
}

func (r *shipperRepo) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&model.Shipper{}, id).Error
}