package http

import (
	"net/http"
	"strconv"

	"github.com/abhishek622/moviedock/movie/internal/controller/movie"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	ctrl *movie.Controller
}

func New(ctrl *movie.Controller) *Handler {
	return &Handler{ctrl}
}

func (h *Handler) GetMovieDetails(c *gin.Context) {
	idStr := c.Query("id")
	if idStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id parameter is required"})
		return
	}

	id, err := strconv.ParseInt(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id format"})
		return
	}

	details, err := h.ctrl.Get(c.Request.Context(), int32(id))
	if err != nil {
		if err == movie.ErrNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "movie not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		}
		return
	}

	c.JSON(http.StatusOK, details)
}

func (h *Handler) RegisterRoutes(router *gin.Engine) {
	router.GET("/api/v1/movie", h.GetMovieDetails)
}
