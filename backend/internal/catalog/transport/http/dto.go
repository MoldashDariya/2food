package http

import "gitlab.com/2foood/backend/internal/domain"

type CitiesResponse struct {
	Cities []*domain.City `json:"cities"`
}

type CityResponse struct {
	City *domain.City `json:"city"`
}

type CreateCityRequest struct {
	Name  string `json:"name" validate:"required,min=2,max=50"`
	Label string `json:"label" validate:"required,min=2,max=50"`
}

type CuisineTypesResponse struct {
	CuisineTypes []*domain.CuisineType `json:"cuisine_types"`
}

type CuisineTypeResponse struct {
	CuisineType *domain.CuisineType `json:"cuisine_type"`
}

type CreateCuisineTypeRequest struct {
	Name  string `json:"name" validate:"required,min=2,max=50"`
	Label string `json:"label" validate:"required,min=2,max=50"`
}

type ProductCategoriesResponse struct {
	ProductCategories []*domain.ProductCategory `json:"product_categories"`
}

type ProductCategoryResponse struct {
	ProductCategory *domain.ProductCategory `json:"product_category"`
}

type CreateProductCategoryRequest struct {
	Name  string `json:"name" validate:"required,min=2,max=50"`
	Label string `json:"label" validate:"required,min=2,max=50"`
}
