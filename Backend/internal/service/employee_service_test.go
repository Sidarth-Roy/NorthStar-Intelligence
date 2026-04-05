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

// MockEmployeeRepo implementation
type MockEmployeeRepo struct {
	mock.Mock
}

func (m *MockEmployeeRepo) Create(ctx context.Context, e *model.Employee) error {
	args := m.Called(ctx, e)
	e.ID = 1 
	return args.Error(0)
}

func (m *MockEmployeeRepo) GetByID(ctx context.Context, id uint) (*model.Employee, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Employee), args.Error(1)
}

func (m *MockEmployeeRepo) GetAll(ctx context.Context) ([]model.Employee, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]model.Employee), args.Error(1)
}

func (m *MockEmployeeRepo) Update(ctx context.Context, e *model.Employee) error {
	args := m.Called(ctx, e)
	return args.Error(0)
}

func (m *MockEmployeeRepo) Delete(ctx context.Context, id uint) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

// --- TESTS ---

func TestEmployeeService_Create(t *testing.T) {
	ctx := context.TODO()

	t.Run("Success Path", func(t *testing.T) {
		mockRepo := new(MockEmployeeRepo)
		svc := NewEmployeeSvc(mockRepo)
		req := dto.EmployeeUpsertReq{EmployeeName: "John Doe"}

		mockRepo.On("Create", ctx, mock.AnythingOfType("*model.Employee")).Return(nil)

		res, err := svc.Create(ctx, req)

		assert.NoError(t, err)
		assert.NotNil(t, res)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Create Error", func(t *testing.T) {
		mockRepo := new(MockEmployeeRepo)
		svc := NewEmployeeSvc(mockRepo)
		
		mockRepo.On("Create", ctx, mock.AnythingOfType("*model.Employee")).Return(errors.New("db error"))

		res, err := svc.Create(ctx, dto.EmployeeUpsertReq{EmployeeName: "Fail"})

		assert.Error(t, err)
		assert.Nil(t, res)
	})
}

func TestEmployeeService_Get(t *testing.T) {
	ctx := context.TODO()

	t.Run("Found Employee", func(t *testing.T) {
		mockRepo := new(MockEmployeeRepo)
		svc := NewEmployeeSvc(mockRepo)
		mockRepo.On("GetByID", ctx, uint(1)).Return(&model.Employee{EmployeeName: "Jane Smith"}, nil)

		res, err := svc.Get(ctx, 1)

		assert.NoError(t, err)
		assert.Equal(t, "Jane Smith", res.EmployeeName)
	})

	t.Run("Get Error", func(t *testing.T) {
		mockRepo := new(MockEmployeeRepo)
		svc := NewEmployeeSvc(mockRepo)
		mockRepo.On("GetByID", ctx, uint(1)).Return(nil, errors.New("not found"))

		res, err := svc.Get(ctx, 1)

		assert.Error(t, err)
		assert.Nil(t, res)
	})
}

func TestEmployeeService_List(t *testing.T) {
	ctx := context.TODO()

	t.Run("Success List", func(t *testing.T) {
		mockRepo := new(MockEmployeeRepo)
		svc := NewEmployeeSvc(mockRepo)
		mockData := []model.Employee{{EmployeeName: "Emp 1"}}
		mockRepo.On("GetAll", ctx).Return(mockData, nil)

		res, err := svc.List(ctx)

		assert.NoError(t, err)
		assert.Len(t, res, 1)
	})

	t.Run("List Error", func(t *testing.T) {
		mockRepo := new(MockEmployeeRepo)
		svc := NewEmployeeSvc(mockRepo)
		mockRepo.On("GetAll", ctx).Return(nil, errors.New("query failed"))

		res, err := svc.List(ctx)

		assert.Error(t, err)
		assert.Nil(t, res)
	})
}

func TestEmployeeService_Update(t *testing.T) {
	ctx := context.TODO()

	t.Run("Success Update", func(t *testing.T) {
		mockRepo := new(MockEmployeeRepo)
		svc := NewEmployeeSvc(mockRepo)
		existing := &model.Employee{Base: model.Base{ID: 2}, EmployeeName: "Old"}
		req := dto.EmployeeUpsertReq{EmployeeName: "New"}

		mockRepo.On("GetByID", ctx, uint(2)).Return(existing, nil)
		mockRepo.On("Update", ctx, mock.AnythingOfType("*model.Employee")).Return(nil)

		res, err := svc.Update(ctx, 2, req)

		assert.NoError(t, err)
		assert.Equal(t, "New", res.EmployeeName)
	})

	t.Run("Update Fail - Find Phase", func(t *testing.T) {
		mockRepo := new(MockEmployeeRepo)
		svc := NewEmployeeSvc(mockRepo)
		mockRepo.On("GetByID", ctx, uint(404)).Return(nil, errors.New("not found"))

		res, err := svc.Update(ctx, 404, dto.EmployeeUpsertReq{})

		assert.Error(t, err)
		assert.Nil(t, res)
	})

	t.Run("Update Fail - Save Phase", func(t *testing.T) {
		mockRepo := new(MockEmployeeRepo)
		svc := NewEmployeeSvc(mockRepo)
		existing := &model.Employee{Base: model.Base{ID: 2}}
		
		mockRepo.On("GetByID", ctx, uint(2)).Return(existing, nil)
		mockRepo.On("Update", ctx, mock.AnythingOfType("*model.Employee")).Return(errors.New("save error"))

		res, err := svc.Update(ctx, 2, dto.EmployeeUpsertReq{})

		assert.Error(t, err)
		assert.Nil(t, res)
	})
}

func TestEmployeeService_Delete(t *testing.T) {
	ctx := context.TODO()

	t.Run("Success Delete", func(t *testing.T) {
		mockRepo := new(MockEmployeeRepo)
		svc := NewEmployeeSvc(mockRepo)
		mockRepo.On("Delete", ctx, uint(10)).Return(nil)

		err := svc.Delete(ctx, 10)
		assert.NoError(t, err)
	})

	t.Run("Delete Fail", func(t *testing.T) {
		mockRepo := new(MockEmployeeRepo)
		svc := NewEmployeeSvc(mockRepo)
		mockRepo.On("Delete", ctx, uint(10)).Return(errors.New("db crash"))

		err := svc.Delete(ctx, 10)
		assert.Error(t, err)
	})
}