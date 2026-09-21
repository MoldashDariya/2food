package domain

import (
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
)

type Scope string

const (
	ScopeAccessToken  Scope = "access_token"
	ScopeRefreshToken Scope = "refresh_token"
	ScopePreAuthToken Scope = "pre_auth_token"
)

type TokenPair struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type Claims struct {
	UserID uuid.UUID `json:"user_id"`
	Role   UserRole  `json:"role"`
	Scope  Scope     `json:"scope"`
	jwt.RegisteredClaims
}

type OAuthInfo struct {
	OAuthID       string `json:"oauth_id"`
	Email         string `json:"email"`
	EmailVerified bool   `json:"email_verified"`
	Name          string `json:"name"`
	AvatarURL     string `json:"avatar_url"`
}

type OAuthResult struct {
	User      *User      `json:"user"`
	Tokens    *TokenPair `json:"tokens,omitempty"`
	NeedPhone bool       `json:"need_phone"`
	SessionID string     `json:"session_id,omitempty"` // when user needs to bind phone number
}
