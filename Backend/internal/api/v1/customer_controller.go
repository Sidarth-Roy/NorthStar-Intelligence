package v1

import (
	"net/http"
	"strconv"
	"github.com/Sidarth-Roy/NorthStar-Intelligence/Backend/internal/dto"
	"github.com/Sidarth-Roy/NorthStar-Intelligence/Backend/internal/service"
	"github.com/gin-gonic/gin"
)

type CustomerController struct{ svc service.CustomerService }

func NewCustomerController(s service.CustomerService) *CustomerController {
	return &CustomerController{svc: s}
}

func (ctrl *CustomerController) GetAll(c *gin.Context) {
	res, err := ctrl.svc.List(c.Request.Context())
	if err != nil { c.Error(err); return }
	c.JSON(http.StatusOK, res)
}

func (ctrl *CustomerController) GetByID(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	res, err := ctrl.svc.Get(c.Request.Context(), uint(id))
	if err != nil { c.Error(err); return }
	c.JSON(http.StatusOK, res)
}

func (ctrl *CustomerController) Create(c *gin.Context) {
	var req dto.CustomerInsertReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	res, err := ctrl.svc.Create(c.Request.Context(), req)
	if err != nil { c.Error(err); return }
	c.JSON(http.StatusCreated, res)
}

func (ctrl *CustomerController) Update(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var req dto.CustomerUpdateReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	res, err := ctrl.svc.Update(c.Request.Context(), uint(id), req)
	if err != nil { c.Error(err); return }
	c.JSON(http.StatusOK, res)
}

func (ctrl *CustomerController) Delete(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	if err := ctrl.svc.Delete(c.Request.Context(), uint(id)); err != nil {
		c.Error(err); return
	}
	c.Status(http.StatusNoContent)
}