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

// GetByPartner godoc
// @Summary      Get partner's restaurants
// @Description  Retrieves a list of restaurants associated with the currently authenticated partner
// @Tags         restaurants
// @Security     Bearer
// @Produce      json
// @Success      200 {object} RestaurantsResponse "List of partner's restaurants"
// @Failure      401 {object} respond.ErrorResponse "Unauthorized"
// @Failure      500 {object} respond.ErrorResponse "Internal server error"
// @Router       /partner/restaurants [get]
func (h *Handler) GetByPartner(w http.ResponseWriter, r *http.Request) {
	const op = "restaurant.Handler.GetByPartner"
	log := h.log.With(slog.String("op", op))

	userID, ok := ctxkey.GetUserID(r.Context())
	if !ok {
		respond.Error(w, r, log, domain.ErrInternalServerError)
		return
	}

	rests, err := h.restService.GetByPartner(r.Context(), userID)
	if err != nil {
		respond.Error(w, r, log, err)
		return
	}

	render.JSON(w, r, RestaurantsResponse{Restaurants: rests})
}

// UpdateInfo godoc
// @Summary      Update restaurant info
// @Description  Updates information for a specific restaurant. Caller must be an owner.
// @Tags         restaurants
// @Security     Bearer
// @Accept       json
// @Produce      json
// @Param        id path string true "Restaurant ID" format(uuid)
// @Param        input body UpdateRestaurantRequest true "Updated restaurant data"
// @Success      200 {object} RestaurantResponse "Updated restaurant info"
// @Failure      400 {object} respond.ErrorResponse "Validation error or invalid UUID/JSON"
// @Failure      401 {object} respond.ErrorResponse "Unauthorized"
// @Failure      403 {object} respond.ErrorResponse "Permission denied"
// @Failure      404 {object} respond.ErrorResponse "Restaurant not found"
// @Failure      500 {object} respond.ErrorResponse "Internal server error"
// @Router       /partner/restaurants/{id} [patch]
func (h *Handler) UpdateInfo(w http.ResponseWriter, r *http.Request) {
	const op = "restaurant.Handler.UpdateInfo"
	log := h.log.With(slog.String("op", op))

	restID, err := decode.UUID(r, params.RestaurantID)
	if err != nil {
		respond.Error(w, r, log, err)
		return
	}

	actor, ok := ctxkey.GetRestaurantMember(r.Context())
	if !ok {
		respond.Error(w, r, log, domain.ErrInternalServerError)
		return
	}

	var req UpdateRestaurantRequest
	if err := render.DecodeJSON(r.Body, &req); err != nil {
		respond.Error(w, r, log, domain.ErrInvalidJSON)
		return
	}

	if err := h.validate.Struct(req); err != nil {
		respond.Error(w, r, log, err)
		return
	}

	input := domain.UpdateRestaurantInput{
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

	rest, err := h.restService.UpdateInfo(r.Context(), actor, restID, input)
	if err != nil {
		respond.Error(w, r, log, err)
		return
	}

	render.JSON(w, r, RestaurantResponse{Restaurant: rest})
}

// UpdateCuisines godoc
// @Summary      Update restaurant cuisine types
// @Description  Replaces the cuisine types associated with a restaurant. Caller must be an owner.
// @Tags         restaurants
// @Security     Bearer
// @Accept       json
// @Produce      json
// @Param        id    path  string               true  "Restaurant ID" format(uuid)
// @Param        input body  UpdateCuisinesRequest true  "List of cuisine type IDs"
// @Success      200   {object} RestaurantResponse      "Updated restaurant"
// @Failure      400   {object} respond.ErrorResponse   "Validation error or invalid UUID/JSON"
// @Failure      401   {object} respond.ErrorResponse   "Unauthorized"
// @Failure      403   {object} respond.ErrorResponse   "Permission denied"
// @Failure      404   {object} respond.ErrorResponse   "Restaurant not found"
// @Failure      500   {object} respond.ErrorResponse   "Internal server error"
// @Router       /partner/restaurants/{id}/cuisines [patch]
func (h *Handler) UpdateCuisines(w http.ResponseWriter, r *http.Request) {
	const op = "restaurant.Handler.UpdateCuisines"
	log := h.log.With(slog.String("op", op))

	restID, err := decode.UUID(r, params.RestaurantID)
	if err != nil {
		respond.Error(w, r, log, err)
		return
	}

	actor, ok := ctxkey.GetRestaurantMember(r.Context())
	if !ok {
		respond.Error(w, r, log, domain.ErrInternalServerError)
		return
	}

	var req UpdateCuisinesRequest
	if err := render.DecodeJSON(r.Body, &req); err != nil {
		respond.Error(w, r, log, domain.ErrInvalidJSON)
		return
	}

	if err := h.validate.Struct(req); err != nil {
		respond.Error(w, r, log, err)
		return
	}

	if err = h.restService.UpdateCuisines(r.Context(), actor, restID, req.CuisineIDs); err != nil {
		respond.Error(w, r, log, err)
		return
	}

	respond.Message(w, r, "cuisines updated")
}

// AddStaff godoc
// @Summary      Add staff to restaurant
// @Description  Adds a new staff member to the restaurant. Caller must be an owner.
// @Tags         restaurants, staff
// @Security     Bearer
// @Accept       json
// @Produce      json
// @Param        id path string true "Restaurant ID" format(uuid)
// @Param        input body AddStaffRequest true "Staff data"
// @Success      201 {object} respond.MessageResponse "Staff added successfully"
// @Failure      400 {object} respond.ErrorResponse "Validation error or invalid UUID/JSON"
// @Failure      401 {object} respond.ErrorResponse "Unauthorized"
// @Failure      403 {object} respond.ErrorResponse "Permission denied"
// @Failure      409 {object} respond.ErrorResponse "Staff member already exists"
// @Failure      500 {object} respond.ErrorResponse "Internal server error"
// @Router       /partner/restaurants/{id}/staff [post]
func (h *Handler) AddStaff(w http.ResponseWriter, r *http.Request) {
	const op = "restaurant.Handler.AddStaff"
	log := h.log.With(slog.String("op", op))

	restID, err := decode.UUID(r, params.RestaurantID)
	if err != nil {
		respond.Error(w, r, log, err)
		return
	}

	actor, ok := ctxkey.GetRestaurantMember(r.Context())
	if !ok {
		respond.Error(w, r, log, domain.ErrInternalServerError)
		return
	}

	var req AddStaffRequest
	if err := render.DecodeJSON(r.Body, &req); err != nil {
		respond.Error(w, r, log, domain.ErrInvalidJSON)
		return
	}

	if err := h.validate.Struct(req); err != nil {
		respond.Error(w, r, log, err)
		return
	}

	if err = h.staffService.AddStaff(r.Context(), actor, restID, req.UserID, req.Role); err != nil {
		respond.Error(w, r, log, err)
		return
	}

	render.Status(r, http.StatusCreated)
	respond.Message(w, r, "staff added")
}

// GetStaff godoc
// @Summary      Get restaurant staff
// @Description  Retrieves a list of all staff members for a specific restaurant
// @Tags         restaurants, staff
// @Security     Bearer
// @Produce      json
// @Param        id path string true "Restaurant ID" format(uuid)
// @Success      200 {object} StaffResponse "List of staff members"
// @Failure      400 {object} respond.ErrorResponse "Invalid UUID"
// @Failure      401 {object} respond.ErrorResponse "Unauthorized"
// @Failure      403 {object} respond.ErrorResponse "Permission denied"
// @Failure      500 {object} respond.ErrorResponse "Internal server error"
// @Router       /partner/restaurants/{id}/staff [get]
func (h *Handler) GetStaff(w http.ResponseWriter, r *http.Request) {
	const op = "restaurant.Handler.GetStaff"
	log := h.log.With(slog.String("op", op))

	restID, err := decode.UUID(r, params.RestaurantID)
	if err != nil {
		respond.Error(w, r, log, err)
		return
	}

	staff, err := h.staffService.GetStaff(r.Context(), restID)
	if err != nil {
		respond.Error(w, r, log, err)
		return
	}

	render.JSON(w, r, StaffResponse{Staff: staff})
}

// UpdateStaffRole godoc
// @Summary      Update staff role
// @Description  Updates the role of a specific staff member. Caller must be an owner.
// @Tags         restaurants, staff
// @Security     Bearer
// @Accept       json
// @Produce      json
// @Param        id path string true "Restaurant ID" format(uuid)
// @Param        userID path string true "User ID" format(uuid)
// @Param        input body UpdateStaffRoleRequest true "New role data"
// @Success      200 {object} map[string]string "Role successfully updated"
// @Failure      400 {object} respond.ErrorResponse "Validation error or invalid UUID/JSON"
// @Failure      401 {object} respond.ErrorResponse "Unauthorized"
// @Failure      403 {object} respond.ErrorResponse "Permission denied"
// @Failure      404 {object} respond.ErrorResponse "Staff member not found"
// @Failure      500 {object} respond.ErrorResponse "Internal server error"
// @Router       /partner/restaurants/{id}/staff/{userID} [patch]
func (h *Handler) UpdateStaffRole(w http.ResponseWriter, r *http.Request) {
	const op = "restaurant.Handler.UpdateStaffRole"
	log := h.log.With(slog.String("op", op))

	restID, err := decode.UUID(r, params.RestaurantID)
	if err != nil {
		respond.Error(w, r, log, err)
		return
	}

	userID, err := decode.UUID(r, params.UserID)
	if err != nil {
		respond.Error(w, r, log, err)
		return
	}

	actor, ok := ctxkey.GetRestaurantMember(r.Context())
	if !ok {
		respond.Error(w, r, log, domain.ErrInternalServerError)
		return
	}

	var req UpdateStaffRoleRequest
	if err := render.DecodeJSON(r.Body, &req); err != nil {
		respond.Error(w, r, log, domain.ErrInvalidJSON)
		return
	}

	if err := h.validate.Struct(req); err != nil {
		respond.Error(w, r, log, err)
		return
	}

	if err = h.staffService.UpdateStaffRole(r.Context(), actor, restID, userID, req.Role); err != nil {
		respond.Error(w, r, log, err)
		return
	}

	respond.Message(w, r, "role updated")
}

// TransferOwnership godoc
// @Summary      Transfer ownership of a restaurant
// @Description  Transfers ownership of a restaurant to a new user. Caller must be an owner.
// @Tags         restaurants, staff
// @Security     Bearer
// @Accept       json
// @Produce      json
// @Param        id path string true "Restaurant ID" format(uuid)
// @Param        userID path string true "User ID" format(uuid)
// @Success      200 {object} respond.MessageResponse "Ownership transferred successfully"
// @Failure      400 {object} respond.ErrorResponse "Validation error or invalid UUID/JSON"
// @Failure      401 {object} respond.ErrorResponse "Unauthorized"
// @Failure      403 {object} respond.ErrorResponse "Permission denied"
// @Failure      404 {object} respond.ErrorResponse "Staff member not found"
// @Failure      500 {object} respond.ErrorResponse "Internal server error"
// @Router       /partner/restaurants/{id}/staff/{userID}/ownership [patch]
func (h *Handler) TransferOwnership(w http.ResponseWriter, r *http.Request) {
	const op = "restaurant.Handler.TransferOwnership"
	log := h.log.With(slog.String("op", op))

	restID, err := decode.UUID(r, params.RestaurantID)
	if err != nil {
		respond.Error(w, r, log, err)
		return
	}

	targetID, err := decode.UUID(r, params.UserID)
	if err != nil {
		respond.Error(w, r, log, err)
		return
	}

	actor, ok := ctxkey.GetRestaurantMember(r.Context())
	if !ok {
		respond.Error(w, r, log, domain.ErrInternalServerError)
		return
	}

	if err = h.staffService.TransferOwnership(r.Context(), actor, restID, targetID); err != nil {
		respond.Error(w, r, log, err)
		return
	}

	respond.Message(w, r, "ownership transferred")
}

// RemoveStaff godoc
// @Summary      Remove staff member
// @Description  Removes a staff member from the restaurant. Caller must be an owner.
// @Tags         restaurants, staff
// @Security     Bearer
// @Produce      json
// @Param        id path string true "Restaurant ID" format(uuid)
// @Param        userID path string true "User ID" format(uuid)
// @Success      204 "No Content - Staff successfully removed"
// @Failure      400 {object} respond.ErrorResponse "Invalid UUID"
// @Failure      401 {object} respond.ErrorResponse "Unauthorized"
// @Failure      403 {object} respond.ErrorResponse "Permission denied"
// @Failure      404 {object} respond.ErrorResponse "Staff member not found"
// @Failure      500 {object} respond.ErrorResponse "Internal server error"
// @Router       /partner/restaurants/{id}/staff/{userID} [delete]
func (h *Handler) RemoveStaff(w http.ResponseWriter, r *http.Request) {
	const op = "restaurant.Handler.RemoveStaff"
	log := h.log.With(slog.String("op", op))

	restID, err := decode.UUID(r, params.RestaurantID)
	if err != nil {
		respond.Error(w, r, log, err)
		return
	}

	userID, err := decode.UUID(r, params.UserID)
	if err != nil {
		respond.Error(w, r, log, err)
		return
	}

	actor, ok := ctxkey.GetRestaurantMember(r.Context())
	if !ok {
		respond.Error(w, r, log, domain.ErrInternalServerError)
		return
	}

	if err = h.staffService.RemoveStaff(r.Context(), actor, restID, userID); err != nil {
		respond.Error(w, r, log, err)
		return
	}

	render.NoContent(w, r)
}
