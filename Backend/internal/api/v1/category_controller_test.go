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
type MockCategoryService struct {
	mock.Mock
}

func (m *MockCategoryService) List(ctx context.Context) ([]dto.CategoryResponse, error) {
	args := m.Called(ctx)
	return args.Get(0).([]dto.CategoryResponse), args.Error(1)
}

func (m *MockCategoryService) Get(ctx context.Context, id uint) (*dto.CategoryResponse, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.CategoryResponse), args.Error(1)
}

func (m *MockCategoryService) Create(ctx context.Context, req dto.CategoryUpsertReq) (*dto.CategoryResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.CategoryResponse), args.Error(1)
}

func (m *MockCategoryService) Update(ctx context.Context, id uint, req dto.CategoryUpsertReq) (*dto.CategoryResponse, error) {
	args := m.Called(ctx, id, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.CategoryResponse), args.Error(1)
}

func (m *MockCategoryService) Delete(ctx context.Context, id uint) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

// ==========================================
// 2. HELPER: SETUP GIN ROUTER
// ==========================================
func setupCategoryTestRouter(ctrl *CategoryController) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()

	// Minimal Error Middleware to catch c.Error(err)
	r.Use(func(c *gin.Context) {
		c.Next()
		if len(c.Errors) > 0 {
			c.JSON(http.StatusInternalServerError, gin.H{"error": c.Errors.String()})
		}
	})

	api := r.Group("/api/v1/categories")
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

func TestCategoryController_GetAll(t *testing.T) {
	mockSvc := new(MockCategoryService)
	ctrl := NewCategoryController(mockSvc)
	r := setupCategoryTestRouter(ctrl)

	expectedRes := []dto.CategoryResponse{{ID: 1, CategoryName: "Beverages"}}
	mockSvc.On("List", mock.Anything).Return(expectedRes, nil)

	req, _ := http.NewRequest(http.MethodGet, "/api/v1/categories", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var response []dto.CategoryResponse
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, 1, len(response))
	assert.Equal(t, "Beverages", response[0].CategoryName)
}

