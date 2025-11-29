package router

import (
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/example/ai-avatar-studio/internal/config"
	adminhandler "github.com/example/ai-avatar-studio/internal/handler/admin"
	authhandler "github.com/example/ai-avatar-studio/internal/handler/auth"
	chathandler "github.com/example/ai-avatar-studio/internal/handler/chat"
	communityhandler "github.com/example/ai-avatar-studio/internal/handler/community"
	creatorhandler "github.com/example/ai-avatar-studio/internal/handler/creator"
	notificationhandler "github.com/example/ai-avatar-studio/internal/handler/notification"
	presethandler "github.com/example/ai-avatar-studio/internal/handler/preset"
	profilehandler "github.com/example/ai-avatar-studio/internal/handler/profile"
	revenuehandler "github.com/example/ai-avatar-studio/internal/handler/revenue"
	rolehandler "github.com/example/ai-avatar-studio/internal/handler/role"
	imagehandler "github.com/example/ai-avatar-studio/internal/handler/image"
	imageadminhandler "github.com/example/ai-avatar-studio/internal/handler/imageadmin"
	storehandler "github.com/example/ai-avatar-studio/internal/handler/store"
	uploadhandler "github.com/example/ai-avatar-studio/internal/handler/upload"
	"github.com/gin-gonic/gin"
	paymenthandler "github.com/example/ai-avatar-studio/internal/handler/payment"
)

// Handlers aggregates feature specific HTTP handlers so the router can wire them.
type Handlers struct {
	Auth         *authhandler.Handler
	Roles        *rolehandler.Handler
	Chat         *chathandler.Handler
	Community    *communityhandler.Handler
	Creator      *creatorhandler.Handler
	Store        *storehandler.Handler
	Notification *notificationhandler.Handler
	Admin        *adminhandler.Handler
	Revenue      *revenuehandler.Handler
	Profile      *profilehandler.Handler
	Upload       *uploadhandler.Handler
	Presets      *presethandler.Handler
	Payment      *paymenthandler.Handler
	Images       *imagehandler.Handler
	ImageAdmin   *imageadminhandler.Handler
}

// New builds the gin router + registers all HTTP routes.
func New(cfg *config.Config, handlers Handlers) *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	r.Use(gin.Logger(), gin.Recovery(), corsMiddleware(cfg.FrontendOrigin, cfg.Env), securityHeaders())
	r.MaxMultipartMemory = 16 << 20 // 16MB max multipart payload

	// Serve uploaded files (avatars, post images).
	if cfg.UploadDir != "" {
		baseDir := filepath.Clean(cfg.UploadDir)
		// Pre-create expected subfolders to avoid missing-disk errors on first upload.
		_ = os.MkdirAll(filepath.Join(baseDir, "avatars"), 0755)
		_ = os.MkdirAll(filepath.Join(baseDir, "posts"), 0755)
		_ = os.MkdirAll(filepath.Join(baseDir, "roles"), 0755)
		r.Static("/uploads", baseDir)
	}

	api := r.Group("/api")
	handlers.Auth.RegisterRoutes(api.Group("/auth"))
	handlers.Roles.RegisterRoutes(api)
	handlers.Chat.RegisterRoutes(api)
	handlers.Community.RegisterRoutes(api)
	handlers.Creator.RegisterRoutes(api)
	handlers.Store.RegisterRoutes(api)
	handlers.Notification.RegisterRoutes(api)
	handlers.Admin.RegisterRoutes(api)
	handlers.Revenue.RegisterRoutes(api)
	if handlers.Profile != nil {
		handlers.Profile.RegisterRoutes(api)
	}
	if handlers.Upload != nil {
		handlers.Upload.RegisterRoutes(api)
	}
	if handlers.Presets != nil {
		handlers.Presets.RegisterRoutes(api)
	}
	if handlers.Payment != nil {
		handlers.Payment.RegisterRoutes(api)
	}
	if handlers.Images != nil {
		handlers.Images.RegisterRoutes(api)
	}
	if handlers.ImageAdmin != nil {
		handlers.ImageAdmin.RegisterRoutes(api)
	}

	api.GET("/health", func(c *gin.Context) { c.JSON(http.StatusOK, gin.H{"status": "ok"}) })
	return r
}

func corsMiddleware(origins []string, env string) gin.HandlerFunc {
	production := strings.ToLower(env) == "production"
	return func(c *gin.Context) {
		origin := c.GetHeader("Origin")
		allowCredentials := false
		allowed := false
		if origin != "" {
			switch {
			case contains(origins, origin):
				allowed = true
			case !production && contains(origins, "*"):
				allowed = true
			case !production:
				allowed = true
			}
			if allowed {
				c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
				c.Writer.Header().Set("Vary", "Origin")
				allowCredentials = true
			}
		} else if !production {
			// Non-browser clients (curl, server-to-server) without Origin header
			c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		}
		if allowCredentials {
			c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		}
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Authorization, Content-Type, X-Admin-Access, x-admin-access")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET,POST,PUT,DELETE,PATCH,OPTIONS")
		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusOK)
			return
		}
		c.Next()
	}
}

func contains(list []string, target string) bool {
	for _, item := range list {
		if item == target || item == "*" {
			return true
		}
	}
	return false
}

// securityHeaders adds basic hardening headers for API responses.
func securityHeaders() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("X-Content-Type-Options", "nosniff")
		c.Header("X-Frame-Options", "DENY")
		c.Header("Referrer-Policy", "no-referrer")
		c.Header("X-XSS-Protection", "1; mode=block")
		c.Next()
	}
}
