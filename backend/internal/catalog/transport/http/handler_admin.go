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

// CreateCity godoc
// @Summary      Create a new city
// @Description  Creates a new city and assigns the specified user as the owner. Caller must be an admin.
// @Tags         catalog
// @Security     Bearer
// @Accept       json
// @Produce      json
// @Param        input body CreateCityRequest true "City data and owner ID"
// @Success      201 {object} CityResponse "Successfully created city"
// @Failure      400 {object} respond.ErrorResponse "Validation error or invalid JSON"
// @Failure      401 {object} respond.ErrorResponse "Unauthorized"
// @Failure      403 {object} respond.ErrorResponse "Permission denied"
// @Failure      500 {object} respond.ErrorResponse "Internal server error"
// @Router       /admin/catalog/cities [post]
func (h *Handler) CreateCity(w http.ResponseWriter, r *http.Request) {
	const op = "catalog.Handler.CreateCity"
	log := h.log.With(slog.String("op", op))

	var req CreateCityRequest
	if err := render.DecodeJSON(r.Body, &req); err != nil {
		respond.Error(w, r, log, domain.ErrInvalidJSON)
		return
	}

	if err := h.validate.Struct(req); err != nil {
		respond.Error(w, r, log, err)
		return
	}

	input := domain.CreateCityInput{
		Name:  req.Name,
		Label: req.Label,
	}

	city, err := h.service.CreateCity(r.Context(), input)
	if err != nil {
		respond.Error(w, r, log, err)
		return
	}

	render.Status(r, http.StatusCreated)
	render.JSON(w, r, CityResponse{City: city})
}

// DeleteCity godoc
// @Summary      Delete a city
// @Description  Deletes a city from the system. Caller must be an admin.
// @Tags         catalog
// @Security     Bearer
// @Produce      json
// @Param        city_id path int true "City ID"
// @Success      204 "No Content - City successfully deleted"
// @Failure      400 {object} respond.ErrorResponse "Invalid ID"
// @Failure      401 {object} respond.ErrorResponse "Unauthorized"
// @Failure      403 {object} respond.ErrorResponse "Permission denied"
// @Failure      404 {object} respond.ErrorResponse "City not found"
// @Failure      500 {object} respond.ErrorResponse "Internal server error"
// @Router       /admin/catalog/cities/{city_id} [delete]
func (h *Handler) DeleteCity(w http.ResponseWriter, r *http.Request) {
	const op = "catalog.Handler.DeleteCity"
	log := h.log.With(slog.String("op", op))

	id, err := decode.Int(r, params.CityID)
	if err != nil {
		respond.Error(w, r, log, domain.ErrInvalidJSON)
		return
	}

	if err := h.service.DeleteCity(r.Context(), id); err != nil {
		respond.Error(w, r, log, err)
		return
	}

	render.NoContent(w, r)
}

// CreateCuisineType godoc
// @Summary      Create a new cuisine type
// @Description  Creates a new cuisine type and assigns the specified user as the owner. Caller must be an admin.
// @Tags         catalog
// @Security     Bearer
// @Accept       json
// @Produce      json
// @Param        input body CreateCuisineTypeRequest true "Cuisine type data and owner ID"
// @Success      201 {object} CuisineTypeResponse "Successfully created cuisine type"
// @Failure      400 {object} respond.ErrorResponse "Validation error or invalid JSON"
// @Failure      401 {object} respond.ErrorResponse "Unauthorized"
// @Failure      403 {object} respond.ErrorResponse "Permission denied"
// @Failure      500 {object} respond.ErrorResponse "Internal server error"
// @Router       /admin/catalog/cuisine-types [post]
func (h *Handler) CreateCuisineType(w http.ResponseWriter, r *http.Request) {
	const op = "catalog.Handler.CreateCuisineType"
	log := h.log.With(slog.String("op", op))

	var req CreateCuisineTypeRequest
	if err := render.DecodeJSON(r.Body, &req); err != nil {
		respond.Error(w, r, log, domain.ErrInvalidJSON)
		return
	}

	if err := h.validate.Struct(req); err != nil {
		respond.Error(w, r, log, err)
		return
	}

	input := domain.CreateCuisineTypeInput{
		Name:  req.Name,
		Label: req.Label,
	}

	cuisineType, err := h.service.CreateCuisineType(r.Context(), input)
	if err != nil {
		respond.Error(w, r, log, err)
		return
	}

	render.Status(r, http.StatusCreated)
	render.JSON(w, r, CuisineTypeResponse{CuisineType: cuisineType})
}

