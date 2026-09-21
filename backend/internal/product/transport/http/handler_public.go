package http

import (
	"log/slog"
	"net/http"

	"github.com/go-chi/render"
	"gitlab.com/2foood/backend/internal/pkg/decode"
	"gitlab.com/2foood/backend/internal/pkg/params"
	"gitlab.com/2foood/backend/internal/pkg/respond"
)

// GetActiveByRestaurant godoc
// @Summary      Get active products
// @Description  Retrieves a list of all active products for a specific restaurant
// @Tags         products
// @Produce      json
// @Param        restaurantID path string true "Restaurant ID" format(uuid)
// @Success      200 {object} ProductsResponse "List of active products"
// @Failure      400 {object} respond.ErrorResponse "Invalid UUID"
// @Failure      500 {object} respond.ErrorResponse "Internal server error"
// @Router       /restaurants/{restaurantID}/products [get]
func (h *Handler) GetActiveByRestaurant(w http.ResponseWriter, r *http.Request) {
	const op = "product.Handler.GetActiveByRestaurant"
	log := h.log.With(slog.String("op", op))

	restID, err := decode.UUID(r, params.RestaurantID)
	if err != nil {
		respond.Error(w, r, log, err)
		return
	}

	products, err := h.productService.GetActiveByRestaurant(r.Context(), restID)
	if err != nil {
		respond.Error(w, r, log, err)
		return
	}

	render.JSON(w, r, ProductsResponse{Products: products})
}

func (h *Handler) GetCategoriesByRestaurant(w http.ResponseWriter, r *http.Request) {
	const op = "product.Handler.GetCategoriesByRestaurant"
	log := h.log.With(slog.String("op", op))

	restID, err := decode.UUID(r, params.RestaurantID)
	if err != nil {
		respond.Error(w, r, log, err)
		return
	}

	categories, err := h.productService.GetCategoriesByRestaurant(r.Context(), restID)
	if err != nil {
		respond.Error(w, r, log, err)
		return
	}

	render.JSON(w, r, CategoriesResponse{Categories: categories})
}
