package service

import (
	"context"
	"errors"
	"testing"

	"github.com/Sidarth-Roy/NorthStar-Intelligence/Backend/internal/dto"
	"github.com/Sidarth-Roy/NorthStar-Intelligence/Backend/pkg/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// --- Mock Definition ---

type MockProductRepo struct {
	mock.Mock
}

func (m *MockProductRepo) Create(ctx context.Context, p *model.Product) error {
	args := m.Called(ctx, p)
	p.ID = 1 
	return args.Error(0)
}

func (m *MockProductRepo) GetByID(ctx context.Context, id uint) (*model.Product, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Product), args.Error(1)
}

func (m *MockProductRepo) GetAll(ctx context.Context) ([]model.Product, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]model.Product), args.Error(1)
}

func (m *MockProductRepo) Update(ctx context.Context, p *model.Product) error {
	args := m.Called(ctx, p)
	return args.Error(0)
}

func (m *MockProductRepo) Delete(ctx context.Context, id uint) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

// --- Service Tests ---

func TestProductService_Create(t *testing.T) {
	ctx := context.TODO()

	t.Run("Create Success", func(t *testing.T) {
		mockRepo := new(MockProductRepo)
		svc := NewProductSvc(mockRepo)
		req := dto.ProductUpsertReq{ProductName: "Laptop", UnitPrice: 1200.50}

		mockRepo.On("Create", ctx, mock.Anything).Return(nil)
		res, err := svc.Create(ctx, req)

		assert.NoError(t, err)
		assert.NotNil(t, res)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Create Repository Error", func(t *testing.T) {
		mockRepo := new(MockProductRepo)
		svc := NewProductSvc(mockRepo)
		mockRepo.On("Create", ctx, mock.Anything).Return(errors.New("db error"))

		res, err := svc.Create(ctx, dto.ProductUpsertReq{})

		assert.Error(t, err)
		assert.Nil(t, res)
	})
}

func TestProductService_Get(t *testing.T) {
	ctx := context.TODO()
	mockRepo := new(MockProductRepo)
	svc := NewProductSvc(mockRepo)

	t.Run("Get Found", func(t *testing.T) {
		mockRepo.On("GetByID", ctx, uint(5)).Return(&model.Product{ProductName: "Mouse"}, nil)
		res, err := svc.Get(ctx, 5)
		assert.NoError(t, err)
		assert.Equal(t, "Mouse", res.ProductName)
	})

	t.Run("Get Not Found", func(t *testing.T) {
		mockRepo.On("GetByID", ctx, uint(404)).Return(nil, errors.New("not found"))
		res, err := svc.Get(ctx, 404)
		assert.Error(t, err)
		assert.Nil(t, res)
	})
}

func TestProductService_List(t *testing.T) {
	ctx := context.TODO()

	t.Run("List Success", func(t *testing.T) {
		mockRepo := new(MockProductRepo)
		svc := NewProductSvc(mockRepo)
		mockRepo.On("GetAll", ctx).Return([]model.Product{{ProductName: "P1"}}, nil)

		res, err := svc.List(ctx)
		assert.NoError(t, err)
		assert.Len(t, res, 1)
	})

	t.Run("List Repository Error", func(t *testing.T) {
		mockRepo := new(MockProductRepo)
		svc := NewProductSvc(mockRepo)
		mockRepo.On("GetAll", ctx).Return(nil, errors.New("list error"))

		res, err := svc.List(ctx)
		assert.Error(t, err)
		assert.Nil(t, res)
	})
}

func TestProductService_Update(t *testing.T) {
	ctx := context.TODO()

	t.Run("Update Success", func(t *testing.T) {
		mockRepo := new(MockProductRepo)
		svc := NewProductSvc(mockRepo)
		existing := &model.Product{ProductName: "Old"}
		existing.ID = 1

		mockRepo.On("GetByID", ctx, uint(1)).Return(existing, nil)
		mockRepo.On("Update", ctx, mock.Anything).Return(nil)

		res, err := svc.Update(ctx, 1, dto.ProductUpsertReq{ProductName: "New"})
		assert.NoError(t, err)
		assert.Equal(t, "New", res.ProductName)
	})

	t.Run("Update GetByID Error", func(t *testing.T) {
		mockRepo := new(MockProductRepo)
		svc := NewProductSvc(mockRepo)
		mockRepo.On("GetByID", ctx, uint(1)).Return(nil, errors.New("not found"))

		res, err := svc.Update(ctx, 1, dto.ProductUpsertReq{})
		assert.Error(t, err)
		assert.Nil(t, res)
	})

	t.Run("Update Repository Save Error", func(t *testing.T) {
		mockRepo := new(MockProductRepo)
		svc := NewProductSvc(mockRepo)
		existing := &model.Product{ProductName: "Old"}

		mockRepo.On("GetByID", ctx, uint(1)).Return(existing, nil)
		mockRepo.On("Update", ctx, mock.Anything).Return(errors.New("save error"))

		res, err := svc.Update(ctx, 1, dto.ProductUpsertReq{ProductName: "New"})
		assert.Error(t, err)
		assert.Nil(t, res)
	})
}

func TestProductService_Delete(t *testing.T) {
	ctx := context.TODO()
	mockRepo := new(MockProductRepo)
	svc := NewProductSvc(mockRepo)

	t.Run("Delete Success", func(t *testing.T) {
		mockRepo.On("Delete", ctx, uint(1)).Return(nil)
		err := svc.Delete(ctx, 1)
		assert.NoError(t, err)
	})

	t.Run("Delete Failure", func(t *testing.T) {
		mockRepo.On("Delete", ctx, uint(99)).Return(errors.New("delete failed"))
		err := svc.Delete(ctx, 99)
		assert.Error(t, err)
	})
}