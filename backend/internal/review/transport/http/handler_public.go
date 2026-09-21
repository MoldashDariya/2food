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

// Create godoc
// @Summary      Create a new review
// @Description  Creates a new review for a specific order
// @Tags         reviews
// @Security     Bearer
// @Accept       json
// @Produce      json
// @Param        input body CreateReviewRequest true "Review details"
// @Success      201 {object} ReviewResponse "Created review"
// @Failure      400 {object} respond.ErrorResponse "Validation error, invalid JSON, or stock issues"
// @Failure      401 {object} respond.ErrorResponse "Unauthorized"
// @Failure      403 {object} respond.ErrorResponse "Forbidden"
// @Failure      404 {object} respond.ErrorResponse "Order not found"
// @Failure      409 {object} respond.ErrorResponse "Review already exists"
// @Failure      500 {object} respond.ErrorResponse "Internal server error"
// @Router       /orders/{orderID}/review [post]
func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	const op = "review.Handler.Create"
	log := h.log.With(slog.String("op", op))

	userID, ok := ctxkey.GetUserID(r.Context())
	if !ok {
		respond.Error(w, r, log, domain.ErrInternalServerError)
		return
	}

	orderID, err := decode.UUID(r, params.OrderID)
	if err != nil {
		respond.Error(w, r, log, err)
		return
	}

	var req CreateReviewRequest
	if err := render.DecodeJSON(r.Body, &req); err != nil {
		respond.Error(w, r, log, domain.ErrInvalidJSON)
		return
	}

	if err := h.validate.Struct(req); err != nil {
		respond.Error(w, r, log, err)
		return
	}

	input := domain.CreateReviewInput{
		UserID:    userID,
		OrderID:   orderID,
		Rating:    req.Rating,
		Comment:   req.Comment,
		PhotoURLs: req.PhotoURLs,
	}

	review, err := h.service.Create(r.Context(), input)
	if err != nil {
		respond.Error(w, r, log, err)
		return
	}

	render.Status(r, http.StatusCreated)
	render.JSON(w, r, ReviewResponse{Review: review})
}

// GetByRestaurant godoc
// @Summary      Get reviews by restaurant
// @Description  Retrieves all reviews for a specific restaurant
// @Tags         reviews
// @Security     Bearer
// @Produce      json
// @Param        restaurantID path string true "Restaurant ID" format(uuid)
// @Success      200 {object} ReviewsResponse "Reviews"
// @Failure      400 {object} respond.ErrorResponse "Invalid UUID"
// @Failure      401 {object} respond.ErrorResponse "Unauthorized"
// @Failure      404 {object} respond.ErrorResponse "Restaurant not found"
// @Failure      500 {object} respond.ErrorResponse "Internal server error"
// @Router       /restaurants/{restaurantID}/reviews [get]
func (h *Handler) GetByRestaurant(w http.ResponseWriter, r *http.Request) {
	const op = "review.Handler.GetByRestaurant"
	log := h.log.With(slog.String("op", op))

	restID, err := decode.UUID(r, params.RestaurantID)
	if err != nil {
		respond.Error(w, r, log, err)
		return
	}

	reviews, err := h.service.GetByRestaurant(r.Context(), restID)
	if err != nil {
		respond.Error(w, r, log, err)
		return
	}

	render.JSON(w, r, ReviewsResponse{Reviews: reviews})
}

// GetByUser godoc
// @Summary      Get reviews by user
// @Description  Retrieves all reviews for a specific user
// @Tags         reviews
// @Security     Bearer
// @Produce      json
// @Param        userID path string true "User ID" format(uuid)
// @Success      200 {object} ReviewsResponse "Reviews"
// @Failure      400 {object} respond.ErrorResponse "Invalid UUID"
// @Failure      401 {object} respond.ErrorResponse "Unauthorized"
// @Failure      500 {object} respond.ErrorResponse "Internal server error"
// @Router       /users/me/reviews [get]
func (h *Handler) GetByUser(w http.ResponseWriter, r *http.Request) {
	const op = "review.Handler.GetByUser"
	log := h.log.With(slog.String("op", op))

	userID, ok := ctxkey.GetUserID(r.Context())
	if !ok {
		respond.Error(w, r, log, domain.ErrInternalServerError)
		return
	}

	reviews, err := h.service.GetByUser(r.Context(), userID)
	if err != nil {
		respond.Error(w, r, log, err)
		return
	}

	render.JSON(w, r, ReviewsResponse{Reviews: reviews})
}

// GetByOrder godoc
// @Summary      Get review by order
// @Description  Retrieves a review for a specific order
// @Tags         reviews
// @Security     Bearer
// @Produce      json
// @Param        orderID path string true "Order ID" format(uuid)
// @Success      200 {object} ReviewResponse "Review"
// @Failure      400 {object} respond.ErrorResponse "Invalid UUID"
// @Failure      401 {object} respond.ErrorResponse "Unauthorized"
// @Failure      403 {object} respond.ErrorResponse "Forbidden"
// @Failure      404 {object} respond.ErrorResponse "Order not found"
// @Failure      500 {object} respond.ErrorResponse "Internal server error"
// @Router       /orders/{orderID}/review [get]
func (h *Handler) GetByOrder(w http.ResponseWriter, r *http.Request) {
	const op = "review.Handler.GetByOrder"
	log := h.log.With(slog.String("op", op))

	orderID, err := decode.UUID(r, params.OrderID)
	if err != nil {
		respond.Error(w, r, log, err)
		return
	}

	userID, ok := ctxkey.GetUserID(r.Context())
	if !ok {
		respond.Error(w, r, log, domain.ErrInternalServerError)
		return
	}

	review, err := h.service.GetByOrder(r.Context(), userID, orderID)
	if err != nil {
		respond.Error(w, r, log, err)
		return
	}

	render.JSON(w, r, ReviewResponse{Review: review})
}

// GetPhotosByRestaurant godoc
// @Summary      Get restaurant review photos
// @Description  Returns a flat list of all photos from active restaurant reviews
// @Tags         restaurants
// @Security     Bearer
// @Produce      json
// @Param        restaurantID path string true "Restaurant ID" format(uuid)
// @Success      200 {object} ReviewPhotosResponse "List of photos"
// @Failure      400 {object} respond.ErrorResponse "Invalid UUID"
// @Failure      401 {object} respond.ErrorResponse "Unauthorized"
// @Failure      500 {object} respond.ErrorResponse "Internal server error"
// @Router       /restaurants/{restaurantID}/reviews/photos [get]
func (h *Handler) GetPhotosByRestaurant(w http.ResponseWriter, r *http.Request) {
	const op = "review.Handler.GetPhotosByRestaurant"
	log := h.log.With(slog.String("op", op))

	restID, err := decode.UUID(r, params.RestaurantID)
	if err != nil {
		respond.Error(w, r, log, err)
		return
	}

	photos, err := h.service.GetPhotosByRestaurant(r.Context(), restID)
	if err != nil {
		respond.Error(w, r, log, err)
		return
	}

	render.JSON(w, r, ReviewPhotosResponse{Photos: photos})
}
