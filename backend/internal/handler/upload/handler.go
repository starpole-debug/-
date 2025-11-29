package upload

import (
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/example/ai-avatar-studio/internal/pkg/middleware"
	"github.com/example/ai-avatar-studio/internal/pkg/response"
	authsvc "github.com/example/ai-avatar-studio/internal/service/auth"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// Handler exposes simple upload endpoints for avatars and post images.
type Handler struct {
	uploadDir string
	auth      *authsvc.Service
	secret    string
}

func NewHandler(uploadDir, secret string, auth *authsvc.Service) *Handler {
	return &Handler{uploadDir: uploadDir, secret: secret, auth: auth}
}

var allowedImageExt = map[string]bool{
	".png":  true,
	".jpg":  true,
	".jpeg": true,
	".gif":  true,
	".webp": true,
}

func (h *Handler) RegisterRoutes(rg *gin.RouterGroup) {
	auth := middleware.Authenticator(h.secret)
	rg.POST("/uploads/avatar", auth, h.uploadAvatar)
	rg.POST("/uploads/posts", auth, h.uploadPostImage)
	rg.POST("/uploads/roles", auth, h.uploadRoleImage)
}

func (h *Handler) uploadAvatar(c *gin.Context) {
	file, header, err := c.Request.FormFile("file")
	if err != nil {
		response.Error(c, http.StatusBadRequest, "missing file")
		return
	}
	defer file.Close()
	if header.Size > 8<<20 {
		response.Error(c, http.StatusBadRequest, "file too large (limit 8MB)")
		return
	}
	ext, ok := normalizeImageExt(header.Filename)
	if !ok {
		response.Error(c, http.StatusBadRequest, "unsupported image type")
		return
	}
	targetDir := filepath.Join(h.uploadDir, "avatars")
	if err := os.MkdirAll(targetDir, 0755); err != nil {
		response.Error(c, http.StatusInternalServerError, "cannot create upload dir")
		return
	}
	filename := uuid.NewString() + ext
	fullPath := filepath.Join(targetDir, filename)
	if err := c.SaveUploadedFile(header, fullPath); err != nil {
		response.Error(c, http.StatusInternalServerError, "save file failed")
		return
	}
	url := "/uploads/avatars/" + filename
	// best effort update user profile avatar
	userID := middleware.CurrentUserID(c)
	if userID != "" && h.auth != nil {
		_, _ = h.auth.UpdateProfile(c.Request.Context(), userID, "", url)
	}
	response.Success(c, gin.H{"url": url})
}

func (h *Handler) uploadPostImage(c *gin.Context) {
	file, header, err := c.Request.FormFile("file")
	if err != nil {
		response.Error(c, http.StatusBadRequest, "missing file")
		return
	}
	defer file.Close()
	if header.Size > 12<<20 {
		response.Error(c, http.StatusBadRequest, "file too large (limit 12MB)")
		return
	}
	ext, ok := normalizeImageExt(header.Filename)
	if !ok {
		response.Error(c, http.StatusBadRequest, "unsupported image type")
		return
	}
	targetDir := filepath.Join(h.uploadDir, "posts")
	if err := os.MkdirAll(targetDir, 0755); err != nil {
		response.Error(c, http.StatusInternalServerError, "cannot create upload dir")
		return
	}
	filename := uuid.NewString() + ext
	fullPath := filepath.Join(targetDir, filename)
	if err := c.SaveUploadedFile(header, fullPath); err != nil {
		response.Error(c, http.StatusInternalServerError, "save file failed")
		return
	}
	url := "/uploads/posts/" + filename
	response.Success(c, gin.H{
		"url":       url,
		"file_name": header.Filename,
		"mime":      header.Header.Get("Content-Type"),
		"size":      header.Size,
	})
}

func (h *Handler) uploadRoleImage(c *gin.Context) {
	file, header, err := c.Request.FormFile("file")
	if err != nil {
		response.Error(c, http.StatusBadRequest, "missing file")
		return
	}
	defer file.Close()
	if header.Size > 8<<20 {
		response.Error(c, http.StatusBadRequest, "file too large (limit 8MB)")
		return
	}
	ext, ok := normalizeImageExt(header.Filename)
	if !ok {
		response.Error(c, http.StatusBadRequest, "unsupported image type")
		return
	}
	targetDir := filepath.Join(h.uploadDir, "roles")
	if err := os.MkdirAll(targetDir, 0755); err != nil {
		response.Error(c, http.StatusInternalServerError, "cannot create upload dir")
		return
	}
	filename := uuid.NewString() + ext
	fullPath := filepath.Join(targetDir, filename)
	if err := c.SaveUploadedFile(header, fullPath); err != nil {
		response.Error(c, http.StatusInternalServerError, "save file failed")
		return
	}
	url := "/uploads/roles/" + filename
	response.Success(c, gin.H{
		"url":       url,
		"file_name": header.Filename,
		"mime":      header.Header.Get("Content-Type"),
		"size":      header.Size,
	})
}

func normalizeImageExt(filename string) (string, bool) {
	ext := strings.ToLower(filepath.Ext(filename))
	if ext == "" {
		ext = ".png"
	}
	_, ok := allowedImageExt[ext]
	return ext, ok
}
