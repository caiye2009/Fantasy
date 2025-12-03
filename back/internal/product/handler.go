package product

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type ProductHandler struct {
	productService *ProductService
}

func NewProductHandler(productService *ProductService) *ProductHandler {
	return &ProductHandler{productService: productService}
}

// Create godoc
// @Summary      创建产品
// @Description  创建新的产品信息
// @Tags         产品管理
// @Accept       json
// @Produce      json
// @Param        request body Product true "产品信息"
// @Success      200 {object} Product "创建成功"
// @Failure      400 {object} map[string]string "请求参数错误"
// @Failure      500 {object} map[string]string "服务器错误"
// @Security     Bearer
// @Router       /product [post]
func (h *ProductHandler) Create(c *gin.Context) {
	var product Product
	if err := c.ShouldBindJSON(&product); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.productService.Create(&product); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, product)
}

// Get godoc
// @Summary      获取产品详情
// @Description  根据产品ID获取产品详细信息
// @Tags         产品管理
// @Accept       json
// @Produce      json
// @Param        id path int true "产品ID"
// @Success      200 {object} Product "获取成功"
// @Failure      404 {object} map[string]string "产品不存在"
// @Security     Bearer
// @Router       /product/{id} [get]
func (h *ProductHandler) Get(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	product, err := h.productService.Get(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}
	c.JSON(http.StatusOK, product)
}

// List godoc
// @Summary      获取产品列表
// @Description  获取所有产品列表
// @Tags         产品管理
// @Accept       json
// @Produce      json
// @Success      200 {array} Product "获取成功"
// @Failure      500 {object} map[string]string "服务器错误"
// @Security     Bearer
// @Router       /product [get]
func (h *ProductHandler) List(c *gin.Context) {
	list, err := h.productService.List()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, list)
}

// Update godoc
// @Summary      更新产品信息
// @Description  根据产品ID更新产品信息
// @Tags         产品管理
// @Accept       json
// @Produce      json
// @Param        id path int true "产品ID"
// @Param        request body map[string]interface{} true "更新的产品信息"
// @Success      200 {object} map[string]string "更新成功"
// @Failure      400 {object} map[string]string "请求参数错误"
// @Failure      500 {object} map[string]string "服务器错误"
// @Security     Bearer
// @Router       /product/{id} [put]
func (h *ProductHandler) Update(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	data := make(map[string]interface{})
	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.productService.Update(uint(id), data); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "updated"})
}

// Delete godoc
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

	if err := h.productService.Delete(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "deleted"})
}

// CalculateCost godoc
// @Summary      计算产品成本
// @Description  根据产品ID和数量计算产品成本
// @Tags         产品管理
// @Accept       json
// @Produce      json
// @Param        request body object{product_id=uint,quantity=float64,use_min_price=bool} true "成本计算参数"
// @Success      200 {object} object "成本计算结果"
// @Failure      400 {object} map[string]string "请求参数错误"
// @Failure      500 {object} map[string]string "服务器错误"
// @Security     Bearer
// @Router       /product/calculate-cost [post]
func (h *ProductHandler) CalculateCost(c *gin.Context) {
	var req struct {
		ProductID   uint    `json:"product_id" binding:"required" example:"1"`
		Quantity    float64 `json:"quantity" binding:"required,gt=0" example:"100"`
		UseMinPrice bool    `json:"use_min_price" example:"true"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := h.productService.CalculateCost(req.ProductID, req.Quantity, req.UseMinPrice)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}