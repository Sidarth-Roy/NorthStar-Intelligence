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

// MockShipperRepo implementation
type MockShipperRepo struct {
	mock.Mock
}

func (m *MockShipperRepo) Create(ctx context.Context, s *model.Shipper) error {
	args := m.Called(ctx, s)
	if args.Error(0) == nil {
		s.ID = 1
	}
	return args.Error(0)
}

func (m *MockShipperRepo) GetByID(ctx context.Context, id uint) (*model.Shipper, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Shipper), args.Error(1)
}

func (m *MockShipperRepo) GetAll(ctx context.Context) ([]model.Shipper, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]model.Shipper), args.Error(1)
}

func (m *MockShipperRepo) Update(ctx context.Context, s *model.Shipper) error {
	args := m.Called(ctx, s)
	return args.Error(0)
}

func (m *MockShipperRepo) Delete(ctx context.Context, id uint) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

// --- TESTS ---

func TestShipperService_Create(t *testing.T) {
	ctx := context.TODO()
	t.Run("Success", func(t *testing.T) {
		mockRepo := new(MockShipperRepo)
		svc := NewShipperSvc(mockRepo)
		mockRepo.On("Create", ctx, mock.Anything).Return(nil)
		res, err := svc.Create(ctx, dto.ShipperUpsertReq{CompanyName: "A"})
		assert.NoError(t, err)
		assert.NotNil(t, res)
	})

	t.Run("Error", func(t *testing.T) {
		mockRepo := new(MockShipperRepo)
		svc := NewShipperSvc(mockRepo)
		mockRepo.On("Create", ctx, mock.Anything).Return(errors.New("fail"))
		res, err := svc.Create(ctx, dto.ShipperUpsertReq{CompanyName: "A"})
		assert.Error(t, err)
		assert.Nil(t, res)
	})
}

func TestShipperService_Get(t *testing.T) {
	ctx := context.TODO()
	t.Run("Success", func(t *testing.T) {
		mockRepo := new(MockShipperRepo)
		svc := NewShipperSvc(mockRepo)
		mockRepo.On("GetByID", ctx, uint(1)).Return(&model.Shipper{CompanyName: "A"}, nil)
		res, err := svc.Get(ctx, 1)
		assert.NoError(t, err)
		assert.Equal(t, "A", res.CompanyName)
	})

	t.Run("Error", func(t *testing.T) {
		mockRepo := new(MockShipperRepo)
		svc := NewShipperSvc(mockRepo)
		mockRepo.On("GetByID", ctx, uint(1)).Return(nil, errors.New("not found"))
		res, err := svc.Get(ctx, 1)
		assert.Error(t, err)
		assert.Nil(t, res)
	})
}

func TestShipperService_List(t *testing.T) {
	ctx := context.TODO()
	t.Run("Success", func(t *testing.T) {
		mockRepo := new(MockShipperRepo)
		svc := NewShipperSvc(mockRepo)
		mockRepo.On("GetAll", ctx).Return([]model.Shipper{{CompanyName: "A"}}, nil)
		res, err := svc.List(ctx)
		assert.NoError(t, err)
		assert.Len(t, res, 1)
	})

	t.Run("Error", func(t *testing.T) {
		mockRepo := new(MockShipperRepo)
		svc := NewShipperSvc(mockRepo)
		mockRepo.On("GetAll", ctx).Return(nil, errors.New("db fail"))
		res, err := svc.List(ctx)
		assert.Error(t, err)
		assert.Nil(t, res)
	})
}

func TestShipperService_Update(t *testing.T) {
	ctx := context.TODO()
	t.Run("Success", func(t *testing.T) {
		mockRepo := new(MockShipperRepo)
		svc := NewShipperSvc(mockRepo)
		existing := &model.Shipper{CompanyName: "Old"}
		mockRepo.On("GetByID", ctx, uint(1)).Return(existing, nil)
		mockRepo.On("Update", ctx, mock.Anything).Return(nil)
		res, err := svc.Update(ctx, 1, dto.ShipperUpsertReq{CompanyName: "New"})
		assert.NoError(t, err)
		assert.Equal(t, "New", res.CompanyName)
	})

	t.Run("Error_On_Get", func(t *testing.T) {
		mockRepo := new(MockShipperRepo)
		svc := NewShipperSvc(mockRepo)
		mockRepo.On("GetByID", ctx, uint(1)).Return(nil, errors.New("get fail"))
		res, err := svc.Update(ctx, 1, dto.ShipperUpsertReq{CompanyName: "New"})
		assert.Error(t, err)
		assert.Nil(t, res)
	})

	t.Run("Error_On_Update", func(t *testing.T) {
		mockRepo := new(MockShipperRepo)
		svc := NewShipperSvc(mockRepo)
		mockRepo.On("GetByID", ctx, uint(1)).Return(&model.Shipper{}, nil)
		mockRepo.On("Update", ctx, mock.Anything).Return(errors.New("update fail"))
		res, err := svc.Update(ctx, 1, dto.ShipperUpsertReq{CompanyName: "New"})
		assert.Error(t, err)
		assert.Nil(t, res)
	})
}

func TestShipperService_Delete(t *testing.T) {
	ctx := context.TODO()
	t.Run("Success", func(t *testing.T) {
		mockRepo := new(MockShipperRepo)
		svc := NewShipperSvc(mockRepo)
		mockRepo.On("Delete", ctx, uint(1)).Return(nil)
		err := svc.Delete(ctx, 1)
		assert.NoError(t, err)
	})

	t.Run("Error", func(t *testing.T) {
		mockRepo := new(MockShipperRepo)
		svc := NewShipperSvc(mockRepo)
		mockRepo.On("Delete", ctx, uint(1)).Return(errors.New("del fail"))
		err := svc.Delete(ctx, 1)
		assert.Error(t, err)
	})
}