package revenue

import (
	"net/http"

	"github.com/example/ai-avatar-studio/internal/pkg/middleware"
	"github.com/example/ai-avatar-studio/internal/pkg/response"
	revenuesvc "github.com/example/ai-avatar-studio/internal/service/revenue"
	"github.com/gin-gonic/gin"
)

// Handler lets creators inspect wallet + request payouts.
type Handler struct {
	service *revenuesvc.Service
	secret  string
}

func NewHandler(service *revenuesvc.Service, secret string) *Handler {
	return &Handler{service: service, secret: secret}
}

func (h *Handler) RegisterRoutes(rg *gin.RouterGroup) {
	auth := middleware.Authenticator(h.secret)
	rg.GET("/me/wallet", auth, h.wallet)
	rg.POST("/me/payouts", auth, h.requestPayout)
}

func (h *Handler) wallet(c *gin.Context) {
	userID := middleware.CurrentUserID(c)
	wallet, events, payouts, err := h.service.Wallet(c.Request.Context(), userID)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.Success(c, gin.H{"wallet": wallet, "events": events, "payouts": payouts})
}

func (h *Handler) requestPayout(c *gin.Context) {
	var req struct {
		Amount  int64  `json:"amount"`
		Channel string `json:"channel"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "invalid body")
		return
	}
	payout, err := h.service.RequestPayout(c.Request.Context(), middleware.CurrentUserID(c), req.Amount, req.Channel)
	if err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	response.Success(c, payout)
}
