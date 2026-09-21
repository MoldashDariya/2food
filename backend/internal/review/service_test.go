package review_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gitlab.com/2foood/backend/internal/domain"
	"gitlab.com/2foood/backend/internal/review"
	"gitlab.com/2foood/backend/internal/review/mocks"
)

func setupTest() (*review.Service, *mocks.ReviewRepositoryMock, *mocks.OrderRepositoryMock, *mocks.RestaurantRepositoryMock) {
	tm := &mocks.TransactionManagerFake{}
	reviewRepo := &mocks.ReviewRepositoryMock{}
	orderRepo := &mocks.OrderRepositoryMock{}
	restRepo := &mocks.RestaurantRepositoryMock{}

	svc := review.NewService(tm, reviewRepo, orderRepo, restRepo)
	return svc, reviewRepo, orderRepo, restRepo
}

func TestService_Create(t *testing.T) {
	userID := uuid.New()
	orderID := uuid.New()
	restID := uuid.New()
	comment := "Great!"

	tests := []struct {
		name    string
		input   domain.CreateReviewInput
		mockFn  func(reviewRepo *mocks.ReviewRepositoryMock, orderRepo *mocks.OrderRepositoryMock, restRepo *mocks.RestaurantRepositoryMock)
		wantErr error
	}{
		{
			name: "success",
			input: domain.CreateReviewInput{
				UserID:  userID,
				OrderID: orderID,
				Rating:  5,
				Comment: &comment,
			},
			mockFn: func(reviewRepo *mocks.ReviewRepositoryMock, orderRepo *mocks.OrderRepositoryMock, restRepo *mocks.RestaurantRepositoryMock) {
				order := &domain.Order{
					ID:           orderID,
					UserID:       userID,
					RestaurantID: restID,
					Status:       domain.OrderStatusCollected,
				}
				orderRepo.On("GetByID", mock.Anything, orderID).Return(order, nil)
				reviewRepo.On("ExistsByOrder", mock.Anything, orderID).Return(false, nil)

				reviewRepo.On("Create", mock.Anything, mock.AnythingOfType("*domain.Review")).Return(nil)
				restRepo.On("IncrementRating", mock.Anything, restID, int8(5)).Return(nil)
				orderRepo.On("SetReviewed", mock.Anything, orderID).Return(nil)
			},
			wantErr: nil,
		},
		{
			name: "error permission denied - not owner",
			input: domain.CreateReviewInput{
				UserID:  uuid.New(),
				OrderID: orderID,
			},
			mockFn: func(reviewRepo *mocks.ReviewRepositoryMock, orderRepo *mocks.OrderRepositoryMock, restRepo *mocks.RestaurantRepositoryMock) {
				order := &domain.Order{
					ID:     orderID,
					UserID: uuid.New(),
				}
				orderRepo.On("GetByID", mock.Anything, orderID).Return(order, nil)
			},
			wantErr: domain.ErrPermissionDenied,
		},
		{
			name: "error order not completed",
			input: domain.CreateReviewInput{
				UserID:  userID,
				OrderID: orderID,
			},
			mockFn: func(reviewRepo *mocks.ReviewRepositoryMock, orderRepo *mocks.OrderRepositoryMock, restRepo *mocks.RestaurantRepositoryMock) {
				order := &domain.Order{
					ID:     orderID,
					UserID: userID,
					Status: domain.OrderStatusPending,
				}
				orderRepo.On("GetByID", mock.Anything, orderID).Return(order, nil)
			},
			wantErr: domain.ErrOrderNotCompleted,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			svc, reviewRepo, orderRepo, restRepo := setupTest()
			tt.mockFn(reviewRepo, orderRepo, restRepo)

			res, err := svc.Create(context.Background(), tt.input)

			if tt.wantErr != nil {
				assert.Error(t, err)
				assert.ErrorIs(t, err, tt.wantErr)
				assert.Nil(t, res)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, res)
				assert.Equal(t, tt.input.Rating, res.Rating)
			}

			reviewRepo.AssertExpectations(t)
			orderRepo.AssertExpectations(t)
			restRepo.AssertExpectations(t)
		})
	}
}

