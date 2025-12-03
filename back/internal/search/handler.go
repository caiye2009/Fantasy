package search

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type SearchHandler struct {
	searchService *SearchService
}

func NewSearchHandler(searchService *SearchService) *SearchHandler {
	return &SearchHandler{searchService: searchService}
}

// Search godoc
// @Summary      通用搜索
// @Description  执行跨多个索引的通用搜索
// @Tags         搜索
// @Accept       json
// @Produce      json
// @Param        request body SearchRequest true "搜索请求参数"
// @Success      200 {object} object "搜索结果"
// @Failure      400 {object} map[string]string "请求参数错误"
// @Failure      500 {object} map[string]string "服务器错误"
// @Security     Bearer
// @Router       /search [post]
func (h *SearchHandler) Search(c *gin.Context) {
	var req SearchRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	results, err := h.searchService.Search(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, results)
}

// SearchByIndex godoc
// @Summary      按索引搜索
// @Description  在指定索引中执行搜索
// @Tags         搜索
// @Accept       json
// @Produce      json
// @Param        index path string true "索引名称" example:"products"
// @Param        request body SearchRequest true "搜索请求参数"
// @Success      200 {object} object "搜索结果"
// @Failure      400 {object} map[string]string "请求参数错误"
// @Failure      500 {object} map[string]string "服务器错误"
// @Security     Bearer
// @Router       /search/{index} [post]
func (h *SearchHandler) SearchByIndex(c *gin.Context) {
	index := c.Param("index")

	var req SearchRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 强制指定索引
	req.Indices = []string{index}

	results, err := h.searchService.Search(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, results)
}

// GetIndices godoc
// @Summary      获取所有可用索引
// @Description  获取系统中所有可搜索的索引列表
// @Tags         搜索
// @Accept       json
// @Produce      json
// @Success      200 {object} map[string]interface{} "索引列表"
// @Security     Bearer
// @Router       /search/indices [get]
func (h *SearchHandler) GetIndices(c *gin.Context) {
	indices := GetAllIndices()
	
	indexList := make([]map[string]interface{}, 0, len(indices))
	for _, idx := range indices {
		if meta, ok := IndexConfig[idx]; ok {
			indexList = append(indexList, map[string]interface{}{
				"name":   meta.Name,
				"type":   meta.Type,
				"fields": meta.DefaultFields,
			})
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"indices": indexList,
	})
}