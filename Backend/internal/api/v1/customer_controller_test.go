package v1

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Sidarth-Roy/NorthStar-Intelligence/Backend/internal/dto"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// ==========================================
// 1. MOCK THE SERVICE LAYER
// ==========================================
type MockCustomerService struct {
	mock.Mock
}

func (m *MockCustomerService) List(ctx context.Context) ([]dto.CustomerResponse, error) {
	args := m.Called(ctx)
	return args.Get(0).([]dto.CustomerResponse), args.Error(1)
}

func (m *MockCustomerService) Get(ctx context.Context, id uint) (*dto.CustomerResponse, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil { return nil, args.Error(1) }
	return args.Get(0).(*dto.CustomerResponse), args.Error(1)
}

func (m *MockCustomerService) Create(ctx context.Context, req dto.CustomerUpsertReq) (*dto.CustomerResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil { return nil, args.Error(1) }
	return args.Get(0).(*dto.CustomerResponse), args.Error(1)
}

func (m *MockCustomerService) Update(ctx context.Context, id uint, req dto.CustomerUpsertReq) (*dto.CustomerResponse, error) {
	args := m.Called(ctx, id, req)
	if args.Get(0) == nil { return nil, args.Error(1) }
	return args.Get(0).(*dto.CustomerResponse), args.Error(1)
}

func (m *MockCustomerService) Delete(ctx context.Context, id uint) error {
	return m.Called(ctx, id).Error(0)
}

// ==========================================
// 2. HELPER: SETUP GIN ROUTER
// ==========================================
func setupCustomerTestRouter(ctrl *CustomerController) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()

	r.Use(func(c *gin.Context) {
		c.Next()
		if len(c.Errors) > 0 {
			c.JSON(http.StatusInternalServerError, gin.H{"error": c.Errors.String()})
		}
	})

	api := r.Group("/api/v1/customers")
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

