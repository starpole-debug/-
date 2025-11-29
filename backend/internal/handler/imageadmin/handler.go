package imageadmin

import (
	"net/http"
	"encoding/json"
	"log"

	"github.com/example/ai-avatar-studio/internal/model"
	"github.com/example/ai-avatar-studio/internal/pkg/middleware"
	"github.com/example/ai-avatar-studio/internal/pkg/response"
	"github.com/example/ai-avatar-studio/internal/repository"
	"github.com/gin-gonic/gin"
)

// Handler manages image providers and presets for admin.
type Handler struct {
	providers *repository.ImageProviderRepository
	presets   *repository.ImagePresetRepository
	secret    string
}

func NewHandler(providers *repository.ImageProviderRepository, presets *repository.ImagePresetRepository, secret string) *Handler {
	return &Handler{providers: providers, presets: presets, secret: secret}
}

func (h *Handler) RegisterRoutes(rg *gin.RouterGroup) {
	admin := rg.Group("/admin", middleware.AdminOnly(h.secret))
	admin.GET("/image-providers", h.listProviders)
	admin.POST("/image-providers", h.saveProvider)
	admin.PUT("/image-providers/:id", h.saveProvider)
	admin.DELETE("/image-providers/:id", h.deleteProvider)

	admin.GET("/image-presets", h.listPresets)
	admin.POST("/image-presets", h.savePreset)
	admin.PUT("/image-presets/:id", h.savePreset)
	admin.DELETE("/image-presets/:id", h.deletePreset)
}

func (h *Handler) listProviders(c *gin.Context) {
	items, err := h.providers.List(c.Request.Context())
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.Success(c, items)
}

func (h *Handler) saveProvider(c *gin.Context) {
	var payload model.ImageProvider
	if err := c.ShouldBindJSON(&payload); err != nil {
		response.Error(c, http.StatusBadRequest, "invalid body")
		return
	}
	if payload.ID == "" {
		payload.ID = c.Param("id")
	}
	if payload.Status == "" {
		payload.Status = "active"
	}
	if payload.ParamsJSON == "" {
		payload.ParamsJSON = "{}"
	}
	// Normalize params_json to canonical JSON so updates persist.
	if err := validateJSON(payload.ParamsJSON); err != nil {
		response.Error(c, http.StatusBadRequest, "params_json must be valid JSON")
		return
	}
	var canonical interface{}
	if err := json.Unmarshal([]byte(payload.ParamsJSON), &canonical); err == nil {
		if compact, err := json.Marshal(canonical); err == nil {
			payload.ParamsJSON = string(compact)
		}
	}
	log.Printf("admin save image provider id=%s params_json=%s", payload.ID, payload.ParamsJSON)
	item, err := h.providers.Save(c.Request.Context(), &payload)
	if err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	response.Success(c, item)
}

func (h *Handler) deleteProvider(c *gin.Context) {
	if err := h.providers.Delete(c.Request.Context(), c.Param("id")); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	response.Success(c, gin.H{"status": "deleted"})
}

func (h *Handler) listPresets(c *gin.Context) {
	items, err := h.presets.List(c.Request.Context())
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.Success(c, items)
}

func (h *Handler) savePreset(c *gin.Context) {
	var payload model.ImagePreset
	if err := c.ShouldBindJSON(&payload); err != nil {
		response.Error(c, http.StatusBadRequest, "invalid body")
		return
	}
	if payload.ID == "" {
		payload.ID = c.Param("id")
	}
	if payload.Status == "" {
		payload.Status = "active"
	}
	if payload.PresetJSON == "" {
		payload.PresetJSON = "{}"
	}
	if err := validateJSON(payload.PresetJSON); err != nil {
		response.Error(c, http.StatusBadRequest, "preset_json must be valid JSON")
		return
	}
	item, err := h.presets.Save(c.Request.Context(), &payload)
	if err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	response.Success(c, item)
}

func (h *Handler) deletePreset(c *gin.Context) {
	if err := h.presets.Delete(c.Request.Context(), c.Param("id")); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	response.Success(c, gin.H{"status": "deleted"})
}

func validateJSON(raw string) error {
	if raw == "" {
		return nil
	}
	var v interface{}
	return json.Unmarshal([]byte(raw), &v)
}
