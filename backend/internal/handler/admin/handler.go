package admin

import (
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/example/ai-avatar-studio/internal/model"
	"github.com/example/ai-avatar-studio/internal/pkg/middleware"
	"github.com/example/ai-avatar-studio/internal/pkg/response"
	adminsvc "github.com/example/ai-avatar-studio/internal/service/admin"
	authsvc "github.com/example/ai-avatar-studio/internal/service/auth"
	revenuesvc "github.com/example/ai-avatar-studio/internal/service/revenue"
	"github.com/gin-gonic/gin"
)

// Handler provides admin authentication + CRUD surfaces.
type Handler struct {
	service    *adminsvc.Service
	auth       *authsvc.Service
	secret     string
	accessKey  string
	revenue    *revenuesvc.Service
}

func NewHandler(service *adminsvc.Service, auth *authsvc.Service, revenue *revenuesvc.Service, secret, accessKey string) *Handler {
	return &Handler{service: service, auth: auth, revenue: revenue, secret: secret, accessKey: accessKey}
}

func (h *Handler) RegisterRoutes(rg *gin.RouterGroup) {
	rg.POST("/admin/login", h.login)
	secured := rg.Group("/admin", middleware.AdminOnly(h.secret))
	secured.GET("/models", h.models)
	secured.POST("/models", h.createModel)
	secured.PUT("/models/:id", h.updateModel)
	secured.DELETE("/models/:id", h.deleteModel)
	secured.GET("/dictionary", h.dictionary)
	secured.POST("/dictionary", h.saveDictionary)
	secured.DELETE("/dictionary/:id", h.deleteDictionary)
	secured.GET("/roles", h.roles)
	secured.POST("/posts/:id/hide", h.hidePost)
	secured.POST("/posts/:id/limit", h.limitPost)
	secured.POST("/posts/:id/unhide", h.unhidePost)
	secured.POST("/posts/:id/unlimit", h.unlimitPost)
	secured.GET("/posts", h.listPosts)
	secured.GET("/users", h.listUsers)
	secured.POST("/users", h.createUser)
	secured.POST("/users/:id/ban", h.banUser)
	secured.POST("/users/:id/unban", h.unbanUser)
	secured.DELETE("/users/:id", h.deleteUser)
	secured.GET("/comments", h.listComments)
	secured.POST("/comments/:id/hide", h.hideComment)
	secured.DELETE("/comments/:id", h.deleteComment)
	secured.GET("/payouts", h.listPayouts)
	secured.POST("/payouts/:id/approve", h.approvePayout)
	secured.POST("/payouts/:id/reject", h.rejectPayout)
	secured.POST("/notifications/broadcast", h.broadcastNotifications)
}

func (h *Handler) login(c *gin.Context) {
	var req struct {
		Username    string `json:"username"`
		Email       string `json:"email"`
		Password    string `json:"password"`
		AdminSecret string `json:"admin_secret"`
		AccessKey   string `json:"access_key"`
	}
	if err := c.ShouldBind(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "invalid body")
		return
	}
	providedAccess := strings.TrimSpace(req.AccessKey)
	if providedAccess == "" {
		providedAccess = strings.TrimSpace(c.GetHeader("X-Admin-Access"))
	}
	if h.accessKey != "" && providedAccess != h.accessKey {
		response.Error(c, http.StatusForbidden, "admin access denied")
		return
	}
	identifier := strings.TrimSpace(req.Username)
	if identifier == "" {
		identifier = strings.TrimSpace(req.Email)
	}
	if identifier == "" {
		response.Error(c, http.StatusBadRequest, "username is required")
		return
	}
	user, token, err := h.auth.AdminLogin(c.Request.Context(), identifier, req.Password, req.AdminSecret)
	if err != nil {
		response.Error(c, http.StatusUnauthorized, err.Error())
		return
	}
	response.Success(c, gin.H{"user": user, "token": token})
}

func (h *Handler) models(c *gin.Context) {
	models, err := h.service.Models(c.Request.Context())
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.Success(c, models)
}

func (h *Handler) createModel(c *gin.Context) {
	payload, err := bindModelPayload(c)
	if err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	model, err := h.service.SaveModel(c.Request.Context(), payload)
	if err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	response.Success(c, model)
}

func (h *Handler) updateModel(c *gin.Context) {
	payload, err := bindModelPayload(c)
	if err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	payload.ID = c.Param("id")
	model, err := h.service.SaveModel(c.Request.Context(), payload)
	if err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	response.Success(c, model)
}

func (h *Handler) deleteModel(c *gin.Context) {
	if err := h.service.DeleteModel(c.Request.Context(), c.Param("id")); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	response.Success(c, gin.H{"status": "deleted"})
}

