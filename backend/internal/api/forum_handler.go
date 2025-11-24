package api

import (
	"net/http"

	"github.com/forumapp/backend/internal/models"
	"github.com/forumapp/backend/internal/service"
	"github.com/labstack/echo/v4"
)

type ForumHandler struct {
	forumService *service.ForumService
}

func NewForumHandler(forumService *service.ForumService) *ForumHandler {
	return &ForumHandler{
		forumService: forumService,
	}
}

func (h *ForumHandler) GetAllForums(c echo.Context) error {
	forums, err := h.forumService.GetAllForums()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, models.NewErrorResponse(
			err.Error(),
			"INTERNAL_ERROR",
			nil,
		))
	}

	return c.JSON(http.StatusOK, forums)
}

