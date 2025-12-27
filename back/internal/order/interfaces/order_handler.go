package interfaces

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"back/pkg/audit"
	"back/pkg/endpoint"
	"back/internal/order/application"
	"back/internal/order/domain"
)

// OrderHandler 订单 Handler
type OrderHandler struct {
	service *application.OrderService
}

// NewOrderHandler 创建 Handler
func NewOrderHandler(service *application.OrderService) *OrderHandler {
	return &OrderHandler{service: service}
}

// getUserInfo 从上下文获取用户信息
func getUserInfo(c *gin.Context) (loginID uint, username string, role string, err error) {
	loginIDVal, exists := c.Get("loginId")
	if !exists {
		err = errors.New("未找到用户登录信息")
		return
	}

	usernameVal, exists := c.Get("username")
	if !exists {
		err = errors.New("未找到用户名信息")
		return
	}

	roleVal, exists := c.Get("role")
	if !exists {
		err = errors.New("未找到用户角色信息")
		return
	}

	// 转换 loginID
	loginIDStr, ok := loginIDVal.(string)
	if !ok {
		err = errors.New("用户登录信息格式错误")
		return
	}
	id, parseErr := strconv.Atoi(loginIDStr)
	if parseErr != nil {
		err = errors.New("用户登录信息格式错误")
		return
	}
	loginID = uint(id)

	// 转换 username
	username, ok = usernameVal.(string)
	if !ok {
		err = errors.New("用户名格式错误")
		return
	}

	// 转换 role
	role, ok = roleVal.(string)
	if !ok {
		err = errors.New("用户角色格式错误")
		return
	}

	return
}

// Create 创建订单
// @Summary      创建订单
// @Description  创建新的订单信息
// @Tags         订单管理
// @Accept       json
// @Produce      json
// @Param        request body application.CreateOrderRequest true "订单信息"
// @Success      200 {object} application.OrderResponse "创建成功"
// @Failure      400 {object} map[string]string "请求参数错误"
// @Failure      500 {object} map[string]string "服务器错误"
// @Security     Bearer
// @Router       /order [post]
func (h *OrderHandler) Create(c *gin.Context) {
	var req application.CreateOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	resp, err := h.service.Create(c.Request.Context(), &req)
	if err != nil {
		if errors.Is(err, domain.ErrOrderNoDuplicate) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "订单编号已存在"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, resp)
}

// Get 获取订单详情
// @Summary      获取订单详情
// @Description  根据订单ID获取订单详细信息
// @Tags         订单管理
// @Accept       json
// @Produce      json
// @Param        id path int true "订单ID"
// @Success      200 {object} application.OrderResponse "获取成功"
// @Failure      404 {object} map[string]string "订单不存在"
// @Security     Bearer
// @Router       /order/{id} [get]
func (h *OrderHandler) Get(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	
	resp, err := h.service.Get(c.Request.Context(), uint(id))
	if err != nil {
		if errors.Is(err, domain.ErrOrderNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "订单不存在"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, resp)
}

// List 获取订单列表
// @Summary      获取订单列表
// @Description  获取所有订单列表
// @Tags         订单管理
// @Accept       json
// @Produce      json
// @Param        limit query int false "每页数量" default(10)
// @Param        offset query int false "偏移量" default(0)
// @Success      200 {object} application.OrderListResponse "获取成功"
// @Failure      500 {object} map[string]string "服务器错误"
// @Security     Bearer
// @Router       /order [get]
func (h *OrderHandler) List(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))
	
	resp, err := h.service.List(c.Request.Context(), limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, resp)
}

// Update 更新订单
// @Summary      更新订单
// @Description  根据订单ID更新订单信息
// @Tags         订单管理
// @Accept       json
// @Produce      json
// @Param        id path int true "订单ID"
// @Param        request body application.UpdateOrderRequest true "更新的订单信息"
// @Success      200 {object} map[string]string "更新成功"
// @Failure      400 {object} map[string]string "请求参数错误"
// @Failure      500 {object} map[string]string "服务器错误"
// @Security     Bearer
// @Router       /order/{id} [post]
func (h *OrderHandler) Update(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	
	var req application.UpdateOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	if err := h.service.Update(c.Request.Context(), uint(id), &req); err != nil {
		if errors.Is(err, domain.ErrOrderNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "订单不存在"})
			return
		}
		if errors.Is(err, domain.ErrCannotConfirm) ||
			errors.Is(err, domain.ErrCannotStartProduction) ||
			errors.Is(err, domain.ErrCannotComplete) ||
			errors.Is(err, domain.ErrCannotUpdateCompleted) {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{"message": "更新成功"})
}

