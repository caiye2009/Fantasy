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

// Search 通用搜索接口
// POST /api/v1/search
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

// SearchByIndex 按索引搜索
// POST /api/v1/search/:index
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

// GetIndices 获取所有可用索引
// GET /api/v1/search/indices
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