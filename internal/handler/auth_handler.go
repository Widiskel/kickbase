package handler

import (
	"net/http"

	"kickbase/internal/interfaces"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	authService interfaces.AuthService
}

func NewAuthHandler(authService interfaces.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

// Register godoc
// @Summary Register a new user
// @Description Register a new user with role (admin, staff, viewer)
// @Tags Auth
// @Accept json
// @Produce json
// @Param user body RegisterRequest true "User registration data"
// @Success 201 {object} SuccessResponse
// @Failure 400 {object} ErrorResponse
// @Failure 409 {object} ErrorResponse
// @Router /api/auth/register [post]
func (h *AuthHandler) Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		RespondError(c, http.StatusBadRequest, "VALIDATION_ERROR", "Invalid request body", nil)
		return
	}

	user, err := h.authService.Register(req.Username, req.Password, req.Name, req.Role)
	if err != nil {
		if err.Error() == "username already exists" {
			RespondError(c, http.StatusConflict, "CONFLICT", err.Error(), nil)
		} else {
			RespondError(c, http.StatusBadRequest, "VALIDATION_ERROR", err.Error(), nil)
		}
		return
	}

	RespondCreated(c, gin.H{
		"id":       user.ID,
		"username": user.Username,
		"name":     user.Name,
		"role":     user.Role,
	})
}

// Login godoc
// @Summary Login user
// @Description Authenticate user and return Access Token + Refresh Token + Permissions
// @Tags Auth
// @Accept json
// @Produce json
// @Param credentials body LoginRequest true "Login credentials"
// @Success 200 {object} SuccessResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Router /api/auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		RespondError(c, http.StatusBadRequest, "VALIDATION_ERROR", "Invalid request body", nil)
		return
	}

	authResp, err := h.authService.Login(req.Username, req.Password)
	if err != nil {
		RespondError(c, http.StatusUnauthorized, "UNAUTHORIZED", err.Error(), nil)
		return
	}

	RespondSuccess(c, authResp, "Login successful")
}

// RefreshToken godoc
// @Summary Refresh access token
// @Description Exchange a valid refresh token for a new token pair
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body RefreshTokenRequest true "Refresh token payload"
// @Success 200 {object} SuccessResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Router /api/auth/refresh [post]
func (h *AuthHandler) RefreshToken(c *gin.Context) {
	var req RefreshTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		RespondError(c, http.StatusBadRequest, "VALIDATION_ERROR", "Invalid request body", nil)
		return
	}

	authResp, err := h.authService.RefreshToken(req.RefreshToken)
	if err != nil {
		RespondError(c, http.StatusUnauthorized, "UNAUTHORIZED", err.Error(), nil)
		return
	}

	RespondSuccess(c, authResp, "Token refreshed successfully")
}
