package oauth

import "time"

type Config struct {
	SessionTTL time.Duration `default:"15m" envconfig:"OAUTH_SESSION_TTL"`
	KeyPrefix  string        `default:"oauth:" envconfig:"OAUTH_KEY_PREFIX"`
}
