package config

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"
	"strconv"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
)

// Config aggregates server level configuration with sane defaults for local dev.
type Config struct {
	Env                  string
	ServerPort           string
	DatabaseURL          string
	DatabaseMaxConns     int32
	RedisAddr            string
	RedisPassword        string
	JWTSecret            string
	AdminSecret          string
	AdminAccessKey       string
	FrontendOrigin       []string
	DefaultModelID       string
	AdminBootstrapUser   string
	AdminBootstrapPass   string
	AdminBootstrapEmail  string
	AdminBootstrapNick   string
	UploadDir            string
	PayGateway           string
	PayMerchantID        string
	PayKey               string
	PayNotifyURL         string
	PayReturnURL         string
	CoinsPerYuan         int64
	SMTPHost             string
	SMTPPort             int
	SMTPUser             string
	SMTPPass             string
	SMTPFrom             string
}

// Load reads environment variables and .env if present.
func Load() (*Config, error) {
	_ = godotenv.Load(".env", "../.env", "backend/.env")
	jwtSecret := strings.TrimSpace(os.Getenv("AUTH_JWT_SECRET"))
	if jwtSecret == "" {
		jwtSecret = getEnv("JWT_SECRET", "dev-secret")
	}
	cfg := &Config{
		Env:                 getEnv("APP_ENV", "local"),
		ServerPort:          getEnv("SERVER_PORT", "8080"),
		DatabaseURL:         getEnv("DB_DSN", "postgres://aiapp:aiapp@localhost:5432/aiapp?sslmode=disable"),
		DatabaseMaxConns:    int32(parseInt(getEnv("DB_MAX_CONNS", "10"), 10)),
		RedisAddr:           getEnv("REDIS_ADDR", "localhost:6379"),
		RedisPassword:       strings.TrimSpace(os.Getenv("REDIS_PASSWORD")),
		JWTSecret:           jwtSecret,
		AdminSecret:         getEnv("ADMIN_SECRET", "admin-secret"),
		AdminAccessKey:      strings.TrimSpace(os.Getenv("ADMIN_ACCESS_KEY")),
		DefaultModelID:      strings.TrimSpace(os.Getenv("DEFAULT_MODEL_ID")),
		AdminBootstrapUser:  strings.TrimSpace(os.Getenv("ADMIN_BOOTSTRAP_USERNAME")),
		AdminBootstrapPass:  strings.TrimSpace(os.Getenv("ADMIN_BOOTSTRAP_PASSWORD")),
		AdminBootstrapEmail: strings.TrimSpace(os.Getenv("ADMIN_BOOTSTRAP_EMAIL")),
		AdminBootstrapNick:  strings.TrimSpace(os.Getenv("ADMIN_BOOTSTRAP_NICKNAME")),
		UploadDir:           getEnv("UPLOAD_DIR", "uploads"),
		PayGateway:          strings.TrimSpace(os.Getenv("PAY_GATEWAY")),
		PayMerchantID:       strings.TrimSpace(os.Getenv("PAY_MERCHANT_ID")),
		PayKey:              strings.TrimSpace(os.Getenv("PAY_KEY")),
		PayNotifyURL:        strings.TrimSpace(os.Getenv("PAY_NOTIFY_URL")),
		PayReturnURL:        strings.TrimSpace(os.Getenv("PAY_RETURN_URL")),
		CoinsPerYuan:        parseInt64(getEnv("COINS_PER_YUAN", "1000"), 1000),
		SMTPHost:            strings.TrimSpace(os.Getenv("SMTP_HOST")),
		SMTPPort:            parseInt(getEnv("SMTP_PORT", "587"), 587),
		SMTPUser:            strings.TrimSpace(os.Getenv("SMTP_USER")),
		SMTPPass:            strings.TrimSpace(os.Getenv("SMTP_PASS")),
		SMTPFrom:            strings.TrimSpace(os.Getenv("SMTP_FROM")),
	}
	origins := getEnv("FRONTEND_ORIGIN", "*")
	for _, o := range strings.Split(origins, ",") {
		o = strings.TrimSpace(o)
		if o != "" {
			cfg.FrontendOrigin = append(cfg.FrontendOrigin, o)
		}
	}
	if len(cfg.FrontendOrigin) == 0 {
		cfg.FrontendOrigin = []string{"*"}
	}
	if cfg.JWTSecret == "" {
		return nil, fmt.Errorf("JWT_SECRET is required")
	}
	if cfg.AdminSecret == "" {
		return nil, fmt.Errorf("ADMIN_SECRET is required")
	}
	if cfg.PayGateway == "" || cfg.PayMerchantID == "" || cfg.PayKey == "" {
		return nil, fmt.Errorf("payment gateway config (PAY_GATEWAY, PAY_MERCHANT_ID, PAY_KEY) is required")
	}

	envLower := strings.ToLower(cfg.Env)
	if envLower == "production" {
		if cfg.JWTSecret == "dev-secret" || len(cfg.JWTSecret) < 24 {
			return nil, fmt.Errorf("production requires a strong JWT_SECRET")
		}
		if cfg.AdminSecret == "admin-secret" || len(cfg.AdminSecret) < 16 {
			return nil, fmt.Errorf("production requires a strong ADMIN_SECRET")
		}
		if cfg.AdminAccessKey == "" {
			return nil, fmt.Errorf("ADMIN_ACCESS_KEY is required in production")
		}
		if contains(cfg.FrontendOrigin, "*") {
			return nil, fmt.Errorf("production FRONTEND_ORIGIN cannot be '*'")
		}
		if cfg.PayMerchantID == "" || cfg.PayMerchantID == "1323" {
			return nil, fmt.Errorf("PAY_MERCHANT_ID must be set for production")
		}
		if cfg.PayKey == "" || cfg.PayKey == "i0AJIXg3Gx4al4N9a0LnLJoN1ad0hg0l" {
			return nil, fmt.Errorf("PAY_KEY must be set for production")
		}
		if cfg.SMTPHost == "" || cfg.SMTPUser == "" || cfg.SMTPPass == "" || cfg.SMTPFrom == "" {
			return nil, fmt.Errorf("SMTP_HOST/SMTP_USER/SMTP_PASS/SMTP_FROM are required in production for email verification")
		}
	}
	return cfg, nil
}

