package middleware

import (
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4" // Replaced Gin with Echo
)

// TODO: Centralize jwtSecret, perhaps in a config package or environment variable.
// This secret should be the same as the one used in the handler package for token generation.
var jwtSecret = []byte("temporary-secret-key-please-change")

// AuthMiddleware creates an echo.HandlerFunc for JWT authentication.
func AuthMiddleware(next echo.HandlerFunc) echo.HandlerFunc { // Changed signature
	return func(c echo.Context) error { // Changed signature
		authHeader := c.Request().Header.Get("Authorization") // Get header from request for Echo
		if authHeader == "" {
			return c.JSON(http.StatusUnauthorized, echo.Map{"error": "Authorization header required"}) // Changed to return c.JSON and echo.Map
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			return c.JSON(http.StatusUnauthorized, echo.Map{"error": "Authorization header format must be Bearer {token}"})
		}

		tokenString := parts[1]

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// Validate the alg is what you expect:
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrSignatureInvalid // Using jwt.ErrSignatureInvalid as specified
			}
			return jwtSecret, nil
		})

		if err != nil {
			// Based on the error type, you can provide more specific messages.
			// For example, if errors.Is(err, jwt.ErrTokenExpired) { ... }
			return c.JSON(http.StatusUnauthorized, echo.Map{"error": "Invalid or expired token"})
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			userIDFloat, okUserID := claims["user_id"].(float64) // JWT standard decodes numbers as float64
			if !okUserID {
				return c.JSON(http.StatusUnauthorized, echo.Map{"error": "Invalid token claims: user_id missing or not a number"})
			}
			// Store the userID in the context for downstream handlers
			c.Set("userID", int64(userIDFloat)) // Convert to int64 for consistency
			return next(c)                      // Call the next handler
		}

		// This case might be redundant if token.Valid already covers it,
		// but it's a good safeguard for malformed claims.
		return c.JSON(http.StatusUnauthorized, echo.Map{"error": "Invalid token claims"})
	}
}
