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

// MockCustomerRepo implementation
type MockCustomerRepo struct {
	mock.Mock
}

func (m *MockCustomerRepo) Create(ctx context.Context, c *model.Customer) error {
	args := m.Called(ctx, c)
	c.ID = 1 
	return args.Error(0)
}

func (m *MockCustomerRepo) GetByID(ctx context.Context, id uint) (*model.Customer, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil { return nil, args.Error(1) }
	return args.Get(0).(*model.Customer), args.Error(1)
}

func (m *MockCustomerRepo) GetAll(ctx context.Context) ([]model.Customer, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil { return nil, args.Error(1) }
	return args.Get(0).([]model.Customer), args.Error(1)
}

func (m *MockCustomerRepo) Update(ctx context.Context, c *model.Customer) error {
	args := m.Called(ctx, c)
	return args.Error(0)
}

func (m *MockCustomerRepo) Delete(ctx context.Context, id uint) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

// --- TESTS ---

func TestCustomerService_Create(t *testing.T) {
	ctx := context.TODO()
	
	t.Run("Success Path", func(t *testing.T) {
		mockRepo := new(MockCustomerRepo)
		svc := NewCustomerSvc(mockRepo)
		req := dto.CustomerUpsertReq{CustomerID: "ALFKI", CompanyName: "Alfreds"}

		mockRepo.On("Create", ctx, mock.AnythingOfType("*model.Customer")).Return(nil)
		res, err := svc.Create(ctx, req)

		assert.NoError(t, err)
		assert.NotNil(t, res)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Create Repo Fail", func(t *testing.T) {
		mockRepo := new(MockCustomerRepo)
		svc := NewCustomerSvc(mockRepo)
		
		mockRepo.On("Create", ctx, mock.AnythingOfType("*model.Customer")).Return(errors.New("db error"))
		res, err := svc.Create(ctx, dto.CustomerUpsertReq{CustomerID: "FAIL"})

		assert.Error(t, err)
		assert.Nil(t, res)
	})
}

func TestCustomerService_Get(t *testing.T) {
	ctx := context.TODO()
	mockRepo := new(MockCustomerRepo)
	svc := NewCustomerSvc(mockRepo)

	t.Run("Found", func(t *testing.T) {
		mockRepo.On("GetByID", ctx, uint(1)).Return(&model.Customer{CustomerID: "BOTTM"}, nil)
		res, err := svc.Get(ctx, 1)
		assert.NoError(t, err)
		assert.Equal(t, "BOTTM", res.CustomerID)
	})

	t.Run("Not Found", func(t *testing.T) {
		mockRepo.On("GetByID", ctx, uint(404)).Return(nil, errors.New("not found"))
		res, err := svc.Get(ctx, 404)
		assert.Error(t, err)
		assert.Nil(t, res)
	})
}

func TestCustomerService_List(t *testing.T) {
	ctx := context.TODO()
	
	t.Run("Success", func(t *testing.T) {
		mockRepo := new(MockCustomerRepo)
		svc := NewCustomerSvc(mockRepo)
		mockRepo.On("GetAll", ctx).Return([]model.Customer{{CustomerID: "A"}}, nil)
		
		res, err := svc.List(ctx)
		assert.NoError(t, err)
		assert.Len(t, res, 1)
	})

	t.Run("List Repo Fail", func(t *testing.T) {
		mockRepo := new(MockCustomerRepo)
		svc := NewCustomerSvc(mockRepo)
		mockRepo.On("GetAll", ctx).Return(nil, errors.New("fetch error"))
		
		res, err := svc.List(ctx)
		assert.Error(t, err)
		assert.Nil(t, res)
	})
}

func TestCustomerService_Update(t *testing.T) {
	ctx := context.TODO()

	t.Run("Success Update", func(t *testing.T) {
		mockRepo := new(MockCustomerRepo)
		svc := NewCustomerSvc(mockRepo)
		existing := &model.Customer{Base: model.Base{ID: 1}, CustomerID: "OLD"}
		
		mockRepo.On("GetByID", ctx, uint(1)).Return(existing, nil)
		mockRepo.On("Update", ctx, mock.AnythingOfType("*model.Customer")).Return(nil)

		res, err := svc.Update(ctx, 1, dto.CustomerUpsertReq{CustomerID: "NEW"})
		assert.NoError(t, err)
		assert.Equal(t, "NEW", res.CustomerID)
	})

	t.Run("Update GetByID Fail", func(t *testing.T) {
		mockRepo := new(MockCustomerRepo)
		svc := NewCustomerSvc(mockRepo)
		mockRepo.On("GetByID", ctx, uint(2)).Return(nil, errors.New("not found"))

		res, err := svc.Update(ctx, 2, dto.CustomerUpsertReq{})
		assert.Error(t, err)
		assert.Nil(t, res)
	})

	t.Run("Update Repo Save Fail", func(t *testing.T) {
		mockRepo := new(MockCustomerRepo)
		svc := NewCustomerSvc(mockRepo)
		existing := &model.Customer{Base: model.Base{ID: 3}}

		mockRepo.On("GetByID", ctx, uint(3)).Return(existing, nil)
		mockRepo.On("Update", ctx, mock.AnythingOfType("*model.Customer")).Return(errors.New("save error"))

		res, err := svc.Update(ctx, 3, dto.CustomerUpsertReq{})
		assert.Error(t, err)
		assert.Nil(t, res)
	})
}

func TestCustomerService_Delete(t *testing.T) {
	ctx := context.TODO()
	
	t.Run("Success Delete", func(t *testing.T) {
		mockRepo := new(MockCustomerRepo)
		svc := NewCustomerSvc(mockRepo)
		mockRepo.On("Delete", ctx, uint(1)).Return(nil)
		err := svc.Delete(ctx, 1)
		assert.NoError(t, err)
	})

	t.Run("Delete Fail", func(t *testing.T) {
		mockRepo := new(MockCustomerRepo)
		svc := NewCustomerSvc(mockRepo)
		mockRepo.On("Delete", ctx, uint(1)).Return(errors.New("db error"))
		err := svc.Delete(ctx, 1)
		assert.Error(t, err)
	})
}