func TestCustomerController_GetAll_Success(t *testing.T) {
	mockSvc := new(MockCustomerService)
	ctrl := NewCustomerController(mockSvc)
	r := setupCustomerTestRouter(ctrl)

	expected := []dto.CustomerResponse{{ID: 1, CustomerID: "ALFKI", CompanyName: "Test Corp"}}
	mockSvc.On("List", mock.Anything).Return(expected, nil)

	req, _ := http.NewRequest(http.MethodGet, "/api/v1/customers", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestCustomerController_GetByID_Success(t *testing.T) {
	mockSvc := new(MockCustomerService)
	ctrl := NewCustomerController(mockSvc)
	r := setupCustomerTestRouter(ctrl)

	expected := &dto.CustomerResponse{ID: 1, CustomerID: "ALFKI", CompanyName: "Test Corp"}
	mockSvc.On("Get", mock.Anything, uint(1)).Return(expected, nil)

	req, _ := http.NewRequest(http.MethodGet, "/api/v1/customers/1", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestCustomerController_Create_Success(t *testing.T) {
	mockSvc := new(MockCustomerService)
	ctrl := NewCustomerController(mockSvc)
	r := setupCustomerTestRouter(ctrl)

	reqDto := dto.CustomerUpsertReq{
		CustomerID:  "NEWCO", 
		CompanyName: "Success Corp",
	}
	resDto := &dto.CustomerResponse{ID: 1, CustomerID: "NEWCO", CompanyName: "Success Corp"}

	mockSvc.On("Create", mock.Anything, reqDto).Return(resDto, nil)

	body, _ := json.Marshal(reqDto)
	req, _ := http.NewRequest(http.MethodPost, "/api/v1/customers", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
}

func TestCustomerController_Update_Success(t *testing.T) {
	mockSvc := new(MockCustomerService)
	ctrl := NewCustomerController(mockSvc)
	r := setupCustomerTestRouter(ctrl)

	reqDto := dto.CustomerUpsertReq{
		CustomerID:  "ALFKI", 
		CompanyName: "Updated Corp",
	}
	resDto := &dto.CustomerResponse{ID: 1, CustomerID: "ALFKI", CompanyName: "Updated Corp"}

	mockSvc.On("Update", mock.Anything, uint(1), reqDto).Return(resDto, nil)

	body, _ := json.Marshal(reqDto)
	req, _ := http.NewRequest(http.MethodPut, "/api/v1/customers/1", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestCustomerController_Delete_Success(t *testing.T) {
	mockSvc := new(MockCustomerService)
	ctrl := NewCustomerController(mockSvc)
	r := setupCustomerTestRouter(ctrl)

	mockSvc.On("Delete", mock.Anything, uint(1)).Return(nil)

	req, _ := http.NewRequest(http.MethodDelete, "/api/v1/customers/1", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNoContent, w.Code)
}

// ==========================================
// 4. ERROR CASE TESTS (BINDING ERRORS)
// ==========================================

func TestCustomerController_Create_BindingError(t *testing.T) {
	mockSvc := new(MockCustomerService)
	ctrl := NewCustomerController(mockSvc)
	r := setupCustomerTestRouter(ctrl)

	// Send completely invalid JSON string
	req, _ := http.NewRequest(http.MethodPost, "/api/v1/customers", bytes.NewBuffer([]byte(`{"company_name": "missing quote}`)))
	req.Header.Set("Content-Type", "application/json")
	
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	mockSvc.AssertNotCalled(t, "Create", mock.Anything, mock.Anything)
}

func TestCustomerController_Update_BindingError(t *testing.T) {
	mockSvc := new(MockCustomerService)
	ctrl := NewCustomerController(mockSvc)
	r := setupCustomerTestRouter(ctrl)

	// Send completely invalid JSON string
	req, _ := http.NewRequest(http.MethodPut, "/api/v1/customers/1", bytes.NewBuffer([]byte(`{"company_name": "missing quote}`)))
	req.Header.Set("Content-Type", "application/json")
	
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	mockSvc.AssertNotCalled(t, "Update", mock.Anything, mock.Anything, mock.Anything)
}

// ==========================================
// 5. ERROR CASE TESTS (SERVICE ERRORS)
// ==========================================

func TestCustomerController_GetAll_ServiceError(t *testing.T) {
	mockSvc := new(MockCustomerService)
	ctrl := NewCustomerController(mockSvc)
	r := setupCustomerTestRouter(ctrl)

	mockSvc.On("List", mock.Anything).Return([]dto.CustomerResponse{}, errors.New("db error"))

	req, _ := http.NewRequest(http.MethodGet, "/api/v1/customers", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestCustomerController_GetByID_ServiceError(t *testing.T) {
	mockSvc := new(MockCustomerService)
	ctrl := NewCustomerController(mockSvc)
	r := setupCustomerTestRouter(ctrl)

	mockSvc.On("Get", mock.Anything, uint(1)).Return(nil, errors.New("not found"))

	req, _ := http.NewRequest(http.MethodGet, "/api/v1/customers/1", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestCustomerController_Create_ServiceError(t *testing.T) {
	mockSvc := new(MockCustomerService)
	ctrl := NewCustomerController(mockSvc)
	r := setupCustomerTestRouter(ctrl)

	reqDto := dto.CustomerUpsertReq{
		CustomerID: "FAIL1", 
		CompanyName: "Fail Corp",
	}
	mockSvc.On("Create", mock.Anything, reqDto).Return(nil, errors.New("db error"))

	body, _ := json.Marshal(reqDto)
	req, _ := http.NewRequest(http.MethodPost, "/api/v1/customers", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestCustomerController_Update_ServiceError(t *testing.T) {
	mockSvc := new(MockCustomerService)
	ctrl := NewCustomerController(mockSvc)
	r := setupCustomerTestRouter(ctrl)

	reqDto := dto.CustomerUpsertReq{
		CustomerID: "ALFKI", 
		CompanyName: "Fail Corp",
	}
	mockSvc.On("Update", mock.Anything, uint(1), reqDto).Return(nil, errors.New("db error"))

	body, _ := json.Marshal(reqDto)
	req, _ := http.NewRequest(http.MethodPut, "/api/v1/customers/1", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestCustomerController_Delete_ServiceError(t *testing.T) {
	mockSvc := new(MockCustomerService)
	ctrl := NewCustomerController(mockSvc)
	r := setupCustomerTestRouter(ctrl)

	mockSvc.On("Delete", mock.Anything, uint(1)).Return(errors.New("delete failed"))

	req, _ := http.NewRequest(http.MethodDelete, "/api/v1/customers/1", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}