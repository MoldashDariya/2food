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
// @Summary      Create product
// @Description  Creates a new product for a specific restaurant. Caller must be a restaurant staff member.
// @Tags         products
// @Security     Bearer
// @Accept       json
// @Produce      json
// @Param        restaurantID path string true "Restaurant ID" format(uuid)
// @Param        input body CreateProductRequest true "Product data"
// @Success      201 {object} ProductResponse "Created product details"
// @Failure      400 {object} respond.ErrorResponse "Validation error or invalid JSON/UUID"
// @Failure      401 {object} respond.ErrorResponse "Unauthorized"
// @Failure      403 {object} respond.ErrorResponse "Permission denied"
// @Failure      500 {object} respond.ErrorResponse "Internal server error"
// @Router       /partner/restaurants/{restaurantID}/products [post]
func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	const op = "product.Handler.Create"
	log := h.log.With(slog.String("op", op))

	restID, err := decode.UUID(r, params.RestaurantID)
	if err != nil {
		respond.Error(w, r, log, err)
		return
	}

	var req CreateProductRequest
	if err := render.DecodeJSON(r.Body, &req); err != nil {
		respond.Error(w, r, log, domain.ErrInvalidJSON)
		return
	}

	if err := h.validate.Struct(req); err != nil {
		respond.Error(w, r, log, err)
		return
	}

	input := domain.CreateProductInput{
		RestaurantID:      restID,
		CategoryID:        req.CategoryID,
		Name:              req.Name,
		Description:       req.Description,
		ImageURL:          req.ImageURL,
		OriginalPrice:     req.OriginalPrice,
		DiscountedPrice:   req.DiscountedPrice,
		PickupTimeStart:   req.PickupTimeStart,
		PickupTimeEnd:     req.PickupTimeEnd,
		QuantityAvailable: req.QuantityAvailable,
	}

	product, err := h.productService.Create(r.Context(), input)
	if err != nil {
		respond.Error(w, r, log, err)
		return
	}

	render.Status(r, http.StatusCreated)
	render.JSON(w, r, ProductResponse{Product: product})
}

// GetByRestaurant godoc
// @Summary      Get all restaurant products
// @Description  Retrieves a list of all products (including inactive) for a specific restaurant. Caller must be a restaurant staff member.
// @Tags         products
// @Security     Bearer
// @Produce      json
// @Param        restaurantID path string true "Restaurant ID" format(uuid)
// @Success      200 {object} ProductsResponse "List of all products"
// @Failure      400 {object} respond.ErrorResponse "Invalid UUID"
// @Failure      401 {object} respond.ErrorResponse "Unauthorized"
// @Failure      403 {object} respond.ErrorResponse "Permission denied"
// @Failure      500 {object} respond.ErrorResponse "Internal server error"
// @Router       /partner/restaurants/{restaurantID}/products [get]
func (h *Handler) GetByRestaurant(w http.ResponseWriter, r *http.Request) {
	const op = "product.Handler.GetByRestaurant"
	log := h.log.With(slog.String("op", op))

	actor, ok := ctxkey.GetRestaurantMember(r.Context())
	if !ok {
		respond.Error(w, r, log, domain.ErrPermissionDenied)
		return
	}

	restID, err := decode.UUID(r, params.RestaurantID)
	if err != nil {
		respond.Error(w, r, log, err)
		return
	}

	products, err := h.productService.GetByRestaurant(r.Context(), actor, restID)
	if err != nil {
		respond.Error(w, r, log, err)
		return
	}

	render.JSON(w, r, ProductsResponse{Products: products})
}

