	package httpapi

	import (
		"errors"
		"log"
		"net/http"

		"github.com/gin-gonic/gin"
		"github.com/guidiguidi/RateMonitorBC/internal/bestchange"
		"github.com/guidiguidi/RateMonitorBC/internal/models"
	)

	type Handler struct {
		service *bestchange.Service
	}

	func NewHandler(service *bestchange.Service) *Handler {
		return &Handler{service: service}
	}

	func (h *Handler) GetBestRate(c *gin.Context) {
		var req models.BestRateRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		log.Printf("GetBestRate request: %+v", req)

		bestRate, err := h.service.GetBestRate(c.Request.Context(), &req)
		if err != nil {
			log.Printf("GetBestRate error: %v", err)
			if errors.Is(err, bestchange.ErrNoSuitableRates) {
				c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
				return
			}
			if errors.Is(err, bestchange.ErrCurrencyNotFound) {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}
			c.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
			return
		}

		log.Printf("GetBestRate response: %+v", bestRate)
		c.JSON(http.StatusOK, bestRate)
	}
