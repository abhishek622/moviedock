package http

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/abhishek622/moviedock/rating/internal/controller/rating"
	"github.com/abhishek622/moviedock/rating/pkg/model"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	ctrl *rating.Controller
}

func New(ctrl *rating.Controller) *Handler {
	return &Handler{ctrl}
}

// RegisterRoutes registers all the routes for the rating service.
func (h *Handler) RegisterRoutes(router *gin.Engine) {
	// Health check endpoint
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	// API v1 routes
	v1 := router.Group("/api/v1/rating")
	{
		v1.POST("", h.PutRating)
		v1.GET("", h.GetAggregatedRating)
	}
}

func (h *Handler) PutRating(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	var req model.Rating
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.ctrl.PutRating(c.Request.Context(), model.RecordID(id), model.RecordType(req.RecordType), &req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

func (h *Handler) GetAggregatedRating(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	recordID := model.RecordID(id)
	recordType := model.RecordType(c.Param("record_type"))
	if recordType == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid record type"})
		return
	}

	v, err := h.ctrl.GetAggregatedRating(c.Request.Context(), recordID, recordType)
	if err != nil && errors.Is(err, rating.ErrNotFound) {
		c.JSON(http.StatusNotFound, gin.H{"error": "rating not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"rating": v})
}

func (h *Handler) DeleteRating(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	if err := h.ctrl.DeleteRating(c.Request.Context(), model.UserID(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}
