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
type MockEmployeeService struct {
	mock.Mock
}

func (m *MockEmployeeService) List(ctx context.Context) ([]dto.EmployeeResponse, error) {
	args := m.Called(ctx)
	return args.Get(0).([]dto.EmployeeResponse), args.Error(1)
}

func (m *MockEmployeeService) Get(ctx context.Context, id uint) (*dto.EmployeeResponse, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil { return nil, args.Error(1) }
	return args.Get(0).(*dto.EmployeeResponse), args.Error(1)
}

func (m *MockEmployeeService) Create(ctx context.Context, req dto.EmployeeUpsertReq) (*dto.EmployeeResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil { return nil, args.Error(1) }
	return args.Get(0).(*dto.EmployeeResponse), args.Error(1)
}

func (m *MockEmployeeService) Update(ctx context.Context, id uint, req dto.EmployeeUpsertReq) (*dto.EmployeeResponse, error) {
	args := m.Called(ctx, id, req)
	if args.Get(0) == nil { return nil, args.Error(1) }
	return args.Get(0).(*dto.EmployeeResponse), args.Error(1)
}

func (m *MockEmployeeService) Delete(ctx context.Context, id uint) error {
	return m.Called(ctx, id).Error(0)
}

// ==========================================
// 2. HELPER: SETUP GIN ROUTER
// ==========================================
func setupEmployeeTestRouter(ctrl *EmployeeController) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()

	// Error middleware to capture c.Error(err)
	r.Use(func(c *gin.Context) {
		c.Next()
		if len(c.Errors) > 0 {
			c.JSON(http.StatusInternalServerError, gin.H{"error": c.Errors.String()})
		}
	})

	api := r.Group("/api/v1/employees")
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

