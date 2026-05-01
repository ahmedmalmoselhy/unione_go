package handlers

import (
	"net/http"

	"github.com/ahmedmalmoselhy/unione_go/internal/apiutil"
	"github.com/ahmedmalmoselhy/unione_go/internal/middlewares"
	"github.com/ahmedmalmoselhy/unione_go/internal/models"
	"github.com/ahmedmalmoselhy/unione_go/internal/services"
	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	authService services.AuthService
}

func NewAuthHandler(authService services.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

type LoginInput struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type RegisterInput struct {
	Email     string      `json:"email" binding:"required,email"`
	Password  string      `json:"password" binding:"required,min=6"`
	FirstName string      `json:"first_name" binding:"required"`
	LastName  string      `json:"last_name" binding:"required"`
	Role      models.Role `json:"role" binding:"required"`
}

// Register godoc
// @Summary Register a new user
// @Description Create a new user with the specified details.
// @Tags auth
// @Accept json
// @Produce json
// @Param input body RegisterInput true "Registration details"
// @Success 201 {object} map[string]string "data: {message: Registration successful}"
// @Failure 400 {object} apiutil.errorResponse
// @Router /auth/register [post]
func (h *AuthHandler) Register(c *gin.Context) {
	var input RegisterInput
	if err := c.ShouldBindJSON(&input); err != nil {
		apiutil.Error(c, http.StatusBadRequest, "validation_error", err.Error())
		return
	}

	user := &models.User{
		Email:     input.Email,
		Password:  input.Password,
		FirstName: input.FirstName,
		LastName:  input.LastName,
		Role:      input.Role,
	}

	if err := h.authService.RegisterUser(user); err != nil {
		apiutil.Error(c, http.StatusBadRequest, "registration_failed", err.Error())
		return
	}

	apiutil.Message(c, http.StatusCreated, "Registration successful")
}

// Login godoc
// @Summary Login and get a JWT token
// @Description Authenticate a user and return a JWT token.
// @Tags auth
// @Accept json
// @Produce json
// @Param input body LoginInput true "Login credentials"
// @Success 200 {object} map[string]interface{} "data: {token: <jwt_token>}"
// @Failure 400 {object} apiutil.errorResponse
// @Failure 401 {object} apiutil.errorResponse
// @Router /auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	var input LoginInput
	if err := c.ShouldBindJSON(&input); err != nil {
		apiutil.Error(c, http.StatusBadRequest, "validation_error", err.Error())
		return
	}

	token, err := h.authService.AuthenticateUser(input.Email, input.Password)
	if err != nil {
		apiutil.Error(c, http.StatusUnauthorized, "authentication_failed", err.Error())
		return
	}

	apiutil.Success(c, http.StatusOK, gin.H{"token": token})
}

// GetMe godoc
// @Summary Get the current authenticated user
// @Description Returns the profile of the user associated with the JWT token.
// @Tags auth
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} models.User
// @Failure 401 {object} apiutil.errorResponse
// @Failure 404 {object} apiutil.errorResponse
// @Router /auth/me [get]
func (h *AuthHandler) GetMe(c *gin.Context) {
	userID, exists := middlewares.GetUserID(c)
	if !exists {
		apiutil.Error(c, http.StatusUnauthorized, "unauthorized", "Authentication required")
		return
	}

	user, err := h.authService.GetUserByID(userID)
	if err != nil {
		apiutil.Error(c, http.StatusNotFound, "not_found", "User not found")
		return
	}

	apiutil.Success(c, http.StatusOK, user)
}
