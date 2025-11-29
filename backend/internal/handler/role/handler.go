package role

import (
	"net/http"
	"strings"

	"github.com/example/ai-avatar-studio/internal/model"
	"github.com/example/ai-avatar-studio/internal/pkg/middleware"
	"github.com/example/ai-avatar-studio/internal/pkg/response"
	rolesvc "github.com/example/ai-avatar-studio/internal/service/role"
	"github.com/gin-gonic/gin"
)

// Handler exposes CRUD endpoints for roles.
type Handler struct {
	service *rolesvc.Service
	secret  string
}

func NewHandler(service *rolesvc.Service, secret string) *Handler {
	return &Handler{service: service, secret: secret}
}

func (h *Handler) RegisterRoutes(rg *gin.RouterGroup) {
	public := rg
	public.Use(middleware.OptionalAuth(h.secret))
	public.GET("/roles", h.list)
	public.GET("/roles/:id", h.get)
	public.GET("/roles/featured", h.featured)
	auth := middleware.Authenticator(h.secret)
	rg.POST("/roles", auth, h.create)
	rg.PUT("/roles/:id", auth, h.update)
	rg.POST("/roles/:id/publish", auth, h.publish)
	rg.POST("/roles/:id/archive", auth, h.archive)
	rg.POST("/roles/:id/prompt", auth, h.snapshotPrompt)
	rg.GET("/roles/favorites", auth, h.favorites)
	rg.POST("/roles/:id/favorite", auth, h.favorite)
	rg.DELETE("/roles/:id/favorite", auth, h.unfavorite)
}

func (h *Handler) list(c *gin.Context) {
	roles, err := h.service.List(c.Request.Context())
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.Success(c, roles)
}

func (h *Handler) featured(c *gin.Context) {
	roles, err := h.service.Featured(c.Request.Context())
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.Success(c, roles)
}

func (h *Handler) get(c *gin.Context) {
	userID := middleware.CurrentUserID(c)
	role, err := h.service.Get(c.Request.Context(), c.Param("id"))
	if err != nil || role == nil {
		response.Error(c, http.StatusNotFound, "role not found")
		return
	}
	h.service.PopulateFavoriteMetadata(c.Request.Context(), userID, role)
	response.Success(c, role)
}

func (h *Handler) create(c *gin.Context) {
	userID := middleware.CurrentUserID(c)
	var payload model.Role
	if err := c.ShouldBindJSON(&payload); err != nil {
		response.Error(c, http.StatusBadRequest, "invalid body")
		return
	}
	role, err := h.service.Save(c.Request.Context(), userID, &payload)
	if err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	response.Created(c, role)
}

func (h *Handler) update(c *gin.Context) {
	userID := middleware.CurrentUserID(c)
	var payload model.Role
	if err := c.ShouldBindJSON(&payload); err != nil {
		response.Error(c, http.StatusBadRequest, "invalid body")
		return
	}
	payload.ID = c.Param("id")
	role, err := h.service.Save(c.Request.Context(), userID, &payload)
	if err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	response.Success(c, role)
}

func (h *Handler) publish(c *gin.Context) {
	if err := h.service.Publish(c.Request.Context(), c.Param("id")); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	response.Success(c, gin.H{"status": "published"})
}

func (h *Handler) archive(c *gin.Context) {
	if err := h.service.Archive(c.Request.Context(), c.Param("id")); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	response.Success(c, gin.H{"status": "archived"})
}

func (h *Handler) snapshotPrompt(c *gin.Context) {
	var req struct {
		Prompt string `json:"prompt"`
	}
	if err := c.ShouldBindJSON(&req); err != nil || strings.TrimSpace(req.Prompt) == "" {
		response.Error(c, http.StatusBadRequest, "prompt required")
		return
	}
	if err := h.service.SnapshotPrompt(c.Request.Context(), c.Param("id"), req.Prompt); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	response.Success(c, gin.H{"status": "saved"})
}

func (h *Handler) favorite(c *gin.Context) {
	userID := middleware.CurrentUserID(c)
	if err := h.service.Favorite(c.Request.Context(), userID, c.Param("id")); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	response.Success(c, gin.H{"status": "favorited"})
}

func (h *Handler) unfavorite(c *gin.Context) {
	userID := middleware.CurrentUserID(c)
	if err := h.service.Unfavorite(c.Request.Context(), userID, c.Param("id")); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	response.Success(c, gin.H{"status": "unfavorited"})
}

func (h *Handler) favorites(c *gin.Context) {
	userID := middleware.CurrentUserID(c)
	roles, err := h.service.ListFavorites(c.Request.Context(), userID)
	if err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	response.Success(c, roles)
}
