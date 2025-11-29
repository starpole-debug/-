package middleware

import (
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

const (
	userIDKey    = "currentUserID"
	userAdminKey = "currentUserAdmin"
)

// AuthClaims encodes the user identity inside JWT tokens.
type AuthClaims struct {
	UserID  string `json:"user_id"`
	IsAdmin bool   `json:"is_admin"`
	jwt.RegisteredClaims
}

// Authenticator enforces a valid bearer token for protected APIs.
func Authenticator(secret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		token := readToken(c)
		if token == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "missing token"})
			return
		}
		claims, err := parseToken(secret, token)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			return
		}
		c.Set(userIDKey, claims.UserID)
		c.Set(userAdminKey, claims.IsAdmin)
		c.Next()
	}
}

// OptionalAuth tries to parse the token but does not reject anonymous requests.
func OptionalAuth(secret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		token := readToken(c)
		if token == "" {
			c.Next()
			return
		}
		claims, err := parseToken(secret, token)
		if err == nil {
			c.Set(userIDKey, claims.UserID)
			c.Set(userAdminKey, claims.IsAdmin)
		}
		c.Next()
	}
}

// AdminOnly ensures the caller holds an admin claim.
func AdminOnly(secret string) gin.HandlerFunc {
	auth := Authenticator(secret)
	return func(c *gin.Context) {
		auth(c)
		if c.IsAborted() {
			return
		}
		if !IsAdmin(c) {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "admin only"})
			return
		}
		c.Next()
	}
}

// IssueToken builds a signed JWT with the provided claims.
func IssueToken(secret, userID string, isAdmin bool, ttl time.Duration) (string, error) {
	if ttl <= 0 {
		ttl = 24 * time.Hour
	}
	claims := AuthClaims{
		UserID:  userID,
		IsAdmin: isAdmin,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(ttl)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

// CurrentUserID extracts the user id stored by the middleware.
func CurrentUserID(c *gin.Context) string {
	if v, ok := c.Get(userIDKey); ok {
		if s, ok := v.(string); ok {
			return s
		}
	}
	return ""
}

// IsAdmin returns whether the current request is issued by an administrator.
func IsAdmin(c *gin.Context) bool {
	if v, ok := c.Get(userAdminKey); ok {
		if b, ok := v.(bool); ok {
			return b
		}
	}
	return false
}

func readToken(c *gin.Context) string {
	token := c.GetHeader("Authorization")
	if strings.HasPrefix(strings.ToLower(token), "bearer ") {
		return strings.TrimSpace(token[7:])
	}
	if t, err := c.Cookie("auth-token"); err == nil && strings.TrimSpace(t) != "" {
		return t
	}
	if t, err := c.Cookie("token"); err == nil {
		return t
	}
	return ""
}

func parseToken(secret, token string) (*AuthClaims, error) {
	parsed, err := jwt.ParseWithClaims(token, &AuthClaims{}, func(t *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})
	if err != nil {
		return nil, err
	}
	if claims, ok := parsed.Claims.(*AuthClaims); ok && parsed.Valid {
		return claims, nil
	}
	return nil, jwt.ErrTokenInvalidClaims
}
