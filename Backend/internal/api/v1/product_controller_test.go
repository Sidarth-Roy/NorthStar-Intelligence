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
type MockProductService struct {
	mock.Mock
}

func (m *MockProductService) List(ctx context.Context) ([]dto.ProductResponse, error) {
	args := m.Called(ctx)
	return args.Get(0).([]dto.ProductResponse), args.Error(1)
}

func (m *MockProductService) Get(ctx context.Context, id uint) (*dto.ProductResponse, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil { return nil, args.Error(1) }
	return args.Get(0).(*dto.ProductResponse), args.Error(1)
}

func (m *MockProductService) Create(ctx context.Context, req dto.ProductUpsertReq) (*dto.ProductResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil { return nil, args.Error(1) }
	return args.Get(0).(*dto.ProductResponse), args.Error(1)
}

func (m *MockProductService) Update(ctx context.Context, id uint, req dto.ProductUpsertReq) (*dto.ProductResponse, error) {
	args := m.Called(ctx, id, req)
	if args.Get(0) == nil { return nil, args.Error(1) }
	return args.Get(0).(*dto.ProductResponse), args.Error(1)
}

func (m *MockProductService) Delete(ctx context.Context, id uint) error {
	return m.Called(ctx, id).Error(0)
}

// ==========================================
// 2. HELPER: SETUP GIN ROUTER
// ==========================================
func setupProductTestRouter(ctrl *ProductController) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()

	// Essential middleware to handle c.Error(err) branches
	r.Use(func(c *gin.Context) {
		c.Next()
		if len(c.Errors) > 0 {
			c.JSON(http.StatusInternalServerError, gin.H{"error": c.Errors.String()})
		}
	})

	api := r.Group("/api/v1/products")
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