func TestEmployeeController_GetAll_Success(t *testing.T) {
	mockSvc := new(MockEmployeeService)
	ctrl := NewEmployeeController(mockSvc)
	r := setupEmployeeTestRouter(ctrl)

	expected := []dto.EmployeeResponse{{ID: 1, EmployeeName: "John Doe"}}
	mockSvc.On("List", mock.Anything).Return(expected, nil)

	req, _ := http.NewRequest(http.MethodGet, "/api/v1/employees", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestEmployeeController_GetByID_Success(t *testing.T) {
	mockSvc := new(MockEmployeeService)
	ctrl := NewEmployeeController(mockSvc)
	r := setupEmployeeTestRouter(ctrl)

	expected := &dto.EmployeeResponse{ID: 1, EmployeeName: "John Doe"}
	mockSvc.On("Get", mock.Anything, uint(1)).Return(expected, nil)

	req, _ := http.NewRequest(http.MethodGet, "/api/v1/employees/1", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestEmployeeController_Create_Success(t *testing.T) {
	mockSvc := new(MockEmployeeService)
	ctrl := NewEmployeeController(mockSvc)
	r := setupEmployeeTestRouter(ctrl)

	reqDto := dto.EmployeeUpsertReq{EmployeeName: "Jane Doe", Title: "Manager"}
	resDto := &dto.EmployeeResponse{ID: 1, EmployeeName: "Jane Doe"}

	mockSvc.On("Create", mock.Anything, reqDto).Return(resDto, nil)

	body, _ := json.Marshal(reqDto)
	req, _ := http.NewRequest(http.MethodPost, "/api/v1/employees", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
}

func TestEmployeeController_Update_Success(t *testing.T) {
	mockSvc := new(MockEmployeeService)
	ctrl := NewEmployeeController(mockSvc)
	r := setupEmployeeTestRouter(ctrl)

	reqDto := dto.EmployeeUpsertReq{EmployeeName: "Jane Smith", Title: "Director"}
	resDto := &dto.EmployeeResponse{ID: 1, EmployeeName: "Jane Smith"}

	mockSvc.On("Update", mock.Anything, uint(1), reqDto).Return(resDto, nil)

	body, _ := json.Marshal(reqDto)
	req, _ := http.NewRequest(http.MethodPut, "/api/v1/employees/1", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestEmployeeController_Delete_Success(t *testing.T) {
	mockSvc := new(MockEmployeeService)
	ctrl := NewEmployeeController(mockSvc)
	r := setupEmployeeTestRouter(ctrl)

	mockSvc.On("Delete", mock.Anything, uint(1)).Return(nil)

	req, _ := http.NewRequest(http.MethodDelete, "/api/v1/employees/1", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNoContent, w.Code)
}

// ==========================================
// 4. ERROR CASE TESTS (BINDING ERRORS - 400)
// ==========================================

func TestEmployeeController_Create_BindingError(t *testing.T) {
	mockSvc := new(MockEmployeeService)
	ctrl := NewEmployeeController(mockSvc)
	r := setupEmployeeTestRouter(ctrl)

	// Case A: Malformed JSON
	req, _ := http.NewRequest(http.MethodPost, "/api/v1/employees", bytes.NewBuffer([]byte(`{invalid}`)))
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)

	// Case B: Validation Fail (Name too short - min=3)
	badDto := dto.EmployeeUpsertReq{EmployeeName: "Jo", Title: "Dev"}
	body, _ := json.Marshal(badDto)
	req, _ = http.NewRequest(http.MethodPost, "/api/v1/employees", bytes.NewBuffer(body))
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestEmployeeController_Update_BindingError(t *testing.T) {
	mockSvc := new(MockEmployeeService)
	ctrl := NewEmployeeController(mockSvc)
	r := setupEmployeeTestRouter(ctrl)

	// Validation Fail (Missing Title)
	badDto := dto.EmployeeUpsertReq{EmployeeName: "John Doe"} 
	body, _ := json.Marshal(badDto)
	req, _ := http.NewRequest(http.MethodPut, "/api/v1/employees/1", bytes.NewBuffer(body))
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

// ==========================================
// 5. ERROR CASE TESTS (SERVICE ERRORS - 500)
// ==========================================

func TestEmployeeController_GetAll_ServiceError(t *testing.T) {
	mockSvc := new(MockEmployeeService)
	ctrl := NewEmployeeController(mockSvc)
	r := setupEmployeeTestRouter(ctrl)

	mockSvc.On("List", mock.Anything).Return([]dto.EmployeeResponse{}, errors.New("list failed"))

	req, _ := http.NewRequest(http.MethodGet, "/api/v1/employees", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestEmployeeController_GetByID_ServiceError(t *testing.T) {
	mockSvc := new(MockEmployeeService)
	ctrl := NewEmployeeController(mockSvc)
	r := setupEmployeeTestRouter(ctrl)

	mockSvc.On("Get", mock.Anything, uint(1)).Return(nil, errors.New("get failed"))

	req, _ := http.NewRequest(http.MethodGet, "/api/v1/employees/1", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestEmployeeController_Create_ServiceError(t *testing.T) {
	mockSvc := new(MockEmployeeService)
	ctrl := NewEmployeeController(mockSvc)
	r := setupEmployeeTestRouter(ctrl)

	reqDto := dto.EmployeeUpsertReq{EmployeeName: "Jane Doe", Title: "Manager"}
	mockSvc.On("Create", mock.Anything, reqDto).Return(nil, errors.New("create failed"))

	body, _ := json.Marshal(reqDto)
	req, _ := http.NewRequest(http.MethodPost, "/api/v1/employees", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestEmployeeController_Update_ServiceError(t *testing.T) {
	mockSvc := new(MockEmployeeService)
	ctrl := NewEmployeeController(mockSvc)
	r := setupEmployeeTestRouter(ctrl)

	reqDto := dto.EmployeeUpsertReq{EmployeeName: "Jane Smith", Title: "Director"}
	mockSvc.On("Update", mock.Anything, uint(1), reqDto).Return(nil, errors.New("update failed"))

	body, _ := json.Marshal(reqDto)
	req, _ := http.NewRequest(http.MethodPut, "/api/v1/employees/1", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestEmployeeController_Delete_ServiceError(t *testing.T) {
	mockSvc := new(MockEmployeeService)
	ctrl := NewEmployeeController(mockSvc)
	r := setupEmployeeTestRouter(ctrl)

	mockSvc.On("Delete", mock.Anything, uint(1)).Return(errors.New("delete failed"))

	req, _ := http.NewRequest(http.MethodDelete, "/api/v1/employees/1", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}