func TestCategoryController_GetByID(t *testing.T) {
	mockSvc := new(MockCategoryService)
	ctrl := NewCategoryController(mockSvc)
	r := setupCategoryTestRouter(ctrl)

	expectedRes := &dto.CategoryResponse{ID: 5, CategoryName: "Grains"}
	mockSvc.On("Get", mock.Anything, uint(5)).Return(expectedRes, nil)

	req, _ := http.NewRequest(http.MethodGet, "/api/v1/categories/5", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	mockSvc.AssertExpectations(t)
}

func TestCategoryController_Create_Success(t *testing.T) {
	mockSvc := new(MockCategoryService)
	ctrl := NewCategoryController(mockSvc)
	r := setupCategoryTestRouter(ctrl)

	reqDto := dto.CategoryUpsertReq{CategoryName: "New Category"}
	resDto := &dto.CategoryResponse{ID: 10, CategoryName: "New Category"}

	mockSvc.On("Create", mock.Anything, reqDto).Return(resDto, nil)

	payload, _ := json.Marshal(reqDto)
	req, _ := http.NewRequest(http.MethodPost, "/api/v1/categories", bytes.NewBuffer(payload))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
	mockSvc.AssertExpectations(t)
}

func TestCategoryController_Update_Success(t *testing.T) {
	mockSvc := new(MockCategoryService)
	ctrl := NewCategoryController(mockSvc)
	r := setupCategoryTestRouter(ctrl)

	reqDto := dto.CategoryUpsertReq{CategoryName: "Updated"}
	resDto := &dto.CategoryResponse{ID: 1, CategoryName: "Updated"}

	mockSvc.On("Update", mock.Anything, uint(1), reqDto).Return(resDto, nil)

	payload, _ := json.Marshal(reqDto)
	req, _ := http.NewRequest(http.MethodPut, "/api/v1/categories/1", bytes.NewBuffer(payload))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	mockSvc.AssertExpectations(t)
}

func TestCategoryController_Delete_Success(t *testing.T) {
	mockSvc := new(MockCategoryService)
	ctrl := NewCategoryController(mockSvc)
	r := setupCategoryTestRouter(ctrl)

	mockSvc.On("Delete", mock.Anything, uint(1)).Return(nil)

	req, _ := http.NewRequest(http.MethodDelete, "/api/v1/categories/1", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNoContent, w.Code)
	mockSvc.AssertExpectations(t)
}

// ==========================================
// 4. ERROR CASE TESTS (BINDING ERRORS)
// ==========================================

func TestCategoryController_Create_BindingError(t *testing.T) {
	mockSvc := new(MockCategoryService)
	ctrl := NewCategoryController(mockSvc)
	r := setupCategoryTestRouter(ctrl)

	// Sending malformed JSON (missing closing bracket)
	payload := []byte(`{"category_name": "Incomplete JSON`) 
	req, _ := http.NewRequest(http.MethodPost, "/api/v1/categories", bytes.NewBuffer(payload))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	// Should return 400 Bad Request
	assert.Equal(t, http.StatusBadRequest, w.Code)
	// Verify service was NEVER called
	mockSvc.AssertNotCalled(t, "Create", mock.Anything, mock.Anything)
}

func TestCategoryController_Update_BindingError(t *testing.T) {
	mockSvc := new(MockCategoryService)
	ctrl := NewCategoryController(mockSvc)
	r := setupCategoryTestRouter(ctrl)

	// Sending malformed JSON
	payload := []byte(`{"category_name": 123}`) // Wrong type if validation is strict, but malformed is better
	payload = []byte(`{invalid-json}`)
	
	req, _ := http.NewRequest(http.MethodPut, "/api/v1/categories/1", bytes.NewBuffer(payload))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	mockSvc.AssertNotCalled(t, "Update", mock.Anything, mock.Anything, mock.Anything)
}

// ==========================================
// 5. ERROR CASE TESTS (SERVICE ERRORS)
// ==========================================

func TestCategoryController_GetAll_ServiceError(t *testing.T) {
	mockSvc := new(MockCategoryService)
	ctrl := NewCategoryController(mockSvc)
	r := setupCategoryTestRouter(ctrl)

	mockSvc.On("List", mock.Anything).Return([]dto.CategoryResponse{}, errors.New("db down"))

	req, _ := http.NewRequest(http.MethodGet, "/api/v1/categories", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestCategoryController_GetByID_ServiceError(t *testing.T) {
	mockSvc := new(MockCategoryService)
	ctrl := NewCategoryController(mockSvc)
	r := setupCategoryTestRouter(ctrl)

	mockSvc.On("Get", mock.Anything, uint(1)).Return(nil, errors.New("service failure"))

	req, _ := http.NewRequest(http.MethodGet, "/api/v1/categories/1", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	mockSvc.AssertExpectations(t)
}

func TestCategoryController_Create_ServiceError(t *testing.T) {
	mockSvc := new(MockCategoryService)
	ctrl := NewCategoryController(mockSvc)
	r := setupCategoryTestRouter(ctrl)

	reqDto := dto.CategoryUpsertReq{CategoryName: "Valid Name"}
	mockSvc.On("Create", mock.Anything, reqDto).Return(nil, errors.New("database insert failed"))

	payload, _ := json.Marshal(reqDto)
	req, _ := http.NewRequest(http.MethodPost, "/api/v1/categories", bytes.NewBuffer(payload))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	mockSvc.AssertExpectations(t)
}

func TestCategoryController_Update_ServiceError(t *testing.T) {
	mockSvc := new(MockCategoryService)
	ctrl := NewCategoryController(mockSvc)
	r := setupCategoryTestRouter(ctrl)

	reqDto := dto.CategoryUpsertReq{CategoryName: "New Name"}
	mockSvc.On("Update", mock.Anything, uint(1), reqDto).Return(nil, errors.New("update conflict"))

	payload, _ := json.Marshal(reqDto)
	req, _ := http.NewRequest(http.MethodPut, "/api/v1/categories/1", bytes.NewBuffer(payload))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	mockSvc.AssertExpectations(t)
}

func TestCategoryController_Delete_ServiceError(t *testing.T) {
	mockSvc := new(MockCategoryService)
	ctrl := NewCategoryController(mockSvc)
	r := setupCategoryTestRouter(ctrl)

	mockSvc.On("Delete", mock.Anything, uint(99)).Return(errors.New("constraint violation"))

	req, _ := http.NewRequest(http.MethodDelete, "/api/v1/categories/99", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}