// GetByID godoc
// @Summary      Get product by ID
// @Description  Retrieves detailed information about a specific product. Caller must be a restaurant staff member.
// @Tags         products
// @Security     Bearer
// @Produce      json
// @Param        productID path string true "Product ID" format(uuid)
// @Success      200 {object} ProductResponse "Product details"
// @Failure      400 {object} respond.ErrorResponse "Invalid UUID"
// @Failure      401 {object} respond.ErrorResponse "Unauthorized"
// @Failure      403 {object} respond.ErrorResponse "Permission denied"
// @Failure      404 {object} respond.ErrorResponse "Product not found"
// @Failure      500 {object} respond.ErrorResponse "Internal server error"
// @Router       /partner/restaurants/{restaurantID}/products/{productID} [get]
func (h *Handler) GetByID(w http.ResponseWriter, r *http.Request) {
	const op = "product.Handler.GetByID"
	log := h.log.With(slog.String("op", op))

	actor, ok := ctxkey.GetRestaurantMember(r.Context())
	if !ok {
		respond.Error(w, r, log, domain.ErrPermissionDenied)
		return
	}

	productID, err := decode.UUID(r, params.ProductID)
	if err != nil {
		respond.Error(w, r, log, err)
		return
	}

	product, err := h.productService.GetByID(r.Context(), actor, productID)
	if err != nil {
		respond.Error(w, r, log, err)
		return
	}

	render.JSON(w, r, ProductResponse{Product: product})
}

// Update godoc
// @Summary      Update product
// @Description  Updates general information of a specific product. Caller must be a restaurant staff member.
// @Tags         products
// @Security     Bearer
// @Accept       json
// @Produce      json
// @Param        productID path string true "Product ID" format(uuid)
// @Param        input body UpdateProductRequest true "Updated product data"
// @Success      200 {object} ProductResponse "Updated product details"
// @Failure      400 {object} respond.ErrorResponse "Validation error or invalid JSON/UUID"
// @Failure      401 {object} respond.ErrorResponse "Unauthorized"
// @Failure      403 {object} respond.ErrorResponse "Permission denied"
// @Failure      404 {object} respond.ErrorResponse "Product not found"
// @Failure      500 {object} respond.ErrorResponse "Internal server error"
// @Router       /partner/restaurants/{restaurantID}/products/{productID} [patch]
func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	const op = "product.Handler.Update"
	log := h.log.With(slog.String("op", op))

	actor, ok := ctxkey.GetRestaurantMember(r.Context())
	if !ok {
		respond.Error(w, r, log, domain.ErrPermissionDenied)
		return
	}

	productID, err := decode.UUID(r, params.ProductID)
	if err != nil {
		respond.Error(w, r, log, err)
		return
	}

	var req UpdateProductRequest
	if err := render.DecodeJSON(r.Body, &req); err != nil {
		respond.Error(w, r, log, domain.ErrInvalidJSON)
		return
	}

	if err := h.validate.Struct(req); err != nil {
		respond.Error(w, r, log, err)
		return
	}

	input := domain.UpdateProductInput{
		CategoryID:      req.CategoryID,
		Name:            req.Name,
		Description:     req.Description,
		ImageURL:        req.ImageURL,
		OriginalPrice:   req.OriginalPrice,
		DiscountedPrice: req.DiscountedPrice,
		PickupTimeStart: req.PickupTimeStart,
		PickupTimeEnd:   req.PickupTimeEnd,
	}

	product, err := h.productService.Update(r.Context(), actor, productID, input)
	if err != nil {
		respond.Error(w, r, log, err)
		return
	}

	render.JSON(w, r, ProductResponse{Product: product})
}

// UpdateStatus godoc
// @Summary      Update product active status
// @Description  Activates or deactivates a product. Caller must be a restaurant staff member.
// @Tags         products
// @Security     Bearer
// @Accept       json
// @Produce      json
// @Param        productID path string true "Product ID" format(uuid)
// @Param        input body UpdateStatusRequest true "New status data"
// @Success      200 {object} respond.MessageResponse "Status successfully updated"
// @Failure      400 {object} respond.ErrorResponse "Validation error or invalid JSON/UUID"
// @Failure      401 {object} respond.ErrorResponse "Unauthorized"
// @Failure      403 {object} respond.ErrorResponse "Permission denied"
// @Failure      404 {object} respond.ErrorResponse "Product not found"
// @Failure      500 {object} respond.ErrorResponse "Internal server error"
// @Router       /partner/restaurants/{restaurantID}/products/{productID}/status [patch]
func (h *Handler) UpdateStatus(w http.ResponseWriter, r *http.Request) {
	const op = "product.Handler.UpdateStatus"
	log := h.log.With(slog.String("op", op))

	actor, ok := ctxkey.GetRestaurantMember(r.Context())
	if !ok {
		respond.Error(w, r, log, domain.ErrPermissionDenied)
		return
	}

	productID, err := decode.UUID(r, params.ProductID)
	if err != nil {
		respond.Error(w, r, log, err)
		return
	}

	var req UpdateStatusRequest
	if err := render.DecodeJSON(r.Body, &req); err != nil {
		respond.Error(w, r, log, domain.ErrInvalidJSON)
		return
	}

	if err := h.validate.Struct(req); err != nil {
		respond.Error(w, r, log, err)
		return
	}

	err = h.productService.UpdateStatus(r.Context(), actor, productID, *req.IsActive)
	if err != nil {
		respond.Error(w, r, log, err)
		return
	}

	respond.Message(w, r, "status updated")
}

