package v1

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/Sidarth-Roy/NorthStar-Intelligence/Backend/internal/dto"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// ==========================================
// 1. MOCK THE SERVICE LAYER
// ==========================================
type MockOrderService struct {
	mock.Mock
}

func (m *MockOrderService) List(ctx context.Context) ([]dto.OrderResponse, error) {
	args := m.Called(ctx)
	return args.Get(0).([]dto.OrderResponse), args.Error(1)
}

func (m *MockOrderService) Get(ctx context.Context, id uint) (*dto.OrderResponse, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil { return nil, args.Error(1) }
	return args.Get(0).(*dto.OrderResponse), args.Error(1)
}

func (m *MockOrderService) Create(ctx context.Context, req dto.OrderUpsertReq) (*dto.OrderResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil { return nil, args.Error(1) }
	return args.Get(0).(*dto.OrderResponse), args.Error(1)
}

func (m *MockOrderService) Update(ctx context.Context, id uint, req dto.OrderUpsertReq) (*dto.OrderResponse, error) {
	args := m.Called(ctx, id, req)
	if args.Get(0) == nil { return nil, args.Error(1) }
	return args.Get(0).(*dto.OrderResponse), args.Error(1)
}

func (m *MockOrderService) Delete(ctx context.Context, id uint) error {
	return m.Called(ctx, id).Error(0)
}

// ==========================================
// 2. HELPER: SETUP GIN ROUTER
// ==========================================
func setupOrderTestRouter(ctrl *OrderController) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()

	r.Use(func(c *gin.Context) {
		c.Next()
		if len(c.Errors) > 0 {
			c.JSON(http.StatusInternalServerError, gin.H{"error": c.Errors.String()})
		}
	})

	api := r.Group("/api/v1/orders")
	{
		api.GET("", ctrl.GetAll)
		api.GET("/:id", ctrl.GetByID)
		api.POST("", ctrl.Create)
		api.PUT("/:id", ctrl.Update)
		api.DELETE("/:id", ctrl.Delete)
	}
	return r
}

// ==========================================
// 3. SUCCESS CASE TESTS
// ==========================================

