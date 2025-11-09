package http

import (
	"errors"
	"log"
	"net/http"
	"strconv"

	"github.com/abhishek622/moviedock/metadata/internal/controller/metadata"
	"github.com/abhishek622/moviedock/metadata/pkg/model"
	"github.com/gin-gonic/gin"
)

// Handler defines a movie metadata HTTP handler.
type Handler struct {
	ctrl *metadata.Controller
}

// New creates a new movie metadata HTTP handler.
func New(ctrl *metadata.Controller) *Handler {
	return &Handler{ctrl}
}

// RegisterRoutes registers all the routes for the metadata service.
func (h *Handler) RegisterRoutes(router *gin.Engine) {
	// Health check endpoint
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	// API v1 routes
	v1 := router.Group("/api/v1/metadata")
	{
		v1.POST("", h.CreateMetadata)
		v1.GET("", h.ListMetadata)
		v1.GET("/:id", h.GetMetadata)
		v1.PUT("/:id", h.UpdateMetadata)
		v1.DELETE("/:id", h.DeleteMetadata)
	}
}

func (h *Handler) GetMetadata(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	m, err := h.ctrl.Get(c.Request.Context(), int32(id))
	if err != nil {
		if errors.Is(err, metadata.ErrNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "metadata not found"})
			return
		}
		log.Printf("Failed to get metadata: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}
	c.JSON(http.StatusOK, m)
}

func (h *Handler) CreateMetadata(c *gin.Context) {
	var req *model.Metadata
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	metadata, err := h.ctrl.Create(c.Request.Context(), req)
	if err != nil {
		log.Printf("Failed to create metadata: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create metadata"})
		return
	}

	c.JSON(http.StatusCreated, metadata)
}

func (h *Handler) UpdateMetadata(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	var req model.Metadata
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	m, err := h.ctrl.Update(c.Request.Context(), int32(id), &req)
	if err != nil {
		if errors.Is(err, metadata.ErrNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "metadata not found"})
			return
		}
		log.Printf("Failed to update metadata: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update metadata"})
		return
	}

	c.JSON(http.StatusOK, m)
}

func (h *Handler) DeleteMetadata(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	if err := h.ctrl.Delete(c.Request.Context(), int32(id)); err != nil {
		if err == metadata.ErrNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "metadata not found"})
			return
		}
		log.Printf("Failed to delete metadata: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete metadata"})
		return
	}

	c.Status(http.StatusNoContent)
}

func (h *Handler) ListMetadata(c *gin.Context) {
	limit, err := strconv.Atoi(c.DefaultQuery("limit", "10"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid limit"})
		return
	}
	offset, err := strconv.Atoi(c.DefaultQuery("offset", "0"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid offset"})
		return
	}
	metadata, err := h.ctrl.List(c.Request.Context(), limit, offset)
	if err != nil {
		log.Printf("Failed to list metadata: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list metadata"})
		return
	}

	c.JSON(http.StatusOK, metadata)
}
