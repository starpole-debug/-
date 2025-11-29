package auth

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"
	"crypto/rand"
	"math/big"

	"github.com/example/ai-avatar-studio/internal/model"
	"github.com/example/ai-avatar-studio/internal/pkg/mailer"
	"github.com/example/ai-avatar-studio/internal/pkg/middleware"
	"github.com/example/ai-avatar-studio/internal/pkg/password"
	"github.com/example/ai-avatar-studio/internal/repository"
)

// Service wires user repository logic for registration/login flows.
type Service struct {
	users          *repository.UserRepository
	verifications  *repository.VerificationRepository
	mailer         mailer.Client
	jwtSecret      string
	adminSecret    string
	tokenDuration  time.Duration
	codeTTL        time.Duration
}

func NewService(users *repository.UserRepository, verifications *repository.VerificationRepository, mailer mailer.Client, jwtSecret, adminSecret string) *Service {
	return &Service{
		users:         users,
		verifications: verifications,
		mailer:        mailer,
		jwtSecret:     jwtSecret,
		adminSecret:   adminSecret,
		tokenDuration: 24 * time.Hour,
		codeTTL:       10 * time.Minute,
	}
}

// Register creates a user and returns a JWT token.
func (s *Service) Register(ctx context.Context, username, email, code, rawPassword, nickname string) (*model.User, string, error) {
	username = strings.TrimSpace(strings.ToLower(username))
	email = strings.TrimSpace(strings.ToLower(email))
	if username == "" || email == "" || rawPassword == "" {
		return nil, "", errors.New("missing credentials")
	}
	if !strings.Contains(email, "@") {
		return nil, "", errors.New("invalid email")
	}
	if err := s.validateCode(ctx, email, "signup", code); err != nil {
		return nil, "", err
	}
	if existing, _ := s.users.FindByUsername(ctx, username); existing != nil {
		return nil, "", errors.New("username already taken")
	}
	if existing, _ := s.users.FindByEmail(ctx, email); existing != nil {
		return nil, "", errors.New("email already registered")
	}
	hash, err := password.Hash(rawPassword)
	if err != nil {
		return nil, "", err
	}
	nickname = strings.TrimSpace(nickname)
	if nickname == "" {
		nickname = username
	}
	user := &model.User{Username: username, Email: email, PasswordHash: hash, Nickname: nickname}
	if err := s.users.Create(ctx, user); err != nil {
		return nil, "", err
	}
	_ = s.consumeCode(ctx, email, "signup")
	token, err := middleware.IssueToken(s.jwtSecret, user.ID, user.IsAdmin, s.tokenDuration)
	if err != nil {
		return nil, "", err
	}
	return user, token, nil
}

// Login validates credentials and issues a JWT.
func (s *Service) Login(ctx context.Context, identifier, rawPassword string) (*model.User, string, error) {
	identifier = strings.TrimSpace(strings.ToLower(identifier))
	if identifier == "" || rawPassword == "" {
		return nil, "", errors.New("invalid credentials")
	}
	var user *model.User
	var err error
	user, err = s.users.FindByUsername(ctx, identifier)
	if err != nil {
		return nil, "", errors.New("invalid credentials")
	}
	if user == nil {
		user, err = s.users.FindByEmail(ctx, identifier)
		if err != nil || user == nil {
			return nil, "", errors.New("invalid credentials")
		}
	}
	if user.IsBanned {
		return nil, "", errors.New("account disabled")
	}
	if err := password.Compare(user.PasswordHash, rawPassword); err != nil {
		return nil, "", errors.New("invalid credentials")
	}
	token, err := middleware.IssueToken(s.jwtSecret, user.ID, user.IsAdmin, s.tokenDuration)
	if err != nil {
		return nil, "", err
	}
	return user, token, nil
}

// AdminLogin enforces the admin secret on top of username/password.
func (s *Service) AdminLogin(ctx context.Context, username, rawPassword, adminSecret string) (*model.User, string, error) {
	user, token, err := s.Login(ctx, username, rawPassword)
	if err != nil {
		return nil, "", err
	}
	if !user.IsAdmin || adminSecret != s.adminSecret {
		return nil, "", errors.New("admin secret mismatch")
	}
	return user, token, nil
}

