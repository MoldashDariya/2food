package catalog_test

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gitlab.com/2foood/backend/internal/catalog"
	"gitlab.com/2foood/backend/internal/catalog/mocks"
	"gitlab.com/2foood/backend/internal/domain"
)

func setupTest() (*mocks.CityRepositoryMock, *mocks.CuisineTypeRepositoryMock, *mocks.ProductCategoryRepositoryMock, *catalog.Service) {
	cityRepo := &mocks.CityRepositoryMock{}
	cuisineRepo := &mocks.CuisineTypeRepositoryMock{}
	productCategoryRepo := &mocks.ProductCategoryRepositoryMock{}
	svc := catalog.NewService(cityRepo, cuisineRepo, productCategoryRepo)
	return cityRepo, cuisineRepo, productCategoryRepo, svc
}

func TestCatalogService_CreateCity(t *testing.T) {
	input := domain.CreateCityInput{Name: "Almaty", Label: "almaty"}

	tests := []struct {
		name    string
		mockFn  func(r *mocks.CityRepositoryMock)
		wantErr bool
	}{
		{
			name: "db error",
			mockFn: func(r *mocks.CityRepositoryMock) {
				r.On("Create", mock.Anything, mock.AnythingOfType("*domain.City")).Return(errors.New("db error"))
			},
			wantErr: true,
		},
		{
			name: "success",
			mockFn: func(r *mocks.CityRepositoryMock) {
				r.On("Create", mock.Anything, mock.AnythingOfType("*domain.City")).Return(nil)
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			cityRepo, _, _, svc := setupTest()
			tt.mockFn(cityRepo)

			city, err := svc.CreateCity(context.Background(), input)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, city)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, input.Name, city.Name)
			}
			cityRepo.AssertExpectations(t)
		})
	}
}

func TestCatalogService_CityGetters(t *testing.T) {
	t.Run("GetCities_success", func(t *testing.T) {
		t.Parallel()
		cityRepo, _, _, svc := setupTest()

		cityRepo.On("GetAll", mock.Anything).Return([]*domain.City{{ID: 1, Name: "Almaty"}}, nil)
		res, err := svc.GetCities(context.Background())

		assert.NoError(t, err)
		assert.Len(t, res, 1)
		cityRepo.AssertExpectations(t)
	})

	t.Run("GetCityByID_success", func(t *testing.T) {
		t.Parallel()
		cityRepo, _, _, svc := setupTest()

		cityRepo.On("GetByID", mock.Anything, 1).Return(&domain.City{ID: 1, Name: "Almaty"}, nil)
		res, err := svc.GetCityByID(context.Background(), 1)

		assert.NoError(t, err)
		assert.Equal(t, 1, res.ID)
		cityRepo.AssertExpectations(t)
	})

	t.Run("GetCityByID_not_found", func(t *testing.T) {
		t.Parallel()
		cityRepo, _, _, svc := setupTest()

		cityRepo.On("GetByID", mock.Anything, 99).Return(nil, domain.ErrNotFound)
		res, err := svc.GetCityByID(context.Background(), 99)

		assert.Error(t, err)
		assert.Nil(t, res)
		cityRepo.AssertExpectations(t)
	})

	t.Run("GetCityByName_success", func(t *testing.T) {
		t.Parallel()
		cityRepo, _, _, svc := setupTest()

		cityRepo.On("GetByName", mock.Anything, "Almaty").Return(&domain.City{ID: 1, Name: "Almaty"}, nil)
		res, err := svc.GetCityByName(context.Background(), "Almaty")

		assert.NoError(t, err)
		assert.Equal(t, "Almaty", res.Name)
		cityRepo.AssertExpectations(t)
	})
}

func TestCatalogService_DeleteCity(t *testing.T) {
	tests := []struct {
		name    string
		mockFn  func(r *mocks.CityRepositoryMock)
		wantErr bool
	}{
		{
			name: "not found",
			mockFn: func(r *mocks.CityRepositoryMock) {
				r.On("Delete", mock.Anything, 1).Return(domain.ErrNotFound)
			},
			wantErr: true,
		},
		{
			name: "success",
			mockFn: func(r *mocks.CityRepositoryMock) {
				r.On("Delete", mock.Anything, 1).Return(nil)
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			cityRepo, _, _, svc := setupTest()
			tt.mockFn(cityRepo)

			err := svc.DeleteCity(context.Background(), 1)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			cityRepo.AssertExpectations(t)
		})
	}
}

