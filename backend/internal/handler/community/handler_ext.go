package community

import (
	"net/http"

	"github.com/example/ai-avatar-studio/internal/pkg/middleware"
	"github.com/example/ai-avatar-studio/internal/pkg/response"
	"github.com/gin-gonic/gin"
)

func (h *Handler) getUserProfile(c *gin.Context) {
	user, err := h.service.GetUserProfile(c.Request.Context(), c.Param("id"), middleware.CurrentUserID(c))
	if err != nil {
		response.Error(c, http.StatusNotFound, err.Error())
		return
	}
	response.Success(c, user)
}

func (h *Handler) getUserPosts(c *gin.Context) {
	posts, err := h.service.GetUserPosts(c.Request.Context(), c.Param("id"))
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.Success(c, posts)
}

func (h *Handler) followUser(c *gin.Context) {
	viewer := middleware.CurrentUserID(c)
	target := c.Param("id")
	following, followerCount, followingCount, err := h.service.ToggleFollow(c.Request.Context(), viewer, target)
	if err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	response.Success(c, gin.H{
		"following":        following,
		"follower_count":   followerCount,
		"following_count":  followingCount,
	})
}
