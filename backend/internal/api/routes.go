package api

import (
	"github.com/forumapp/backend/internal/config"
	"github.com/forumapp/backend/internal/service"
	"github.com/labstack/echo/v4"
)

func SetupRoutes(
	e *echo.Echo,
	authService *service.AuthService,
	userService *service.UserService,
	forumService *service.ForumService,
	threadService *service.ThreadService,
	postService *service.PostService,
	cfg *config.Config,
) {
	// Initialize handlers
	authHandler := NewAuthHandler(authService)
	userHandler := NewUserHandler(userService)
	forumHandler := NewForumHandler(forumService)
	threadHandler := NewThreadHandler(threadService)
	postHandler := NewPostHandler(postService)

	// Middleware
	authMiddleware := NewAuthMiddleware(authService)

	// Public routes
	api := e.Group("/api")

	// Auth routes
	api.POST("/users/register", authHandler.Register)
	api.POST("/auth/login", authHandler.Login)

	// Forum routes
	api.GET("/forums", forumHandler.GetAllForums)
	api.GET("/forums/:slug/threads", threadHandler.GetThreadsByForumSlug)

	// Thread routes
	api.GET("/threads/:id", threadHandler.GetThreadByID)
	api.GET("/threads/:id/posts", postHandler.GetPostsByThreadID)

	// Protected routes
	protected := api.Group("")
	protected.Use(authMiddleware.RequireAuth)

	// User routes (protected)
	protected.GET("/users/me", userHandler.GetCurrentUser)
	protected.PATCH("/users/me", userHandler.UpdateCurrentUser)

	// Thread routes (protected)
	protected.POST("/threads", threadHandler.CreateThread)

	// Post routes (protected)
	protected.POST("/threads/:id/posts", postHandler.CreatePost)
	protected.PATCH("/posts/:id", postHandler.UpdatePost)
	protected.DELETE("/posts/:id", postHandler.DeletePost)

	// Health check
	e.GET("/health", func(c echo.Context) error {
		return c.JSON(200, map[string]string{"status": "ok"})
	})
}

