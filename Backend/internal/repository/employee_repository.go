package repository

import (
	"context"
	"github.com/Sidarth-Roy/NorthStar-Intelligence/Backend/pkg/model"
	"gorm.io/gorm"
)

type EmployeeRepository interface {
	Create(ctx context.Context, e *model.Employee) error
	GetByID(ctx context.Context, id uint) (*model.Employee, error)
	GetAll(ctx context.Context) ([]model.Employee, error)
	Update(ctx context.Context, e *model.Employee) error
	Delete(ctx context.Context, id uint) error
}

type employeeRepo struct{ db *gorm.DB }

func NewEmployeeRepo(db *gorm.DB) EmployeeRepository { return &employeeRepo{db: db} }

func (r *employeeRepo) Create(ctx context.Context, e *model.Employee) error {
	return r.db.WithContext(ctx).Create(e).Error
}

func (r *employeeRepo) GetByID(ctx context.Context, id uint) (*model.Employee, error) {
	var e model.Employee
	err := r.db.WithContext(ctx).
				Preload("Manager").
				Preload("Orders").
				Preload("Orders.Customer").
				Preload("Orders.Shipper").
				First(&e, id).Error
	return &e, err
}

func (r *employeeRepo) GetAll(ctx context.Context) ([]model.Employee, error) {
	var employees []model.Employee
	err := r.db.WithContext(ctx).
				Preload("Manager").
				Preload("Orders").
				Preload("Orders.Customer").
				Preload("Orders.Shipper").
				Find(&employees).Error
	return employees, err
}

func (r *employeeRepo) Update(ctx context.Context, e *model.Employee) error {
    return r.db.WithContext(ctx).Model(e).Select("*").Updates(e).Error
}

func (r *employeeRepo) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&model.Employee{}, id).Error
}