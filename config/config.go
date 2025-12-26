package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

const (
	DefaultHTTPAddr            = ":8080"
	DefaultPublicKeyPath       = "public.pem"
	DefaultPrivateKeyPath      = "private.pem"
	DefaultJWTIssuer           = "favourites-api"
	DefaultJWTAudience         = "web"
	DefaultPaginationPage      = 1
	DefaultPaginationLimit     = 20
	DefaultPaginationMaxLimit  = 100
)

type AppConfig struct {
	HTTPAddr   string
	Auth       AuthConfig
	Pagination PaginationConfig
}

type AuthConfig struct {
	PublicKeyPath  string
	PrivateKeyPath string
	Issuer         string
	Audience       string
}

type PaginationConfig struct {
	DefaultPage  int
	DefaultLimit int
	MaxLimit     int
}

func Load() (AppConfig, error) {
	cfg := AppConfig{
		HTTPAddr: envOrDefault("HTTP_ADDR", DefaultHTTPAddr),
		Auth: AuthConfig{
			PublicKeyPath:  envOrDefault("PUBLIC_KEY_PATH", DefaultPublicKeyPath),
			PrivateKeyPath: envOrDefault("PRIVATE_KEY_PATH", DefaultPrivateKeyPath),
			Issuer:         envOrDefault("JWT_ISSUER", DefaultJWTIssuer),
			Audience:       envOrDefault("JWT_AUDIENCE", DefaultJWTAudience),
		},
		Pagination: PaginationConfig{
			DefaultPage:  DefaultPaginationPage,
			DefaultLimit: DefaultPaginationLimit,
			MaxLimit:     DefaultPaginationMaxLimit,
		},
	}

	if v, ok, err := intFromEnv("PAGINATION_DEFAULT_LIMIT"); err != nil {
		return AppConfig{}, err
	} else if ok {
		cfg.Pagination.DefaultLimit = v
	}
	if v, ok, err := intFromEnv("PAGINATION_MAX_LIMIT"); err != nil {
		return AppConfig{}, err
	} else if ok {
		cfg.Pagination.MaxLimit = v
	}
	if v, ok, err := intFromEnv("PAGINATION_DEFAULT_PAGE"); err != nil {
		return AppConfig{}, err
	} else if ok {
		cfg.Pagination.DefaultPage = v
	}

	if err := validate(cfg); err != nil {
		return AppConfig{}, err
	}

	return cfg, nil
}

func validate(cfg AppConfig) error {
	if strings.TrimSpace(cfg.Auth.Issuer) == "" {
		return fmt.Errorf("jwt issuer is required")
	}
	if strings.TrimSpace(cfg.Auth.Audience) == "" {
		return fmt.Errorf("jwt audience is required")
	}
	if strings.TrimSpace(cfg.Auth.PublicKeyPath) == "" {
		return fmt.Errorf("public key path is required")
	}
	if cfg.Pagination.DefaultPage <= 0 {
		return fmt.Errorf("pagination default page must be > 0")
	}
	if cfg.Pagination.DefaultLimit <= 0 {
		return fmt.Errorf("pagination default limit must be > 0")
	}
	if cfg.Pagination.MaxLimit <= 0 {
		return fmt.Errorf("pagination max limit must be > 0")
	}
	if cfg.Pagination.DefaultLimit > cfg.Pagination.MaxLimit {
		return fmt.Errorf("pagination default limit (%d) cannot exceed max limit (%d)", cfg.Pagination.DefaultLimit, cfg.Pagination.MaxLimit)
	}
	return nil
}

func envOrDefault(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}

// intFromEnv returns (value, true, nil) if key exists and is a valid int, (0, false, nil) if unset.
func intFromEnv(key string) (int, bool, error) {
	v := os.Getenv(key)
	if v == "" {
		return 0, false, nil
	}
	num, err := strconv.Atoi(v)
	if err != nil {
		return 0, false, fmt.Errorf("%s must be an integer: %w", key, err)
	}
	return num, true, nil
}
