package preset

import (
	"net/http"

	"github.com/example/ai-avatar-studio/internal/model"
	"github.com/example/ai-avatar-studio/internal/pkg/middleware"
	"github.com/example/ai-avatar-studio/internal/service/preset"
	"github.com/gin-gonic/gin"
	"log"
)

type Handler struct {
	service *preset.Service
	secret  string
}

func NewPresetHandler(service *preset.Service, secret string) *Handler {
	return &Handler{service: service, secret: secret}
}

func (h *Handler) RegisterRoutes(rg *gin.RouterGroup) {
	auth := middleware.Authenticator(h.secret)
	optional := middleware.OptionalAuth(h.secret)
	presets := rg.Group("/presets", optional)

	// Public marketplace & fetch by id (enforces is_public inside service)
	presets.GET("/public", h.ListPublic)
	presets.GET("/:id", h.Get)

	// Protected endpoints
	presets.POST("", auth, h.Create)
	presets.PUT("/:id", auth, h.Update)
	presets.DELETE("/:id", auth, h.Delete)
	presets.POST("/:id/publish", auth, h.Publish)
	// For compatibility, keep owner list behind auth
	presets.GET("", auth, h.List)
}

func (h *Handler) Create(c *gin.Context) {
	userID := middleware.CurrentUserID(c)
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	var req model.Preset
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	p, err := h.service.CreatePreset(c.Request.Context(), userID, req)
	if err != nil {
		log.Printf("preset create failed user=%s err=%v", userID, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"data": p})
}

func (h *Handler) Update(c *gin.Context) {
	userID := middleware.CurrentUserID(c)
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	id := c.Param("id")
	var req model.Preset
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	p, err := h.service.UpdatePreset(c.Request.Context(), userID, id, req)
	if err != nil {
		log.Printf("preset update failed user=%s id=%s err=%v", userID, id, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": p})
}

func (h *Handler) Delete(c *gin.Context) {
	userID := middleware.CurrentUserID(c)
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	id := c.Param("id")
	if err := h.service.DeletePreset(c.Request.Context(), userID, id); err != nil {
		log.Printf("preset delete failed user=%s id=%s err=%v", userID, id, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "deleted"})
}

func (h *Handler) Get(c *gin.Context) {
	userID := middleware.CurrentUserID(c)
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	id := c.Param("id")
	p, err := h.service.GetPreset(c.Request.Context(), userID, id)
	if err != nil {
		log.Printf("preset get failed user=%s id=%s err=%v", userID, id, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": p})
}

func (h *Handler) List(c *gin.Context) {
	userID := middleware.CurrentUserID(c)
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	presets, err := h.service.ListUserPresets(c.Request.Context(), userID)
	if err != nil {
		log.Printf("preset list failed user=%s err=%v", userID, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": presets})
}

func (h *Handler) Publish(c *gin.Context) {
	userID := middleware.CurrentUserID(c)
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	id := c.Param("id")
	var body struct {
		IsPublic bool `json:"is_public"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid body"})
		return
	}
	p, err := h.service.PublishPreset(c.Request.Context(), userID, id, body.IsPublic)
	if err != nil {
		log.Printf("preset publish failed user=%s id=%s err=%v", userID, id, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": p})
}

func (h *Handler) ListPublic(c *gin.Context) {
	presets, err := h.service.ListPublicPresets(c.Request.Context(), 50)
	if err != nil {
		log.Printf("preset list public failed err=%v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": presets})
}
