package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"
	"fmt"

	"github.com/example/ai-avatar-studio/internal/config"
	adminhandler "github.com/example/ai-avatar-studio/internal/handler/admin"
	authhandler "github.com/example/ai-avatar-studio/internal/handler/auth"
	chathandler "github.com/example/ai-avatar-studio/internal/handler/chat"
	communityhandler "github.com/example/ai-avatar-studio/internal/handler/community"
	creatorhandler "github.com/example/ai-avatar-studio/internal/handler/creator"
	notificationhandler "github.com/example/ai-avatar-studio/internal/handler/notification"
	presethandler "github.com/example/ai-avatar-studio/internal/handler/preset"
	profilehandler "github.com/example/ai-avatar-studio/internal/handler/profile"
	revenuehandler "github.com/example/ai-avatar-studio/internal/handler/revenue"
	rolehandler "github.com/example/ai-avatar-studio/internal/handler/role"
	storehandler "github.com/example/ai-avatar-studio/internal/handler/store"
	uploadhandler "github.com/example/ai-avatar-studio/internal/handler/upload"
	paymenthandler "github.com/example/ai-avatar-studio/internal/handler/payment"
	"github.com/example/ai-avatar-studio/internal/model"
	"github.com/example/ai-avatar-studio/internal/pkg/llm"
	"github.com/example/ai-avatar-studio/internal/pkg/mailer"
	"github.com/example/ai-avatar-studio/internal/pkg/password"
	"github.com/example/ai-avatar-studio/internal/pkg/redisclient"
	"github.com/example/ai-avatar-studio/internal/repository"
	"github.com/example/ai-avatar-studio/internal/router"
	imageadminhandler "github.com/example/ai-avatar-studio/internal/handler/imageadmin"
	imagehandler "github.com/example/ai-avatar-studio/internal/handler/image"
	adminsvc "github.com/example/ai-avatar-studio/internal/service/admin"
	authsvc "github.com/example/ai-avatar-studio/internal/service/auth"
	chatsvc "github.com/example/ai-avatar-studio/internal/service/chat"
	communitysvc "github.com/example/ai-avatar-studio/internal/service/community"
	creatorsvc "github.com/example/ai-avatar-studio/internal/service/creator"
	memorysvc "github.com/example/ai-avatar-studio/internal/service/memory"
	notificationsvc "github.com/example/ai-avatar-studio/internal/service/notification"
	imagesvc "github.com/example/ai-avatar-studio/internal/service/image"
	presetsvc "github.com/example/ai-avatar-studio/internal/service/preset"
	profilesvc "github.com/example/ai-avatar-studio/internal/service/profile"
	ragservice "github.com/example/ai-avatar-studio/internal/service/rag"
	revenuesvc "github.com/example/ai-avatar-studio/internal/service/revenue"
	rolesvc "github.com/example/ai-avatar-studio/internal/service/role"
	storesvc "github.com/example/ai-avatar-studio/internal/service/store"
	paymentsvc "github.com/example/ai-avatar-studio/internal/service/payment"
	"github.com/example/ai-avatar-studio/internal/task"
	"github.com/gin-gonic/gin"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("load config: %v", err)
	}

	pool, err := cfg.ConnectPostgres(ctx)
	if err != nil {
		log.Fatalf("connect postgres: %v", err)
	}
	defer pool.Close()

	redisNative := cfg.NewRedisClient()
	defer redisNative.Close()
	cache := redisclient.New(redisNative)

	if err := task.RunMigrations(ctx, pool, "migrations"); err != nil {
		log.Fatalf("run migrations: %v", err)
	}

	dispatcher := task.NewDispatcher(32)
	defer dispatcher.Stop()

	notifyURL := cfg.PayNotifyURL
	if notifyURL == "" {
		notifyURL = fmt.Sprintf("http://localhost:%s/api/store/payments/notify", cfg.ServerPort)
	}
	returnURL := cfg.PayReturnURL
	if returnURL == "" {
		returnURL = fmt.Sprintf("http://localhost:%s/api/store/payments/return", cfg.ServerPort)
	}

	// repositories
	userRepo := repository.NewUserRepository(pool)
	roleRepo := repository.NewRoleRepository(pool)
	worldRepo := repository.NewWorldbookRepository(pool)
	chatRepo := repository.NewChatRepository(pool)
	communityRepo := repository.NewCommunityRepository(pool)
	configRepo := repository.NewConfigRepository(pool)
	revenueRepo := repository.NewRevenueRepository(pool)
	notificationRepo := repository.NewNotificationRepository(pool)
	memoryRepo := repository.NewMemoryRepository(pool)
	documentRepo := repository.NewDocumentRepository(pool)
	assetRepo := repository.NewUserAssetRepository(pool)
	paymentRepo := repository.NewPaymentRepository(pool)
	presetRepo := repository.NewPresetRepository(pool)
	verificationRepo := repository.NewVerificationRepository(pool)
	imageProviderRepo := repository.NewImageProviderRepository(pool)
	imagePresetRepo := repository.NewImagePresetRepository(pool)
	imageJobRepo := repository.NewImageJobRepository(pool)

	seedAdminUser(ctx, userRepo, cfg)

	// services
	emailer := mailer.NewSMTPClient(cfg.SMTPHost, cfg.SMTPPort, cfg.SMTPUser, cfg.SMTPPass, cfg.SMTPFrom)
	authService := authsvc.NewService(userRepo, verificationRepo, emailer, cfg.JWTSecret, cfg.AdminSecret)
	roleService := rolesvc.NewService(roleRepo)
	ragService := ragservice.NewService(documentRepo)
	memoryService := memorysvc.NewService(memoryRepo)
	revenueService := revenuesvc.NewService(revenueRepo)
	llmClient := llm.NewRouterClient(&http.Client{Timeout: 60 * time.Second})
	chatService := chatsvc.NewService(chatRepo, roleRepo, worldRepo, configRepo, ragService, memoryService, cache, llmClient, cfg.DefaultModelID, assetRepo, revenueService)
	communityService := communitysvc.NewService(communityRepo, userRepo, notificationRepo, dispatcher, configRepo)
	storeService := storesvc.NewService(roleRepo, revenueService)
	creatorService := creatorsvc.NewService(roleRepo, revenueRepo)
	notificationService := notificationsvc.NewService(notificationRepo)
	adminService := adminsvc.NewService(configRepo, roleRepo, communityRepo, userRepo, notificationRepo)
	profileService := profilesvc.NewService(userRepo, communityRepo, chatRepo, roleRepo, assetRepo, revenueRepo)
	presetService := presetsvc.NewService(presetRepo)
	paymentService := paymentsvc.NewService(paymentRepo, assetRepo, notificationRepo, cfg.PayMerchantID, cfg.PayKey, cfg.PayGateway, notifyURL, returnURL, cfg.CoinsPerYuan)
	imageService := imagesvc.NewService(imageProviderRepo, imagePresetRepo, imageJobRepo, chatRepo, configRepo, llmClient)

	handlers := router.Handlers{
		Auth:         authhandler.NewHandler(authService, cfg.JWTSecret),
		Roles:        rolehandler.NewHandler(roleService, cfg.JWTSecret),
		Chat:         chathandler.NewHandler(chatService, cfg.JWTSecret),
		Community:    communityhandler.NewHandler(communityService, cfg.JWTSecret),
		Creator:      creatorhandler.NewHandler(creatorService, roleService, cfg.JWTSecret),
		Store:        storehandler.NewHandler(storeService, cfg.JWTSecret),
		Notification: notificationhandler.NewHandler(notificationService, cfg.JWTSecret),
		Admin:        adminhandler.NewHandler(adminService, authService, revenueService, cfg.JWTSecret, cfg.AdminAccessKey),
		Revenue:      revenuehandler.NewHandler(revenueService, cfg.JWTSecret),
		Profile:      profilehandler.NewHandler(profileService, cfg.JWTSecret),
		Upload:       uploadhandler.NewHandler(cfg.UploadDir, cfg.JWTSecret, authService),
		Presets:      presethandler.NewPresetHandler(presetService, cfg.JWTSecret),
		Payment:      paymenthandler.NewHandler(paymentService, cfg.JWTSecret),
		Images:       imagehandler.NewHandler(imageService, cfg.JWTSecret),
		ImageAdmin:   imageadminhandler.NewHandler(imageProviderRepo, imagePresetRepo, cfg.JWTSecret),
	}

	engine := router.New(cfg, handlers)
	engine.GET("/api/ready", func(c *gin.Context) {
		rctx, cancel := context.WithTimeout(c.Request.Context(), 2*time.Second)
		defer cancel()
		if err := pool.Ping(rctx); err != nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{"status": "down", "db": err.Error()})
			return
		}
		if _, err := redisNative.Ping(rctx).Result(); err != nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{"status": "down", "redis": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})
	srv := &http.Server{
		Addr:              ":" + cfg.ServerPort,
		Handler:           engine,
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       60 * time.Second,
		WriteTimeout:      180 * time.Second, // allow long LLM calls / streaming
		IdleTimeout:       180 * time.Second,
		MaxHeaderBytes:    1 << 20, // 1MB
	}

	go func() {
		log.Printf("server listening on %s", srv.Addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %v", err)
		}
	}()

	<-ctx.Done()
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Printf("server shutdown error: %v", err)
	}
}

func seedAdminUser(ctx context.Context, repo *repository.UserRepository, cfg *config.Config) {
	username := strings.TrimSpace(cfg.AdminBootstrapUser)
	passwordRaw := strings.TrimSpace(cfg.AdminBootstrapPass)
	if username == "" || passwordRaw == "" {
		return
	}
	email := strings.TrimSpace(cfg.AdminBootstrapEmail)
	if email == "" {
		email = username + "@admin.local"
	}
	nickname := strings.TrimSpace(cfg.AdminBootstrapNick)
	if nickname == "" {
		nickname = "Admin"
	}
	hash, err := password.Hash(passwordRaw)
	if err != nil {
		log.Printf("admin bootstrap hash error: %v", err)
		return
	}
	admin := &model.User{
		Username:     strings.ToLower(username),
		Email:        strings.ToLower(email),
		PasswordHash: hash,
		Nickname:     nickname,
		IsAdmin:      true,
	}
	if err := repo.UpsertAdmin(ctx, admin); err != nil {
		log.Printf("admin bootstrap failed: %v", err)
	} else {
		log.Printf("admin account ensured for %s", admin.Username)
	}
}