func TestCatalogService_CreateCuisineType(t *testing.T) {
	input := domain.CreateCuisineTypeInput{Name: "Italian", Label: "italian"}

	tests := []struct {
		name    string
		mockFn  func(r *mocks.CuisineTypeRepositoryMock)
		wantErr bool
	}{
		{
			name: "db error",
			mockFn: func(r *mocks.CuisineTypeRepositoryMock) {
				r.On("Create", mock.Anything, mock.AnythingOfType("*domain.CuisineType")).Return(errors.New("db error"))
			},
			wantErr: true,
		},
		{
			name: "success",
			mockFn: func(r *mocks.CuisineTypeRepositoryMock) {
				r.On("Create", mock.Anything, mock.AnythingOfType("*domain.CuisineType")).Return(nil)
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			_, cuisineRepo, _, svc := setupTest()
			tt.mockFn(cuisineRepo)

			cuisine, err := svc.CreateCuisineType(context.Background(), input)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, cuisine)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, input.Name, cuisine.Name)
			}
			cuisineRepo.AssertExpectations(t)
		})
	}
}

func TestCatalogService_CuisineTypeGetters(t *testing.T) {
	t.Run("GetCuisineTypes_success", func(t *testing.T) {
		t.Parallel()
		_, cuisineRepo, _, svc := setupTest()

		cuisineRepo.On("GetAll", mock.Anything).Return([]*domain.CuisineType{{ID: 1, Name: "Italian"}}, nil)
		res, err := svc.GetCuisineTypes(context.Background())

		assert.NoError(t, err)
		assert.Len(t, res, 1)
		cuisineRepo.AssertExpectations(t)
	})

	t.Run("GetCuisineTypeByID_success", func(t *testing.T) {
		t.Parallel()
		_, cuisineRepo, _, svc := setupTest()

		cuisineRepo.On("GetByID", mock.Anything, 1).Return(&domain.CuisineType{ID: 1, Name: "Italian"}, nil)
		res, err := svc.GetCuisineTypeByID(context.Background(), 1)

		assert.NoError(t, err)
		assert.Equal(t, 1, res.ID)
		cuisineRepo.AssertExpectations(t)
	})

	t.Run("GetCuisineTypeByID_not_found", func(t *testing.T) {
		t.Parallel()
		_, cuisineRepo, _, svc := setupTest()

		cuisineRepo.On("GetByID", mock.Anything, 99).Return(nil, domain.ErrNotFound)
		res, err := svc.GetCuisineTypeByID(context.Background(), 99)

		assert.Error(t, err)
		assert.Nil(t, res)
		cuisineRepo.AssertExpectations(t)
	})

	t.Run("GetCuisineTypeByName_success", func(t *testing.T) {
		t.Parallel()
		_, cuisineRepo, _, svc := setupTest()

		cuisineRepo.On("GetByName", mock.Anything, "Italian").Return(&domain.CuisineType{ID: 1, Name: "Italian"}, nil)
		res, err := svc.GetCuisineTypeByName(context.Background(), "Italian")

		assert.NoError(t, err)
		assert.Equal(t, "Italian", res.Name)
		cuisineRepo.AssertExpectations(t)
	})
}

func TestCatalogService_DeleteCuisineType(t *testing.T) {
	tests := []struct {
		name    string
		mockFn  func(r *mocks.CuisineTypeRepositoryMock)
		wantErr bool
	}{
		{
			name: "not found",
			mockFn: func(r *mocks.CuisineTypeRepositoryMock) {
				r.On("Delete", mock.Anything, 1).Return(domain.ErrNotFound)
			},
			wantErr: true,
		},
		{
			name: "success",
			mockFn: func(r *mocks.CuisineTypeRepositoryMock) {
				r.On("Delete", mock.Anything, 1).Return(nil)
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			_, cuisineRepo, _, svc := setupTest()
			tt.mockFn(cuisineRepo)

			err := svc.DeleteCuisineType(context.Background(), 1)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			cuisineRepo.AssertExpectations(t)
		})
	}
}

