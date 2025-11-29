package community

import (
	"net/http"

	"github.com/example/ai-avatar-studio/internal/model"
	"github.com/example/ai-avatar-studio/internal/pkg/middleware"
	"github.com/example/ai-avatar-studio/internal/pkg/response"
	communitysvc "github.com/example/ai-avatar-studio/internal/service/community"
	"github.com/gin-gonic/gin"
)

// Handler wires feeds/comments/reactions.
type Handler struct {
	service *communitysvc.Service
	secret  string
}

func NewHandler(service *communitysvc.Service, secret string) *Handler {
	return &Handler{service: service, secret: secret}
}

func (h *Handler) RegisterRoutes(rg *gin.RouterGroup) {
	optAuth := middleware.OptionalAuth(h.secret)
	rg.GET("/community", optAuth, h.feed)
	rg.GET("/community/:id", optAuth, h.detail)
	rg.GET("/community/dictionary", h.dictionary)
	auth := middleware.Authenticator(h.secret)
	rg.POST("/community", auth, h.createPost)
	rg.POST("/community/:id/comments", auth, h.comment)
	rg.POST("/community/:id/reactions", auth, h.react)
	rg.GET("/me/favorites", auth, h.favorites)

	rg.GET("/community/users/:id", optAuth, h.getUserProfile)
	rg.GET("/community/users/:id/posts", optAuth, h.getUserPosts)
	rg.POST("/community/users/:id/follow", auth, h.followUser)
}

func (h *Handler) feed(c *gin.Context) {
	sort := c.Query("sort")
	filter := c.Query("filter")
	search := c.Query("search")
	userID := middleware.CurrentUserID(c)

	posts, err := h.service.Feed(c.Request.Context(), sort, filter, search, userID)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.Success(c, posts)
}

func (h *Handler) detail(c *gin.Context) {
	post, comments, reactions, err := h.service.Detail(c.Request.Context(), c.Param("id"), middleware.CurrentUserID(c))
	if err != nil {
		response.Error(c, http.StatusNotFound, err.Error())
		return
	}
	response.Success(c, gin.H{"post": post, "comments": comments, "reactions": reactions})
}

func (h *Handler) createPost(c *gin.Context) {
	var payload struct {
		model.CommunityPost
		Attachments []string `json:"attachments"`
	}
	if err := c.ShouldBindJSON(&payload); err != nil {
		response.Error(c, http.StatusBadRequest, "invalid body")
		return
	}
	if payload.CommunityPost.TopicIDs == nil {
		payload.CommunityPost.TopicIDs = []string{}
	}
	payload.CommunityPost.AuthorID = middleware.CurrentUserID(c)
	post, err := h.service.CreatePost(c.Request.Context(), payload.AuthorID, &payload.CommunityPost, payload.Attachments)
	if err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	response.Created(c, post)
}

func (h *Handler) comment(c *gin.Context) {
	var req struct {
		Content string `json:"content"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "invalid body")
		return
	}
	userID := middleware.CurrentUserID(c)
	comment, err := h.service.Comment(c.Request.Context(), userID, c.Param("id"), req.Content)
	if err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	response.Created(c, comment)
}

func (h *Handler) react(c *gin.Context) {
	var req struct {
		Type string `json:"type"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "invalid body")
		return
	}
	userID := middleware.CurrentUserID(c)
	reactions, active, err := h.service.React(c.Request.Context(), userID, c.Param("id"), req.Type)
	if err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	response.Success(c, gin.H{"status": "ok", "reactions": reactions, "activated": active})
}

func (h *Handler) dictionary(c *gin.Context) {
	data, err := h.service.Dictionary(c.Request.Context())
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.Success(c, data)
}

func (h *Handler) favorites(c *gin.Context) {
	userID := middleware.CurrentUserID(c)
	posts, err := h.service.Favorites(c.Request.Context(), userID)
	if err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	response.Success(c, posts)
}
