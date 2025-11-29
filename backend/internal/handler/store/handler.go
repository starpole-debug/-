package store

import (
	"net/http"

	"github.com/example/ai-avatar-studio/internal/pkg/middleware"
	"github.com/example/ai-avatar-studio/internal/pkg/response"
	storesvc "github.com/example/ai-avatar-studio/internal/service/store"
	"github.com/gin-gonic/gin"
)

// Handler exposes tipping endpoints consumed by the store UI.
type Handler struct {
	service *storesvc.Service
	secret  string
}

func NewHandler(service *storesvc.Service, secret string) *Handler {
	return &Handler{service: service, secret: secret}
}

func (h *Handler) RegisterRoutes(rg *gin.RouterGroup) {
	rg.GET("/store/options", h.options)
	auth := middleware.Authenticator(h.secret)
	rg.POST("/store/tips", auth, h.tip)
}

func (h *Handler) options(c *gin.Context) {
	response.Success(c, gin.H{
		"amounts":      []int{5, 10, 20},
		"descriptions": []string{"Buy a coffee", "Cheer loudly", "Monthly pass"},
	})
}

func (h *Handler) tip(c *gin.Context) {
	var req struct {
		RoleID string `json:"role_id"`
		Amount int64  `json:"amount"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "invalid body")
		return
	}
	userID := middleware.CurrentUserID(c)
	event, wallet, err := h.service.TipRole(c.Request.Context(), userID, req.RoleID, req.Amount)
	if err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	response.Success(c, gin.H{"event": event, "wallet": wallet})
}
