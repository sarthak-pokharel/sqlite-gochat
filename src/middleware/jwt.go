package middleware

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

// JWTConfig holds JWT middleware configuration
type JWTConfig struct {
	Secret     string
	SkipRoutes []string
}

// CustomClaims represents JWT claims structure
type CustomClaims struct {
	UserID         string `json:"user_id"`
	OrganizationID string `json:"organization_id"`
	Role           string `json:"role"`
	jwt.RegisteredClaims
}

// SetupJWT configures JWT authentication middleware
func SetupJWT(config JWTConfig) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Check if route should skip JWT validation
			path := c.Path()
			for _, skipRoute := range config.SkipRoutes {
				if strings.HasPrefix(path, skipRoute) {
					return next(c)
				}
			}

			// Extract token from Authorization header
			authHeader := c.Request().Header.Get("Authorization")
			if authHeader == "" {
				return echo.NewHTTPError(http.StatusUnauthorized, "missing authorization header")
			}

			// Expected format: "Bearer <token>"
			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || parts[0] != "Bearer" {
				return echo.NewHTTPError(http.StatusUnauthorized, "invalid authorization format")
			}

			tokenString := parts[1]

			// Parse and validate token
			token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
				}
				return []byte(config.Secret), nil
			})

			if err != nil {
				return echo.NewHTTPError(http.StatusUnauthorized, "invalid token")
			}

			if !token.Valid {
				return echo.NewHTTPError(http.StatusUnauthorized, "token is not valid")
			}

			// Extract claims
			claims, ok := token.Claims.(*CustomClaims)
			if !ok {
				return echo.NewHTTPError(http.StatusUnauthorized, "invalid token claims")
			}

			// Set claims in context for handlers to use
			c.Set("user_id", claims.UserID)
			c.Set("organization_id", claims.OrganizationID)
			c.Set("role", claims.Role)

			return next(c)
		}
	}
}

// GetUserID retrieves user ID from context
func GetUserID(c echo.Context) string {
	userID, _ := c.Get("user_id").(string)
	return userID
}

// GetOrganizationID retrieves organization ID from context
func GetOrganizationID(c echo.Context) string {
	orgID, _ := c.Get("organization_id").(string)
	return orgID
}

// GetRole retrieves user role from context
func GetRole(c echo.Context) string {
	role, _ := c.Get("role").(string)
	return role
}
