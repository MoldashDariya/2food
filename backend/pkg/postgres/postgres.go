package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"time"
)

type Config struct {
	User     string `envconfig:"POSTGRES_USER"     required:"true"`
	Password string `envconfig:"POSTGRES_PASSWORD" required:"true"`
	Host     string `envconfig:"POSTGRES_HOST"     required:"true"`
	Port     string `envconfig:"POSTGRES_PORT"     required:"true"`
	DBName   string `envconfig:"POSTGRES_DB_NAME"  required:"true"`
	SSLMode  string `default:"disable" envconfig:"POSTGRES_SSLMODE"`
	PoolConfig
}

type PoolConfig struct {
	MaxOpenConns    int           `default:"25" envconfig:"POSTGRES_MAX_OPEN_CONNS"`
	MaxIdleConns    int           `default:"25" envconfig:"POSTGRES_MAX_IDLE_CONNS"`
	ConnMaxLifeTime time.Duration `default:"5m" envconfig:"POSTGRES_CONN_MAX_LIFE_TIME"`
}

func (c *Config) DSN() string {
	return fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		c.Host, c.User, c.Password, c.DBName, c.Port, c.SSLMode,
	)
}

func New(cfg Config) (*sql.DB, error) {
	db, err := sql.Open("postgres", cfg.DSN())
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// pool settings
	db.SetMaxOpenConns(cfg.MaxOpenConns)
	db.SetMaxIdleConns(cfg.MaxIdleConns)
	db.SetConnMaxIdleTime(cfg.ConnMaxLifeTime)

	// check connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := db.PingContext(ctx); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return db, nil
}
