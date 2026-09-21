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

// GetByRestaurant godoc
// @Summary      Get all restaurant orders
// @Description  Retrieves all orders for a specific restaurant. Caller must be a restaurant staff member.
// @Tags         restaurant-orders
// @Security     Bearer
// @Produce      json
// @Param        restaurantID path string true "Restaurant ID" format(uuid)
// @Success      200 {object} OrdersResponse "List of all restaurant orders"
// @Failure      400 {object} respond.ErrorResponse "Invalid UUID"
// @Failure      401 {object} respond.ErrorResponse "Unauthorized"
// @Failure      403 {object} respond.ErrorResponse "Permission denied"
// @Failure      500 {object} respond.ErrorResponse "Internal server error"
// @Router       /partner/restaurants/{restaurantID}/orders [get]
func (h *Handler) GetByRestaurant(w http.ResponseWriter, r *http.Request) {
	const op = "order.Handler.GetByRestaurant"
	log := h.log.With(slog.String("op", op))

	restID, err := decode.UUID(r, params.RestaurantID)
	if err != nil {
		respond.Error(w, r, log, err)
		return
	}

	orders, err := h.service.GetByRestaurant(r.Context(), restID)
	if err != nil {
		respond.Error(w, r, log, err)
		return
	}

	render.JSON(w, r, OrdersResponse{Orders: orders})
}

// GetActiveByRestaurant godoc
// @Summary      Get active restaurant orders
// @Description  Retrieves only active (not completed/cancelled) orders for a specific restaurant
// @Tags         restaurant-orders
// @Security     Bearer
// @Produce      json
// @Param        restaurantID path string true "Restaurant ID" format(uuid)
// @Success      200 {object} OrdersResponse "List of active orders"
// @Failure      400 {object} respond.ErrorResponse "Invalid UUID"
// @Failure      401 {object} respond.ErrorResponse "Unauthorized"
// @Failure      403 {object} respond.ErrorResponse "Permission denied"
// @Failure      500 {object} respond.ErrorResponse "Internal server error"
// @Router       /partner/restaurants/{restaurantID}/orders/active [get]
func (h *Handler) GetActiveByRestaurant(w http.ResponseWriter, r *http.Request) {
	const op = "order.Handler.GetActiveByRestaurant"
	log := h.log.With(slog.String("op", op))

	restID, err := decode.UUID(r, params.RestaurantID)
	if err != nil {
		respond.Error(w, r, log, err)
		return
	}

	orders, err := h.service.GetActiveByRestaurant(r.Context(), restID)
	if err != nil {
		respond.Error(w, r, log, err)
		return
	}

	render.JSON(w, r, OrdersResponse{Orders: orders})
}

// GetByPickupCode godoc
// @Summary      Get order by pickup code
// @Description  Looks up an order using the unique pickup code provided by the customer
// @Tags         restaurant-orders
// @Security     Bearer
// @Produce      json
// @Param        restaurantID path string true "Restaurant ID" format(uuid)
// @Param        pickup_code query string true "Pickup code from customer"
// @Success      200 {object} OrderResponse "Order details"
// @Failure      400 {object} respond.ErrorResponse "Missing code or invalid UUID"
// @Failure      401 {object} respond.ErrorResponse "Unauthorized"
// @Failure      403 {object} respond.ErrorResponse "Permission denied"
// @Failure      404 {object} respond.ErrorResponse "Order not found or code invalid"
// @Failure      500 {object} respond.ErrorResponse "Internal server error"
// @Router       /partner/restaurants/{restaurantID}/orders/pickup [get]
func (h *Handler) GetByPickupCode(w http.ResponseWriter, r *http.Request) {
	const op = "order.Handler.GetByPickupCode"
	log := h.log.With(slog.String("op", op))

	restID, err := decode.UUID(r, params.RestaurantID)
	if err != nil {
		respond.Error(w, r, log, err)
		return
	}

	code := r.URL.Query().Get(params.PickupCode)
	if code == "" {
		respond.Error(w, r, log, domain.ErrMissingParam)
		return
	}

	order, err := h.service.GetByPickupCode(r.Context(), restID, code)
	if err != nil {
		respond.Error(w, r, log, err)
		return
	}

	render.JSON(w, r, OrderResponse{Order: order})
}

