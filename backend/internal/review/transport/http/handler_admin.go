package http

import (
	"log/slog"
	"net/http"

	"github.com/go-chi/render"
	"gitlab.com/2foood/backend/internal/pkg/decode"
	"gitlab.com/2foood/backend/internal/pkg/params"
	"gitlab.com/2foood/backend/internal/pkg/respond"
)

// GetAll godoc
// @Summary      Get all reviews
// @Description  Retrieves all reviews
// @Tags         reviews
// @Security     Bearer
// @Produce      json
// @Success      200 {object} ReviewsResponse "Reviews"
// @Failure      400 {object} respond.ErrorResponse "Invalid UUID"
// @Failure      401 {object} respond.ErrorResponse "Unauthorized"
// @Failure      403 {object} respond.ErrorResponse "Forbidden"
// @Failure      500 {object} respond.ErrorResponse "Internal server error"
// @Router       /admin/reviews [get]
func (h *Handler) GetAll(w http.ResponseWriter, r *http.Request) {
	const op = "review.Handler.GetAll"
	log := h.log.With(slog.String("op", op))

	reviews, err := h.service.GetAll(r.Context())
	if err != nil {
		respond.Error(w, r, log, err)
		return
	}

	render.JSON(w, r, ReviewsResponse{Reviews: reviews})
}

// Hide godoc
// @Summary      Hide a review
// @Description  Hides a specific review
// @Tags         reviews
// @Security     Bearer
// @Produce      json
// @Param        reviewID path string true "Review ID" format(uuid)
// @Success      200 {object} ReviewResponse "Review hidden"
// @Failure      400 {object} respond.ErrorResponse "Invalid UUID"
// @Failure      401 {object} respond.ErrorResponse "Unauthorized"
// @Failure      403 {object} respond.ErrorResponse "Forbidden"
// @Failure      404 {object} respond.ErrorResponse "Review not found"
// @Failure      500 {object} respond.ErrorResponse "Internal server error"
// @Router       /admin/reviews/{reviewID}/hide [patch]
func (h *Handler) Hide(w http.ResponseWriter, r *http.Request) {
	const op = "review.Handler.Hide"
	log := h.log.With(slog.String("op", op))

	reviewID, err := decode.UUID(r, params.ReviewID)
	if err != nil {
		respond.Error(w, r, log, err)
		return
	}

	if err := h.service.Hide(r.Context(), reviewID); err != nil {
		respond.Error(w, r, log, err)
		return
	}

	render.NoContent(w, r)
}

// Unhide godoc
// @Summary      Unhide a review
// @Description  Unhides a specific review
// @Tags         reviews
// @Security     Bearer
// @Produce      json
// @Param        reviewID path string true "Review ID" format(uuid)
// @Success      200 {object} ReviewResponse "Review unhidden"
// @Failure      400 {object} respond.ErrorResponse "Invalid UUID"
// @Failure      401 {object} respond.ErrorResponse "Unauthorized"
// @Failure      403 {object} respond.ErrorResponse "Forbidden"
// @Failure      404 {object} respond.ErrorResponse "Review not found"
// @Failure      500 {object} respond.ErrorResponse "Internal server error"
// @Router       /admin/reviews/{reviewID}/unhide [patch]
func (h *Handler) Unhide(w http.ResponseWriter, r *http.Request) {
	const op = "review.Handler.Unhide"
	log := h.log.With(slog.String("op", op))

	reviewID, err := decode.UUID(r, params.ReviewID)
	if err != nil {
		respond.Error(w, r, log, err)
		return
	}

	if err := h.service.Unhide(r.Context(), reviewID); err != nil {
		respond.Error(w, r, log, err)
		return
	}

	render.NoContent(w, r)
}
