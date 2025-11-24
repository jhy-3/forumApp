package api

import (
	"net/http"
	"strings"

	"github.com/forumapp/backend/internal/models"
	"github.com/forumapp/backend/internal/service"
	"github.com/labstack/echo/v4"
)

type AuthMiddleware struct {
	authService *service.AuthService
}

func NewAuthMiddleware(authService *service.AuthService) *AuthMiddleware {
	return &AuthMiddleware{
		authService: authService,
	}
}

func (m *AuthMiddleware) RequireAuth(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		// Get authorization header
		authHeader := c.Request().Header.Get("Authorization")
		if authHeader == "" {
			return c.JSON(http.StatusUnauthorized, models.NewErrorResponse(
				"Missing authorization header",
				"UNAUTHORIZED",
				nil,
			))
		}

		// Extract token
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			return c.JSON(http.StatusUnauthorized, models.NewErrorResponse(
				"Invalid authorization header format",
				"UNAUTHORIZED",
				nil,
			))
		}

		token := parts[1]

		// Validate token
		userID, err := m.authService.ValidateToken(token)
		if err != nil {
			return c.JSON(http.StatusUnauthorized, models.NewErrorResponse(
				"Invalid or expired token",
				"UNAUTHORIZED",
				nil,
			))
		}

		// Store user ID in context
		c.Set("userID", userID)

		return next(c)
	}
}

func GetUserIDFromContext(c echo.Context) (uint64, error) {
	userID, ok := c.Get("userID").(uint64)
	if !ok {
		return 0, echo.NewHTTPError(http.StatusUnauthorized, "User ID not found in context")
	}
	return userID, nil
}

