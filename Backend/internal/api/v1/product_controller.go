package v1

import (
	"net/http"
	"strconv"
	"github.com/Sidarth-Roy/NorthStar-Intelligence/Backend/internal/dto"
	"github.com/Sidarth-Roy/NorthStar-Intelligence/Backend/internal/service"
	"github.com/gin-gonic/gin"
)

type ProductController struct{ svc service.ProductService }

func NewProductController(s service.ProductService) *ProductController {
	return &ProductController{svc: s}
}

func (ctrl *ProductController) GetAll(c *gin.Context) {
	res, err := ctrl.svc.List(c.Request.Context())
	if err != nil { c.Error(err); return }
	c.JSON(http.StatusOK, res)
}

func (ctrl *ProductController) GetByID(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	res, err := ctrl.svc.Get(c.Request.Context(), uint(id))
	if err != nil { c.Error(err); return }
	c.JSON(http.StatusOK, res)
}

func (ctrl *ProductController) Create(c *gin.Context) {
	var req dto.ProductUpsertReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	res, err := ctrl.svc.Create(c.Request.Context(), req)
	if err != nil { c.Error(err); return }
	c.JSON(http.StatusCreated, res)
}

func (ctrl *ProductController) Update(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var req dto.ProductUpsertReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	res, err := ctrl.svc.Update(c.Request.Context(), uint(id), req)
	if err != nil { c.Error(err); return }
	c.JSON(http.StatusOK, res)
}

func (ctrl *ProductController) Delete(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	if err := ctrl.svc.Delete(c.Request.Context(), uint(id)); err != nil {
		c.Error(err); return
	}
	c.Status(http.StatusNoContent)
}