// MarkAsAccepted godoc
// @Summary      Mark order as accepted (in work)
// @Description  Updates order status to 'Accepted/Preparing' and notifies the customer.
// @Tags         restaurant-orders
// @Security     Bearer
// @Produce      json
// @Param        restaurantID path string true "Restaurant ID" format(uuid)
// @Param        orderID path string true "Order ID" format(uuid)
// @Success      204 "No Content - Order marked as accepted"
// @Failure      400 {object} respond.ErrorResponse "Invalid UUID"
// @Failure      401 {object} respond.ErrorResponse "Unauthorized"
// @Failure      403 {object} respond.ErrorResponse "Permission denied"
// @Failure      409 {object} respond.ErrorResponse "Invalid status transition"
// @Failure      500 {object} respond.ErrorResponse "Internal server error"
// @Router       /partner/restaurants/{restaurantID}/orders/{orderID}/accept [patch]
func (h *Handler) MarkAsAccepted(w http.ResponseWriter, r *http.Request) {
	const op = "order.Handler.MarkAsAccepted"
	log := h.log.With(slog.String("op", op))

	actor, ok := ctxkey.GetRestaurantMember(r.Context())
	if !ok {
		respond.Error(w, r, log, domain.ErrPermissionDenied)
		return
	}

	orderID, err := decode.UUID(r, params.OrderID)
	if err != nil {
		respond.Error(w, r, log, err)
		return
	}

	if err := h.service.MarkAsAccepted(r.Context(), actor, orderID); err != nil {
		respond.Error(w, r, log, err)
		return
	}

	render.NoContent(w, r)
}

// MarkAsReady godoc
// @Summary      Mark order as ready
// @Description  Updates order status to 'Ready for Pickup'. Caller must be a restaurant staff member.
// @Tags         restaurant-orders
// @Security     Bearer
// @Produce      json
// @Param        orderID path string true "Order ID" format(uuid)
// @Success      204 "No Content - Order marked as ready"
// @Failure      400 {object} respond.ErrorResponse "Invalid UUID"
// @Failure      401 {object} respond.ErrorResponse "Unauthorized"
// @Failure      403 {object} respond.ErrorResponse "Permission denied"
// @Failure      409 {object} respond.ErrorResponse "Invalid status transition"
// @Failure      500 {object} respond.ErrorResponse "Internal server error"
// @Router       /partner/restaurants/{restaurantID}/orders/{orderID}/ready [patch]
func (h *Handler) MarkAsReady(w http.ResponseWriter, r *http.Request) {
	const op = "order.Handler.MarkAsReady"
	log := h.log.With(slog.String("op", op))

	actor, ok := ctxkey.GetRestaurantMember(r.Context())
	if !ok {
		respond.Error(w, r, log, domain.ErrPermissionDenied)
		return
	}

	orderID, err := decode.UUID(r, params.OrderID)
	if err != nil {
		respond.Error(w, r, log, err)
		return
	}

	if err := h.service.MarkAsReady(r.Context(), actor, orderID); err != nil {
		respond.Error(w, r, log, err)
		return
	}

	render.NoContent(w, r)
}

// MarkAsCollected godoc
// @Summary      Mark order as collected
// @Description  Updates order status to 'Completed/Collected'. Requires valid PIN code unless force=true.
// @Tags         restaurant-orders
// @Security     Bearer
// @Accept       json
// @Produce      json
// @Param        restaurantID path string true "Restaurant ID" format(uuid)
// @Param        orderID path string true "Order ID" format(uuid)
// @Param        input body MarkAsCollectedRequest true "Pickup code and force flag"
// @Success      204 "No Content - Order marked as collected"
// @Failure      400 {object} respond.ErrorResponse "Invalid JSON/UUID"
// @Failure      401 {object} respond.ErrorResponse "Unauthorized"
// @Failure      403 {object} respond.ErrorResponse "Permission denied or invalid pickup code"
// @Failure      409 {object} respond.ErrorResponse "Invalid status transition"
// @Failure      500 {object} respond.ErrorResponse "Internal server error"
// @Router       /partner/restaurants/{restaurantID}/orders/{orderID}/collect [patch]
func (h *Handler) MarkAsCollected(w http.ResponseWriter, r *http.Request) {
	const op = "order.Handler.MarkAsCollected"
	log := h.log.With(slog.String("op", op))

	actor, ok := ctxkey.GetRestaurantMember(r.Context())
	if !ok {
		respond.Error(w, r, log, domain.ErrPermissionDenied)
		return
	}

	orderID, err := decode.UUID(r, params.OrderID)
	if err != nil {
		respond.Error(w, r, log, err)
		return
	}

	var req MarkAsCollectedRequest
	if err := render.DecodeJSON(r.Body, &req); err != nil {
		respond.Error(w, r, log, domain.ErrInvalidJSON)
		return
	}

	if err := h.service.MarkAsCollected(r.Context(), actor, orderID, req.PickupCode, req.Force); err != nil {
		respond.Error(w, r, log, err)
		return
	}

	render.NoContent(w, r)
}

