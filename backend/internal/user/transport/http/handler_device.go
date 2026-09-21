package http

import (
	"log/slog"
	"net/http"

	"github.com/go-chi/render"
	"gitlab.com/2foood/backend/internal/domain"
	"gitlab.com/2foood/backend/internal/pkg/ctxkey"
	"gitlab.com/2foood/backend/internal/pkg/respond"
)

// RegisterDeviceToken godoc
// @Summary      Register device token
// @Description  Registers a push notification device token for the current user
// @Tags         users
// @Security     Bearer
// @Accept       json
// @Produce      json
// @Param        input body RegisterDeviceTokenRequest true "Device token data"
// @Success      200 {object} respond.MessageResponse "Device token registered"
// @Failure      400 {object} respond.ErrorResponse "Validation error or invalid JSON"
// @Failure      401 {object} respond.ErrorResponse "Unauthorized"
// @Failure      500 {object} respond.ErrorResponse "Internal server error"
// @Router       /users/me/devices [post]
func (h *Handler) RegisterDeviceToken(w http.ResponseWriter, r *http.Request) {
	const op = "user.Handler.RegisterDeviceToken"
	log := h.log.With(slog.String("op", op))

	userID, ok := ctxkey.GetUserID(r.Context())
	if !ok {
		respond.Error(w, r, log, domain.ErrInternalServerError)
		return
	}

	var req RegisterDeviceTokenRequest
	if err := render.DecodeJSON(r.Body, &req); err != nil {
		respond.Error(w, r, log, domain.ErrInvalidJSON)
		return
	}

	if err := h.validate.Struct(req); err != nil {
		respond.Error(w, r, log, err)
		return
	}

	input := domain.RegisterDeviceTokenInput{
		UserID:   userID,
		Token:    req.Token,
		Platform: req.Platform,
	}

	if err := h.service.RegisterDeviceToken(r.Context(), input); err != nil {
		respond.Error(w, r, log, err)
		return
	}

	respond.Message(w, r, "device token registered")
}

// DeleteDeviceToken godoc
// @Summary      Delete device token
// @Description  Removes a push notification device token for the current user
// @Tags         users
// @Security     Bearer
// @Accept       json
// @Produce      json
// @Param        input body DeleteDeviceTokenRequest true "Device token to delete"
// @Success      204 "No Content - Token successfully deleted"
// @Failure      400 {object} respond.ErrorResponse "Validation error or invalid JSON"
// @Failure      401 {object} respond.ErrorResponse "Unauthorized"
// @Failure      404 {object} respond.ErrorResponse "Token not found"
// @Failure      500 {object} respond.ErrorResponse "Internal server error"
// @Router       /users/me/devices [delete]
func (h *Handler) DeleteDeviceToken(w http.ResponseWriter, r *http.Request) {
	const op = "user.Handler.DeleteDeviceToken"
	log := h.log.With(slog.String("op", op))

	userID, ok := ctxkey.GetUserID(r.Context())
	if !ok {
		respond.Error(w, r, log, domain.ErrInternalServerError)
		return
	}

	var req DeleteDeviceTokenRequest
	if err := render.DecodeJSON(r.Body, &req); err != nil {
		respond.Error(w, r, log, domain.ErrInvalidJSON)
		return
	}

	if err := h.validate.Struct(req); err != nil {
		respond.Error(w, r, log, err)
		return
	}

	if err := h.service.DeleteDeviceToken(r.Context(), userID, req.Token); err != nil {
		respond.Error(w, r, log, err)
		return
	}

	render.NoContent(w, r)
}
