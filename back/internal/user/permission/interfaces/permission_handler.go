package interfaces

import (
	"net/http"

	"back/pkg/endpoint"
	"back/internal/user/permission/application"
	"github.com/gin-gonic/gin"
)

// PermissionHandler 权限管理 Handler
type PermissionHandler struct {
	service *application.PermissionService
}

// NewPermissionHandler 创建 Handler
func NewPermissionHandler(service *application.PermissionService) *PermissionHandler {
	return &PermissionHandler{service: service}
}

// AddUserPermission 给用户添加权限
// @Summary      给用户添加权限
// @Description  给指定用户添加个性化权限
// @Tags         权限管理
// @Accept       json
// @Produce      json
// @Param        request body application.AddUserPermissionRequest true "请求参数"
// @Success      200 {object} map[string]interface{} "添加成功"
// @Failure      400 {object} map[string]interface{} "参数错误"
// @Security     Bearer
// @Router       /permission/user/add [post]
func (h *PermissionHandler) AddUserPermission(c *gin.Context) {
	var req application.AddUserPermissionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
		return
	}

	if err := h.service.AddUserPermission(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "权限添加成功",
		"data": gin.H{
			"loginId":    req.LoginID,
			"permission": req.Permission,
		},
	})
}

// RemoveUserPermission 删除用户权限
// @Summary      删除用户权限
// @Description  删除指定用户的个性化权限
// @Tags         权限管理
// @Accept       json
// @Produce      json
// @Param        request body application.RemoveUserPermissionRequest true "请求参数"
// @Success      200 {object} map[string]interface{} "删除成功"
// @Failure      400 {object} map[string]interface{} "参数错误"
// @Security     Bearer
// @Router       /permission/user/remove [post]
func (h *PermissionHandler) RemoveUserPermission(c *gin.Context) {
	var req application.RemoveUserPermissionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
		return
	}

	if err := h.service.RemoveUserPermission(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "权限删除成功",
		"data": gin.H{
			"loginId":    req.LoginID,
			"permission": req.Permission,
		},
	})
}

// GetUserPermissions 获取用户权限列表
// @Summary      获取用户权限列表
// @Description  获取指定用户的所有个性化权限
// @Tags         权限管理
// @Accept       json
// @Produce      json
// @Param        loginId query string true "用户 LoginID"
// @Success      200 {object} application.UserPermissionsResponse "权限列表"
// @Failure      400 {object} map[string]interface{} "参数错误"
// @Security     Bearer
// @Router       /permission/user [get]
func (h *PermissionHandler) GetUserPermissions(c *gin.Context) {
	loginID := c.Query("loginId")
	if loginID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "loginId 不能为空"})
		return
	}

	resp, err := h.service.GetUserPermissions(&application.GetUserPermissionsRequest{
		LoginID: loginID,
	})
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}

// ListAllPermissions 列出所有可用权限
// @Summary      列出所有可用权限
// @Description  列出系统中所有可分配的权限
// @Tags         权限管理
// @Accept       json
// @Produce      json
// @Success      200 {object} application.PermissionsListResponse "权限列表"
// @Security     Bearer
// @Router       /permission/list [get]
func (h *PermissionHandler) ListAllPermissions(c *gin.Context) {
	resp := h.service.ListAllPermissions()
	c.JSON(http.StatusOK, resp)
}

// ListPermissionsByDomain 按域分组列出权限
// @Summary      按域分组列出权限
// @Description  列出系统中所有权限，按业务域分组
// @Tags         权限管理
// @Accept       json
// @Produce      json
// @Success      200 {object} application.PermissionsByDomainResponse "权限列表"
// @Security     Bearer
// @Router       /permission/list-by-domain [get]
func (h *PermissionHandler) ListPermissionsByDomain(c *gin.Context) {
	resp := h.service.ListPermissionsByDomain()
	c.JSON(http.StatusOK, resp)
}

// GetRoutes 返回路由定义
func (h *PermissionHandler) GetRoutes() []endpoint.RouteDefinition {
	return []endpoint.RouteDefinition{
		{Method: "POST", Path: "/permission/user/add", Handler: h.AddUserPermission, Domain: "", Action: ""},
		{Method: "POST", Path: "/permission/user/remove", Handler: h.RemoveUserPermission, Domain: "", Action: ""},
		{Method: "GET", Path: "/permission/user", Handler: h.GetUserPermissions, Domain: "", Action: ""},
		{Method: "GET", Path: "/permission/list", Handler: h.ListAllPermissions, Domain: "", Action: ""},
		{Method: "GET", Path: "/permission/list-by-domain", Handler: h.ListPermissionsByDomain, Domain: "", Action: ""},
	}
}
