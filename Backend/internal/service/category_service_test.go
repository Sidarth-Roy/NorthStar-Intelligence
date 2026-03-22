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

// MockCategoryRepo implementation
type MockCategoryRepo struct {
	mock.Mock
}

func (m *MockCategoryRepo) Create(ctx context.Context, c *model.Category) error {
	args := m.Called(ctx, c)
	if args.Error(0) == nil {
		c.ID = 1
	}
	return args.Error(0)
}

func (m *MockCategoryRepo) GetByID(ctx context.Context, id uint) (*model.Category, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Category), args.Error(1)
}

func (m *MockCategoryRepo) GetAll(ctx context.Context) ([]model.Category, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]model.Category), args.Error(1)
}

func (m *MockCategoryRepo) Update(ctx context.Context, c *model.Category) error {
	args := m.Called(ctx, c)
	return args.Error(0)
}

func (m *MockCategoryRepo) Delete(ctx context.Context, id uint) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

// --- TESTS ---

func TestCategoryService_Create(t *testing.T) {
	ctx := context.TODO()

	t.Run("Success", func(t *testing.T) {
		mockRepo := new(MockCategoryRepo)
		svc := NewCategorySvc(mockRepo)
		req := dto.CategoryUpsertReq{CategoryName: "Electronics"}
		mockRepo.On("Create", ctx, mock.AnythingOfType("*model.Category")).Return(nil)

		res, err := svc.Create(ctx, req)

		assert.NoError(t, err)
		assert.NotNil(t, res)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Error_RepoFailure", func(t *testing.T) {
		mockRepo := new(MockCategoryRepo)
		svc := NewCategorySvc(mockRepo)
		req := dto.CategoryUpsertReq{CategoryName: "Electronics"}
		mockRepo.On("Create", ctx, mock.AnythingOfType("*model.Category")).Return(errors.New("db error"))

		res, err := svc.Create(ctx, req)

		assert.Error(t, err)
		assert.Nil(t, res)
		mockRepo.AssertExpectations(t)
	})
}

func TestCategoryService_Get(t *testing.T) {
	ctx := context.TODO()

	t.Run("Success", func(t *testing.T) {
		mockRepo := new(MockCategoryRepo)
		svc := NewCategorySvc(mockRepo)
		mockRepo.On("GetByID", ctx, uint(1)).Return(&model.Category{CategoryName: "Test"}, nil)

		res, err := svc.Get(ctx, 1)

		assert.NoError(t, err)
		assert.Equal(t, "Test", res.CategoryName)
	})

	t.Run("Error_NotFound", func(t *testing.T) {
		mockRepo := new(MockCategoryRepo)
		svc := NewCategorySvc(mockRepo)
		mockRepo.On("GetByID", ctx, uint(1)).Return(nil, errors.New("not found"))

		res, err := svc.Get(ctx, 1)

		assert.Error(t, err)
		assert.Nil(t, res)
	})
}

func TestCategoryService_List(t *testing.T) {
	ctx := context.TODO()

	t.Run("Success", func(t *testing.T) {
		mockRepo := new(MockCategoryRepo)
		svc := NewCategorySvc(mockRepo)
		mockData := []model.Category{{CategoryName: "Cat 1"}}
		mockRepo.On("GetAll", ctx).Return(mockData, nil)

		res, err := svc.List(ctx)

		assert.NoError(t, err)
		assert.Len(t, res, 1)
	})

	t.Run("Error_RepoFailure", func(t *testing.T) {
		mockRepo := new(MockCategoryRepo)
		svc := NewCategorySvc(mockRepo)
		mockRepo.On("GetAll", ctx).Return(nil, errors.New("db error"))

		res, err := svc.List(ctx)

		assert.Error(t, err)
		assert.Nil(t, res)
	})
}

func TestCategoryService_Update(t *testing.T) {
	ctx := context.TODO()

	t.Run("Success", func(t *testing.T) {
		mockRepo := new(MockCategoryRepo)
		svc := NewCategorySvc(mockRepo)
		existing := &model.Category{Base: model.Base{ID: 1}, CategoryName: "Old"}
		req := dto.CategoryUpsertReq{CategoryName: "New"}

		mockRepo.On("GetByID", ctx, uint(1)).Return(existing, nil)
		mockRepo.On("Update", ctx, mock.AnythingOfType("*model.Category")).Return(nil)

		res, err := svc.Update(ctx, 1, req)

		assert.NoError(t, err)
		assert.Equal(t, "New", res.CategoryName)
	})

	t.Run("Error_GetByIDFailure", func(t *testing.T) {
		mockRepo := new(MockCategoryRepo)
		svc := NewCategorySvc(mockRepo)
		mockRepo.On("GetByID", ctx, uint(1)).Return(nil, errors.New("not found"))

		res, err := svc.Update(ctx, 1, dto.CategoryUpsertReq{})

		assert.Error(t, err)
		assert.Nil(t, res)
	})

	t.Run("Error_UpdateFailure", func(t *testing.T) {
		mockRepo := new(MockCategoryRepo)
		svc := NewCategorySvc(mockRepo)
		existing := &model.Category{Base: model.Base{ID: 1}}
		
		mockRepo.On("GetByID", ctx, uint(1)).Return(existing, nil)
		mockRepo.On("Update", ctx, mock.AnythingOfType("*model.Category")).Return(errors.New("update failed"))

		res, err := svc.Update(ctx, 1, dto.CategoryUpsertReq{CategoryName: "New"})

		assert.Error(t, err)
		assert.Nil(t, res)
	})
}

func TestCategoryService_Delete(t *testing.T) {
	ctx := context.TODO()

	t.Run("Success", func(t *testing.T) {
		mockRepo := new(MockCategoryRepo)
		svc := NewCategorySvc(mockRepo)
		mockRepo.On("Delete", ctx, uint(1)).Return(nil)

		err := svc.Delete(ctx, 1)

		assert.NoError(t, err)
	})

	t.Run("Error_DeleteFailure", func(t *testing.T) {
		mockRepo := new(MockCategoryRepo)
		svc := NewCategorySvc(mockRepo)
		mockRepo.On("Delete", ctx, uint(1)).Return(errors.New("db error"))

		err := svc.Delete(ctx, 1)

		assert.Error(t, err)
	})
}