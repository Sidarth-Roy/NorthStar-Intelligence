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
type MockShipperService struct {
	mock.Mock
}

func (m *MockShipperService) List(ctx context.Context) ([]dto.ShipperResponse, error) {
	args := m.Called(ctx)
	return args.Get(0).([]dto.ShipperResponse), args.Error(1)
}

func (m *MockShipperService) Get(ctx context.Context, id uint) (*dto.ShipperResponse, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil { return nil, args.Error(1) }
	return args.Get(0).(*dto.ShipperResponse), args.Error(1)
}

func (m *MockShipperService) Create(ctx context.Context, req dto.ShipperUpsertReq) (*dto.ShipperResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil { return nil, args.Error(1) }
	return args.Get(0).(*dto.ShipperResponse), args.Error(1)
}

func (m *MockShipperService) Update(ctx context.Context, id uint, req dto.ShipperUpsertReq) (*dto.ShipperResponse, error) {
	args := m.Called(ctx, id, req)
	if args.Get(0) == nil { return nil, args.Error(1) }
	return args.Get(0).(*dto.ShipperResponse), args.Error(1)
}

func (m *MockShipperService) Delete(ctx context.Context, id uint) error {
	return m.Called(ctx, id).Error(0)
}

// ==========================================
// 2. HELPER: SETUP GIN ROUTER
// ==========================================
func setupShipperTestRouter(ctrl *ShipperController) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()

	// Error middleware to handle the c.Error(err) branches
	r.Use(func(c *gin.Context) {
		c.Next()
		if len(c.Errors) > 0 {
			c.JSON(http.StatusInternalServerError, gin.H{"error": c.Errors.String()})
		}
	})

	api := r.Group("/api/v1/shippers")
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

