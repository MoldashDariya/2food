package http

import (
	"log/slog"
	"net/http"

	"github.com/go-chi/render"
	"gitlab.com/2foood/backend/internal/domain"
	"gitlab.com/2foood/backend/internal/pkg/ctxkey"
	"gitlab.com/2foood/backend/internal/pkg/decode"
	"gitlab.com/2foood/backend/internal/pkg/params"
	"gitlab.com/2foood/backend/internal/pkg/respond"
)

// GetByRestaurantAll godoc
// @Summary      Get all reviews by restaurant
// @Description  Retrieves all reviews for a specific restaurant
// @Tags         reviews
// @Security     Bearer
// @Produce      json
// @Param        restaurantID path string true "Restaurant ID" format(uuid)
// @Success      200 {object} ReviewsResponse "Reviews"
// @Failure      400 {object} respond.ErrorResponse "Invalid UUID"
// @Failure      401 {object} respond.ErrorResponse "Unauthorized"
// @Failure      403 {object} respond.ErrorResponse "Forbidden"
// @Failure      404 {object} respond.ErrorResponse "Restaurant not found"
// @Failure      500 {object} respond.ErrorResponse "Internal server error"
// @Router       /partner/restaurants/{restaurantID}/reviews [get]
func (h *Handler) GetByRestaurantAll(w http.ResponseWriter, r *http.Request) {
	const op = "review.Handler.GetByRestaurantAll"
	log := h.log.With(slog.String("op", op))

	restID, err := decode.UUID(r, params.RestaurantID)
	if err != nil {
		respond.Error(w, r, log, err)
		return
	}

	reviews, err := h.service.GetByRestaurantAll(r.Context(), restID)
	if err != nil {
		respond.Error(w, r, log, err)
		return
	}

	render.JSON(w, r, ReviewsResponse{Reviews: reviews})
}

// AddReply godoc
// @Summary      Add a reply to a review
// @Description  Adds a reply to a specific review
// @Tags         reviews
// @Security     Bearer
// @Produce      json
// @Param        reviewID path string true "Review ID" format(uuid)
// @Param        input body AddReplyRequest true "Reply details"
// @Success      200 {object} ReviewResponse "Reply added"
// @Failure      400 {object} respond.ErrorResponse "Validation error or invalid JSON"
// @Failure      401 {object} respond.ErrorResponse "Unauthorized"
// @Failure      403 {object} respond.ErrorResponse "Forbidden"
// @Failure      404 {object} respond.ErrorResponse "Review not found"
// @Failure      500 {object} respond.ErrorResponse "Internal server error"
// @Router       /partner/restaurants/{restaurantID}/reviews/{reviewID}/reply [post]
func (h *Handler) AddReply(w http.ResponseWriter, r *http.Request) {
	const op = "review.Handler.AddReply"
	log := h.log.With(slog.String("op", op))

	actor, ok := ctxkey.GetRestaurantMember(r.Context())
	if !ok {
		respond.Error(w, r, log, domain.ErrInternalServerError)
		return
	}

	reviewID, err := decode.UUID(r, params.ReviewID)
	if err != nil {
		respond.Error(w, r, log, err)
		return
	}

	var req AddReplyRequest
	if err := render.DecodeJSON(r.Body, &req); err != nil {
		respond.Error(w, r, log, domain.ErrInvalidJSON)
		return
	}

	if err := h.validate.Struct(req); err != nil {
		respond.Error(w, r, log, err)
		return
	}

	review, err := h.service.AddReply(r.Context(), actor, reviewID, req.Reply)
	if err != nil {
		respond.Error(w, r, log, err)
		return
	}

	render.JSON(w, r, ReviewResponse{Review: review})
}
