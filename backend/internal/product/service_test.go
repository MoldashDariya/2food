package product_test

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gitlab.com/2foood/backend/internal/domain"
	"gitlab.com/2foood/backend/internal/product"
	"gitlab.com/2foood/backend/internal/product/mocks"
)

func setupTest() (*mocks.ProductRepositoryMock, *product.Service) {
	tm := &mocks.TransactionManagerFake{}
	productRepo := &mocks.ProductRepositoryMock{}
	svc := product.NewService(tm, productRepo)
	return productRepo, svc
}

func TestProductService_Create(t *testing.T) {
	restID := uuid.New()
	catID := 1

	tests := []struct {
		name    string
		input   domain.CreateProductInput
		mockFn  func(productRepo *mocks.ProductRepositoryMock)
		wantErr bool
	}{
		{
			name: "db error",
			input: domain.CreateProductInput{
				RestaurantID: restID,
				Name:         "Test Product",
			},
			mockFn: func(productRepo *mocks.ProductRepositoryMock) {
				productRepo.On("Create", mock.Anything, mock.AnythingOfType("*domain.Product")).
					Return(errors.New("db error"))
			},
			wantErr: true,
		},
		{
			name: "success without category",
			input: domain.CreateProductInput{
				RestaurantID: restID,
				CategoryID:   nil,
				Name:         "No Category Product",
			},
			mockFn: func(productRepo *mocks.ProductRepositoryMock) {
				productRepo.On("Create", mock.Anything, mock.AnythingOfType("*domain.Product")).Return(nil)
			},
			wantErr: false,
		},
		{
			name: "success with category",
			input: domain.CreateProductInput{
				RestaurantID: restID,
				CategoryID:   &catID,
				Name:         "Burger",
			},
			mockFn: func(productRepo *mocks.ProductRepositoryMock) {
				productRepo.On("Create", mock.Anything, mock.AnythingOfType("*domain.Product")).Return(nil)
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			productRepo, svc := setupTest()
			tt.mockFn(productRepo)

			res, err := svc.Create(context.Background(), tt.input)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, res)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, res)
				assert.Equal(t, tt.input.Name, res.Name)
			}
			productRepo.AssertExpectations(t)
		})
	}
}

func TestProductService_Getters(t *testing.T) {
	restID := uuid.New()
	staff := &domain.RestaurantStaff{RestaurantID: restID}

	t.Run("GetByRestaurant_PermissionDenied", func(t *testing.T) {
		t.Parallel()
		productRepo, svc := setupTest()

		otherRestID := uuid.New()
		res, err := svc.GetByRestaurant(context.Background(), staff, otherRestID)

		assert.Error(t, err)
		assert.ErrorIs(t, err, domain.ErrPermissionDenied)
		assert.Nil(t, res)
		productRepo.AssertExpectations(t)
	})

	t.Run("GetByRestaurant_Success", func(t *testing.T) {
		t.Parallel()
		productRepo, svc := setupTest()

		productRepo.On("GetByRestaurant", mock.Anything, restID).
			Return([]*domain.Product{{RestaurantID: restID}}, nil)
		res, err := svc.GetByRestaurant(context.Background(), staff, restID)

		assert.NoError(t, err)
		assert.Len(t, res, 1)
		productRepo.AssertExpectations(t)
	})

	t.Run("GetActiveByRestaurant_Success", func(t *testing.T) {
		t.Parallel()
		productRepo, svc := setupTest()

		productRepo.On("GetActiveByRestaurant", mock.Anything, restID).
			Return([]*domain.Product{{IsActive: true}}, nil)
		res, err := svc.GetActiveByRestaurant(context.Background(), restID)

		assert.NoError(t, err)
		assert.Len(t, res, 1)
		productRepo.AssertExpectations(t)
	})

	t.Run("GetByID_PermissionDenied", func(t *testing.T) {
		t.Parallel()
		productID := uuid.New()
		productRepo, svc := setupTest()

		productRepo.On("GetByID", mock.Anything, productID).
			Return(&domain.Product{RestaurantID: uuid.New()}, nil)

		res, err := svc.GetByID(context.Background(), staff, productID)

		assert.Error(t, err)
		assert.ErrorIs(t, err, domain.ErrPermissionDenied)
		assert.Nil(t, res)
		productRepo.AssertExpectations(t)
	})

	t.Run("GetByID_Success", func(t *testing.T) {
		t.Parallel()
		productID := uuid.New()
		productRepo, svc := setupTest()

		productRepo.On("GetByID", mock.Anything, productID).
			Return(&domain.Product{ID: productID, RestaurantID: restID}, nil)

		res, err := svc.GetByID(context.Background(), staff, productID)

		assert.NoError(t, err)
		assert.Equal(t, productID, res.ID)
		productRepo.AssertExpectations(t)
	})
}

