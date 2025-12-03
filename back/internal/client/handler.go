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

func (h *ClientHandler) Get(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	client, err := h.clientService.Get(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}
	c.JSON(http.StatusOK, client)
}

func (h *ClientHandler) List(c *gin.Context) {
	list, err := h.clientService.List()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, list)
}

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

func (h *ClientHandler) Delete(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	if err := h.clientService.Delete(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "deleted"})
}