func (s *Service) Profile(ctx context.Context, id string) (*model.User, error) {
	return s.users.FindByID(ctx, id)
}

// UpdateProfile lets users change nickname and avatar.
func (s *Service) UpdateProfile(ctx context.Context, id, nickname, avatarURL string) (*model.User, error) {
	if strings.TrimSpace(id) == "" {
		return nil, errors.New("invalid user")
	}
	return s.users.UpdateProfile(ctx, id, strings.TrimSpace(nickname), strings.TrimSpace(avatarURL))
}

// SendVerificationCode sends a signup or reset code to the email.
func (s *Service) SendVerificationCode(ctx context.Context, email, purpose string) error {
	if s.verifications == nil || s.mailer == nil {
		return errors.New("email service unavailable")
	}
	email = strings.TrimSpace(strings.ToLower(email))
	if email == "" || !strings.Contains(email, "@") {
		return errors.New("invalid email")
	}
	purpose = strings.ToLower(strings.TrimSpace(purpose))
	if purpose != "signup" && purpose != "reset" {
		return errors.New("invalid purpose")
	}
	if purpose == "signup" {
		if existing, _ := s.users.FindByEmail(ctx, email); existing != nil {
			return errors.New("email already registered")
		}
	} else if purpose == "reset" {
		if existing, _ := s.users.FindByEmail(ctx, email); existing == nil {
			return errors.New("user not found")
		}
	}
	code := generateCode()
	if err := s.verifications.Upsert(ctx, email, purpose, code, time.Now().Add(s.codeTTL)); err != nil {
		return err
	}
	subject := "Your verification code"
	if purpose == "reset" {
		subject = "Password reset code"
	}
	body := fmt.Sprintf("您的验证码是 %s ，有效期10分钟。如非本人操作请忽略。", code)
	return s.mailer.Send(email, subject, body)
}

// ResetPassword allows users to set a new password using a verified code.
func (s *Service) ResetPassword(ctx context.Context, email, code, newPassword string) error {
	email = strings.TrimSpace(strings.ToLower(email))
	if email == "" || newPassword == "" {
		return errors.New("missing email or password")
	}
	if err := s.validateCode(ctx, email, "reset", code); err != nil {
		return err
	}
	user, err := s.users.FindByEmail(ctx, email)
	if err != nil || user == nil {
		return errors.New("user not found")
	}
	hash, err := password.Hash(newPassword)
	if err != nil {
		return err
	}
	if err := s.users.UpdatePassword(ctx, user.ID, hash); err != nil {
		return err
	}
	return s.consumeCode(ctx, email, "reset")
}

func (s *Service) validateCode(ctx context.Context, email, purpose, code string) error {
	if s.verifications == nil {
		return errors.New("verification unavailable")
	}
	code = strings.TrimSpace(code)
	if code == "" {
		return errors.New("verification code required")
	}
	record, err := s.verifications.Find(ctx, email, purpose)
	if err != nil {
		return err
	}
	if record == nil || record.Code != code {
		return errors.New("invalid verification code")
	}
	if record.ConsumedAt != nil {
		return errors.New("verification code already used")
	}
	if time.Now().After(record.ExpiresAt) {
		return errors.New("verification code expired")
	}
	return nil
}

func (s *Service) consumeCode(ctx context.Context, email, purpose string) error {
	record, err := s.verifications.Find(ctx, email, purpose)
	if err != nil || record == nil {
		return err
	}
	return s.verifications.Consume(ctx, record.ID)
}

func generateCode() string {
	const digits = "0123456789"
	b := make([]byte, 6)
	for i := range b {
		n, err := rand.Int(rand.Reader, big.NewInt(int64(len(digits))))
		if err != nil {
			b[i] = digits[0]
			continue
		}
		b[i] = digits[n.Int64()]
	}
	return string(b)
}