// UpdateStock godoc
// @Summary      Update product stock
// @Description  Updates the available quantity of a product. Caller must be a restaurant staff member.
// @Tags         products
// @Security     Bearer
// @Accept       json
// @Produce      json
// @Param        productID path string true "Product ID" format(uuid)
// @Param        input body UpdateStockRequest true "New stock quantity data"
// @Success      200 {object} respond.MessageResponse "Stock successfully updated"
// @Failure      400 {object} respond.ErrorResponse "Validation error or invalid JSON/UUID"
// @Failure      401 {object} respond.ErrorResponse "Unauthorized"
// @Failure      403 {object} respond.ErrorResponse "Permission denied"
// @Failure      404 {object} respond.ErrorResponse "Product not found"
// @Failure      500 {object} respond.ErrorResponse "Internal server error"
// @Router       /partner/restaurants/{restaurantID}/products/{productID}/stock [patch]
func (h *Handler) UpdateStock(w http.ResponseWriter, r *http.Request) {
	const op = "product.Handler.UpdateStock"
	log := h.log.With(slog.String("op", op))

	actor, ok := ctxkey.GetRestaurantMember(r.Context())
	if !ok {
		respond.Error(w, r, log, domain.ErrPermissionDenied)
		return
	}

	productID, err := decode.UUID(r, params.ProductID)
	if err != nil {
		respond.Error(w, r, log, err)
		return
	}

	var req UpdateStockRequest
	if err := render.DecodeJSON(r.Body, &req); err != nil {
		respond.Error(w, r, log, domain.ErrInvalidJSON)
		return
	}

	if err := h.validate.Struct(req); err != nil {
		respond.Error(w, r, log, err)
		return
	}

	err = h.productService.UpdateStock(r.Context(), actor, productID, *req.Quantity)
	if err != nil {
		respond.Error(w, r, log, err)
		return
	}

	respond.Message(w, r, "stock updated")
}

// Delete godoc
// @Summary      Delete product
// @Description  Deletes a product from the system. Caller must be a restaurant staff member.
// @Tags         products
// @Security     Bearer
// @Produce      json
// @Param        productID path string true "Product ID" format(uuid)
// @Success      200 {object} respond.MessageResponse "Product successfully deleted"
// @Failure      400 {object} respond.ErrorResponse "Invalid UUID"
// @Failure      401 {object} respond.ErrorResponse "Unauthorized"
// @Failure      403 {object} respond.ErrorResponse "Permission denied"
// @Failure      404 {object} respond.ErrorResponse "Product not found"
// @Failure      500 {object} respond.ErrorResponse "Internal server error"
// @Router       /partner/restaurants/{restaurantID}/products/{productID} [delete]
func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	const op = "product.Handler.Delete"
	log := h.log.With(slog.String("op", op))

	actor, ok := ctxkey.GetRestaurantMember(r.Context())
	if !ok {
		respond.Error(w, r, log, domain.ErrPermissionDenied)
		return
	}

	productID, err := decode.UUID(r, params.ProductID)
	if err != nil {
		respond.Error(w, r, log, err)
		return
	}

	err = h.productService.Delete(r.Context(), actor, productID)
	if err != nil {
		respond.Error(w, r, log, err)
		return
	}

	respond.Message(w, r, "product deleted")
}
