package middleware

import (
	"net/http"
	"slices"

	"gitlab.com/2foood/backend/internal/domain"
	"gitlab.com/2foood/backend/internal/pkg/ctxkey"
	"gitlab.com/2foood/backend/internal/pkg/respond"
)

// has to run after Auth (it reads the role from the context)
func (m *Manager) RequireRole(requiredRoles ...domain.UserRole) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			role, ok := ctxkey.GetUserRole(r.Context())
			if !ok {
				respond.Error(w, r, m.log, domain.ErrInternalServerError)
				return
			}

			// admin has access to all routes
			if role == domain.UserRoleAdmin {
				next.ServeHTTP(w, r)
				return
			}

			if slices.Contains(requiredRoles, role) {
				next.ServeHTTP(w, r)
				return
			}

			respond.Error(w, r, m.log, domain.ErrPermissionDenied)
		})
	}
}