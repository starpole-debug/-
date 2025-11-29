package image

import (
	"net/http"

	"github.com/example/ai-avatar-studio/internal/pkg/middleware"
	"github.com/example/ai-avatar-studio/internal/pkg/response"
	imagesvc "github.com/example/ai-avatar-studio/internal/service/image"
	"github.com/gin-gonic/gin"
)

// Handler exposes endpoints for chat image generation.
type Handler struct {
	service *imagesvc.Service
	secret  string
}

func NewHandler(service *imagesvc.Service, secret string) *Handler {
	return &Handler{service: service, secret: secret}
}

func (h *Handler) RegisterRoutes(rg *gin.RouterGroup) {
	auth := middleware.Authenticator(h.secret)
	rg.POST("/chat/images", auth, h.create)
	rg.GET("/chat/images/:id", auth, h.detail)
}

func (h *Handler) create(c *gin.Context) {
	userID := middleware.CurrentUserID(c)
	var req struct {
		SessionID string `json:"session_id"`
		MessageID string `json:"message_id"`
		Prompt    string `json:"prompt"`
	}
	if err := c.ShouldBindJSON(&req); err != nil || req.SessionID == "" {
		response.Error(c, http.StatusBadRequest, "invalid body")
		return
	}
	job, err := h.service.RequestImage(c.Request.Context(), userID, req.SessionID, req.MessageID, req.Prompt)
	if err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	response.Created(c, job)
}

func (h *Handler) detail(c *gin.Context) {
	job, err := h.service.GetJob(c.Request.Context(), c.Param("id"))
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	if job == nil {
		response.Error(c, http.StatusNotFound, "not found")
		return
	}
	response.Success(c, job)
}
