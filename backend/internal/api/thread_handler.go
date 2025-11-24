package api

import (
	"net/http"
	"strconv"

	"github.com/forumapp/backend/internal/models"
	"github.com/forumapp/backend/internal/service"
	"github.com/labstack/echo/v4"
)

type ThreadHandler struct {
	threadService *service.ThreadService
}

func NewThreadHandler(threadService *service.ThreadService) *ThreadHandler {
	return &ThreadHandler{
		threadService: threadService,
	}
}

func (h *ThreadHandler) CreateThread(c echo.Context) error {
	userID, err := GetUserIDFromContext(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, models.NewErrorResponse(
			"Unauthorized",
			"UNAUTHORIZED",
			nil,
		))
	}

	var req models.ThreadCreateRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, models.NewErrorResponse(
			"Invalid request body",
			"BAD_REQUEST",
			err.Error(),
		))
	}

	// Basic validation
	if req.ForumID == 0 || req.Title == "" || req.Content == "" {
		return c.JSON(http.StatusBadRequest, models.NewErrorResponse(
			"ForumID, title, and content are required",
			"VALIDATION_ERROR",
			nil,
		))
	}

	thread, err := h.threadService.CreateThread(userID, &req)
	if err != nil {
		return c.JSON(http.StatusBadRequest, models.NewErrorResponse(
			err.Error(),
			"CREATE_FAILED",
			nil,
		))
	}

	return c.JSON(http.StatusCreated, thread)
}

func (h *ThreadHandler) GetThreadByID(c echo.Context) error {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, models.NewErrorResponse(
			"Invalid thread ID",
			"BAD_REQUEST",
			nil,
		))
	}

	thread, err := h.threadService.GetThreadByID(id)
	if err != nil {
		return c.JSON(http.StatusNotFound, models.NewErrorResponse(
			err.Error(),
			"THREAD_NOT_FOUND",
			nil,
		))
	}

	return c.JSON(http.StatusOK, thread)
}

func (h *ThreadHandler) GetThreadsByForumSlug(c echo.Context) error {
	slug := c.Param("slug")

	// Parse pagination parameters
	page, _ := strconv.Atoi(c.QueryParam("page"))
	if page < 1 {
		page = 1
	}

	limit, _ := strconv.Atoi(c.QueryParam("limit"))
	if limit < 1 || limit > 100 {
		limit = 25
	}

	response, err := h.threadService.GetThreadsByForumSlug(slug, page, limit)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, models.NewErrorResponse(
			err.Error(),
			"INTERNAL_ERROR",
			nil,
		))
	}

	return c.JSON(http.StatusOK, response)
}

