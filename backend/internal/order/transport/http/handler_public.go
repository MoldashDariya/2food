package http

import (
	"io"
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
// @Summary      Create a new order
// @Description  Creates a new order from a cart and returns a Stripe client secret for payment execution
// @Tags         orders
// @Security     Bearer
// @Accept       json
// @Produce      json
// @Param        input body CreateOrderRequest true "Order details"
// @Success      201 {object} CreateOrderResponse "Created order and payment intent client secret"
// @Failure      400 {object} respond.ErrorResponse "Validation error, invalid JSON, or stock issues"
// @Failure      401 {object} respond.ErrorResponse "Unauthorized"
// @Failure      500 {object} respond.ErrorResponse "Internal server error"
// @Router       /orders [post]
func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	const op = "order.Handler.Create"
	log := h.log.With(slog.String("op", op))

	userID, ok := ctxkey.GetUserID(r.Context())
	if !ok {
		respond.Error(w, r, log, domain.ErrInternalServerError)
		return
	}

	var req CreateOrderRequest
	if err := render.DecodeJSON(r.Body, &req); err != nil {
		respond.Error(w, r, log, domain.ErrInvalidJSON)
		return
	}

	if err := h.validate.Struct(req); err != nil {
		respond.Error(w, r, log, err)
		return
	}

	items := make([]domain.CreateOrderItemInput, len(req.Items))
	for i, item := range req.Items {
		items[i] = domain.CreateOrderItemInput{
			ProductID: item.ProductID,
			Quantity:  item.Quantity,
		}
	}

	input := domain.CreateOrderInput{
		UserID:            userID,
		RestaurantID:      req.RestaurantID,
		PickupWindowStart: req.PickupWindowStart,
		PickupWindowEnd:   req.PickupWindowEnd,
		Items:             items,
	}

	order, clientSecret, err := h.service.Create(r.Context(), input)
	if err != nil {
		respond.Error(w, r, log, err)
		return
	}

	render.Status(r, http.StatusCreated)
	render.JSON(w, r, CreateOrderResponse{Order: order, ClientSecret: clientSecret})
}

// GetByID godoc
// @Summary      Get order by ID
// @Description  Retrieves detailed information about a specific order belonging to the current user
// @Tags         orders
// @Security     Bearer
// @Produce      json
// @Param        orderID path string true "Order ID" format(uuid)
// @Success      200 {object} OrderResponse "Order details"
// @Failure      400 {object} respond.ErrorResponse "Invalid UUID"
// @Failure      401 {object} respond.ErrorResponse "Unauthorized"
// @Failure      403 {object} respond.ErrorResponse "Permission denied"
// @Failure      404 {object} respond.ErrorResponse "Order not found"
// @Failure      500 {object} respond.ErrorResponse "Internal server error"
// @Router       /orders/{orderID} [get]
func (h *Handler) GetByID(w http.ResponseWriter, r *http.Request) {
	const op = "order.Handler.GetByID"
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

	order, err := h.service.GetByID(r.Context(), orderID)
	if err != nil {
		respond.Error(w, r, log, err)
		return
	}

	if order.UserID != userID {
		respond.Error(w, r, log, domain.ErrPermissionDenied)
		return
	}

	render.JSON(w, r, OrderResponse{Order: order})
}

// GetByUser godoc
// @Summary      Get user orders
// @Description  Retrieves a list of all orders placed by the current authenticated user
// @Tags         orders
// @Security     Bearer
// @Produce      json
// @Success      200 {object} OrdersResponse "List of user's orders"
// @Failure      401 {object} respond.ErrorResponse "Unauthorized"
// @Failure      500 {object} respond.ErrorResponse "Internal server error"
// @Router       /orders [get]
func (h *Handler) GetByUser(w http.ResponseWriter, r *http.Request) {
	const op = "order.Handler.GetByUser"
	log := h.log.With(slog.String("op", op))

	userID, ok := ctxkey.GetUserID(r.Context())
	if !ok {
		respond.Error(w, r, log, domain.ErrInternalServerError)
		return
	}

	orders, err := h.service.GetByUser(r.Context(), userID)
	if err != nil {
		respond.Error(w, r, log, err)
		return
	}

	render.JSON(w, r, OrdersResponse{Orders: orders})
}

