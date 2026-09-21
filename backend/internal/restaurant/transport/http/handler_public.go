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
// @Summary      Get all restaurants
// @Description  Retrieves a list of all registered restaurants
// @Tags         restaurants
// @Produce      json
// @Success      200 {object} RestaurantsResponse "List of restaurants"
// @Failure      500 {object} respond.ErrorResponse "Internal server error"
// @Router       /restaurants [get]
func (h *Handler) GetAll(w http.ResponseWriter, r *http.Request) {
	const op = "restaurant.Handler.GetAll"
	log := h.log.With(slog.String("op", op))

	restaurants, err := h.restService.GetAll(r.Context())
	if err != nil {
		respond.Error(w, r, log, err)
		return
	}

	render.JSON(w, r, RestaurantsResponse{Restaurants: restaurants})
}

// GetByID godoc
// @Summary      Get restaurant by ID
// @Description  Retrieves detailed information about a specific restaurant
// @Tags         restaurants
// @Produce      json
// @Param        id path string true "Restaurant ID" format(uuid)
// @Success      200 {object} RestaurantResponse "Restaurant details"
// @Failure      400 {object} respond.ErrorResponse "Invalid UUID"
// @Failure      404 {object} respond.ErrorResponse "Restaurant not found"
// @Failure      500 {object} respond.ErrorResponse "Internal server error"
// @Router       /restaurants/{id} [get]
func (h *Handler) GetByID(w http.ResponseWriter, r *http.Request) {
	const op = "restaurant.Handler.GetByID"
	log := h.log.With(slog.String("op", op))

	restID, err := decode.UUID(r, params.RestaurantID)
	if err != nil {
		respond.Error(w, r, log, err)
		return
	}

	restaurant, err := h.restService.GetByID(r.Context(), restID)
	if err != nil {
		respond.Error(w, r, log, err)
		return
	}

	render.JSON(w, r, RestaurantResponse{Restaurant: restaurant})
}