// DeleteCuisineType godoc
// @Summary      Delete a cuisine type
// @Description  Deletes a cuisine type from the system. Caller must be an admin.
// @Tags         catalog
// @Security     Bearer
// @Produce      json
// @Param        cuisine_type_id path int true "Cuisine type ID"
// @Success      204 "No Content - Cuisine type successfully deleted"
// @Failure      400 {object} respond.ErrorResponse "Invalid ID"
// @Failure      401 {object} respond.ErrorResponse "Unauthorized"
// @Failure      403 {object} respond.ErrorResponse "Permission denied"
// @Failure      404 {object} respond.ErrorResponse "Cuisine type not found"
// @Failure      500 {object} respond.ErrorResponse "Internal server error"
// @Router       /admin/catalog/cuisine-types/{cuisine_type_id} [delete]
func (h *Handler) DeleteCuisineType(w http.ResponseWriter, r *http.Request) {
	const op = "catalog.Handler.DeleteCuisineType"
	log := h.log.With(slog.String("op", op))

	id, err := decode.Int(r, params.CuisineTypeID)
	if err != nil {
		respond.Error(w, r, log, err)
		return
	}

	if err := h.service.DeleteCuisineType(r.Context(), id); err != nil {
		respond.Error(w, r, log, err)
		return
	}

	render.NoContent(w, r)
}

// CreateProductCategory godoc
// @Summary      Create a new product category
// @Description  Creates a new product category and assigns the specified user as the owner. Caller must be an admin.
// @Tags         catalog
// @Security     Bearer
// @Accept       json
// @Produce      json
// @Param        input body CreateProductCategoryRequest true "Product category data and owner ID"
// @Success      201 {object} ProductCategoryResponse "Successfully created product category"
// @Failure      400 {object} respond.ErrorResponse "Validation error or invalid JSON"
// @Failure      401 {object} respond.ErrorResponse "Unauthorized"
// @Failure      403 {object} respond.ErrorResponse "Permission denied"
// @Failure      500 {object} respond.ErrorResponse "Internal server error"
// @Router       /admin/catalog/product-categories [post]
func (h *Handler) CreateProductCategory(w http.ResponseWriter, r *http.Request) {
	const op = "catalog.Handler.CreateProductCategory"
	log := h.log.With(slog.String("op", op))

	var req CreateProductCategoryRequest
	if err := render.DecodeJSON(r.Body, &req); err != nil {
		respond.Error(w, r, log, domain.ErrInvalidJSON)
		return
	}

	if err := h.validate.Struct(req); err != nil {
		respond.Error(w, r, log, err)
		return
	}

	input := domain.CreateProductCategoryInput{
		Name:  req.Name,
		Label: req.Label,
	}

	productCategory, err := h.service.CreateProductCategory(r.Context(), input)
	if err != nil {
		respond.Error(w, r, log, err)
		return
	}

	render.Status(r, http.StatusCreated)
	render.JSON(w, r, ProductCategoryResponse{ProductCategory: productCategory})
}

// DeleteProductCategory godoc
// @Summary      Delete a product category
// @Description  Deletes a product category from the system. Caller must be an admin.
// @Tags         catalog
// @Security     Bearer
// @Produce      json
// @Param        product_category_id path int true "Product category ID"
// @Success      204 "No Content - Product category successfully deleted"
// @Failure      400 {object} respond.ErrorResponse "Invalid ID"
// @Failure      401 {object} respond.ErrorResponse "Unauthorized"
// @Failure      403 {object} respond.ErrorResponse "Permission denied"
// @Failure      404 {object} respond.ErrorResponse "Product category not found"
// @Failure      500 {object} respond.ErrorResponse "Internal server error"
// @Router       /admin/catalog/product-categories/{product_category_id} [delete]
func (h *Handler) DeleteProductCategory(w http.ResponseWriter, r *http.Request) {
	const op = "catalog.Handler.DeleteProductCategory"
	log := h.log.With(slog.String("op", op))

	id, err := decode.Int(r, params.ProductCategoryID)
	if err != nil {
		respond.Error(w, r, log, err)
		return
	}

	if err := h.service.DeleteProductCategory(r.Context(), id); err != nil {
		respond.Error(w, r, log, err)
		return
	}

	render.NoContent(w, r)
}
