package http

import (
	"log/slog"
	"net/http"

	"github.com/go-chi/render"
	"gitlab.com/2foood/backend/internal/pkg/ctxkey"
	"gitlab.com/2foood/backend/internal/pkg/decode"
	"gitlab.com/2foood/backend/internal/pkg/params"
	"gitlab.com/2foood/backend/internal/pkg/respond"
)

// GetCities godoc
// @Summary      Get all cities
// @Description  Returns a list of all available cities
// @Tags         catalog
// @Produce      json
// @Success      200 {object} CitiesResponse          "List of cities"
// @Failure      500 {object} respond.ErrorResponse   "Internal server error"
// @Router       /catalog/cities [get]
func (h *Handler) GetCities(w http.ResponseWriter, r *http.Request) {
	const op = "catalog.Handler.GetCities"
	log := h.log.With(slog.String("op", op))

	cities, err := h.service.GetCities(r.Context())
	if err != nil {
		respond.Error(w, r, log, err)
		return
	}

	lang := ctxkey.GetLang(r.Context())
	for _, city := range cities {
		city.Label = h.tr.TFallback(lang, "City_"+city.Name, city.Label, nil)
	}

	render.JSON(w, r, CitiesResponse{Cities: cities})
}

// GetCityByID godoc
// @Summary      Get city by ID
// @Description  Returns a single city by its ID
// @Tags         catalog
// @Produce      json
// @Param        city_id path int true "City ID"
// @Success      200 {object} CityResponse            "City"
// @Failure      400 {object} respond.ErrorResponse   "Invalid ID"
// @Failure      404 {object} respond.ErrorResponse   "City not found"
// @Failure      500 {object} respond.ErrorResponse   "Internal server error"
// @Router       /catalog/cities/{city_id} [get]
func (h *Handler) GetCityByID(w http.ResponseWriter, r *http.Request) {
	const op = "catalog.Handler.GetCityByID"
	log := h.log.With(slog.String("op", op))

	id, err := decode.Int(r, params.CityID)
	if err != nil {
		respond.Error(w, r, log, err)
		return
	}

	city, err := h.service.GetCityByID(r.Context(), id)
	if err != nil {
		respond.Error(w, r, log, err)
		return
	}

	lang := ctxkey.GetLang(r.Context())
	city.Label = h.tr.TFallback(lang, "City_"+city.Name, city.Label, nil)

	render.JSON(w, r, CityResponse{City: city})
}

// GetCuisineTypes godoc
// @Summary      Get all cuisine types
// @Description  Returns a list of all available cuisine types
// @Tags         catalog
// @Produce      json
// @Success      200 {object} CuisineTypesResponse    "List of cuisine types"
// @Failure      500 {object} respond.ErrorResponse   "Internal server error"
// @Router       /catalog/cuisine-types [get]
func (h *Handler) GetCuisineTypes(w http.ResponseWriter, r *http.Request) {
	const op = "catalog.Handler.GetCuisineTypes"
	log := h.log.With(slog.String("op", op))

	cuisineTypes, err := h.service.GetCuisineTypes(r.Context())
	if err != nil {
		respond.Error(w, r, log, err)
		return
	}

	lang := ctxkey.GetLang(r.Context())
	for _, cuisine := range cuisineTypes {
		cuisine.Label = h.tr.TFallback(lang, "Cuisine_"+cuisine.Name, cuisine.Label, nil)
	}

	render.JSON(w, r, CuisineTypesResponse{CuisineTypes: cuisineTypes})
}

// GetCuisineTypeByID godoc
// @Summary      Get cuisine type by ID
// @Description  Returns a single cuisine type by its ID
// @Tags         catalog
// @Produce      json
// @Param        cuisine_type_id path int true "Cuisine type ID"
// @Success      200 {object} CuisineTypeResponse     "Cuisine type"
// @Failure      400 {object} respond.ErrorResponse   "Invalid ID"
// @Failure      404 {object} respond.ErrorResponse   "Cuisine type not found"
// @Failure      500 {object} respond.ErrorResponse   "Internal server error"
// @Router       /catalog/cuisine-types/{cuisine_type_id} [get]
func (h *Handler) GetCuisineTypeByID(w http.ResponseWriter, r *http.Request) {
	const op = "catalog.Handler.GetCuisineTypeByID"
	log := h.log.With(slog.String("op", op))

	id, err := decode.Int(r, params.CuisineTypeID)
	if err != nil {
		respond.Error(w, r, log, err)
		return
	}

	cuisineType, err := h.service.GetCuisineTypeByID(r.Context(), id)
	if err != nil {
		respond.Error(w, r, log, err)
		return
	}

	lang := ctxkey.GetLang(r.Context())
	cuisineType.Label = h.tr.TFallback(lang, "Cuisine_"+cuisineType.Name, cuisineType.Label, nil)

	render.JSON(w, r, CuisineTypeResponse{CuisineType: cuisineType})
}

// GetProductCategories godoc
// @Summary      Get all product categories
// @Description  Returns a list of all available product categories
// @Tags         catalog
// @Produce      json
// @Success      200 {object} ProductCategoriesResponse "List of product categories"
// @Failure      500 {object} respond.ErrorResponse "Internal server error"
// @Router       /catalog/product-categories [get]
func (h *Handler) GetProductCategories(w http.ResponseWriter, r *http.Request) {
	const op = "catalog.Handler.GetProductCategories"
	log := h.log.With(slog.String("op", op))

	productCategories, err := h.service.GetProductCategories(r.Context())
	if err != nil {
		respond.Error(w, r, log, err)
		return
	}

	lang := ctxkey.GetLang(r.Context())
	for _, pc := range productCategories {
		pc.Label = h.tr.TFallback(lang, "ProdCat_"+pc.Name, pc.Label, nil)
	}

	render.JSON(w, r, ProductCategoriesResponse{ProductCategories: productCategories})
}

// GetProductCategoryByID godoc
// @Summary      Get product category by ID
// @Description  Returns a single product category by its ID
// @Tags         catalog
// @Produce      json
// @Param        product_category_id path int true "Product category ID"
// @Success      200 {object} ProductCategoryResponse "Product category"
// @Failure      400 {object} respond.ErrorResponse "Invalid ID"
// @Failure      404 {object} respond.ErrorResponse "Product category not found"
// @Failure      500 {object} respond.ErrorResponse "Internal server error"
// @Router       /catalog/product-categories/{product_category_id} [get]
func (h *Handler) GetProductCategoryByID(w http.ResponseWriter, r *http.Request) {
	const op = "catalog.Handler.GetProductCategoryByID"
	log := h.log.With(slog.String("op", op))

	id, err := decode.Int(r, params.ProductCategoryID)
	if err != nil {
		respond.Error(w, r, log, err)
		return
	}

	productCategory, err := h.service.GetProductCategoryByID(r.Context(), id)
	if err != nil {
		respond.Error(w, r, log, err)
		return
	}

	lang := ctxkey.GetLang(r.Context())
	productCategory.Label = h.tr.TFallback(lang, "ProdCat_"+productCategory.Name, productCategory.Label, nil)

	render.JSON(w, r, ProductCategoryResponse{ProductCategory: productCategory})
}
