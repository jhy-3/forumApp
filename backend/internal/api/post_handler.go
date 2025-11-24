package api

import (
	"net/http"
	"strconv"

	"github.com/forumapp/backend/internal/models"
	"github.com/forumapp/backend/internal/service"
	"github.com/labstack/echo/v4"
)

type PostHandler struct {
	postService *service.PostService
}

func NewPostHandler(postService *service.PostService) *PostHandler {
	return &PostHandler{
		postService: postService,
	}
}

func (h *PostHandler) CreatePost(c echo.Context) error {
	userID, err := GetUserIDFromContext(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, models.NewErrorResponse(
			"Unauthorized",
			"UNAUTHORIZED",
			nil,
		))
	}

	idParam := c.Param("id")
	threadID, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, models.NewErrorResponse(
			"Invalid thread ID",
			"BAD_REQUEST",
			nil,
		))
	}

	var req models.PostCreateRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, models.NewErrorResponse(
			"Invalid request body",
			"BAD_REQUEST",
			err.Error(),
		))
	}

	// Basic validation
	if req.Content == "" {
		return c.JSON(http.StatusBadRequest, models.NewErrorResponse(
			"Content is required",
			"VALIDATION_ERROR",
			nil,
		))
	}

	post, err := h.postService.CreatePost(userID, threadID, &req)
	if err != nil {
		return c.JSON(http.StatusBadRequest, models.NewErrorResponse(
			err.Error(),
			"CREATE_FAILED",
			nil,
		))
	}

	return c.JSON(http.StatusCreated, post)
}

func (h *PostHandler) GetPostsByThreadID(c echo.Context) error {
	idParam := c.Param("id")
	threadID, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, models.NewErrorResponse(
			"Invalid thread ID",
			"BAD_REQUEST",
			nil,
		))
	}

	// Parse pagination parameters
	page, _ := strconv.Atoi(c.QueryParam("page"))
	if page < 1 {
		page = 1
	}

	limit, _ := strconv.Atoi(c.QueryParam("limit"))
	if limit < 1 || limit > 100 {
		limit = 25
	}

	response, err := h.postService.GetPostsByThreadID(threadID, page, limit)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, models.NewErrorResponse(
			err.Error(),
			"INTERNAL_ERROR",
			nil,
		))
	}

	return c.JSON(http.StatusOK, response)
}

func (h *PostHandler) UpdatePost(c echo.Context) error {
	userID, err := GetUserIDFromContext(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, models.NewErrorResponse(
			"Unauthorized",
			"UNAUTHORIZED",
			nil,
		))
	}

	idParam := c.Param("id")
	postID, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, models.NewErrorResponse(
			"Invalid post ID",
			"BAD_REQUEST",
			nil,
		))
	}

	var req models.PostUpdateRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, models.NewErrorResponse(
			"Invalid request body",
			"BAD_REQUEST",
			err.Error(),
		))
	}

	// Basic validation
	if req.Content == "" {
		return c.JSON(http.StatusBadRequest, models.NewErrorResponse(
			"Content is required",
			"VALIDATION_ERROR",
			nil,
		))
	}

	post, err := h.postService.UpdatePost(postID, userID, &req)
	if err != nil {
		return c.JSON(http.StatusBadRequest, models.NewErrorResponse(
			err.Error(),
			"UPDATE_FAILED",
			nil,
		))
	}

	return c.JSON(http.StatusOK, post)
}

func (h *PostHandler) DeletePost(c echo.Context) error {
	userID, err := GetUserIDFromContext(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, models.NewErrorResponse(
			"Unauthorized",
			"UNAUTHORIZED",
			nil,
		))
	}

	idParam := c.Param("id")
	postID, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, models.NewErrorResponse(
			"Invalid post ID",
			"BAD_REQUEST",
			nil,
		))
	}

	err = h.postService.DeletePost(postID, userID)
	if err != nil {
		return c.JSON(http.StatusBadRequest, models.NewErrorResponse(
			err.Error(),
			"DELETE_FAILED",
			nil,
		))
	}

	return c.NoContent(http.StatusNoContent)
}

