package http

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"gitlab.com/2foood/backend/internal/domain"
	"gitlab.com/2foood/backend/internal/pkg/ctxkey"
	"gitlab.com/2foood/backend/internal/pkg/respond"
)

type UserService interface {
	Get(ctx context.Context, id uuid.UUID) (*domain.User, error)
	Update(ctx context.Context, id uuid.UUID, input domain.UpdateUserInput) (*domain.User, error)
	UpdatePhone(ctx context.Context, id uuid.UUID, phone string, code string) error
	UpdateEmail(ctx context.Context, id uuid.UUID, email string, code string) error
	Delete(ctx context.Context, id uuid.UUID) error
	// Device token
	RegisterDeviceToken(ctx context.Context, input domain.RegisterDeviceTokenInput) error
	DeleteDeviceToken(ctx context.Context, userID uuid.UUID, token string) error
}

type Handler struct {
	log      *slog.Logger
	service  UserService
	validate *validator.Validate
}

func NewHandler(log *slog.Logger, service UserService, validate *validator.Validate) *Handler {
	return &Handler{
		log:      log,
		service:  service,
		validate: validate,
	}
}

// Get godoc
// @Summary      Get current user
// @Description  Retrieves the profile information of the currently authenticated user
// @Tags         users
// @Security     Bearer
// @Produce      json
// @Success      200 {object} UserResponse "User profile data"
// @Failure      401 {object} respond.ErrorResponse "Unauthorized"
// @Failure      404 {object} respond.ErrorResponse "User not found"
// @Failure      500 {object} respond.ErrorResponse "Internal server error"
// @Router       /users/me [get]
func (h *Handler) Get(w http.ResponseWriter, r *http.Request) {
	const op = "user.Handler.Get"
	log := h.log.With(slog.String("op", op))

	userID, ok := ctxkey.GetUserID(r.Context())
	if !ok {
		respond.Error(w, r, log, domain.ErrInternalServerError)
		return
	}

	user, err := h.service.Get(r.Context(), userID)
	if err != nil {
		respond.Error(w, r, log, err)
		return
	}

	render.JSON(w, r, UserResponse{User: user})
}

// Update godoc
// @Summary      Update current user
// @Description  Updates the profile information of the currently authenticated user
// @Tags         users
// @Security     Bearer
// @Accept       json
// @Produce      json
// @Param        input body UpdateUserRequest true "Data for updating user profile"
// @Success      200 {object} UserResponse "Updated user profile data"
// @Failure      400 {object} respond.ErrorResponse "Validation error or invalid JSON"
// @Failure      401 {object} respond.ErrorResponse "Unauthorized"
// @Failure      404 {object} respond.ErrorResponse "User not found"
// @Failure      500 {object} respond.ErrorResponse "Internal server error"
// @Router       /users/me [patch]
func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	const op = "user.Handler.Update"
	log := h.log.With(slog.String("op", op))

	var req UpdateUserRequest
	if err := render.DecodeJSON(r.Body, &req); err != nil {
		respond.Error(w, r, log, domain.ErrInvalidJSON)
		return
	}

	if err := h.validate.Struct(req); err != nil {
		respond.Error(w, r, log, err)
		return
	}

	userID, ok := ctxkey.GetUserID(r.Context())
	if !ok {
		respond.Error(w, r, log, domain.ErrInternalServerError)
		return
	}

	input := domain.UpdateUserInput{
		Name:                 req.Name,
		AvatarURL:            req.AvatarURL,
		CityID:               req.CityID,
		Language:             req.Language,
		NotificationSettings: req.NotificationSettings,
		DietaryPreferences:   req.DietaryPreferences,
	}

	user, err := h.service.Update(r.Context(), userID, input)
	if err != nil {
		respond.Error(w, r, log, err)
		return
	}

	render.JSON(w, r, UserResponse{User: user})
}

// UpdatePhone godoc
// @Summary      Update current user phone
// @Description  Updates the phone number of the currently authenticated user
// @Tags         users
// @Security     Bearer
// @Accept       json
// @Produce      json
// @Param        input body UpdatePhoneRequest true "Data for updating user phone"
// @Success      200 {object} respond.MessageResponse "Phone updated"
// @Failure      400 {object} respond.ErrorResponse "Validation error or invalid JSON"
// @Failure      401 {object} respond.ErrorResponse "Unauthorized"
// @Failure      500 {object} respond.ErrorResponse "Internal server error"
// @Router       /users/me/phone [patch]
func (h *Handler) UpdatePhone(w http.ResponseWriter, r *http.Request) {
	const op = "user.Handler.UpdatePhone"
	log := h.log.With(slog.String("op", op))

	var req UpdatePhoneRequest
	if err := render.DecodeJSON(r.Body, &req); err != nil {
		respond.Error(w, r, log, domain.ErrInvalidJSON)
		return
	}

	if err := h.validate.Struct(req); err != nil {
		respond.Error(w, r, log, err)
		return
	}

	userID, ok := ctxkey.GetUserID(r.Context())
	if !ok {
		respond.Error(w, r, log, domain.ErrInternalServerError)
		return
	}

	if err := h.service.UpdatePhone(r.Context(), userID, req.Phone, req.Code); err != nil {
		respond.Error(w, r, log, err)
		return
	}

	respond.Message(w, r, "phone updated")
}

// UpdateEmail godoc
// @Summary      Update current user email
// @Description  Updates the email of the currently authenticated user
// @Tags         users
// @Security     Bearer
// @Accept       json
// @Produce      json
// @Param        input body UpdateEmailRequest true "Data for updating user email"
// @Success      200 {object} respond.MessageResponse "Email updated"
// @Failure      400 {object} respond.ErrorResponse "Validation error or invalid JSON"
// @Failure      401 {object} respond.ErrorResponse "Unauthorized"
// @Failure      500 {object} respond.ErrorResponse "Internal server error"
// @Router       /users/me/email [patch]
func (h *Handler) UpdateEmail(w http.ResponseWriter, r *http.Request) {
	const op = "user.Handler.UpdateEmail"
	log := h.log.With(slog.String("op", op))

	var req UpdateEmailRequest
	if err := render.DecodeJSON(r.Body, &req); err != nil {
		respond.Error(w, r, log, domain.ErrInvalidJSON)
		return
	}

	if err := h.validate.Struct(req); err != nil {
		respond.Error(w, r, log, err)
		return
	}

	userID, ok := ctxkey.GetUserID(r.Context())
	if !ok {
		respond.Error(w, r, log, domain.ErrInternalServerError)
		return
	}

	if err := h.service.UpdateEmail(r.Context(), userID, req.Email, req.Code); err != nil {
		respond.Error(w, r, log, err)
		return
	}

	respond.Message(w, r, "email updated")
}

// Delete godoc
// @Summary      Delete current user
// @Description  Deletes the account of the currently authenticated user
// @Tags         users
// @Security     Bearer
// @Success      204 "No Content - User successfully deleted"
// @Failure      401 {object} respond.ErrorResponse "Unauthorized"
// @Failure      500 {object} respond.ErrorResponse "Internal server error"
// @Router       /users/me [delete]
func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	const op = "user.Handler.Delete"
	log := h.log.With(slog.String("op", op))

	userID, ok := ctxkey.GetUserID(r.Context())
	if !ok {
		respond.Error(w, r, log, domain.ErrInternalServerError)
		return
	}

	if err := h.service.Delete(r.Context(), userID); err != nil {
		respond.Error(w, r, log, err)
		return
	}

	render.NoContent(w, r)
}