// Delete 删除订单
// @Summary      删除订单
// @Description  根据订单ID删除订单
// @Tags         订单管理
// @Accept       json
// @Produce      json
// @Param        id path int true "订单ID"
// @Success      200 {object} map[string]string "删除成功"
// @Failure      500 {object} map[string]string "服务器错误"
// @Security     Bearer
// @Router       /order/{id} [delete]
func (h *OrderHandler) Delete(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	
	if err := h.service.Delete(c.Request.Context(), uint(id)); err != nil {
		if errors.Is(err, domain.ErrOrderNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "订单不存在"})
			return
		}
		if errors.Is(err, domain.ErrCannotDeleteCompleted) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "已完成的订单不能删除"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{"message": "删除成功"})
}

// ==================== 新版协作工作流Handler ====================

// CreateV2 创建订单（新版，支持协作工作流）
// @Summary      创建订单（新版）
// @Description  创建新订单，记录参与者、初始化进度、记录事件
// @Tags         订单管理
// @Accept       json
// @Produce      json
// @Param        request body application.CreateOrderRequestV2 true "订单信息"
// @Success      200 {object} application.OrderDetailResponse "创建成功"
// @Failure      400 {object} map[string]string "请求参数错误"
// @Failure      500 {object} map[string]string "服务器错误"
// @Security     Bearer
// @Router       /order [post]
func (h *OrderHandler) CreateV2(c *gin.Context) {
	var req application.CreateOrderRequestV2
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 从上下文获取用户信息
	creatorID, username, role, err := getUserInfo(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	resp, err := h.service.CreateV2(c.Request.Context(), &req, creatorID, username, role)
	if err != nil {
		if errors.Is(err, domain.ErrOrderNoDuplicate) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "订单编号已存在"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}

// AssignDepartment 分配部门
// @Summary      分配部门
// @Description  生产总监分配订单到部门
// @Tags         订单管理
// @Accept       json
// @Produce      json
// @Param        id path int true "订单ID"
// @Param        request body application.AssignDepartmentRequest true "部门信息"
// @Success      200 {object} map[string]string "分配成功"
// @Failure      400 {object} map[string]string "请求参数错误"
// @Failure      500 {object} map[string]string "服务器错误"
// @Security     Bearer
// @Router       /order/{id}/assign-department [post]
func (h *OrderHandler) AssignDepartment(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	var req application.AssignDepartmentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// ===== Audit: 记录旧值 =====
	recorder := audit.Get(c)
	if recorder != nil {
		recorder.SetResourceID(id) // 设置被操作的资源ID

		// 获取旧值（分配前的部门信息）
		oldOrder, err := h.service.Get(c.Request.Context(), uint(id))
		if err == nil {
			recorder.SetOld(map[string]interface{}{
				"order_id":            oldOrder.ID,
				"order_no":            oldOrder.OrderNo,
				"assigned_department": "", // 旧的部门为空
			})
		}
	}
	// ===== Audit End =====

	// 从上下文获取用户信息
	directorID, directorName, role, err := getUserInfo(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	if err := h.service.AssignDepartment(c.Request.Context(), uint(id), &req, directorID, directorName, role); err != nil {
		if errors.Is(err, domain.ErrOrderNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "订单不存在"})
			return
		}
		if errors.Is(err, domain.ErrInvalidOrderStatus) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "订单状态不允许此操作"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// ===== Audit: 记录新值（分配后的部门信息） =====
	if recorder != nil {
		recorder.SetNew(map[string]interface{}{
			"order_id":            id,
			"assigned_department": req.Department,
		})
	}
	// ===== Audit End =====

	c.JSON(http.StatusOK, gin.H{"message": "分配部门成功"})
}

// AssignPersonnel 分配人员
// @Summary      分配人员
// @Description  生产助理分配生产专员和跟单，并设定胚布目标数量
// @Tags         订单管理
// @Accept       json
// @Produce      json
// @Param        id path int true "订单ID"
// @Param        request body application.AssignPersonnelRequest true "人员信息"
// @Success      200 {object} map[string]string "分配成功"
// @Failure      400 {object} map[string]string "请求参数错误"
// @Failure      500 {object} map[string]string "服务器错误"
// @Security     Bearer
// @Router       /order/{id}/assign-personnel [post]
func (h *OrderHandler) AssignPersonnel(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	var req application.AssignPersonnelRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 从上下文获取用户信息
	assistantID, assistantName, role, err := getUserInfo(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	if err := h.service.AssignPersonnel(c.Request.Context(), uint(id), &req, assistantID, assistantName, role); err != nil {
		if errors.Is(err, domain.ErrOrderNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "订单不存在"})
			return
		}
		if errors.Is(err, domain.ErrInvalidOrderStatus) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "订单状态不允许此操作"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "分配人员成功"})
}

