package otp

import "time"

type Config struct {
	Length    int           `default:"6" envconfig:"OTP_LENGTH"`
	TTL       time.Duration `default:"5m" envconfig:"OTP_TTL"`
	KeyPrefix string        `default:"otp:" envconfig:"OTP_KEY_PREFIX"`

	// rate limitting
	RateLimit     int64         `default:"3"    envconfig:"OTP_RATE_LIMIT"`
	RatePeriod    time.Duration `default:"1h"   envconfig:"OTP_RATE_PERIOD"`
	RateKeyPrefix string        `default:"rate:" envconfig:"OTP_RATE_KEY_PREFIX"`
}
