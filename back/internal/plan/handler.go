package plan

import (
    "net/http"
    "strconv"

    "github.com/gin-gonic/gin"
)

type PlanHandler struct {
    planService *PlanService
}

func NewPlanHandler(planService *PlanService) *PlanHandler {
    return &PlanHandler{planService: planService}
}

func (h *PlanHandler) Create(c *gin.Context) {
    var p Plan
    if err := c.ShouldBindJSON(&p); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    if err := h.planService.Create(&p); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, p)
}

func (h *PlanHandler) Get(c *gin.Context) {
    id, _ := strconv.Atoi(c.Param("id"))
    p, err := h.planService.Get(uint(id))
    if err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
        return
    }
    c.JSON(http.StatusOK, p)
}

func (h *PlanHandler) List(c *gin.Context) {
    list, err := h.planService.List()
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    c.JSON(http.StatusOK, list)
}

func (h *PlanHandler) Update(c *gin.Context) {
    id, _ := strconv.Atoi(c.Param("id"))

    data := make(map[string]interface{})
    if err := c.ShouldBindJSON(&data); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    if err := h.planService.Update(uint(id), data); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    c.JSON(http.StatusOK, gin.H{"status": "updated"})
}

func (h *PlanHandler) Delete(c *gin.Context) {
    id, _ := strconv.Atoi(c.Param("id"))

    if err := h.planService.Delete(uint(id)); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    c.JSON(http.StatusOK, gin.H{"status": "deleted"})
}
