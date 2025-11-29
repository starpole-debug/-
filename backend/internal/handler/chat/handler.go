package chat

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/example/ai-avatar-studio/internal/model"
	"github.com/example/ai-avatar-studio/internal/pkg/middleware"
	"github.com/example/ai-avatar-studio/internal/pkg/response"
	chatsvc "github.com/example/ai-avatar-studio/internal/service/chat"
	"github.com/gin-gonic/gin"
)

// Handler powers chat session + messaging endpoints.
type Handler struct {
	service *chatsvc.Service
	secret  string
}

func NewHandler(service *chatsvc.Service, secret string) *Handler {
	return &Handler{service: service, secret: secret}
}

func (h *Handler) RegisterRoutes(rg *gin.RouterGroup) {
	auth := middleware.Authenticator(h.secret)
	rg.POST("/chat/sessions", auth, h.createSession)
	rg.GET("/chat/sessions", auth, h.listSessions)
	rg.GET("/chat/sessions/:id", auth, h.overview)
	rg.POST("/chat/sessions/:id/messages", auth, h.sendMessage)
	rg.PATCH("/chat/messages/:id", auth, h.updateMessage)
	rg.DELETE("/chat/messages/:id", auth, h.deleteMessage)
	rg.DELETE("/chat/sessions/:id", auth, h.deleteSession)
	rg.POST("/chat/messages/:id/retry", auth, h.retryMessage)
	rg.DELETE("/chat/sessions/:id/messages", auth, h.clearSession)
	rg.PATCH("/chat/sessions/:id/settings", auth, h.updateSettings)
	rg.GET("/chat/models", auth, h.listModels)
}

func (h *Handler) createSession(c *gin.Context) {
	userID := middleware.CurrentUserID(c)
	var req struct {
		RoleID   string `json:"role_id"`
		ModelKey string `json:"model_key"`
		Title    string `json:"title"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "invalid body")
		return
	}
	session, err := h.service.StartSession(c.Request.Context(), userID, req.RoleID, req.ModelKey, req.Title)
	if err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	response.Created(c, session)
}

func (h *Handler) listSessions(c *gin.Context) {
	userID := middleware.CurrentUserID(c)
	sessions, err := h.service.ListSessions(c.Request.Context(), userID)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.Success(c, sessions)
}

func (h *Handler) overview(c *gin.Context) {
	userID := middleware.CurrentUserID(c)
	view, err := h.service.SessionOverview(c.Request.Context(), userID, c.Param("id"))
	if err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	response.Success(c, view)
}

func (h *Handler) sendMessage(c *gin.Context) {
	userID := middleware.CurrentUserID(c)
	var req struct {
		Content string       `json:"content"`
		Preset  *model.Preset `json:"preset"`
		Stream  bool         `json:"stream"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "invalid body")
		return
	}
	if req.Stream {
		flusher, ok := c.Writer.(http.Flusher)
		if !ok {
			response.Error(c, http.StatusInternalServerError, "stream not supported")
			return
		}
		c.Writer.Header().Set("Content-Type", "application/json")
		c.Writer.Header().Set("Transfer-Encoding", "chunked")
		_, _ = c.Writer.Write([]byte{}) // ensure headers are sent
		flusher.Flush()
		_, err := h.service.SendMessageStream(c.Request.Context(), userID, c.Param("id"), req.Content, req.Preset, func(delta, reasoning string) {
			payload, _ := json.Marshal(gin.H{
				"content":   delta,
				"reasoning": reasoning,
			})
			_, _ = c.Writer.Write(append(payload, '\n'))
			flusher.Flush()
		})
		if err != nil {
			log.Printf("chat: stream send failed user=%s session=%s err=%v", userID, c.Param("id"), err)
			return
		}
		_, _ = c.Writer.Write([]byte(`{"done":true}` + "\n"))
		flusher.Flush()
	} else {
		msgs, err := h.service.SendMessage(c.Request.Context(), userID, c.Param("id"), req.Content, req.Preset)
		if err != nil {
			// Log the error with session/user context for easier troubleshooting.
			log.Printf("chat: send message failed user=%s session=%s err=%v", userID, c.Param("id"), err)
			response.Error(c, http.StatusBadRequest, err.Error())
			return
		}
		response.Success(c, msgs)
	}
}

func (h *Handler) updateSettings(c *gin.Context) {
	userID := middleware.CurrentUserID(c)
	var req struct {
		Mode           string   `json:"mode"`
		ModelKey       string   `json:"model_key"`
		Temperature    *float64 `json:"temperature"`
		MaxTokens      *int     `json:"max_tokens"`
		NarrativeFocus *string  `json:"narrative_focus"`
		ActionRichness *string  `json:"action_richness"`
		SFWMode        *bool    `json:"sfw_mode"`
		Immersive      *bool    `json:"immersive"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "invalid body")
		return
	}
	session, err := h.service.UpdateSettings(
		c.Request.Context(),
		userID,
		c.Param("id"),
		req.Mode,
		req.ModelKey,
		chatsvc.SettingsPatch{
			Temperature:    req.Temperature,
			MaxTokens:      req.MaxTokens,
			NarrativeFocus: req.NarrativeFocus,
			ActionRichness: req.ActionRichness,
			SFWMode:        req.SFWMode,
			Immersive:      req.Immersive,
		},
	)
	if err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	response.Success(c, session)
}

func (h *Handler) updateMessage(c *gin.Context) {
	userID := middleware.CurrentUserID(c)
	var req struct {
		Content string `json:"content"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "invalid body")
		return
	}
	if err := h.service.UpdateMessage(c.Request.Context(), userID, c.Param("id"), req.Content); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	response.Success(c, gin.H{"status": "updated"})
}

func (h *Handler) deleteMessage(c *gin.Context) {
	userID := middleware.CurrentUserID(c)
	if err := h.service.DeleteMessage(c.Request.Context(), userID, c.Param("id")); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	response.Success(c, gin.H{"status": "deleted"})
}

func (h *Handler) deleteSession(c *gin.Context) {
	userID := middleware.CurrentUserID(c)
	if err := h.service.DeleteSession(c.Request.Context(), userID, c.Param("id")); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	response.Success(c, gin.H{"status": "deleted"})
}

func (h *Handler) retryMessage(c *gin.Context) {
	userID := middleware.CurrentUserID(c)
	msgs, err := h.service.RetryAssistantMessage(c.Request.Context(), userID, c.Param("id"))
	if err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	response.Success(c, msgs)
}

func (h *Handler) clearSession(c *gin.Context) {
	userID := middleware.CurrentUserID(c)
	if err := h.service.ClearSession(c.Request.Context(), userID, c.Param("id")); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	response.Success(c, gin.H{"status": "cleared"})
}

func (h *Handler) listModels(c *gin.Context) {
	models, err := h.service.ListModels(c.Request.Context())
	if err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	response.Success(c, models)
}
