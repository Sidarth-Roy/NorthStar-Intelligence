package repository

import (
	"context"
	"github.com/Sidarth-Roy/NorthStar-Intelligence/Backend/pkg/model"
	"gorm.io/gorm"
)

type ProductRepository interface {
	Create(ctx context.Context, p *model.Product) error
	GetByID(ctx context.Context, id uint) (*model.Product, error)
	GetAll(ctx context.Context) ([]model.Product, error)
	Update(ctx context.Context, p *model.Product) error
	Delete(ctx context.Context, id uint) error
}

type productRepo struct{ db *gorm.DB }

func NewProductRepo(db *gorm.DB) ProductRepository { return &productRepo{db} }

// func (r *productRepo) Create(ctx context.Context, p *model.Product) error {
// 	return r.db.WithContext(ctx).Create(p).Error
// }

func (r *productRepo) Create(ctx context.Context, p *model.Product) error {
	err := r.db.WithContext(ctx).Create(p).Error
	if err != nil {
		return err
	}
	// Re-fetch with Preload to get CategoryName
	return r.db.WithContext(ctx).Preload("Category").First(p, p.ID).Error
}

func (r *productRepo) GetByID(ctx context.Context, id uint) (*model.Product, error) {
	var p model.Product
	err := r.db.WithContext(ctx).Preload("Category").First(&p, id).Error
	return &p, err
}

func (r *productRepo) GetAll(ctx context.Context) ([]model.Product, error) {
	var products []model.Product
	err := r.db.WithContext(ctx).Preload("Category").Find(&products).Error
	return products, err
}

// func (r *productRepo) Update(ctx context.Context, p *model.Product) error {
// 	return r.db.WithContext(ctx).Save(p).Error
// }

func (r *productRepo) Update(ctx context.Context, p *model.Product) error {
	err := r.db.WithContext(ctx).Save(p).Error
	if err != nil {
		return err
	}
	// Re-fetch with Preload to get CategoryName
	return r.db.WithContext(ctx).Preload("Category").First(p, p.ID).Error
}

func (r *productRepo) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&model.Product{}, id).Error
}