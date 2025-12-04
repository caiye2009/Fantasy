package interfaces

import (
	"errors"
	"net/http"
	"strconv"
	
	"github.com/gin-gonic/gin"
	
	"back/internal/product/application"
	"back/internal/product/domain"
)

// ProductHandler 产品 Handler
type ProductHandler struct {
	service    *application.ProductService
	calculator *application.CostCalculator
}

// NewProductHandler 创建 Handler
func NewProductHandler(service *application.ProductService, calculator *application.CostCalculator) *ProductHandler {
	return &ProductHandler{
		service:    service,
		calculator: calculator,
	}
}

// Create 创建产品
// @Summary      创建产品
// @Description  创建新的产品信息
// @Tags         产品管理
// @Accept       json
// @Produce      json
// @Param        request body application.CreateProductRequest true "产品信息"
// @Success      200 {object} application.ProductResponse "创建成功"
// @Failure      400 {object} map[string]string "请求参数错误"
// @Failure      500 {object} map[string]string "服务器错误"
// @Security     Bearer
// @Router       /product [post]
func (h *ProductHandler) Create(c *gin.Context) {
	var req application.CreateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	resp, err := h.service.Create(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, resp)
}

// Get 获取产品详情
// @Summary      获取产品详情
// @Description  根据产品ID获取产品详细信息
// @Tags         产品管理
// @Accept       json
// @Produce      json
// @Param        id path int true "产品ID"
// @Success      200 {object} application.ProductResponse "获取成功"
// @Failure      404 {object} map[string]string "产品不存在"
// @Security     Bearer
// @Router       /product/{id} [get]
func (h *ProductHandler) Get(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	
	resp, err := h.service.Get(c.Request.Context(), uint(id))
	if err != nil {
		if errors.Is(err, domain.ErrProductNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "产品不存在"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, resp)
}

// List 获取产品列表
// @Summary      获取产品列表
// @Description  获取所有产品列表
// @Tags         产品管理
// @Accept       json
// @Produce      json
// @Param        limit query int false "每页数量" default(10)
// @Param        offset query int false "偏移量" default(0)
// @Success      200 {object} application.ProductListResponse "获取成功"
// @Failure      500 {object} map[string]string "服务器错误"
// @Security     Bearer
// @Router       /product [get]
func (h *ProductHandler) List(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))
	
	resp, err := h.service.List(c.Request.Context(), limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, resp)
}

// Update 更新产品信息
// @Summary      更新产品信息
// @Description  根据产品ID更新产品信息
// @Tags         产品管理
// @Accept       json
// @Produce      json
// @Param        id path int true "产品ID"
// @Param        request body application.UpdateProductRequest true "更新的产品信息"
// @Success      200 {object} map[string]string "更新成功"
// @Failure      400 {object} map[string]string "请求参数错误"
// @Failure      500 {object} map[string]string "服务器错误"
// @Security     Bearer
// @Router       /product/{id} [put]
func (h *ProductHandler) Update(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	
	var req application.UpdateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	if err := h.service.Update(c.Request.Context(), uint(id), &req); err != nil {
		if errors.Is(err, domain.ErrProductNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "产品不存在"})
			return
		}
		if errors.Is(err, domain.ErrCannotUpdateApprovedProduct) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "已审批的产品不能修改"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{"message": "更新成功"})
}

// Delete 删除产品
// @Summary      删除产品
// @Description  根据产品ID删除产品
// @Tags         产品管理
// @Accept       json
// @Produce      json
// @Param        id path int true "产品ID"
// @Success      200 {object} map[string]string "删除成功"
// @Failure      500 {object} map[string]string "服务器错误"
// @Security     Bearer
// @Router       /product/{id} [delete]
func (h *ProductHandler) Delete(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	
	if err := h.service.Delete(c.Request.Context(), uint(id)); err != nil {
		if errors.Is(err, domain.ErrProductNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "产品不存在"})
			return
		}
		if errors.Is(err, domain.ErrCannotDeleteApprovedProduct) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "已审批的产品不能删除"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{"message": "删除成功"})
}

// CalculateCost 计算产品成本
// @Summary      计算产品成本
// @Description  根据产品ID和数量计算产品成本
// @Tags         产品管理
// @Accept       json
// @Produce      json
// @Param        request body application.CalculateCostRequest true "成本计算参数"
// @Success      200 {object} domain.CostResult "成本计算结果"
// @Failure      400 {object} map[string]string "请求参数错误"
// @Failure      500 {object} map[string]string "服务器错误"
// @Security     Bearer
// @Router       /product/calculate-cost [post]
func (h *ProductHandler) CalculateCost(c *gin.Context) {
	var req application.CalculateCostRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	result, err := h.calculator.Calculate(c.Request.Context(), &req)
	if err != nil {
		if errors.Is(err, domain.ErrProductNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "产品不存在"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, result)
}

// RegisterProductHandlers 注册路由
func RegisterProductHandlers(rg *gin.RouterGroup, service *application.ProductService, calculator *application.CostCalculator) {
	handler := NewProductHandler(service, calculator)
	
	rg.POST("/product", handler.Create)
	rg.GET("/product/:id", handler.Get)
	rg.GET("/product", handler.List)
	rg.PUT("/product/:id", handler.Update)
	rg.DELETE("/product/:id", handler.Delete)
	rg.POST("/product/calculate-cost", handler.CalculateCost)
}