// GetStatusHistoryByRestaurant godoc
// @Summary      Get order history (Restaurant)
// @Description  Retrieves the complete status history for a specific order. Caller must be a restaurant staff member.
// @Tags         restaurant-orders
// @Security     Bearer
// @Produce      json
// @Param        orderID path string true "Order ID" format(uuid)
// @Success      200 {object} StatusHistoryResponse "List of status changes"
// @Failure      400 {object} respond.ErrorResponse "Invalid UUID"
// @Failure      401 {object} respond.ErrorResponse "Unauthorized"
// @Failure      403 {object} respond.ErrorResponse "Permission denied"
// @Failure      500 {object} respond.ErrorResponse "Internal server error"
// @Router       /restaurants/{restaurantID}/orders/{orderID}/history [get]
func (h *Handler) GetStatusHistoryByRestaurant(w http.ResponseWriter, r *http.Request) {
	const op = "order.Handler.GetStatusHistoryByRestaurant"
	log := h.log.With(slog.String("op", op))

	actor, ok := ctxkey.GetRestaurantMember(r.Context())
	if !ok {
		respond.Error(w, r, log, domain.ErrPermissionDenied)
		return
	}

	orderID, err := decode.UUID(r, params.OrderID)
	if err != nil {
		respond.Error(w, r, log, err)
		return
	}

	history, err := h.service.GetStatusHistoryByRestaurant(r.Context(), actor, orderID)
	if err != nil {
		respond.Error(w, r, log, err)
		return
	}

	render.JSON(w, r, StatusHistoryResponse{History: history})
}

// CancelByRestaurant godoc
// @Summary      Cancel order (Restaurant)
// @Description  Allows a restaurant staff member to cancel an order and provide a reason
// @Tags         restaurant-orders
// @Security     Bearer
// @Accept       json
// @Produce      json
// @Param        orderID path string true "Order ID" format(uuid)
// @Param        input body CancelByRestaurantRequest true "Cancellation reason"
// @Success      204 "No Content - Order successfully cancelled"
// @Failure      400 {object} respond.ErrorResponse "Invalid JSON/UUID or validation error"
// @Failure      401 {object} respond.ErrorResponse "Unauthorized"
// @Failure      403 {object} respond.ErrorResponse "Permission denied"
// @Failure      409 {object} respond.ErrorResponse "Invalid status transition"
// @Failure      500 {object} respond.ErrorResponse "Internal server error"
// @Router       /partner/restaurants/{restaurantID}/orders/{orderID} [delete]
func (h *Handler) CancelByRestaurant(w http.ResponseWriter, r *http.Request) {
	const op = "order.Handler.CancelByRestaurant"
	log := h.log.With(slog.String("op", op))

	actor, ok := ctxkey.GetRestaurantMember(r.Context())
	if !ok {
		respond.Error(w, r, log, domain.ErrPermissionDenied)
		return
	}

	orderID, err := decode.UUID(r, params.OrderID)
	if err != nil {
		respond.Error(w, r, log, err)
		return
	}

	var req CancelByRestaurantRequest
	if err := render.DecodeJSON(r.Body, &req); err != nil {
		respond.Error(w, r, log, domain.ErrInvalidJSON)
		return
	}

	if err := h.validate.Struct(req); err != nil {
		respond.Error(w, r, log, err)
		return
	}
	if err := h.service.CancelByRestaurant(r.Context(), actor, orderID, &req.Reason); err != nil {
		respond.Error(w, r, log, err)
		return
	}

	render.NoContent(w, r)
}
