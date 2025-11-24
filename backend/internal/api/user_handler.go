package api

import (
	"net/http"

	"github.com/forumapp/backend/internal/models"
	"github.com/forumapp/backend/internal/service"
	"github.com/labstack/echo/v4"
)

type UserHandler struct {
	userService *service.UserService
}

func NewUserHandler(userService *service.UserService) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

func (h *UserHandler) GetCurrentUser(c echo.Context) error {
	userID, err := GetUserIDFromContext(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, models.NewErrorResponse(
			"Unauthorized",
			"UNAUTHORIZED",
			nil,
		))
	}

	user, err := h.userService.GetUserByID(userID)
	if err != nil {
		return c.JSON(http.StatusNotFound, models.NewErrorResponse(
			err.Error(),
			"USER_NOT_FOUND",
			nil,
		))
	}

	return c.JSON(http.StatusOK, user)
}

func (h *UserHandler) UpdateCurrentUser(c echo.Context) error {
	userID, err := GetUserIDFromContext(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, models.NewErrorResponse(
			"Unauthorized",
			"UNAUTHORIZED",
			nil,
		))
	}

	var req models.UserUpdateRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, models.NewErrorResponse(
			"Invalid request body",
			"BAD_REQUEST",
			err.Error(),
		))
	}

	user, err := h.userService.UpdateUser(userID, &req)
	if err != nil {
		return c.JSON(http.StatusBadRequest, models.NewErrorResponse(
			err.Error(),
			"UPDATE_FAILED",
			nil,
		))
	}

	return c.JSON(http.StatusOK, user)
}

