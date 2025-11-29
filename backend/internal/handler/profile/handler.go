package profile

import (
	"net/http"

	"github.com/example/ai-avatar-studio/internal/pkg/middleware"
	"github.com/example/ai-avatar-studio/internal/pkg/response"
	profilesvc "github.com/example/ai-avatar-studio/internal/service/profile"
	"github.com/gin-gonic/gin"
)

// Handler exposes personal homepage endpoints.
type Handler struct {
	service *profilesvc.Service
	secret  string
}

func NewHandler(service *profilesvc.Service, secret string) *Handler {
	return &Handler{service: service, secret: secret}
}

func (h *Handler) RegisterRoutes(rg *gin.RouterGroup) {
	auth := middleware.Authenticator(h.secret)
	rg.GET("/me/home", auth, h.home)
}

func (h *Handler) home(c *gin.Context) {
	userID := middleware.CurrentUserID(c)
	data, err := h.service.Overview(c.Request.Context(), userID)
	if err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	response.Success(c, data)
}