// UpdateFabricInput 更新胚布投入进度
// @Summary      更新胚布投入进度
// @Description  跟单更新胚布投入进度
// @Tags         订单管理
// @Accept       json
// @Produce      json
// @Param        id path int true "订单ID"
// @Param        request body application.UpdateProgressRequest true "进度信息"
// @Success      200 {object} map[string]string "更新成功"
// @Failure      400 {object} map[string]string "请求参数错误"
// @Failure      500 {object} map[string]string "服务器错误"
// @Security     Bearer
// @Router       /order/{id}/progress/fabric-input [post]
func (h *OrderHandler) UpdateFabricInput(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	var req application.UpdateProgressRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 从上下文获取用户信息
	operatorID, operatorName, role, err := getUserInfo(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	if err := h.service.UpdateProgress(c.Request.Context(), uint(id), domain.ProgressTypeFabricInput, &req, operatorID, operatorName, role); err != nil {
		if errors.Is(err, domain.ErrProgressNotExists) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "该进度项不存在"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "更新胚布投入进度成功"})
}

// UpdateProduction 更新生产进度
// @Summary      更新生产进度
// @Description  跟单更新生产进度
// @Tags         订单管理
// @Accept       json
// @Produce      json
// @Param        id path int true "订单ID"
// @Param        request body application.UpdateProgressRequest true "进度信息"
// @Success      200 {object} map[string]string "更新成功"
// @Failure      400 {object} map[string]string "请求参数错误"
// @Failure      500 {object} map[string]string "服务器错误"
// @Security     Bearer
// @Router       /order/{id}/progress/production [post]
func (h *OrderHandler) UpdateProduction(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	var req application.UpdateProgressRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 从上下文获取用户信息
	operatorID, operatorName, role, err := getUserInfo(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	if err := h.service.UpdateProgress(c.Request.Context(), uint(id), domain.ProgressTypeProduction, &req, operatorID, operatorName, role); err != nil {
		if errors.Is(err, domain.ErrProgressNotExists) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "该进度项不存在"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "更新生产进度成功"})
}

// UpdateWarehouseCheck 更新验货进度
// @Summary      更新验货进度
// @Description  仓管更新验货进度
// @Tags         订单管理
// @Accept       json
// @Produce      json
// @Param        id path int true "订单ID"
// @Param        request body application.UpdateProgressRequest true "进度信息"
// @Success      200 {object} map[string]string "更新成功"
// @Failure      400 {object} map[string]string "请求参数错误"
// @Failure      500 {object} map[string]string "服务器错误"
// @Security     Bearer
// @Router       /order/{id}/progress/warehouse-check [post]
func (h *OrderHandler) UpdateWarehouseCheck(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	var req application.UpdateProgressRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 从上下文获取用户信息
	operatorID, operatorName, role, err := getUserInfo(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	if err := h.service.UpdateProgress(c.Request.Context(), uint(id), domain.ProgressTypeWarehouseCheck, &req, operatorID, operatorName, role); err != nil {
		if errors.Is(err, domain.ErrProgressNotExists) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "该进度项不存在"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "更新验货进度成功"})
}

// UpdateRework 更新回修进度
// @Summary      更新回修进度
// @Description  跟单更新回修进度
// @Tags         订单管理
// @Accept       json
// @Produce      json
// @Param        id path int true "订单ID"
// @Param        request body application.UpdateProgressRequest true "进度信息"
// @Success      200 {object} map[string]string "更新成功"
// @Failure      400 {object} map[string]string "请求参数错误"
// @Failure      500 {object} map[string]string "服务器错误"
// @Security     Bearer
// @Router       /order/{id}/progress/rework [post]
func (h *OrderHandler) UpdateRework(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	var req application.UpdateProgressRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 从上下文获取用户信息
	operatorID, operatorName, role, err := getUserInfo(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	if err := h.service.UpdateProgress(c.Request.Context(), uint(id), domain.ProgressTypeRework, &req, operatorID, operatorName, role); err != nil {
		if errors.Is(err, domain.ErrProgressNotExists) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "该进度项不存在"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "更新回修进度成功"})
}

// AddDefect 录入次品
// @Summary      录入次品
// @Description  仓管录入次品数量
// @Tags         订单管理
// @Accept       json
// @Produce      json
// @Param        id path int true "订单ID"
// @Param        request body application.AddDefectRequest true "次品信息"
// @Success      200 {object} map[string]string "录入成功"
// @Failure      400 {object} map[string]string "请求参数错误"
// @Failure      500 {object} map[string]string "服务器错误"
// @Security     Bearer
// @Router       /order/{id}/defect [post]
func (h *OrderHandler) AddDefect(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	var req application.AddDefectRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 从上下文获取用户信息
	warehouseID, warehouseName, role, err := getUserInfo(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	if err := h.service.AddDefect(c.Request.Context(), uint(id), &req, warehouseID, warehouseName, role); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "录入次品成功"})
}

