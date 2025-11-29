package auth

import (
	"net/http"

	"github.com/example/ai-avatar-studio/internal/model"
	"github.com/example/ai-avatar-studio/internal/pkg/middleware"
	"github.com/example/ai-avatar-studio/internal/pkg/response"
	authsvc "github.com/example/ai-avatar-studio/internal/service/auth"
	"github.com/gin-gonic/gin"
)

// Handler wires HTTP endpoints for registration/login/profile.
type Handler struct {
	service *authsvc.Service
	secret  string
}

func NewHandler(service *authsvc.Service, secret string) *Handler {
	return &Handler{service: service, secret: secret}
}

func (h *Handler) RegisterRoutes(rg *gin.RouterGroup) {
	rg.POST("/register", h.register)
	rg.POST("/login", h.login)
	rg.POST("/send-code", h.sendCode)
	rg.POST("/password/reset", h.resetPassword)
	auth := middleware.Authenticator(h.secret)
	rg.GET("/me", auth, h.profile)
	rg.PUT("/me", auth, h.updateProfile)
}

func (h *Handler) register(c *gin.Context) {
	var req struct {
		Username string `json:"username"`
		Email    string `json:"email"`
		Code     string `json:"code"`
		Password string `json:"password"`
		Nickname string `json:"nickname"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "invalid body")
		return
	}
	user, token, err := h.service.Register(c.Request.Context(), req.Username, req.Email, req.Code, req.Password, req.Nickname)
	if err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	response.Created(c, gin.H{"user": summarizeUser(user), "token": token})
}

func (h *Handler) login(c *gin.Context) {
	var req struct {
		Identifier string `json:"identifier"` // username or email
		Password   string `json:"password"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "invalid body")
		return
	}
	user, token, err := h.service.Login(c.Request.Context(), req.Identifier, req.Password)
	if err != nil {
		response.Error(c, http.StatusUnauthorized, err.Error())
		return
	}
	response.Success(c, gin.H{"user": summarizeUser(user), "token": token})
}

func (h *Handler) sendCode(c *gin.Context) {
	var req struct {
		Email   string `json:"email"`
		Purpose string `json:"purpose"` // signup | reset
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "invalid body")
		return
	}
	if err := h.service.SendVerificationCode(c.Request.Context(), req.Email, req.Purpose); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	response.Success(c, gin.H{"status": "sent"})
}

func (h *Handler) resetPassword(c *gin.Context) {
	var req struct {
		Email    string `json:"email"`
		Code     string `json:"code"`
		Password string `json:"password"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "invalid body")
		return
	}
	if err := h.service.ResetPassword(c.Request.Context(), req.Email, req.Code, req.Password); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	response.Success(c, gin.H{"status": "reset"})
}

func (h *Handler) profile(c *gin.Context) {
	userID := middleware.CurrentUserID(c)
	user, err := h.service.Profile(c.Request.Context(), userID)
	if err != nil || user == nil {
		response.Error(c, http.StatusNotFound, "user not found")
		return
	}
	response.Success(c, summarizeUser(user))
}

func (h *Handler) updateProfile(c *gin.Context) {
	userID := middleware.CurrentUserID(c)
	var req struct {
		Nickname  string `json:"nickname"`
		AvatarURL string `json:"avatar_url"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "invalid body")
		return
	}
	user, err := h.service.UpdateProfile(c.Request.Context(), userID, req.Nickname, req.AvatarURL)
	if err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	response.Success(c, summarizeUser(user))
}

func summarizeUser(u *model.User) gin.H {
	return gin.H{
		"id":         u.ID,
		"username":   u.Username,
		"email":      u.Email,
		"nickname":   u.Nickname,
		"is_admin":   u.IsAdmin,
		"avatar_url": u.AvatarURL,
	}
}