// GetStatusHistoryByUser godoc
// @Summary      Get order status history
// @Description  Retrieves the complete history of status changes for a specific order
// @Tags         orders
// @Security     Bearer
// @Produce      json
// @Param        orderID path string true "Order ID" format(uuid)
// @Success      200 {object} StatusHistoryResponse "List of status changes"
// @Failure      400 {object} respond.ErrorResponse "Invalid UUID"
// @Failure      401 {object} respond.ErrorResponse "Unauthorized"
// @Failure      403 {object} respond.ErrorResponse "Permission denied"
// @Failure      500 {object} respond.ErrorResponse "Internal server error"
// @Router       /orders/{orderID}/history [get]
func (h *Handler) GetStatusHistoryByUser(w http.ResponseWriter, r *http.Request) {
	const op = "order.Handler.GetStatusHistoryByUser"
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

	history, err := h.service.GetStatusHistoryByUser(r.Context(), userID, orderID)
	if err != nil {
		respond.Error(w, r, log, err)
		return
	}

	render.JSON(w, r, StatusHistoryResponse{History: history})
}

func (h *Handler) GetPaymentHistoryByOrder(w http.ResponseWriter, r *http.Request) {
	const op = "order.Handler.GetPaymentHistoryByOrder"
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

	history, err := h.service.GetPaymentHistoryByOrder(r.Context(), userID, orderID)
	if err != nil {
		respond.Error(w, r, log, err)
		return
	}

	render.JSON(w, r, PaymentHistoryResponse{History: history})
}

func (h *Handler) GetPaymentHistoryByUser(w http.ResponseWriter, r *http.Request) {
	const op = "order.Handler.GetPaymentHistoryByUser"
	log := h.log.With(slog.String("op", op))

	userID, ok := ctxkey.GetUserID(r.Context())
	if !ok {
		respond.Error(w, r, log, domain.ErrInternalServerError)
		return
	}

	history, err := h.service.GetPaymentHistoryByUser(r.Context(), userID)
	if err != nil {
		respond.Error(w, r, log, err)
		return
	}

	render.JSON(w, r, PaymentHistoryResponse{History: history})
}

// CancelByUser godoc
// @Summary      Cancel order
// @Description  Allows a user to cancel their order (if permitted by current order status)
// @Tags         orders
// @Security     Bearer
// @Accept       json
// @Produce      json
// @Param        orderID path string true "Order ID" format(uuid)
// @Param        input body CancelByUserRequest false "Cancellation reason (optional)"
// @Success      204 "No Content - Order successfully cancelled"
// @Failure      400 {object} respond.ErrorResponse "Invalid JSON/UUID or invalid status transition"
// @Failure      401 {object} respond.ErrorResponse "Unauthorized"
// @Failure      403 {object} respond.ErrorResponse "Permission denied"
// @Failure      409 {object} respond.ErrorResponse "Conflict - Order cannot be cancelled anymore"
// @Failure      500 {object} respond.ErrorResponse "Internal server error"
// @Router       /orders/{orderID} [delete]
func (h *Handler) CancelByUser(w http.ResponseWriter, r *http.Request) {
	const op = "order.Handler.CancelByUser"
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

	var req CancelByUserRequest
	err = render.DecodeJSON(r.Body, &req)
	if err != nil && err != io.EOF {
		respond.Error(w, r, log, domain.ErrInvalidJSON)
		return
	}

	if err := h.validate.Struct(req); err != nil {
		respond.Error(w, r, log, err)
		return
	}

	if err := h.service.CancelByUser(r.Context(), userID, orderID, req.Reason); err != nil {
		respond.Error(w, r, log, err)
		return
	}

	render.NoContent(w, r)
}