func TestShipperController_GetAll_Success(t *testing.T) {
	mockSvc := new(MockShipperService)
	ctrl := NewShipperController(mockSvc)
	r := setupShipperTestRouter(ctrl)

	expected := []dto.ShipperResponse{{ID: 1, CompanyName: "Speedy Express"}}
	mockSvc.On("List", mock.Anything).Return(expected, nil)

	req, _ := http.NewRequest(http.MethodGet, "/api/v1/shippers", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestShipperController_GetByID_Success(t *testing.T) {
	mockSvc := new(MockShipperService)
	ctrl := NewShipperController(mockSvc)
	r := setupShipperTestRouter(ctrl)

	expected := &dto.ShipperResponse{ID: 1, CompanyName: "United Package"}
	mockSvc.On("Get", mock.Anything, uint(1)).Return(expected, nil)

	req, _ := http.NewRequest(http.MethodGet, "/api/v1/shippers/1", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestShipperController_Create_Success(t *testing.T) {
	mockSvc := new(MockShipperService)
	ctrl := NewShipperController(mockSvc)
	r := setupShipperTestRouter(ctrl)

	reqDto := dto.ShipperUpsertReq{CompanyName: "Federal Shipping"}
	resDto := &dto.ShipperResponse{ID: 3, CompanyName: "Federal Shipping"}

	mockSvc.On("Create", mock.Anything, reqDto).Return(resDto, nil)

	body, _ := json.Marshal(reqDto)
	req, _ := http.NewRequest(http.MethodPost, "/api/v1/shippers", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
}

func TestShipperController_Update_Success(t *testing.T) {
	mockSvc := new(MockShipperService)
	ctrl := NewShipperController(mockSvc)
	r := setupShipperTestRouter(ctrl)

	reqDto := dto.ShipperUpsertReq{CompanyName: "Updated Shipping"}
	resDto := &dto.ShipperResponse{ID: 1, CompanyName: "Updated Shipping"}

	mockSvc.On("Update", mock.Anything, uint(1), reqDto).Return(resDto, nil)

	body, _ := json.Marshal(reqDto)
	req, _ := http.NewRequest(http.MethodPut, "/api/v1/shippers/1", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestShipperController_Delete_Success(t *testing.T) {
	mockSvc := new(MockShipperService)
	ctrl := NewShipperController(mockSvc)
	r := setupShipperTestRouter(ctrl)

	mockSvc.On("Delete", mock.Anything, uint(1)).Return(nil)

	req, _ := http.NewRequest(http.MethodDelete, "/api/v1/shippers/1", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNoContent, w.Code)
}

// ==========================================
// 4. ERROR CASE TESTS (BINDING ERRORS - 400)
// ==========================================

func TestShipperController_Create_BindingError(t *testing.T) {
	mockSvc := new(MockShipperService)
	ctrl := NewShipperController(mockSvc)
	r := setupShipperTestRouter(ctrl)

	// Case: Name too short (min=2 validation fails)
	badDto := dto.ShipperUpsertReq{CompanyName: "A"}
	body, _ := json.Marshal(badDto)
	req, _ := http.NewRequest(http.MethodPost, "/api/v1/shippers", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestShipperController_Update_BindingError(t *testing.T) {
	mockSvc := new(MockShipperService)
	ctrl := NewShipperController(mockSvc)
	r := setupShipperTestRouter(ctrl)

	// Case: Missing required field
	req, _ := http.NewRequest(http.MethodPut, "/api/v1/shippers/1", bytes.NewBuffer([]byte(`{}`)))
	req.Header.Set("Content-Type", "application/json")
	
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

// ==========================================
// 5. ERROR CASE TESTS (SERVICE ERRORS - 500)
// ==========================================

func TestShipperController_GetAll_ServiceError(t *testing.T) {
	mockSvc := new(MockShipperService)
	ctrl := NewShipperController(mockSvc)
	r := setupShipperTestRouter(ctrl)

	mockSvc.On("List", mock.Anything).Return([]dto.ShipperResponse{}, errors.New("db fail"))

	req, _ := http.NewRequest(http.MethodGet, "/api/v1/shippers", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestShipperController_GetByID_ServiceError(t *testing.T) {
	mockSvc := new(MockShipperService)
	ctrl := NewShipperController(mockSvc)
	r := setupShipperTestRouter(ctrl)

	mockSvc.On("Get", mock.Anything, uint(1)).Return(nil, errors.New("not found"))

	req, _ := http.NewRequest(http.MethodGet, "/api/v1/shippers/1", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestShipperController_Create_ServiceError(t *testing.T) {
	mockSvc := new(MockShipperService)
	ctrl := NewShipperController(mockSvc)
	r := setupShipperTestRouter(ctrl)

	reqDto := dto.ShipperUpsertReq{CompanyName: "Fail Express"}
	mockSvc.On("Create", mock.Anything, reqDto).Return(nil, errors.New("failed create"))

	body, _ := json.Marshal(reqDto)
	req, _ := http.NewRequest(http.MethodPost, "/api/v1/shippers", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestShipperController_Update_ServiceError(t *testing.T) {
	mockSvc := new(MockShipperService)
	ctrl := NewShipperController(mockSvc)
	r := setupShipperTestRouter(ctrl)

	reqDto := dto.ShipperUpsertReq{CompanyName: "Fail Express"}
	mockSvc.On("Update", mock.Anything, uint(1), reqDto).Return(nil, errors.New("failed update"))

	body, _ := json.Marshal(reqDto)
	req, _ := http.NewRequest(http.MethodPut, "/api/v1/shippers/1", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestShipperController_Delete_ServiceError(t *testing.T) {
	mockSvc := new(MockShipperService)
	ctrl := NewShipperController(mockSvc)
	r := setupShipperTestRouter(ctrl)

	mockSvc.On("Delete", mock.Anything, uint(1)).Return(errors.New("failed delete"))

	req, _ := http.NewRequest(http.MethodDelete, "/api/v1/shippers/1", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}