package middleware

import (
	"net/http"
	"strings"

	"gitlab.com/2foood/backend/internal/domain"
	"gitlab.com/2foood/backend/internal/pkg/ctxkey"
	"gitlab.com/2foood/backend/internal/pkg/respond"
)

// expects `Authorization: Bearer <access token>`.
// puts user id + role into the request context, handlers read them from there (ctxkey)
func (m *Manager) Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			respond.Error(w, r, m.log, domain.ErrInvalidAuthorizationHeader)
			return
		}

		token := strings.TrimPrefix(authHeader, "Bearer ")
		if token == "" {
			respond.Error(w, r, m.log, domain.ErrInvalidAuthorizationHeader)
			return
		}

		claims, err := m.tokenValidator.ValidateToken(token)
		if err != nil {
			respond.Error(w, r, m.log, err)
			return
		}

		// refresh tokens are signed the same way, so check the scope - they can't be used as access tokens
		if claims.Scope != domain.ScopeAccessToken {
			respond.Error(w, r, m.log, domain.ErrInvalidToken)
			return
		}

		ctx := ctxkey.WithUserID(r.Context(), claims.UserID)
		ctx = ctxkey.WithUserRole(ctx, claims.Role)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