func (h *Handler) dictionary(c *gin.Context) {
	items, err := h.service.Dictionary(c.Request.Context(), c.Query("group"))
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.Success(c, items)
}

func (h *Handler) saveDictionary(c *gin.Context) {
	var payload model.DictionaryItem
	if err := c.ShouldBindJSON(&payload); err != nil {
		response.Error(c, http.StatusBadRequest, "invalid body")
		return
	}
	item, err := h.service.SaveDictionary(c.Request.Context(), &payload)
	if err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	response.Success(c, item)
}

func (h *Handler) deleteDictionary(c *gin.Context) {
	if err := h.service.DeleteDictionary(c.Request.Context(), c.Param("id")); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	response.Success(c, gin.H{"status": "deleted"})
}

func (h *Handler) roles(c *gin.Context) {
	roles, err := h.service.Roles(c.Request.Context())
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.Success(c, roles)
}

func (h *Handler) hidePost(c *gin.Context) {
	if err := h.service.HidePost(c.Request.Context(), c.Param("id")); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	response.Success(c, gin.H{"status": "hidden"})
}

func (h *Handler) limitPost(c *gin.Context) {
	if err := h.service.LimitPost(c.Request.Context(), c.Param("id")); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	response.Success(c, gin.H{"status": "limited"})
}

func (h *Handler) unhidePost(c *gin.Context) {
	if err := h.service.RestorePost(c.Request.Context(), c.Param("id")); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	response.Success(c, gin.H{"status": "visible"})
}

func (h *Handler) unlimitPost(c *gin.Context) {
	if err := h.service.RestorePost(c.Request.Context(), c.Param("id")); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	response.Success(c, gin.H{"status": "visible"})
}

func (h *Handler) listPosts(c *gin.Context) {
	limit := parseIntDefault(c.Query("limit"), 50)
	offset := parseIntDefault(c.Query("offset"), 0)
	posts, err := h.service.ListPosts(c.Request.Context(), c.Query("query"), c.Query("visibility"), limit, offset)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.Success(c, posts)
}

func (h *Handler) listUsers(c *gin.Context) {
	limit := parseIntDefault(c.Query("limit"), 50)
	offset := parseIntDefault(c.Query("offset"), 0)
	users, err := h.service.ListUsers(c.Request.Context(), c.Query("query"), limit, offset)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.Success(c, users)
}

func (h *Handler) createUser(c *gin.Context) {
	var payload struct {
		Username string `json:"username"`
		Email    string `json:"email"`
		Password string `json:"password"`
		Nickname string `json:"nickname"`
		IsAdmin  bool   `json:"is_admin"`
	}
	if err := c.ShouldBindJSON(&payload); err != nil {
		response.Error(c, http.StatusBadRequest, "invalid body")
		return
	}
	user, err := h.service.CreateUser(c.Request.Context(), payload.Username, payload.Email, payload.Password, payload.Nickname, payload.IsAdmin)
	if err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	response.Created(c, user)
}

func (h *Handler) banUser(c *gin.Context) {
	if err := h.service.BanUser(c.Request.Context(), c.Param("id")); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	response.Success(c, gin.H{"status": "banned"})
}

func (h *Handler) unbanUser(c *gin.Context) {
	if err := h.service.UnbanUser(c.Request.Context(), c.Param("id")); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	response.Success(c, gin.H{"status": "active"})
}

func (h *Handler) deleteUser(c *gin.Context) {
	if err := h.service.DeleteUser(c.Request.Context(), c.Param("id")); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	response.Success(c, gin.H{"status": "deleted"})
}

func (h *Handler) listComments(c *gin.Context) {
	limit := parseIntDefault(c.Query("limit"), 50)
	offset := parseIntDefault(c.Query("offset"), 0)
	comments, err := h.service.ListComments(c.Request.Context(), c.Query("post_id"), c.Query("visibility"), limit, offset)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.Success(c, comments)
}

func (h *Handler) hideComment(c *gin.Context) {
	if err := h.service.HideComment(c.Request.Context(), c.Param("id")); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	response.Success(c, gin.H{"status": "hidden"})
}

func (h *Handler) deleteComment(c *gin.Context) {
	if err := h.service.DeleteComment(c.Request.Context(), c.Param("id")); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	response.Success(c, gin.H{"status": "deleted"})
}

func (h *Handler) listPayouts(c *gin.Context) {
	limit := parseIntDefault(c.Query("limit"), 50)
	offset := parseIntDefault(c.Query("offset"), 0)
	status := strings.TrimSpace(c.Query("status"))
	payouts, err := h.revenue.AdminListPayouts(c.Request.Context(), status, limit, offset)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.Success(c, payouts)
}

