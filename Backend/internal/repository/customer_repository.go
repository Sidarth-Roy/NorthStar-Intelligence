package repository

import (
	"context"
	"github.com/Sidarth-Roy/NorthStar-Intelligence/Backend/pkg/model"
	"gorm.io/gorm"
)

type CustomerRepository interface {
	Create(ctx context.Context, c *model.Customer) error
	GetByID(ctx context.Context, id uint) (*model.Customer, error)
	GetAll(ctx context.Context) ([]model.Customer, error)
	Update(ctx context.Context, c *model.Customer) error
	Delete(ctx context.Context, id uint) error
}

type customerRepo struct{ db *gorm.DB }

func NewCustomerRepo(db *gorm.DB) CustomerRepository { return &customerRepo{db: db} }

func (r *customerRepo) Create(ctx context.Context, c *model.Customer) error {
	return r.db.WithContext(ctx).Create(c).Error
}

func (r *customerRepo) GetByID(ctx context.Context, id uint) (*model.Customer, error) {
	var c model.Customer
	err := r.db.WithContext(ctx).First(&c, id).Error
	return &c, err
}

func (r *customerRepo) GetAll(ctx context.Context) ([]model.Customer, error) {
	var customers []model.Customer
	err := r.db.WithContext(ctx).Find(&customers).Error
	return customers, err
}

func (r *customerRepo) Update(ctx context.Context, c *model.Customer) error {
	return r.db.WithContext(ctx).Save(c).Error
}

func (r *customerRepo) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&model.Customer{}, id).Error
}