func TestService_GetPhotosByRestaurant(t *testing.T) {
	restID := uuid.New()

	tests := []struct {
		name    string
		mockFn  func(reviewRepo *mocks.ReviewRepositoryMock)
		wantErr bool
		wantLen int
	}{
		{
			name: "success empty",
			mockFn: func(reviewRepo *mocks.ReviewRepositoryMock) {
				reviewRepo.On("GetPhotosByRestaurant", mock.Anything, restID).
					Return([]*domain.ReviewPhoto{}, nil)
			},
			wantErr: false,
			wantLen: 0,
		},
		{
			name: "success with photos",
			mockFn: func(reviewRepo *mocks.ReviewRepositoryMock) {
				reviewRepo.On("GetPhotosByRestaurant", mock.Anything, restID).
					Return([]*domain.ReviewPhoto{
						{URL: "img1.jpg", ReviewID: uuid.New(), Rating: 5},
						{URL: "img2.jpg", ReviewID: uuid.New(), Rating: 3},
					}, nil)
			},
			wantErr: false,
			wantLen: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			svc, reviewRepo, _, _ := setupTest()
			tt.mockFn(reviewRepo)

			res, err := svc.GetPhotosByRestaurant(context.Background(), restID)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, res)
			} else {
				assert.NoError(t, err)
				assert.Len(t, res, tt.wantLen)
			}

			reviewRepo.AssertExpectations(t)
		})
	}
}

func TestService_GetByOrder(t *testing.T) {
	actorID := uuid.New()
	orderID := uuid.New()

	tests := []struct {
		name    string
		mockFn  func(reviewRepo *mocks.ReviewRepositoryMock, orderRepo *mocks.OrderRepositoryMock)
		wantErr error
	}{
		{
			name: "success",
			mockFn: func(reviewRepo *mocks.ReviewRepositoryMock, orderRepo *mocks.OrderRepositoryMock) {
				order := &domain.Order{
					ID:     orderID,
					UserID: actorID,
				}
				expectedReview := &domain.Review{
					ID:      uuid.New(),
					OrderID: orderID,
				}

				orderRepo.On("GetByID", mock.Anything, orderID).Return(order, nil)
				reviewRepo.On("GetByOrder", mock.Anything, orderID).Return(expectedReview, nil)
			},
			wantErr: nil,
		},
		{
			name: "error permission denied - wrong user",
			mockFn: func(reviewRepo *mocks.ReviewRepositoryMock, orderRepo *mocks.OrderRepositoryMock) {
				order := &domain.Order{
					ID:     orderID,
					UserID: uuid.New(),
				}
				orderRepo.On("GetByID", mock.Anything, orderID).Return(order, nil)
			},
			wantErr: domain.ErrPermissionDenied,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			svc, reviewRepo, orderRepo, _ := setupTest()
			tt.mockFn(reviewRepo, orderRepo)

			res, err := svc.GetByOrder(context.Background(), actorID, orderID)

			if tt.wantErr != nil {
				assert.Error(t, err)
				assert.ErrorIs(t, err, tt.wantErr)
				assert.Nil(t, res)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, res)
			}

			reviewRepo.AssertExpectations(t)
			orderRepo.AssertExpectations(t)
		})
	}
}