func TestOrderController_GetAll_Success(t *testing.T) {
	mockSvc := new(MockOrderService)
	ctrl := NewOrderController(mockSvc)
	r := setupOrderTestRouter(ctrl)

	expected := []dto.OrderResponse{{ID: 1, CustomerID: "ALFKI"}}
	mockSvc.On("List", mock.Anything).Return(expected, nil)

	req, _ := http.NewRequest(http.MethodGet, "/api/v1/orders", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestOrderController_GetByID_Success(t *testing.T) {
	mockSvc := new(MockOrderService)
	ctrl := NewOrderController(mockSvc)
	r := setupOrderTestRouter(ctrl)

	expected := &dto.OrderResponse{ID: 1, CustomerID: "ALFKI"}
	mockSvc.On("Get", mock.Anything, uint(1)).Return(expected, nil)

	req, _ := http.NewRequest(http.MethodGet, "/api/v1/orders/1", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestOrderController_Create_Success(t *testing.T) {
	mockSvc := new(MockOrderService)
	ctrl := NewOrderController(mockSvc)
	r := setupOrderTestRouter(ctrl)

	testTime := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	reqDto := dto.OrderUpsertReq{
		CustomerID:   "ALFKI",
		EmployeeID:   1,
		OrderDate:    testTime,
		RequiredDate: testTime,
		ShipperID:    1,
		Freight:      10.0,
		OrderDetails: []dto.OrderDetailReq{
			{ProductID: 1, UnitPrice: 10.0, Quantity: 2},
		},
	}
	resDto := &dto.OrderResponse{ID: 1, CustomerID: "ALFKI"}

	mockSvc.On("Create", mock.Anything, reqDto).Return(resDto, nil)

	body, _ := json.Marshal(reqDto)
	req, _ := http.NewRequest(http.MethodPost, "/api/v1/orders", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
}

func TestOrderController_Update_Success(t *testing.T) {
	mockSvc := new(MockOrderService)
	ctrl := NewOrderController(mockSvc)
	r := setupOrderTestRouter(ctrl)

	testTime := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	reqDto := dto.OrderUpsertReq{
		CustomerID:   "ALFKI",
		EmployeeID:   1,
		OrderDate:    testTime,
		RequiredDate: testTime,
		ShipperID:    1,
		OrderDetails: []dto.OrderDetailReq{
			{ProductID: 1, UnitPrice: 10.0, Quantity: 2},
		},
	}
	resDto := &dto.OrderResponse{ID: 1, CustomerID: "ALFKI"}

	mockSvc.On("Update", mock.Anything, uint(1), reqDto).Return(resDto, nil)

	body, _ := json.Marshal(reqDto)
	req, _ := http.NewRequest(http.MethodPut, "/api/v1/orders/1", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestOrderController_Delete_Success(t *testing.T) {
	mockSvc := new(MockOrderService)
	ctrl := NewOrderController(mockSvc)
	r := setupOrderTestRouter(ctrl)

	mockSvc.On("Delete", mock.Anything, uint(1)).Return(nil)

	req, _ := http.NewRequest(http.MethodDelete, "/api/v1/orders/1", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNoContent, w.Code)
}

// ==========================================
// 4. ERROR CASE TESTS (BINDING ERRORS - 400)
// ==========================================

func TestOrderController_Create_BindingError(t *testing.T) {
	mockSvc := new(MockOrderService)
	ctrl := NewOrderController(mockSvc)
	r := setupOrderTestRouter(ctrl)

	// Case: Negative Quantity in OrderDetails (gt=0 fails)
	badDto := dto.OrderUpsertReq{
		CustomerID: "ALFKI",
		OrderDetails: []dto.OrderDetailReq{
			{ProductID: 1, UnitPrice: 10.0, Quantity: -5},
		},
	}
	body, _ := json.Marshal(badDto)
	req, _ := http.NewRequest(http.MethodPost, "/api/v1/orders", bytes.NewBuffer(body))
	
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestOrderController_Update_BindingError(t *testing.T) {
	mockSvc := new(MockOrderService)
	ctrl := NewOrderController(mockSvc)
	r := setupOrderTestRouter(ctrl)

	// Case: Missing Required OrderDetails
	badDto := dto.OrderUpsertReq{CustomerID: "ALFKI"}
	body, _ := json.Marshal(badDto)
	req, _ := http.NewRequest(http.MethodPut, "/api/v1/orders/1", bytes.NewBuffer(body))
	
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

// ==========================================
// 5. ERROR CASE TESTS (SERVICE ERRORS - 500)
// ==========================================

func TestOrderController_GetAll_ServiceError(t *testing.T) {
	mockSvc := new(MockOrderService)
	ctrl := NewOrderController(mockSvc)
	r := setupOrderTestRouter(ctrl)

	mockSvc.On("List", mock.Anything).Return([]dto.OrderResponse{}, errors.New("db error"))

	req, _ := http.NewRequest(http.MethodGet, "/api/v1/orders", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestOrderController_GetByID_ServiceError(t *testing.T) {
	mockSvc := new(MockOrderService)
	ctrl := NewOrderController(mockSvc)
	r := setupOrderTestRouter(ctrl)

	mockSvc.On("Get", mock.Anything, uint(1)).Return(nil, errors.New("not found"))

	req, _ := http.NewRequest(http.MethodGet, "/api/v1/orders/1", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestOrderController_Create_ServiceError(t *testing.T) {
	mockSvc := new(MockOrderService)
	ctrl := NewOrderController(mockSvc)
	r := setupOrderTestRouter(ctrl)

	testTime := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	reqDto := dto.OrderUpsertReq{
		CustomerID: "ALFKI", EmployeeID: 1, OrderDate: testTime, RequiredDate: testTime, ShipperID: 1,
		OrderDetails: []dto.OrderDetailReq{{ProductID: 1, UnitPrice: 10.0, Quantity: 1}},
	}
	mockSvc.On("Create", mock.Anything, reqDto).Return(nil, errors.New("create fail"))

	body, _ := json.Marshal(reqDto)
	req, _ := http.NewRequest(http.MethodPost, "/api/v1/orders", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestOrderController_Update_ServiceError(t *testing.T) {
	mockSvc := new(MockOrderService)
	ctrl := NewOrderController(mockSvc)
	r := setupOrderTestRouter(ctrl)

	testTime := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	reqDto := dto.OrderUpsertReq{
		CustomerID: "ALFKI", EmployeeID: 1, OrderDate: testTime, RequiredDate: testTime, ShipperID: 1,
		OrderDetails: []dto.OrderDetailReq{{ProductID: 1, UnitPrice: 10.0, Quantity: 1}},
	}
	mockSvc.On("Update", mock.Anything, uint(1), reqDto).Return(nil, errors.New("update fail"))

	body, _ := json.Marshal(reqDto)
	req, _ := http.NewRequest(http.MethodPut, "/api/v1/orders/1", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestOrderController_Delete_ServiceError(t *testing.T) {
	mockSvc := new(MockOrderService)
	ctrl := NewOrderController(mockSvc)
	r := setupOrderTestRouter(ctrl)

	mockSvc.On("Delete", mock.Anything, uint(1)).Return(errors.New("delete fail"))

	req, _ := http.NewRequest(http.MethodDelete, "/api/v1/orders/1", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}