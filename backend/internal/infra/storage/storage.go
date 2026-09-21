package storage

import "time"

type Config struct {
	PresignTTL time.Duration `envconfig:"STORAGE_PRESIGN_TTL" required:"true"`
}