func TestProductService_GetCategoriesByRestaurant(t *testing.T) {
	restID := uuid.New()

	tests := []struct {
		name    string
		mockFn  func(productRepo *mocks.ProductRepositoryMock)
		wantErr bool
		wantLen int
	}{
		{
			name: "db error",
			mockFn: func(productRepo *mocks.ProductRepositoryMock) {
				productRepo.On("GetCategoriesByRestaurant", mock.Anything, restID).
					Return(nil, errors.New("db error"))
			},
			wantErr: true,
		},
		{
			name: "success empty",
			mockFn: func(productRepo *mocks.ProductRepositoryMock) {
				productRepo.On("GetCategoriesByRestaurant", mock.Anything, restID).
					Return([]*domain.ProductCategory{}, nil)
			},
			wantErr: false,
			wantLen: 0,
		},
		{
			name: "success with categories",
			mockFn: func(productRepo *mocks.ProductRepositoryMock) {
				productRepo.On("GetCategoriesByRestaurant", mock.Anything, restID).
					Return([]*domain.ProductCategory{
						{ID: 1, Name: "bakery", Label: "Bakery"},
						{ID: 2, Name: "drinks", Label: "Drinks"},
					}, nil)
			},
			wantErr: false,
			wantLen: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			productRepo, svc := setupTest()
			tt.mockFn(productRepo)

			res, err := svc.GetCategoriesByRestaurant(context.Background(), restID)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, res)
			} else {
				assert.NoError(t, err)
				assert.Len(t, res, tt.wantLen)
			}
			productRepo.AssertExpectations(t)
		})
	}
}

func TestProductService_UpdateModifiers(t *testing.T) {
	productID := uuid.New()
	restID := uuid.New()
	staff := &domain.RestaurantStaff{RestaurantID: restID}

	tests := []struct {
		name    string
		method  string
		mockFn  func(productRepo *mocks.ProductRepositoryMock)
		wantErr bool
	}{
		{
			name:   "product not found",
			method: "UpdateStatus",
			mockFn: func(productRepo *mocks.ProductRepositoryMock) {
				productRepo.On("GetByID", mock.Anything, productID).Return(nil, domain.ErrNotFound)
			},
			wantErr: true,
		},
		{
			name:   "permission denied",
			method: "UpdateStock",
			mockFn: func(productRepo *mocks.ProductRepositoryMock) {
				productRepo.On("GetByID", mock.Anything, productID).
					Return(&domain.Product{RestaurantID: uuid.New()}, nil)
			},
			wantErr: true,
		},
		{
			name:   "success UpdateStatus",
			method: "UpdateStatus",
			mockFn: func(productRepo *mocks.ProductRepositoryMock) {
				productRepo.On("GetByID", mock.Anything, productID).
					Return(&domain.Product{ID: productID, RestaurantID: restID}, nil)
				productRepo.On("UpdateStatus", mock.Anything, productID, false).Return(nil)
			},
			wantErr: false,
		},
		{
			name:   "success UpdateStock",
			method: "UpdateStock",
			mockFn: func(productRepo *mocks.ProductRepositoryMock) {
				productRepo.On("GetByID", mock.Anything, productID).
					Return(&domain.Product{ID: productID, RestaurantID: restID}, nil)
				productRepo.On("UpdateStock", mock.Anything, productID, 10).Return(nil)
			},
			wantErr: false,
		},
		{
			name:   "success Update",
			method: "Update",
			mockFn: func(productRepo *mocks.ProductRepositoryMock) {
				productRepo.On("GetByID", mock.Anything, productID).
					Return(&domain.Product{ID: productID, RestaurantID: restID}, nil)
				productRepo.On("Update", mock.Anything, mock.AnythingOfType("*domain.Product")).Return(nil)
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			productRepo, svc := setupTest()
			tt.mockFn(productRepo)

			var err error
			switch tt.method {
			case "UpdateStatus":
				err = svc.UpdateStatus(context.Background(), staff, productID, false)
			case "UpdateStock":
				err = svc.UpdateStock(context.Background(), staff, productID, 10)
			case "Update":
				_, err = svc.Update(context.Background(), staff, productID, domain.UpdateProductInput{})
			}

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			productRepo.AssertExpectations(t)
		})
	}
}

func TestProductService_Delete(t *testing.T) {
	productID := uuid.New()
	restID := uuid.New()
	staff := &domain.RestaurantStaff{RestaurantID: restID}

	tests := []struct {
		name    string
		mockFn  func(productRepo *mocks.ProductRepositoryMock)
		wantErr bool
	}{
		{
			name: "product not found",
			mockFn: func(productRepo *mocks.ProductRepositoryMock) {
				productRepo.On("GetByID", mock.Anything, productID).Return(nil, domain.ErrNotFound)
			},
			wantErr: true,
		},
		{
			name: "permission denied",
			mockFn: func(productRepo *mocks.ProductRepositoryMock) {
				productRepo.On("GetByID", mock.Anything, productID).
					Return(&domain.Product{ID: productID, RestaurantID: uuid.New()}, nil)
			},
			wantErr: true,
		},
		{
			name: "success",
			mockFn: func(productRepo *mocks.ProductRepositoryMock) {
				productRepo.On("GetByID", mock.Anything, productID).
					Return(&domain.Product{ID: productID, RestaurantID: restID}, nil)
				productRepo.On("Delete", mock.Anything, productID).Return(nil)
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			productRepo, svc := setupTest()
			tt.mockFn(productRepo)

			err := svc.Delete(context.Background(), staff, productID)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			productRepo.AssertExpectations(t)
		})
	}
}
