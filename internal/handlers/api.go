package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/zakwanzambri/Gotrader/internal/database"
	"github.com/zakwanzambri/Gotrader/internal/models"
)

type Handler struct {
	db *database.DB
}

// NewHandler creates a new handler instance
func NewHandler(db *database.DB) *Handler {
	return &Handler{db: db}
}

// GetSignals returns trading signals with pagination and filtering
func (h *Handler) GetSignals(c *gin.Context) {
	// Parse query parameters
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	symbol := c.Query("symbol")
	status := c.Query("status")

	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}

	offset := (page - 1) * limit

	// Get signals from database
	signals, err := h.db.GetSignals(limit, offset, symbol, status)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch signals"})
		return
	}

	// Get total count for pagination
	total, err := h.db.GetSignalCount(symbol, status)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get signal count"})
		return
	}

	response := models.SignalResponse{
		Signals: signals,
		Total:   total,
		Page:    page,
		Limit:   limit,
	}

	c.JSON(http.StatusOK, response)
}

// UpdateSignalStatus updates the status of a trading signal
func (h *Handler) UpdateSignalStatus(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid signal ID"})
		return
	}

	var req struct {
		Status string `json:"status" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate status
	validStatuses := map[string]bool{
		"ACTIVE":  true,
		"CLOSED":  true,
		"EXPIRED": true,
	}

	if !validStatuses[req.Status] {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid status"})
		return
	}

	err = h.db.UpdateSignalStatus(id, req.Status)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update signal"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Signal updated successfully"})
}

// GetStats returns statistics about trading signals
func (h *Handler) GetStats(c *gin.Context) {
	// Get total signals
	total, err := h.db.GetSignalCount("", "")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get stats"})
		return
	}

	// Get active signals
	active, err := h.db.GetSignalCount("", "ACTIVE")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get stats"})
		return
	}

	// Get closed signals
	closed, err := h.db.GetSignalCount("", "CLOSED")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get stats"})
		return
	}

	stats := gin.H{
		"total":  total,
		"active": active,
		"closed": closed,
	}

	c.JSON(http.StatusOK, stats)
}

// HealthCheck returns the health status of the API
func (h *Handler) HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status": "ok",
		"service": "gotrader-api",
	})
}