func TestCatalogService_CreateProductCategory(t *testing.T) {
	input := domain.CreateProductCategoryInput{Name: "Burgers", Label: "burgers"}

	tests := []struct {
		name    string
		mockFn  func(r *mocks.ProductCategoryRepositoryMock)
		wantErr bool
	}{
		{
			name: "db error",
			mockFn: func(r *mocks.ProductCategoryRepositoryMock) {
				r.On("Create", mock.Anything, mock.AnythingOfType("*domain.ProductCategory")).Return(errors.New("db error"))
			},
			wantErr: true,
		},
		{
			name: "already exists",
			mockFn: func(r *mocks.ProductCategoryRepositoryMock) {
				r.On("Create", mock.Anything, mock.AnythingOfType("*domain.ProductCategory")).Return(domain.ErrAlreadyExists)
			},
			wantErr: true,
		},
		{
			name: "success",
			mockFn: func(r *mocks.ProductCategoryRepositoryMock) {
				r.On("Create", mock.Anything, mock.AnythingOfType("*domain.ProductCategory")).Return(nil)
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			_, _, productCategoryRepo, svc := setupTest()
			tt.mockFn(productCategoryRepo)

			category, err := svc.CreateProductCategory(context.Background(), input)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, category)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, input.Name, category.Name)
				assert.Equal(t, input.Label, category.Label)
			}
			productCategoryRepo.AssertExpectations(t)
		})
	}
}

func TestCatalogService_ProductCategoryGetters(t *testing.T) {
	t.Run("GetProductCategories_success", func(t *testing.T) {
		t.Parallel()
		_, _, productCategoryRepo, svc := setupTest()

		productCategoryRepo.On("GetAll", mock.Anything).Return([]*domain.ProductCategory{{ID: 1, Name: "Burgers"}}, nil)
		res, err := svc.GetProductCategories(context.Background())

		assert.NoError(t, err)
		assert.Len(t, res, 1)
		productCategoryRepo.AssertExpectations(t)
	})

	t.Run("GetProductCategories_db_error", func(t *testing.T) {
		t.Parallel()
		_, _, productCategoryRepo, svc := setupTest()

		productCategoryRepo.On("GetAll", mock.Anything).Return(nil, errors.New("db error"))
		res, err := svc.GetProductCategories(context.Background())

		assert.Error(t, err)
		assert.Nil(t, res)
		productCategoryRepo.AssertExpectations(t)
	})

	t.Run("GetProductCategoryByID_success", func(t *testing.T) {
		t.Parallel()
		_, _, productCategoryRepo, svc := setupTest()

		productCategoryRepo.On("GetByID", mock.Anything, 1).Return(&domain.ProductCategory{ID: 1, Name: "Burgers"}, nil)
		res, err := svc.GetProductCategoryByID(context.Background(), 1)

		assert.NoError(t, err)
		assert.Equal(t, 1, res.ID)
		assert.Equal(t, "Burgers", res.Name)
		productCategoryRepo.AssertExpectations(t)
	})

	t.Run("GetProductCategoryByID_not_found", func(t *testing.T) {
		t.Parallel()
		_, _, productCategoryRepo, svc := setupTest()

		productCategoryRepo.On("GetByID", mock.Anything, 99).Return(nil, domain.ErrNotFound)
		res, err := svc.GetProductCategoryByID(context.Background(), 99)

		assert.Error(t, err)
		assert.Nil(t, res)
		productCategoryRepo.AssertExpectations(t)
	})
}

func TestCatalogService_DeleteProductCategory(t *testing.T) {
	tests := []struct {
		name    string
		mockFn  func(r *mocks.ProductCategoryRepositoryMock)
		wantErr bool
	}{
		{
			name: "not found",
			mockFn: func(r *mocks.ProductCategoryRepositoryMock) {
				r.On("Delete", mock.Anything, 1).Return(domain.ErrNotFound)
			},
			wantErr: true,
		},
		{
			name: "success",
			mockFn: func(r *mocks.ProductCategoryRepositoryMock) {
				r.On("Delete", mock.Anything, 1).Return(nil)
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			_, _, productCategoryRepo, svc := setupTest()
			tt.mockFn(productCategoryRepo)

			err := svc.DeleteProductCategory(context.Background(), 1)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			productCategoryRepo.AssertExpectations(t)
		})
	}
}
