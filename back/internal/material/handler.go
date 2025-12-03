package material

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type MaterialHandler struct {
	materialService *MaterialService
}

func NewMaterialHandler(materialService *MaterialService) *MaterialHandler {
	return &MaterialHandler{materialService: materialService}
}

func (h *MaterialHandler) Create(c *gin.Context) {
	var req CreateMaterialRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	material, err := h.materialService.Create(&req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, material)
}

func (h *MaterialHandler) Get(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	m, err := h.materialService.Get(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}
	c.JSON(http.StatusOK, m)
}

func (h *MaterialHandler) List(c *gin.Context) {
	list, err := h.materialService.List()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, list)
}

func (h *MaterialHandler) Update(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	var req UpdateMaterialRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.materialService.Update(uint(id), &req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "updated"})
}

func (h *MaterialHandler) Delete(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	if err := h.materialService.Delete(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "deleted"})
}