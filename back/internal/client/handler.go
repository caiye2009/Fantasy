package client

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type ClientHandler struct {
	clientService *ClientService
}

func NewClientHandler(clientService *ClientService) *ClientHandler {
	return &ClientHandler{clientService: clientService}
}

// Create godoc
// @Summary      创建客户
// @Description  创建新的客户信息
// @Tags         客户管理
// @Accept       json
// @Produce      json
// @Param        request body CreateClientRequest true "客户信息"
// @Success      200 {object} Client "创建成功"
// @Failure      400 {object} map[string]string "请求参数错误"
// @Failure      500 {object} map[string]string "服务器错误"
// @Security     Bearer
// @Router       /client [post]
func (h *ClientHandler) Create(c *gin.Context) {
	var req CreateClientRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	client, err := h.clientService.Create(&req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, client)
}

// Get godoc
// @Summary      获取客户详情
// @Description  根据客户ID获取客户详细信息
// @Tags         客户管理
// @Accept       json
// @Produce      json
// @Param        id path int true "客户ID"
// @Success      200 {object} Client "获取成功"
// @Failure      404 {object} map[string]string "客户不存在"
// @Security     Bearer
// @Router       /client/{id} [get]
func (h *ClientHandler) Get(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	client, err := h.clientService.Get(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}
	c.JSON(http.StatusOK, client)
}

// List godoc
// @Summary      获取客户列表
// @Description  获取所有客户列表
// @Tags         客户管理
// @Accept       json
// @Produce      json
// @Success      200 {array} Client "获取成功"
// @Failure      500 {object} map[string]string "服务器错误"
// @Security     Bearer
// @Router       /client [get]
func (h *ClientHandler) List(c *gin.Context) {
	list, err := h.clientService.List()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, list)
}

// Update godoc
// @Summary      更新客户信息
// @Description  根据客户ID更新客户信息
// @Tags         客户管理
// @Accept       json
// @Produce      json
// @Param        id path int true "客户ID"
// @Param        request body UpdateClientRequest true "更新的客户信息"
// @Success      200 {object} map[string]string "更新成功"
// @Failure      400 {object} map[string]string "请求参数错误"
// @Failure      500 {object} map[string]string "服务器错误"
// @Security     Bearer
// @Router       /client/{id} [put]
func (h *ClientHandler) Update(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	var req UpdateClientRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.clientService.Update(uint(id), &req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "updated"})
}

// Delete godoc
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

	if err := h.clientService.Delete(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "deleted"})
}