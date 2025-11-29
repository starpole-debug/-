package creator

import (
	"net/http"

	"github.com/example/ai-avatar-studio/internal/pkg/middleware"
	"github.com/example/ai-avatar-studio/internal/pkg/response"
	creatorsvc "github.com/example/ai-avatar-studio/internal/service/creator"
	rolesvc "github.com/example/ai-avatar-studio/internal/service/role"
	"github.com/gin-gonic/gin"
)

// Handler exposes APIs for the creator control center.
type Handler struct {
	service *creatorsvc.Service
	roles   *rolesvc.Service
	secret  string
}

func NewHandler(service *creatorsvc.Service, roles *rolesvc.Service, secret string) *Handler {
	return &Handler{service: service, roles: roles, secret: secret}
}

func (h *Handler) RegisterRoutes(rg *gin.RouterGroup) {
	auth := middleware.Authenticator(h.secret)
	rg.GET("/creator/dashboard", auth, h.dashboard)
	rg.GET("/creator/roles", auth, h.rolesList)
	rg.GET("/creator/roles/:id", auth, h.roleDetail)
}

func (h *Handler) dashboard(c *gin.Context) {
	userID := middleware.CurrentUserID(c)
	data, err := h.service.Dashboard(c.Request.Context(), userID)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.Success(c, data)
}

func (h *Handler) rolesList(c *gin.Context) {
	userID := middleware.CurrentUserID(c)
	roles, err := h.roles.ListByCreator(c.Request.Context(), userID)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.Success(c, roles)
}

func (h *Handler) roleDetail(c *gin.Context) {
	userID := middleware.CurrentUserID(c)
	role, err := h.roles.Get(c.Request.Context(), c.Param("id"))
	if err != nil || role == nil {
		response.Error(c, http.StatusNotFound, "role not found")
		return
	}
	if role.CreatorID != userID {
		response.Error(c, http.StatusForbidden, "forbidden")
		return
	}
	response.Success(c, role)
}
