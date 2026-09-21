package http

import (
	"log/slog"
	"net/http"

	"github.com/go-chi/render"
	"gitlab.com/2foood/backend/internal/domain"
	"gitlab.com/2foood/backend/internal/pkg/decode"
	"gitlab.com/2foood/backend/internal/pkg/params"
	"gitlab.com/2foood/backend/internal/pkg/respond"
)

// Create godoc
// @Summary      Create a new restaurant
// @Description  Creates a new restaurant and assigns the specified user as the owner. Caller must be an admin.
// @Tags         restaurants
// @Security     Bearer
// @Accept       json
// @Produce      json
// @Param        input body CreateRestaurantRequest true "Restaurant data and owner ID"
// @Success      201 {object} RestaurantResponse "Successfully created restaurant"
// @Failure      400 {object} respond.ErrorResponse "Validation error or invalid JSON"
// @Failure      401 {object} respond.ErrorResponse "Unauthorized"
// @Failure      403 {object} respond.ErrorResponse "Permission denied"
// @Failure      500 {object} respond.ErrorResponse "Internal server error"
// @Router       /admin/restaurants [post]
func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	const op = "restaurant.Handler.Create"
	log := h.log.With(slog.String("op", op))

	var req CreateRestaurantRequest
	if err := render.DecodeJSON(r.Body, &req); err != nil {
		respond.Error(w, r, log, domain.ErrInvalidJSON)
		return
	}

	if err := h.validate.Struct(req); err != nil {
		respond.Error(w, r, log, err)
		return
	}

	input := domain.CreateRestaurantInput{
		Name:         req.Name,
		Description:  req.Description,
		Address:      req.Address,
		Phone:        req.Phone,
		Website:      req.Website,
		ImageURL:     req.ImageURL,
		Latitude:     req.Latitude,
		Longitude:    req.Longitude,
		WorkingHours: req.WorkingHours,
		SocialLinks:  req.SocialLinks,
		CityID:       req.CityID,
	}

	rest, err := h.restService.Create(r.Context(), req.OwnerID, input)
	if err != nil {
		respond.Error(w, r, log, err)
		return
	}

	render.Status(r, http.StatusCreated)
	render.JSON(w, r, RestaurantResponse{Restaurant: rest})
}

// Delete godoc
// @Summary      Delete a restaurant
// @Description  Deletes a restaurant from the system. Caller must be an admin.
// @Tags         restaurants
// @Security     Bearer
// @Produce      json
// @Param        id path string true "Restaurant ID" format(uuid)
// @Success      204 "No Content - Restaurant successfully deleted"
// @Failure      400 {object} respond.ErrorResponse "Invalid UUID"
// @Failure      401 {object} respond.ErrorResponse "Unauthorized"
// @Failure      403 {object} respond.ErrorResponse "Permission denied"
// @Failure      404 {object} respond.ErrorResponse "Restaurant not found"
// @Failure      500 {object} respond.ErrorResponse "Internal server error"
// @Router       /admin/restaurants/{id} [delete]
func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	const op = "restaurant.Handler.Delete"
	log := h.log.With(slog.String("op", op))

	restID, err := decode.UUID(r, params.RestaurantID)
	if err != nil {
		respond.Error(w, r, log, err)
		return
	}

	if err := h.restService.Delete(r.Context(), restID); err != nil {
		respond.Error(w, r, log, err)
		return
	}

	render.NoContent(w, r)
}

// AssignOwner godoc
// @Summary      Assign owner to a restaurant
// @Description  Assigns an owner to a restaurant. Caller must be an admin.
// @Tags         restaurants
// @Security     Bearer
// @Produce      json
// @Param        id path string true "Restaurant ID" format(uuid)
// @Param        input body AssignOwnerRequest true "Owner data"
// @Success      200 {object} respond.MessageResponse "Owner assigned successfully"
// @Failure      400 {object} respond.ErrorResponse "Validation error or invalid JSON"
// @Failure      401 {object} respond.ErrorResponse "Unauthorized"
// @Failure      403 {object} respond.ErrorResponse "Permission denied"
// @Failure      404 {object} respond.ErrorResponse "Restaurant not found"
// @Failure      500 {object} respond.ErrorResponse "Internal server error"
// @Router       /admin/restaurants/{id}/owner [patch]
func (h *Handler) AssignOwner(w http.ResponseWriter, r *http.Request) {
	const op = "restaurant.Handler.AssignOwner"
	log := h.log.With(slog.String("op", op))

	restID, err := decode.UUID(r, params.RestaurantID)
	if err != nil {
		respond.Error(w, r, log, err)
		return
	}

	var req AssignOwnerRequest
	if err := render.DecodeJSON(r.Body, &req); err != nil {
		respond.Error(w, r, log, domain.ErrInvalidJSON)
		return
	}

	if err := h.validate.Struct(req); err != nil {
		respond.Error(w, r, log, err)
		return
	}

	if err := h.staffService.AssignOwner(r.Context(), restID, req.UserID); err != nil {
		respond.Error(w, r, log, err)
		return
	}

	respond.Message(w, r, "owner assigned")
}
