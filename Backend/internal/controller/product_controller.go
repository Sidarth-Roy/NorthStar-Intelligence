package controller

import (
	"net/http"
	"strconv"

	"github.com/Sidarth-Roy/NorthStar-Intelligence/Backend/internal/dto"
	"github.com/Sidarth-Roy/NorthStar-Intelligence/Backend/internal/service"
	"github.com/gin-gonic/gin"
)

type ProductController struct {
	svc *service.ProductService
}

func NewProductController(svc *service.ProductService) *ProductController {
	return &ProductController{svc: svc}
}

func (pc *ProductController) Create(c *gin.Context) {
	var req dto.ProductCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	res, err := pc.svc.Create(req)
	if err != nil {
		c.Error(err) // Passes to Global Exception Handler
		return
	}

	c.Header("X-Content-Type-Options", "nosniff")
	c.JSON(http.StatusCreated, res)
}

func (pc *ProductController) GetByID(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	res, err := pc.svc.GetByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		return
	}

	c.Header("Cache-Control", "public, max-age=3600")
	c.JSON(http.StatusOK, res)
}