package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/guidiguidi/RateMonitorBC/internal/bestchange"
	"github.com/guidiguidi/RateMonitorBC/internal/models"
)

type Handler struct {
	bc *bestchange.Client
}

func NewHandler(bc *bestchange.Client) *Handler {
	return &Handler{bc: bc}
}

func (h *Handler) GetBestExchange(c *gin.Context) {
	var req models.BestRateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	best, err := h.bc.GetBestRateWithFilters(
		c.Request.Context(),
		req.FromID,
		req.ToID,
		req.Amount,
		req.Marks,
	)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resp := models.BestRateResponse{
		Best:   *best,
		Source: "bestchange",
	}
	c.JSON(http.StatusOK, resp)
}
