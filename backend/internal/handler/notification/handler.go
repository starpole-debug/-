package notification

import (
	"net/http"

	"github.com/example/ai-avatar-studio/internal/pkg/middleware"
	"github.com/example/ai-avatar-studio/internal/pkg/response"
	notifsvc "github.com/example/ai-avatar-studio/internal/service/notification"
	"github.com/gin-gonic/gin"
)

// Handler surfaces notification list + mark-read APIs.
type Handler struct {
	service *notifsvc.Service
	secret  string
}

func NewHandler(service *notifsvc.Service, secret string) *Handler {
	return &Handler{service: service, secret: secret}
}

func (h *Handler) RegisterRoutes(rg *gin.RouterGroup) {
	auth := middleware.Authenticator(h.secret)
	rg.GET("/notifications", auth, h.list)
	rg.POST("/notifications/:id/read", auth, h.mark)
	rg.POST("/notifications/all/read", auth, h.markAll)
}

func (h *Handler) list(c *gin.Context) {
	userID := middleware.CurrentUserID(c)
	notifications, err := h.service.List(c.Request.Context(), userID)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.Success(c, notifications)
}

func (h *Handler) mark(c *gin.Context) {
	userID := middleware.CurrentUserID(c)
	if err := h.service.MarkRead(c.Request.Context(), userID, c.Param("id")); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	response.Success(c, gin.H{"status": "ok"})
}

func (h *Handler) markAll(c *gin.Context) {
	userID := middleware.CurrentUserID(c)
	if err := h.service.MarkRead(c.Request.Context(), userID, "all"); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	response.Success(c, gin.H{"status": "ok"})
}
