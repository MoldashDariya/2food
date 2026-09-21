package middleware

import (
	"errors"
	"net/http"

	"gitlab.com/2foood/backend/internal/domain"
	"gitlab.com/2foood/backend/internal/pkg/ctxkey"
	"gitlab.com/2foood/backend/internal/pkg/decode"
	"gitlab.com/2foood/backend/internal/pkg/params"
	"gitlab.com/2foood/backend/internal/pkg/respond"
)

// checks that the user is a staff member of the restaurant in the url ({restaurantID}).
// the staff record is put into the context so handlers know the user's role there
func (m *Manager) RequireRestaurantAccess(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		const op = "middleware.RequireRestaurantAccess"

		userID, ok := ctxkey.GetUserID(r.Context())
		if !ok {
			respond.Error(w, r, m.log, domain.ErrInternalServerError)
			return
		}

		restID, err := decode.UUID(r, params.RestaurantID)
		if err != nil {
			respond.Error(w, r, m.log, err)
			return
		}

		if role, ok := ctxkey.GetUserRole(r.Context()); ok && role == domain.UserRoleAdmin {
			// admins aren't in the staff table, so we make a fake owner record for them
			// create virtual member for admin
			virtualMember := &domain.RestaurantStaff{
				UserID:       userID,
				RestaurantID: restID,
				Role:         domain.StaffRoleOwner,
			}

			// add virtual member to context and pass it to the next handler
			ctx := ctxkey.WithRestaurantMember(r.Context(), virtualMember)
			next.ServeHTTP(w, r.WithContext(ctx))
			return
		}

		member, err := m.staffRepo.GetMember(r.Context(), restID, userID)
		if err != nil {
			if errors.Is(err, domain.ErrNotFound) {
				respond.Error(w, r, m.log, domain.ErrPermissionDenied)
				return
			}
			respond.Error(w, r, m.log, err)
			return
		}

		ctx := ctxkey.WithRestaurantMember(r.Context(), member)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
