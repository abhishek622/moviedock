package http

import (
	"net/http"

	"github.com/abhishek622/moviedock/user/internal/controller/user"
	"github.com/abhishek622/moviedock/user/pkg/model"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"golang.org/x/crypto/bcrypt"
)

var validate *validator.Validate

func init() {
	validate = validator.New()
}

type Handler struct {
	ctrl *user.Controller
}

func New(ctrl *user.Controller) *Handler {
	return &Handler{ctrl}
}

// RegisterRoutes registers all the routes for the user service.
func (h *Handler) RegisterRoutes(router *gin.Engine) {
	// Health check endpoint
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	// API v1 routes
	v1 := router.Group("/api/v1")
	{
		// Public routes
		auth := v1.Group("/auth")
		{
			auth.POST("/register", h.RegisterUser)
			auth.POST("/login", h.LoginUser)
		}

		// Protected routes
		// user := v1.Group("/user")
		// user.Use(h.AuthMiddleware())
		// {
		// 	user.GET("/profile", h.GetProfile)
		// 	user.PUT("/profile", h.UpdateProfile)
		// 	user.POST("/logout", h.LogoutUser)
		// 	user.POST("/refresh", h.RefreshToken)
		// }
	}
}

// RegisterRequest represents the JSON payload for user registration
type RegisterRequest struct {
	Email    string  `json:"email" binding:"required,email"`
	Password string  `json:"password" binding:"required,min=8"`
	FullName string  `json:"full_name" binding:"required"`
	Timezone *string `json:"timezone,omitempty"`
}

// RegisterUser handles user registration
func (h *Handler) RegisterUser(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body", "details": err.Error()})
		return
	}

	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process password"})
		return
	}

	// Create user model
	user := &model.User{
		Email:             req.Email,
		EncryptedPassword: string(hashedPassword),
		FullName:          req.FullName,
		Role:              model.RoleUser, // Default role
		IsActive:          true,
		Timezone:          req.Timezone,
	}

	// Call controller to register user
	registeredUser, err := h.ctrl.RegisterUser(c.Request.Context(), user)
	if err != nil {
		// Check if it's a duplicate email error
		if err.Error() == "user with this email already exists" {
			c.JSON(http.StatusConflict, gin.H{"error": "User with this email already exists"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to register user"})
		return
	}

	// For now, just return the created user without tokens
	// In a real implementation, you would generate JWT tokens here
	response := model.UserResponse{
		UserId:   registeredUser.UserID,
		FullName: registeredUser.FullName,
		Email:    registeredUser.Email,
		Role:     string(registeredUser.Role),
	}

	c.JSON(http.StatusCreated, response)
}

// LoginUser handles user login
func (h *Handler) LoginUser(c *gin.Context) {
	var loginReq model.UserLogin
	if err := c.ShouldBindJSON(&loginReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	// Validate request
	if err := validate.Struct(loginReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Validation failed", "details": err.Error()})
		return
	}

	// Create user model with login credentials
	user := &model.User{
		Email:             loginReq.Email,
		EncryptedPassword: loginReq.Password,
	}

	// Call controller to authenticate user
	loggedInUser, err := h.ctrl.LoginUser(c.Request.Context(), user)
	if err != nil {
		if err.Error() == "user not found" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process login"})
		}
		return
	}

	// Verify password
	err = bcrypt.CompareHashAndPassword([]byte(loggedInUser.EncryptedPassword), []byte(loginReq.Password))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return
	}

	// For now, just return the user without tokens
	// In a real implementation, you would generate JWT tokens here
	response := model.UserResponse{
		UserId:   loggedInUser.UserID,
		FullName: loggedInUser.FullName,
		Email:    loggedInUser.Email,
		Role:     string(loggedInUser.Role),
	}

	c.JSON(http.StatusOK, response)
}
