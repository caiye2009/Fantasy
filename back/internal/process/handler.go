package process

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type ProcessHandler struct {
	processService *ProcessService
}

func NewProcessHandler(processService *ProcessService) *ProcessHandler {
	return &ProcessHandler{processService: processService}
}

func (h *ProcessHandler) Create(c *gin.Context) {
	var req CreateProcessRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	process, err := h.processService.Create(&req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, process)
}

func (h *ProcessHandler) Get(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	p, err := h.processService.Get(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}
	c.JSON(http.StatusOK, p)
}

func (h *ProcessHandler) List(c *gin.Context) {
	list, err := h.processService.List()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, list)
}

func (h *ProcessHandler) Update(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	var req UpdateProcessRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.processService.Update(uint(id), &req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "updated"})
}

func (h *ProcessHandler) Delete(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	if err := h.processService.Delete(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "deleted"})
}