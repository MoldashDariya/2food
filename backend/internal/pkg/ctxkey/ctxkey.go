package ctxkey

import (
	"context"

	"github.com/google/uuid"
	"gitlab.com/2foood/backend/internal/domain"
)

type ctxKey int

const (
	userIDKey ctxKey = iota
	userRoleKey
	restaurantMemberKey
	langKey
)

const DefaultLanguage = "en"

func WithUserID(ctx context.Context, id uuid.UUID) context.Context {
	return context.WithValue(ctx, userIDKey, id)
}

func GetUserID(ctx context.Context) (uuid.UUID, bool) {
	id, ok := ctx.Value(userIDKey).(uuid.UUID)
	return id, ok
}

func WithUserRole(ctx context.Context, role domain.UserRole) context.Context {
	return context.WithValue(ctx, userRoleKey, role)
}

func GetUserRole(ctx context.Context) (domain.UserRole, bool) {
	role, ok := ctx.Value(userRoleKey).(domain.UserRole)
	return role, ok
}

func WithRestaurantMember(ctx context.Context, staff *domain.RestaurantStaff) context.Context {
	return context.WithValue(ctx, restaurantMemberKey, staff)
}

func GetRestaurantMember(ctx context.Context) (*domain.RestaurantStaff, bool) {
	staff, ok := ctx.Value(restaurantMemberKey).(*domain.RestaurantStaff)
	return staff, ok
}

func WithLang(ctx context.Context, lang string) context.Context {
	return context.WithValue(ctx, langKey, lang)
}

func GetLang(ctx context.Context) string {
	lang, ok := ctx.Value(langKey).(string)
	if !ok || lang == "" {
		return DefaultLanguage
	}
	return lang
}
