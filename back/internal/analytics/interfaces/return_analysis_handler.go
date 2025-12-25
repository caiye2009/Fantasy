// back/internal/analytics/interfaces/return_analysis_handler.go
package interfaces

import (
	"net/http"

	"back/pkg/endpoint"
	"back/internal/analytics/application"
	"github.com/gin-gonic/gin"
)

// ReturnAnalysisHandler 退货分析HTTP处理器
type ReturnAnalysisHandler struct {
	service *application.ReturnAnalysisService
}

// NewReturnAnalysisHandler 创建退货分析处理器
func NewReturnAnalysisHandler(service *application.ReturnAnalysisService) *ReturnAnalysisHandler {
	return &ReturnAnalysisHandler{
		service: service,
	}
}

// GetRoutes 返回路由定义
func (h *ReturnAnalysisHandler) GetRoutes() []endpoint.RouteDefinition {
	return []endpoint.RouteDefinition{
		{
			Method:  "GET",
			Path:    "/return-analysis/customers",
			Handler: h.GetCustomerList,
			Name:    "获取客户列表",
		},
		{
			Method:  "POST",
			Path:    "/return-analysis/analysis",
			Handler: h.GetReturnAnalysis,
			Name:    "获取退货分析",
		},
	}
}

// GetCustomerList 获取客户列表
// @Summary 获取客户下拉列表
// @Description 获取所有有完成订单的客户列表，用于前端下拉选择
// @Tags 退货分析
// @Accept json
// @Produce json
// @Success 200 {array} application.CustomerOptionResponse "客户列表"
// @Failure 500 {object} map[string]string "服务器错误"
// @Router /return-analysis/customers [get]
// @Security Bearer
func (h *ReturnAnalysisHandler) GetCustomerList(c *gin.Context) {
	resp, err := h.service.GetCustomerList(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}

// GetReturnAnalysis 获取退货分析
// @Summary 获取退货分析数据
// @Description 获取指定客户和时间范围内的退货分析数据，包括米数、重量、金额三个维度的统计。customerNo 为空时查询所有客户；dateRange 的 start 和 end 都不传时查询全部时间，必须同时传或同时不传
// @Tags 退货分析
// @Accept json
// @Produce json
// @Param request body application.ReturnAnalysisRequest true "查询参数"
// @Success 200 {object} application.ReturnAnalysisResponse "退货分析结果"
// @Failure 400 {object} map[string]string "参数错误"
// @Failure 500 {object} map[string]string "服务器错误"
// @Router /return-analysis/analysis [post]
// @Security Bearer
func (h *ReturnAnalysisHandler) GetReturnAnalysis(c *gin.Context) {
	var req application.ReturnAnalysisRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误: " + err.Error()})
		return
	}

	resp, err := h.service.GetReturnAnalysis(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}