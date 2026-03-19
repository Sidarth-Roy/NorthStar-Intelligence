package repository

import (
	"github.com/Sidarth-Roy/NorthStar-Intelligence/Backend/pkg/model"
	"gorm.io/gorm"
)

type ProductRepository interface {
	Create(p *model.Product) error
	GetByID(id uint) (*model.Product, error)
	GetAll() ([]model.Product, error)
	Update(p *model.Product) error
	Delete(id uint) error
}

type productRepo struct {
	db *gorm.DB
}

func NewProductRepo(db *gorm.DB) ProductRepository {
	return &productRepo{db: db}
}

func (r *productRepo) Create(p *model.Product) error {
	return r.db.Create(p).Error
}

func (r *productRepo) GetByID(id uint) (*model.Product, error) {
	var p model.Product
	if err := r.db.First(&p, id).Error; err != nil {
		return nil, err
	}
	return &p, nil
}

func (r *productRepo) GetAll() ([]model.Product, error) {
	var products []model.Product
	err := r.db.Where("active = ?", true).Find(&products).Error
	return products, err
}

func (r *productRepo) Update(p *model.Product) error {
	return r.db.Save(p).Error
}

func (r *productRepo) Delete(id uint) error {
	// GORM handles soft delete automatically because of gorm.DeletedAt in Base
	return r.db.Delete(&model.Product{}, id).Error
}