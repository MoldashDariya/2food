package jwt

import (
	"time"
)

type Config struct {
	Secret          string        `envconfig:"JWT_SECRET" required:"true"`
	AccessTokenTTL  time.Duration `default:"15m" envconfig:"JWT_ACCESS_TOKEN_TTL"`
	RefreshTokenTTL time.Duration `default:"720h" envconfig:"JWT_REFRESH_TOKEN_TTL"`
}