func (h *Handler) approvePayout(c *gin.Context) {
	payout, wallet, err := h.revenue.AdminUpdatePayout(c.Request.Context(), c.Param("id"), "approved")
	if err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	response.Success(c, gin.H{"payout": payout, "wallet": wallet})
}

func (h *Handler) rejectPayout(c *gin.Context) {
	payout, wallet, err := h.revenue.AdminUpdatePayout(c.Request.Context(), c.Param("id"), "rejected")
	if err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	response.Success(c, gin.H{"payout": payout, "wallet": wallet})
}

func (h *Handler) broadcastNotifications(c *gin.Context) {
	var payload struct {
		Title     string   `json:"title"`
		Content   string   `json:"content"`
		UserIDs   []string `json:"user_ids"`
		Broadcast bool     `json:"broadcast"`
		Query     string   `json:"query"`
		Limit     int      `json:"limit"`
	}
	if err := c.ShouldBindJSON(&payload); err != nil {
		response.Error(c, http.StatusBadRequest, "invalid body")
		return
	}
	if strings.TrimSpace(payload.Title) == "" || strings.TrimSpace(payload.Content) == "" {
		response.Error(c, http.StatusBadRequest, "title and content required")
		return
	}
	created := 0
	var err error
	if payload.Broadcast && len(payload.UserIDs) == 0 {
		created, err = h.service.BroadcastNotification(c.Request.Context(), payload.Title, payload.Content, payload.Query, payload.Limit)
	} else {
		created, err = h.service.SendNotification(c.Request.Context(), payload.Title, payload.Content, payload.UserIDs)
	}
	if err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	response.Success(c, gin.H{"created": created})
}

type modelPayload struct {
	ID          string   `json:"id"`
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Provider    string   `json:"provider"`
	BaseURL     string   `json:"base_url"`
	ModelName   string   `json:"model_name"`
	APIKey      string   `json:"api_key"`
	Temperature *float64 `json:"temperature"`
	MaxTokens   *int     `json:"max_tokens"`
	Status      string   `json:"status"`
	IsDefault   bool     `json:"is_default"`
	IsEnabled   *bool    `json:"is_enabled"`
	PriceCoins  *int64   `json:"price_coins"`
	PriceHint   string   `json:"price_hint"`
	ShareRole   *float64 `json:"share_role_pct"`
	SharePreset *float64 `json:"share_preset_pct"`
}

func bindModelPayload(c *gin.Context) (*model.ModelConfig, error) {
	var body modelPayload
	if err := c.ShouldBindJSON(&body); err != nil {
		return nil, err
	}
	// clamp share percentages
	clampShare := func(v *float64) float64 {
		if v == nil {
			return 0
		}
		if *v < 0 {
			return 0
		}
		if *v > 1 {
			return 1
		}
		return *v
	}
	cfg := &model.ModelConfig{
		ID:          strings.TrimSpace(body.ID),
		Name:        strings.TrimSpace(body.Name),
		Description: strings.TrimSpace(body.Description),
		Provider:    strings.TrimSpace(body.Provider),
		BaseURL:     strings.TrimSpace(body.BaseURL),
		ModelName:   strings.TrimSpace(body.ModelName),
		APIKey:      strings.TrimSpace(body.APIKey),
		Status:      strings.ToLower(strings.TrimSpace(body.Status)),
		IsDefault:   body.IsDefault,
	}
	if cfg.Provider == "" {
		cfg.Provider = "openai"
	}
	if cfg.BaseURL == "" {
		cfg.BaseURL = "https://api.openai.com/v1"
	}
	if cfg.Status == "" {
		cfg.Status = "active"
	}
	if cfg.Name == "" {
		return nil, errors.New("name is required")
	}
	if cfg.ModelName == "" {
		return nil, errors.New("model_name is required")
	}
	if body.Temperature != nil {
		cfg.Temperature = *body.Temperature
	}
	if body.MaxTokens != nil {
		cfg.MaxTokens = *body.MaxTokens
	}
	if body.PriceCoins != nil {
		cfg.PriceCoins = *body.PriceCoins
	}
	cfg.ShareRolePct = clampShare(body.ShareRole)
	cfg.SharePresetPct = clampShare(body.SharePreset)
	cfg.PriceHint = strings.TrimSpace(body.PriceHint)
	if cfg.ShareRolePct+cfg.SharePresetPct > 1.0 {
		return nil, errors.New("角色和预设分成之和不能超过 1")
	}
	if body.IsEnabled != nil {
		cfg.IsEnabled = *body.IsEnabled
	} else {
		cfg.IsEnabled = strings.EqualFold(cfg.Status, "active")
	}
	return cfg, nil
}

func parseIntDefault(raw string, fallback int) int {
	if v, err := strconv.Atoi(raw); err == nil && v >= 0 {
		return v
	}
	return fallback
}
