package api

import (
	"net/http"

	"github.com/forumapp/backend/internal/models"
	"github.com/forumapp/backend/internal/service"
	"github.com/labstack/echo/v4"
)

type AuthHandler struct {
	authService *service.AuthService
}

func NewAuthHandler(authService *service.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

func (h *AuthHandler) Register(c echo.Context) error {
	var req models.UserRegistrationRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, models.NewErrorResponse(
			"Invalid request body",
			"BAD_REQUEST",
			err.Error(),
		))
	}

	// Basic validation
	if req.Username == "" || req.Email == "" || req.Password == "" {
		return c.JSON(http.StatusBadRequest, models.NewErrorResponse(
			"Username, email, and password are required",
			"VALIDATION_ERROR",
			nil,
		))
	}

	if len(req.Password) < 6 {
		return c.JSON(http.StatusBadRequest, models.NewErrorResponse(
			"Password must be at least 6 characters",
			"VALIDATION_ERROR",
			nil,
		))
	}

	resp, err := h.authService.Register(&req)
	if err != nil {
		return c.JSON(http.StatusBadRequest, models.NewErrorResponse(
			err.Error(),
			"REGISTRATION_FAILED",
			nil,
		))
	}

	return c.JSON(http.StatusCreated, resp)
}

func (h *AuthHandler) Login(c echo.Context) error {
	var req models.UserLoginRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, models.NewErrorResponse(
			"Invalid request body",
			"BAD_REQUEST",
			err.Error(),
		))
	}

	// Basic validation
	if req.Username == "" || req.Password == "" {
		return c.JSON(http.StatusBadRequest, models.NewErrorResponse(
			"Username and password are required",
			"VALIDATION_ERROR",
			nil,
		))
	}

	resp, err := h.authService.Login(&req)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, models.NewErrorResponse(
			err.Error(),
			"LOGIN_FAILED",
			nil,
		))
	}

	return c.JSON(http.StatusOK, resp)
}

