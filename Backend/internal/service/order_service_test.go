package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/Sidarth-Roy/NorthStar-Intelligence/Backend/internal/dto"
	"github.com/Sidarth-Roy/NorthStar-Intelligence/Backend/pkg/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockOrderRepo implementation
type MockOrderRepo struct {
	mock.Mock
}

func (m *MockOrderRepo) Create(ctx context.Context, o *model.Order) error {
	args := m.Called(ctx, o)
	o.ID = 100
	return args.Error(0)
}

func (m *MockOrderRepo) GetByID(ctx context.Context, id uint) (*model.Order, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Order), args.Error(1)
}

func (m *MockOrderRepo) GetAll(ctx context.Context) ([]model.Order, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]model.Order), args.Error(1)
}

func (m *MockOrderRepo) Update(ctx context.Context, o *model.Order) error {
	args := m.Called(ctx, o)
	return args.Error(0)
}

func (m *MockOrderRepo) Delete(ctx context.Context, id uint) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

// --- TESTS ---

func TestOrderService_Create(t *testing.T) {
	ctx := context.TODO()
	now := time.Now()

	t.Run("Success Path with Details", func(t *testing.T) {
		mockRepo := new(MockOrderRepo)
		svc := NewOrderSvc(mockRepo)
		
		// Adding OrderDetails ensures the mapping loops are covered
		req := dto.OrderUpsertReq{
			CustomerID: "ALFKI", 
			OrderDate:  now,
			OrderDetails: []dto.OrderDetailReq{
				{ProductID: 1, UnitPrice: 10.5, Quantity: 2},
			},
		}
		mockRepo.On("Create", ctx, mock.AnythingOfType("*model.Order")).Return(nil)

		res, err := svc.Create(ctx, req)

		assert.NoError(t, err)
		assert.NotNil(t, res)
		assert.Len(t, res.OrderDetails, 1) // Verify mapping worked
		mockRepo.AssertExpectations(t)
	})

	t.Run("Repo Error", func(t *testing.T) {
		mockRepo := new(MockOrderRepo)
		svc := NewOrderSvc(mockRepo)
		mockRepo.On("Create", ctx, mock.Anything).Return(errors.New("creation failed"))

		res, err := svc.Create(ctx, dto.OrderUpsertReq{})

		assert.Error(t, err)
		assert.Nil(t, res)
	})
}

func TestOrderService_Get(t *testing.T) {
	ctx := context.TODO()

	t.Run("Success with Details", func(t *testing.T) {
		mockRepo := new(MockOrderRepo)
		svc := NewOrderSvc(mockRepo)
		
		// Mock returning an order that HAS details
		orderWithDetails := &model.Order{
			Base: model.Base{ID: 1},
			OrderDetails: []model.OrderDetail{
				{ProductID: 5},
			},
		}
		mockRepo.On("GetByID", ctx, uint(1)).Return(orderWithDetails, nil)

		res, err := svc.Get(ctx, 1)
		assert.NoError(t, err)
		assert.Len(t, res.OrderDetails, 1)
	})

	t.Run("Not Found Error", func(t *testing.T) {
		mockRepo := new(MockOrderRepo)
		svc := NewOrderSvc(mockRepo)
		mockRepo.On("GetByID", ctx, uint(404)).Return(nil, errors.New("not found"))

		res, err := svc.Get(ctx, 404)
		assert.Error(t, err)
		assert.Nil(t, res)
	})
}

func TestOrderService_List(t *testing.T) {
	ctx := context.TODO()

	t.Run("Success", func(t *testing.T) {
		mockRepo := new(MockOrderRepo)
		svc := NewOrderSvc(mockRepo)
		mockRepo.On("GetAll", ctx).Return([]model.Order{{CustomerID: "A", OrderDetails: []model.OrderDetail{{ProductID: 1}}}}, nil)

		res, err := svc.List(ctx)
		assert.NoError(t, err)
		assert.Len(t, res, 1)
		assert.Len(t, res[0].OrderDetails, 1)
	})

	t.Run("Repo Error", func(t *testing.T) {
		mockRepo := new(MockOrderRepo)
		svc := NewOrderSvc(mockRepo)
		mockRepo.On("GetAll", ctx).Return(nil, errors.New("db error"))

		res, err := svc.List(ctx)
		assert.Error(t, err)
		assert.Nil(t, res)
	})
}

func TestOrderService_Update(t *testing.T) {
	ctx := context.TODO()

	t.Run("Success Update with Details", func(t *testing.T) {
		mockRepo := new(MockOrderRepo)
		svc := NewOrderSvc(mockRepo)
		existing := &model.Order{Base: model.Base{ID: 1}}
		
		mockRepo.On("GetByID", ctx, uint(1)).Return(existing, nil)
		mockRepo.On("Update", ctx, mock.Anything).Return(nil)

		req := dto.OrderUpsertReq{
			CustomerID: "UPDATED",
			OrderDetails: []dto.OrderDetailReq{
				{ProductID: 99},
			},
		}
		res, err := svc.Update(ctx, 1, req)
		assert.NoError(t, err)
		assert.Equal(t, "UPDATED", res.CustomerID)
		assert.Len(t, res.OrderDetails, 1)
	})

	t.Run("Update GetByID Fail", func(t *testing.T) {
		mockRepo := new(MockOrderRepo)
		svc := NewOrderSvc(mockRepo)
		mockRepo.On("GetByID", ctx, uint(1)).Return(nil, errors.New("search error"))

		res, err := svc.Update(ctx, 1, dto.OrderUpsertReq{})
		assert.Error(t, err)
		assert.Nil(t, res)
	})

	t.Run("Update Save Fail", func(t *testing.T) {
		mockRepo := new(MockOrderRepo)
		svc := NewOrderSvc(mockRepo)
		mockRepo.On("GetByID", ctx, uint(1)).Return(&model.Order{}, nil)
		mockRepo.On("Update", ctx, mock.Anything).Return(errors.New("save error"))

		res, err := svc.Update(ctx, 1, dto.OrderUpsertReq{})
		assert.Error(t, err)
		assert.Nil(t, res)
	})
}

func TestOrderService_Delete(t *testing.T) {
	ctx := context.TODO()

	t.Run("Success", func(t *testing.T) {
		mockRepo := new(MockOrderRepo)
		svc := NewOrderSvc(mockRepo)
		mockRepo.On("Delete", ctx, uint(1)).Return(nil)

		err := svc.Delete(ctx, 1)
		assert.NoError(t, err)
	})

	t.Run("Fail", func(t *testing.T) {
		mockRepo := new(MockOrderRepo)
		svc := NewOrderSvc(mockRepo)
		mockRepo.On("Delete", ctx, uint(1)).Return(errors.New("delete failed"))

		err := svc.Delete(ctx, 1)
		assert.Error(t, err)
	})
}