func TestProductController_GetAll_Success(t *testing.T) {
	mockSvc := new(MockProductService)
	ctrl := NewProductController(mockSvc)
	r := setupProductTestRouter(ctrl)

	expected := []dto.ProductResponse{{ID: 1, ProductName: "Chai", UnitPrice: 18.0}}
	mockSvc.On("List", mock.Anything).Return(expected, nil)

	req, _ := http.NewRequest(http.MethodGet, "/api/v1/products", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestProductController_GetByID_Success(t *testing.T) {
	mockSvc := new(MockProductService)
	ctrl := NewProductController(mockSvc)
	r := setupProductTestRouter(ctrl)

	expected := &dto.ProductResponse{ID: 1, ProductName: "Chai"}
	mockSvc.On("Get", mock.Anything, uint(1)).Return(expected, nil)

	req, _ := http.NewRequest(http.MethodGet, "/api/v1/products/1", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestProductController_Create_Success(t *testing.T) {
	mockSvc := new(MockProductService)
	ctrl := NewProductController(mockSvc)
	r := setupProductTestRouter(ctrl)

	reqDto := dto.ProductUpsertReq{
		ProductName: "New Tea", 
		UnitPrice:   10.5, 
		CategoryID:  1,
	}
	resDto := &dto.ProductResponse{ID: 2, ProductName: "New Tea"}

	mockSvc.On("Create", mock.Anything, reqDto).Return(resDto, nil)

	body, _ := json.Marshal(reqDto)
	req, _ := http.NewRequest(http.MethodPost, "/api/v1/products", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
}

func TestProductController_Update_Success(t *testing.T) {
	mockSvc := new(MockProductService)
	ctrl := NewProductController(mockSvc)
	r := setupProductTestRouter(ctrl)

	reqDto := dto.ProductUpsertReq{
		ProductName: "Updated Chai", 
		UnitPrice:   20.0, 
		CategoryID:  1,
	}
	resDto := &dto.ProductResponse{ID: 1, ProductName: "Updated Chai"}

	mockSvc.On("Update", mock.Anything, uint(1), reqDto).Return(resDto, nil)

	body, _ := json.Marshal(reqDto)
	req, _ := http.NewRequest(http.MethodPut, "/api/v1/products/1", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestProductController_Delete_Success(t *testing.T) {
	mockSvc := new(MockProductService)
	ctrl := NewProductController(mockSvc)
	r := setupProductTestRouter(ctrl)

	mockSvc.On("Delete", mock.Anything, uint(1)).Return(nil)

	req, _ := http.NewRequest(http.MethodDelete, "/api/v1/products/1", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNoContent, w.Code)
}

// ==========================================
// 4. ERROR CASE TESTS (BINDING ERRORS - 400)
// ==========================================

func TestProductController_Create_BindingError(t *testing.T) {
	mockSvc := new(MockProductService)
	ctrl := NewProductController(mockSvc)
	r := setupProductTestRouter(ctrl)

	// Case: UnitPrice <= 0 (gt=0 validation fails)
	badDto := dto.ProductUpsertReq{
		ProductName: "Bad Price", 
		UnitPrice:   -1.0, 
		CategoryID:  1,
	}
	body, _ := json.Marshal(badDto)
	req, _ := http.NewRequest(http.MethodPost, "/api/v1/products", bytes.NewBuffer(body))
	
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestProductController_Update_BindingError(t *testing.T) {
	mockSvc := new(MockProductService)
	ctrl := NewProductController(mockSvc)
	r := setupProductTestRouter(ctrl)

	// Case: Name too short (min=3 validation fails)
	badDto := dto.ProductUpsertReq{
		ProductName: "Ab", 
		UnitPrice:   10.0, 
		CategoryID:  1,
	}
	body, _ := json.Marshal(badDto)
	req, _ := http.NewRequest(http.MethodPut, "/api/v1/products/1", bytes.NewBuffer(body))
	
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

// ==========================================
// 5. ERROR CASE TESTS (SERVICE ERRORS - 500)
// ==========================================

func TestProductController_GetAll_ServiceError(t *testing.T) {
	mockSvc := new(MockProductService)
	ctrl := NewProductController(mockSvc)
	r := setupProductTestRouter(ctrl)

	mockSvc.On("List", mock.Anything).Return([]dto.ProductResponse{}, errors.New("db error"))

	req, _ := http.NewRequest(http.MethodGet, "/api/v1/products", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestProductController_GetByID_ServiceError(t *testing.T) {
	mockSvc := new(MockProductService)
	ctrl := NewProductController(mockSvc)
	r := setupProductTestRouter(ctrl)

	mockSvc.On("Get", mock.Anything, uint(1)).Return(nil, errors.New("not found"))

	req, _ := http.NewRequest(http.MethodGet, "/api/v1/products/1", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestProductController_Create_ServiceError(t *testing.T) {
	mockSvc := new(MockProductService)
	ctrl := NewProductController(mockSvc)
	r := setupProductTestRouter(ctrl)

	reqDto := dto.ProductUpsertReq{ProductName: "Chai", UnitPrice: 10.0, CategoryID: 1}
	mockSvc.On("Create", mock.Anything, reqDto).Return(nil, errors.New("failed to save"))

	body, _ := json.Marshal(reqDto)
	req, _ := http.NewRequest(http.MethodPost, "/api/v1/products", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestProductController_Update_ServiceError(t *testing.T) {
	mockSvc := new(MockProductService)
	ctrl := NewProductController(mockSvc)
	r := setupProductTestRouter(ctrl)

	reqDto := dto.ProductUpsertReq{ProductName: "Chai", UnitPrice: 10.0, CategoryID: 1}
	mockSvc.On("Update", mock.Anything, uint(1), reqDto).Return(nil, errors.New("update failed"))

	body, _ := json.Marshal(reqDto)
	req, _ := http.NewRequest(http.MethodPut, "/api/v1/products/1", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestProductController_Delete_ServiceError(t *testing.T) {
	mockSvc := new(MockProductService)
	ctrl := NewProductController(mockSvc)
	r := setupProductTestRouter(ctrl)

	mockSvc.On("Delete", mock.Anything, uint(1)).Return(errors.New("delete failed"))

	req, _ := http.NewRequest(http.MethodDelete, "/api/v1/products/1", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}