func TestService_Getters(t *testing.T) {
	restID := uuid.New()
	userID := uuid.New()

	t.Run("GetAll", func(t *testing.T) {
		t.Parallel()
		svc, reviewRepo, _, _ := setupTest()

		expected := []*domain.Review{{ID: uuid.New()}}
		reviewRepo.On("GetAll", mock.Anything).Return(expected, nil)

		res, err := svc.GetAll(context.Background())

		assert.NoError(t, err)
		assert.Equal(t, expected, res)
		reviewRepo.AssertExpectations(t)
	})

	t.Run("GetByRestaurant", func(t *testing.T) {
		t.Parallel()
		svc, reviewRepo, _, _ := setupTest()

		expected := []*domain.Review{{ID: uuid.New()}}
		reviewRepo.On("GetByRestaurant", mock.Anything, restID).Return(expected, nil)

		res, err := svc.GetByRestaurant(context.Background(), restID)

		assert.NoError(t, err)
		assert.Equal(t, expected, res)
		reviewRepo.AssertExpectations(t)
	})

	t.Run("GetByRestaurantAll", func(t *testing.T) {
		t.Parallel()
		svc, reviewRepo, _, _ := setupTest()

		expected := []*domain.Review{{ID: uuid.New()}}
		reviewRepo.On("GetByRestaurantAll", mock.Anything, restID).Return(expected, nil)

		res, err := svc.GetByRestaurantAll(context.Background(), restID)

		assert.NoError(t, err)
		assert.Equal(t, expected, res)
		reviewRepo.AssertExpectations(t)
	})

	t.Run("GetByUser", func(t *testing.T) {
		t.Parallel()
		svc, reviewRepo, _, _ := setupTest()

		expected := []*domain.Review{{ID: uuid.New()}}
		reviewRepo.On("GetByUser", mock.Anything, userID).Return(expected, nil)

		res, err := svc.GetByUser(context.Background(), userID)

		assert.NoError(t, err)
		assert.Equal(t, expected, res)
		reviewRepo.AssertExpectations(t)
	})
}

func TestService_AddReply(t *testing.T) {
	reviewID := uuid.New()
	restID := uuid.New()
	actor := &domain.RestaurantStaff{RestaurantID: restID}
	reply := "Thank you!"

	tests := []struct {
		name    string
		mockFn  func(reviewRepo *mocks.ReviewRepositoryMock)
		wantErr error
	}{
		{
			name: "success",
			mockFn: func(reviewRepo *mocks.ReviewRepositoryMock) {
				review := &domain.Review{
					ID:           reviewID,
					RestaurantID: restID,
				}
				reviewRepo.On("GetByID", mock.Anything, reviewID).Return(review, nil)
				reviewRepo.On("SetReply", mock.Anything, mock.AnythingOfType("*domain.Review")).Return(nil)
			},
			wantErr: nil,
		},
		{
			name: "error permission denied - wrong restaurant",
			mockFn: func(reviewRepo *mocks.ReviewRepositoryMock) {
				review := &domain.Review{
					ID:           reviewID,
					RestaurantID: uuid.New(),
				}
				reviewRepo.On("GetByID", mock.Anything, reviewID).Return(review, nil)
			},
			wantErr: domain.ErrPermissionDenied,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			svc, reviewRepo, _, _ := setupTest()
			tt.mockFn(reviewRepo)

			res, err := svc.AddReply(context.Background(), actor, reviewID, reply)

			if tt.wantErr != nil {
				assert.Error(t, err)
				assert.ErrorIs(t, err, tt.wantErr)
				assert.Nil(t, res)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, res)
				assert.Equal(t, &reply, res.Reply)
			}

			reviewRepo.AssertExpectations(t)
		})
	}
}

func TestService_Visibility(t *testing.T) {
	reviewID := uuid.New()

	tests := []struct {
		name    string
		method  string
		mockFn  func(reviewRepo *mocks.ReviewRepositoryMock)
		wantErr bool
	}{
		{
			name:   "success Hide",
			method: "Hide",
			mockFn: func(reviewRepo *mocks.ReviewRepositoryMock) {
				reviewRepo.On("SetVisibility", mock.Anything, reviewID, true).Return(nil)
			},
			wantErr: false,
		},
		{
			name:   "success Unhide",
			method: "Unhide",
			mockFn: func(reviewRepo *mocks.ReviewRepositoryMock) {
				reviewRepo.On("SetVisibility", mock.Anything, reviewID, false).Return(nil)
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			svc, reviewRepo, _, _ := setupTest()
			tt.mockFn(reviewRepo)

			var err error
			switch tt.method {
			case "Hide":
				err = svc.Hide(context.Background(), reviewID)
			case "Unhide":
				err = svc.Unhide(context.Background(), reviewID)
			}

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			reviewRepo.AssertExpectations(t)
		})
	}
}
