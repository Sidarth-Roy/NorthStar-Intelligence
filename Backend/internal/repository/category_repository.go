package repository

import (
	"context"
	"github.com/Sidarth-Roy/NorthStar-Intelligence/Backend/pkg/model"
	"gorm.io/gorm"
)

type CategoryRepository interface {
	Create(ctx context.Context, c *model.Category) error
	GetByID(ctx context.Context, id uint) (*model.Category, error)
	GetAll(ctx context.Context) ([]model.Category, error)
	GetProductsByCategoryID(ctx context.Context, categoryID uint) ([]model.Product, error)
	Update(ctx context.Context, c *model.Category) error
	Delete(ctx context.Context, id uint) error
}

type categoryRepo struct{ db *gorm.DB }

func NewCategoryRepo(db *gorm.DB) CategoryRepository { return &categoryRepo{db: db} }

func (r *categoryRepo) Create(ctx context.Context, c *model.Category) error {
	return r.db.WithContext(ctx).Create(c).Error
}

func (r *categoryRepo) GetByID(ctx context.Context, id uint) (*model.Category, error) {
	var c model.Category
	err := r.db.WithContext(ctx).First(&c, id).Error
	return &c, err
}

func (r *categoryRepo) GetAll(ctx context.Context) ([]model.Category, error) {
	var categories []model.Category
	err := r.db.WithContext(ctx).Find(&categories).Error
	return categories, err
}

func (r *categoryRepo) GetProductsByCategoryID(ctx context.Context, categoryID uint) ([]model.Product, error) {
	var products []model.Product
	err := r.db.WithContext(ctx).Where("category_id = ?", categoryID).Find(&products).Error
	return products, err
}

func (r *categoryRepo) Update(ctx context.Context, c *model.Category) error {
	return r.db.WithContext(ctx).Save(c).Error
}

func (r *categoryRepo) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&model.Category{}, id).Error
}