// GetDetail 获取订单详情
// @Summary      获取订单详情（完整）
// @Description  获取订单详情，含参与者、进度、事件流
// @Tags         订单管理
// @Accept       json
// @Produce      json
// @Param        id path int true "订单ID"
// @Success      200 {object} application.OrderDetailResponse "获取成功"
// @Failure      404 {object} map[string]string "订单不存在"
// @Failure      500 {object} map[string]string "服务器错误"
// @Security     Bearer
// @Router       /order/{id}/detail [get]
func (h *OrderHandler) GetDetail(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	// 从上下文获取用户信息
	loginID, _ := c.Get("loginId")
	userID, _ := strconv.Atoi(loginID.(string))

	resp, err := h.service.GetDetail(c.Request.Context(), uint(id), uint(userID))
	if err != nil {
		if errors.Is(err, domain.ErrOrderNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "订单不存在"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}

// GetEvents 获取订单事件流
// @Summary      获取订单事件流
// @Description  获取订单的完整操作历史事件流
// @Tags         订单管理
// @Accept       json
// @Produce      json
// @Param        id path int true "订单ID"
// @Success      200 {array} application.EventResponse "获取成功"
// @Failure      500 {object} map[string]string "服务器错误"
// @Security     Bearer
// @Router       /order/{id}/events [get]
func (h *OrderHandler) GetEvents(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	resp, err := h.service.GetEvents(c.Request.Context(), uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}

// ListWithDetail 获取订单列表（含完整详情）
// @Summary      获取订单列表（含完整详情）
// @Description  获取订单列表，包含客户名、产品名、进度项、操作日志等完整信息
// @Tags         订单管理
// @Accept       json
// @Produce      json
// @Param        limit query int false "每页数量" default(10)
// @Param        offset query int false "偏移量" default(0)
// @Success      200 {object} application.OrderListDetailResponse "获取成功"
// @Failure      500 {object} map[string]string "服务器错误"
// @Security     Bearer
// @Router       /order/list-detail [get]
func (h *OrderHandler) ListWithDetail(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	resp, err := h.service.ListWithDetail(c.Request.Context(), limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}

// GetRoutes 返回路由定义
func (h *OrderHandler) GetRoutes() []endpoint.RouteDefinition {
	return []endpoint.RouteDefinition{
		// 基本CRUD
		{Method: "POST", Path: "/order", Handler: h.Create, Domain: "order", Action: "create"},
		{Method: "GET", Path: "/order/:id", Handler: h.Get, Domain: "", Action: ""},
		{Method: "GET", Path: "/order", Handler: h.List, Domain: "", Action: ""},
		{Method: "PUT", Path: "/order/:id", Handler: h.Update, Domain: "order", Action: "update"},
		{Method: "DELETE", Path: "/order/:id", Handler: h.Delete, Domain: "order", Action: "delete"},
		// 工作流端点
		{Method: "POST", Path: "/order/:id/assign-department", Handler: h.AssignDepartment, Domain: "order", Action: "assignDepartment"},
		{Method: "POST", Path: "/order/:id/assign-personnel", Handler: h.AssignPersonnel, Domain: "order", Action: "assignPersonnel"},
		{Method: "POST", Path: "/order/:id/progress/fabric-input", Handler: h.UpdateFabricInput, Domain: "order", Action: "fabricInput"},
		{Method: "POST", Path: "/order/:id/progress/production", Handler: h.UpdateProduction, Domain: "order", Action: "production"},
		{Method: "POST", Path: "/order/:id/progress/warehouse-check", Handler: h.UpdateWarehouseCheck, Domain: "order", Action: "warehouseCheck"},
		{Method: "POST", Path: "/order/:id/progress/rework", Handler: h.UpdateRework, Domain: "order", Action: "rework"},
		{Method: "POST", Path: "/order/:id/defect", Handler: h.AddDefect, Domain: "order", Action: "defect"},
		// 详情端点
		{Method: "GET", Path: "/order/list-detail", Handler: h.ListWithDetail, Domain: "", Action: ""},
		{Method: "GET", Path: "/order/:id/detail", Handler: h.GetDetail, Domain: "", Action: ""},
		{Method: "GET", Path: "/order/:id/events", Handler: h.GetEvents, Domain: "", Action: ""},
	}
}