// ConnectPostgres builds a pgx Pool that pings the DB for readiness.
func (c *Config) ConnectPostgres(ctx context.Context) (*pgxpool.Pool, error) {
	cfg, err := pgxpool.ParseConfig(c.DatabaseURL)
	if err != nil {
		return nil, fmt.Errorf("parse db url: %w", err)
	}
	if c.DatabaseMaxConns > 0 {
		cfg.MaxConns = c.DatabaseMaxConns
	} else {
		cfg.MaxConns = 10
	}
	cfg.MaxConnIdleTime = time.Minute
	pool, err := pgxpool.NewWithConfig(ctx, cfg)
	if err != nil {
		return nil, fmt.Errorf("connect postgres: %w", err)
	}
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	if err := pool.Ping(ctx); err != nil {
		return nil, fmt.Errorf("ping postgres: %w", err)
	}
	return pool, nil
}

// NewRedisClient returns a configured redis client instance.
func (c *Config) NewRedisClient() *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     c.RedisAddr,
		Password: c.RedisPassword,
	})
}

func getEnv(key, fallback string) string {
	if v := strings.TrimSpace(os.Getenv(key)); v != "" {
		return v
	}
	return fallback
}

func parseInt64(v string, fallback int64) int64 {
	if v == "" {
		return fallback
	}
	if n, err := strconv.ParseInt(v, 10, 64); err == nil {
		return n
	}
	return fallback
}

func parseInt(v string, fallback int) int {
	if v == "" {
		return fallback
	}
	if n, err := strconv.Atoi(v); err == nil {
		return n
	}
	return fallback
}

func contains(list []string, target string) bool {
	for _, item := range list {
		if item == target || item == "*" {
			return true
		}
	}
	return false
}
