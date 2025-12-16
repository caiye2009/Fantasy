package interfaces

import (
	"errors"
	"net/http"
	"strconv"
	
	"github.com/gin-gonic/gin"
	
	"back/internal/client/application"
	"back/internal/client/domain"
)

// ClientHandler 客户 Handler
type ClientHandler struct {
	service *application.ClientService
}

// NewClientHandler 创建 Handler
func NewClientHandler(service *application.ClientService) *ClientHandler {
	return &ClientHandler{service: service}
}

// Create 创建客户
// @Summary      创建客户
// @Description  创建新的客户信息
// @Tags         客户管理
// @Accept       json
// @Produce      json
// @Param        request body application.CreateClientRequest true "客户信息"
// @Success      200 {object} application.ClientResponse "创建成功"
// @Failure      400 {object} map[string]string "请求参数错误"
// @Failure      500 {object} map[string]string "服务器错误"
// @Security     Bearer
// @Router       /client [post]
func (h *ClientHandler) Create(c *gin.Context) {
	var req application.CreateClientRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	resp, err := h.service.Create(c.Request.Context(), &req)
	if err != nil {
		if errors.Is(err, domain.ErrClientNameExists) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "客户名称已存在"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, resp)
}

// Get 获取客户详情
// @Summary      获取客户详情
// @Description  根据客户ID获取客户详细信息
// @Tags         客户管理
// @Accept       json
// @Produce      json
// @Param        id path int true "客户ID"
// @Success      200 {object} application.ClientResponse "获取成功"
// @Failure      404 {object} map[string]string "客户不存在"
// @Security     Bearer
// @Router       /client/{id} [get]
func (h *ClientHandler) Get(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	
	resp, err := h.service.Get(c.Request.Context(), uint(id))
	if err != nil {
		if errors.Is(err, domain.ErrClientNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "客户不存在"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, resp)
}


// func (h *ClientHandler) List(c *gin.Context) {
// 	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
// 	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))
	
// 	resp, err := h.service.List(c.Request.Context(), limit, offset)
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
// 		return
// 	}
	
// 	c.JSON(http.StatusOK, resp)
// }

// Update 更新客户信息
// @Summary      更新客户信息
// @Description  根据客户ID更新客户信息
// @Tags         客户管理
// @Accept       json
// @Produce      json
// @Param        id path int true "客户ID"
// @Param        request body application.UpdateClientRequest true "更新的客户信息"
// @Success      200 {object} map[string]string "更新成功"
// @Failure      400 {object} map[string]string "请求参数错误"
// @Failure      500 {object} map[string]string "服务器错误"
// @Security     Bearer
// @Router       /client/{id} [post]
func (h *ClientHandler) Update(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	
	var req application.UpdateClientRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	if err := h.service.Update(c.Request.Context(), uint(id), &req); err != nil {
		if errors.Is(err, domain.ErrClientNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "客户不存在"})
			return
		}
		if errors.Is(err, domain.ErrClientNameExists) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "客户名称已存在"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{"message": "更新成功"})
}

// Delete 删除客户
// @Summary      删除客户
// @Description  根据客户ID删除客户
// @Tags         客户管理
// @Accept       json
// @Produce      json
// @Param        id path int true "客户ID"
// @Success      200 {object} map[string]string "删除成功"
// @Failure      500 {object} map[string]string "服务器错误"
// @Security     Bearer
// @Router       /client/{id} [delete]
func (h *ClientHandler) Delete(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	
	if err := h.service.Delete(c.Request.Context(), uint(id)); err != nil {
		if errors.Is(err, domain.ErrClientNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "客户不存在"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{"message": "删除成功"})
}

// RegisterClientHandlers 注册路由
func RegisterClientHandlers(rg *gin.RouterGroup, service *application.ClientService) {
	handler := NewClientHandler(service)
	
	rg.POST("/client", handler.Create)
	rg.GET("/client/:id", handler.Get)
	//rg.GET("/client", handler.List)
	rg.POST("/client/:id", handler.Update)
	rg.DELETE("/client/:id", handler.Delete)
}