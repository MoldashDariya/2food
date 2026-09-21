package order_test

import (
	"context"
	"errors"
	"log/slog"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gitlab.com/2foood/backend/internal/domain"
	"gitlab.com/2foood/backend/internal/order"
	"gitlab.com/2foood/backend/internal/order/mocks"
)

const testAuditDelay = 5 * time.Minute

func setupTest() (
	*mocks.TransactionManagerFake,
	*mocks.OrderRepositoryMock,
	*mocks.ProductRepositoryMock,
	*mocks.PaymentGatewayMock,
	*mocks.PaymentDispatcherMock,
	*mocks.PushDispatcherMock,
	*order.Service,
) {
	tm := &mocks.TransactionManagerFake{}
	orderRepo := &mocks.OrderRepositoryMock{}
	productRepo := &mocks.ProductRepositoryMock{}
	paymentGateway := &mocks.PaymentGatewayMock{}
	paymentDispatcher := &mocks.PaymentDispatcherMock{}
	pushDispatcher := &mocks.PushDispatcherMock{}

	svc := order.NewService(tm, slog.Default(), orderRepo, productRepo, paymentGateway, paymentDispatcher, pushDispatcher, 6, 6, testAuditDelay, "usd")
	return tm, orderRepo, productRepo, paymentGateway, paymentDispatcher, pushDispatcher, svc
}

func TestOrderService_Create(t *testing.T) {
	restID := uuid.New()
	prodID := uuid.New()
	today := time.Now().UTC().Truncate(24 * time.Hour)
	input := domain.CreateOrderInput{
		UserID:            uuid.New(),
		RestaurantID:      restID,
		Items:             []domain.CreateOrderItemInput{{ProductID: prodID, Quantity: 2}},
		PickupWindowStart: today.Add(10 * time.Hour),
		PickupWindowEnd:   today.Add(12 * time.Hour),
	}

	tests := []struct {
		name    string
		input   domain.CreateOrderInput
		mockFn  func(productRepo *mocks.ProductRepositoryMock, paymentGateway *mocks.PaymentGatewayMock, orderRepo *mocks.OrderRepositoryMock)
		wantErr bool
	}{
		{
			name:  "product not found",
			input: input,
			mockFn: func(productRepo *mocks.ProductRepositoryMock, paymentGateway *mocks.PaymentGatewayMock, orderRepo *mocks.OrderRepositoryMock) {
				productRepo.On("GetByID", mock.Anything, prodID).Return(nil, errors.New("db error"))
			},
			wantErr: true,
		},
		{
			name:  "product not active",
			input: input,
			mockFn: func(productRepo *mocks.ProductRepositoryMock, paymentGateway *mocks.PaymentGatewayMock, orderRepo *mocks.OrderRepositoryMock) {
				productRepo.On("GetByID", mock.Anything, prodID).Return(&domain.Product{IsActive: false, Name: "Test"}, nil)
			},
			wantErr: true,
		},
		{
			name:  "insufficient stock",
			input: input,
			mockFn: func(productRepo *mocks.ProductRepositoryMock, paymentGateway *mocks.PaymentGatewayMock, orderRepo *mocks.OrderRepositoryMock) {
				productRepo.On("GetByID", mock.Anything, prodID).Return(&domain.Product{
					IsActive: true, RestaurantID: restID, QuantityAvailable: 1, Name: "Test",
				}, nil)
			},
			wantErr: true,
		},
		{
			name:  "success",
			input: input,
			mockFn: func(productRepo *mocks.ProductRepositoryMock, paymentGateway *mocks.PaymentGatewayMock, orderRepo *mocks.OrderRepositoryMock) {
				productRepo.On("GetByID", mock.Anything, prodID).Return(&domain.Product{
					IsActive: true, RestaurantID: restID, QuantityAvailable: 5,
					Name: "Surprise Box", DiscountedPrice: 1500,
					PickupTimeStart: "00:00", PickupTimeEnd: "23:59",
				}, nil)
				paymentGateway.On("CreatePaymentIntent", 3000.0, mock.Anything).Return("pi_123", "secret_456", nil)
				orderRepo.On("AddPaymentHistory", mock.Anything, mock.AnythingOfType("*domain.PaymentHistoryEntry")).Return(nil)
				orderRepo.On("Create", mock.Anything, mock.AnythingOfType("*domain.Order")).Return(nil)
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			_, orderRepo, productRepo, paymentGateway, _, _, svc := setupTest()
			tt.mockFn(productRepo, paymentGateway, orderRepo)

			orderRes, secret, err := svc.Create(context.Background(), tt.input)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, orderRes)
				assert.Empty(t, secret)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, orderRes)
				assert.Equal(t, "secret_456", secret)
				assert.Equal(t, "pi_123", *orderRes.PaymentID)
			}
			productRepo.AssertExpectations(t)
			paymentGateway.AssertExpectations(t)
			orderRepo.AssertExpectations(t)
		})
	}
}

func TestOrderService_ConfirmPayment(t *testing.T) {
	orderID := uuid.New()
	userID := uuid.New()
	paymentID := "pi_123"

	tests := []struct {
		name    string
		mockFn  func(paymentGateway *mocks.PaymentGatewayMock, paymentDispatcher *mocks.PaymentDispatcherMock, orderRepo *mocks.OrderRepositoryMock, productRepo *mocks.ProductRepositoryMock, pushDispatcher *mocks.PushDispatcherMock)
		wantErr bool
	}{
		{
			name: "invalid payment ID",
			mockFn: func(paymentGateway *mocks.PaymentGatewayMock, paymentDispatcher *mocks.PaymentDispatcherMock, orderRepo *mocks.OrderRepositoryMock, productRepo *mocks.ProductRepositoryMock, pushDispatcher *mocks.PushDispatcherMock) {
				badPaymentID := "pi_wrong"
				orderRepo.On("GetByID", mock.Anything, orderID).Return(&domain.Order{
					ID: orderID, UserID: userID, PaymentID: &badPaymentID,
				}, nil)
			},
			wantErr: true,
		},
		{
			name: "out of stock during payment",
			mockFn: func(paymentGateway *mocks.PaymentGatewayMock, paymentDispatcher *mocks.PaymentDispatcherMock, orderRepo *mocks.OrderRepositoryMock, productRepo *mocks.ProductRepositoryMock, pushDispatcher *mocks.PushDispatcherMock) {
				prodID := uuid.New()
				o := &domain.Order{
					ID: orderID, UserID: userID, PaymentID: &paymentID,
					Status: domain.OrderStatusPending,
					Items:  []domain.OrderItem{{ProductID: prodID, Quantity: 5}},
				}
				orderRepo.On("GetByID", mock.Anything, orderID).Return(o, nil)
				productRepo.On("GetByIDForUpdate", mock.Anything, prodID).Return(&domain.Product{QuantityAvailable: 2}, nil)
				orderRepo.On("UpdateStatus", mock.Anything, orderID, domain.OrderStatusCancelled).Return(nil)
				orderRepo.On("AddStatusHistory", mock.Anything, mock.AnythingOfType("*domain.OrderStatusHistoryEntry")).Return(nil)
				paymentDispatcher.On("CancelAsync", mock.Anything, paymentID, orderID).Return(nil)
				pushDispatcher.On("SendAsync", mock.Anything, userID, mock.AnythingOfType("*domain.PushMessage")).Return(nil)
			},
			wantErr: false,
		},
		{
			name: "success",
			mockFn: func(paymentGateway *mocks.PaymentGatewayMock, paymentDispatcher *mocks.PaymentDispatcherMock, orderRepo *mocks.OrderRepositoryMock, productRepo *mocks.ProductRepositoryMock, pushDispatcher *mocks.PushDispatcherMock) {
				prodID := uuid.New()
				o := &domain.Order{
					ID: orderID, UserID: userID, PaymentID: &paymentID,
					Status: domain.OrderStatusPending,
					Items:  []domain.OrderItem{{ProductID: prodID, Quantity: 2}},
				}
				orderRepo.On("GetByID", mock.Anything, orderID).Return(o, nil)
				productRepo.On("GetByIDForUpdate", mock.Anything, prodID).Return(&domain.Product{QuantityAvailable: 5}, nil)
				productRepo.On("UpdateStock", mock.Anything, prodID, 3).Return(nil)
				orderRepo.On("UpdateStatus", mock.Anything, orderID, domain.OrderStatusPaid).Return(nil)
				orderRepo.On("AddStatusHistory", mock.Anything, mock.AnythingOfType("*domain.OrderStatusHistoryEntry")).Return(nil)
				paymentDispatcher.On("ConfirmAsync", mock.Anything, paymentID, orderID).Return(nil)
				paymentDispatcher.On("AuditAsync", mock.Anything, paymentID, orderID, testAuditDelay).Return(nil)
				pushDispatcher.On("SendAsync", mock.Anything, userID, mock.AnythingOfType("*domain.PushMessage")).Return(nil)
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			_, orderRepo, productRepo, paymentGateway, paymentDispatcher, pushDispatcher, svc := setupTest()
			tt.mockFn(paymentGateway, paymentDispatcher, orderRepo, productRepo, pushDispatcher)

			err := svc.ConfirmPayment(context.Background(), orderID, paymentID)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			paymentGateway.AssertExpectations(t)
			paymentDispatcher.AssertExpectations(t)
			orderRepo.AssertExpectations(t)
			productRepo.AssertExpectations(t)
			pushDispatcher.AssertExpectations(t)
		})
	}
}

func TestOrderService_ConfirmPaymentWebhook(t *testing.T) {
	orderID := uuid.New()
	userID := uuid.New()
	orderIDStr := orderID.String()
	paymentID := "pi_123"

	tests := []struct {
		name    string
		payload []byte
		sig     string
		mockFn  func(paymentGateway *mocks.PaymentGatewayMock, paymentDispatcher *mocks.PaymentDispatcherMock, orderRepo *mocks.OrderRepositoryMock, productRepo *mocks.ProductRepositoryMock, pushDispatcher *mocks.PushDispatcherMock)
		wantErr bool
	}{
		{
			name:    "signature verification fails",
			payload: []byte("payload"),
			sig:     "bad_sig",
			mockFn: func(paymentGateway *mocks.PaymentGatewayMock, paymentDispatcher *mocks.PaymentDispatcherMock, orderRepo *mocks.OrderRepositoryMock, productRepo *mocks.ProductRepositoryMock, pushDispatcher *mocks.PushDispatcherMock) {
				paymentGateway.On("VerifyWebhookSignature", []byte("payload"), "bad_sig").Return("", "", errors.New("invalid signature"))
			},
			wantErr: true,
		},
		{
			name:    "ignored webhook event",
			payload: []byte("payload"),
			sig:     "good_sig",
			mockFn: func(paymentGateway *mocks.PaymentGatewayMock, paymentDispatcher *mocks.PaymentDispatcherMock, orderRepo *mocks.OrderRepositoryMock, productRepo *mocks.ProductRepositoryMock, pushDispatcher *mocks.PushDispatcherMock) {
				paymentGateway.On("VerifyWebhookSignature", []byte("payload"), "good_sig").Return("", "", domain.ErrIgnoredWebhook)
			},
			wantErr: false,
		},
		{
			name:    "success payment confirmation",
			payload: []byte("payload"),
			sig:     "good_sig",
			mockFn: func(paymentGateway *mocks.PaymentGatewayMock, paymentDispatcher *mocks.PaymentDispatcherMock, orderRepo *mocks.OrderRepositoryMock, productRepo *mocks.ProductRepositoryMock, pushDispatcher *mocks.PushDispatcherMock) {
				paymentGateway.On("VerifyWebhookSignature", []byte("payload"), "good_sig").Return(paymentID, orderIDStr, nil)

				prodID := uuid.New()
				o := &domain.Order{
					ID: orderID, UserID: userID, PaymentID: &paymentID,
					Status: domain.OrderStatusPending,
					Items:  []domain.OrderItem{{ProductID: prodID, Quantity: 1}},
				}
				orderRepo.On("GetByID", mock.Anything, orderID).Return(o, nil)
				productRepo.On("GetByIDForUpdate", mock.Anything, prodID).Return(&domain.Product{QuantityAvailable: 2}, nil)
				productRepo.On("UpdateStock", mock.Anything, prodID, 1).Return(nil)
				orderRepo.On("UpdateStatus", mock.Anything, orderID, domain.OrderStatusPaid).Return(nil)
				orderRepo.On("AddStatusHistory", mock.Anything, mock.AnythingOfType("*domain.OrderStatusHistoryEntry")).Return(nil)
				paymentDispatcher.On("ConfirmAsync", mock.Anything, paymentID, orderID).Return(nil)
				paymentDispatcher.On("AuditAsync", mock.Anything, paymentID, orderID, testAuditDelay).Return(nil)
				pushDispatcher.On("SendAsync", mock.Anything, userID, mock.AnythingOfType("*domain.PushMessage")).Return(nil)
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			_, orderRepo, productRepo, paymentGateway, paymentDispatcher, pushDispatcher, svc := setupTest()
			tt.mockFn(paymentGateway, paymentDispatcher, orderRepo, productRepo, pushDispatcher)

			err := svc.ConfirmPaymentWebhook(context.Background(), tt.payload, tt.sig)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			paymentGateway.AssertExpectations(t)
			paymentDispatcher.AssertExpectations(t)
			orderRepo.AssertExpectations(t)
			productRepo.AssertExpectations(t)
			pushDispatcher.AssertExpectations(t)
		})
	}
}

func TestOrderService_AuditPayment(t *testing.T) {
	orderID := uuid.New()
	userID := uuid.New()
	paymentID := "pi_123"

	tests := []struct {
		name    string
		mockFn  func(paymentGateway *mocks.PaymentGatewayMock, orderRepo *mocks.OrderRepositoryMock, productRepo *mocks.ProductRepositoryMock, paymentDispatcher *mocks.PaymentDispatcherMock, pushDispatcher *mocks.PushDispatcherMock)
		wantErr bool
	}{
		{
			name: "payment not failed — skip audit",
			mockFn: func(paymentGateway *mocks.PaymentGatewayMock, orderRepo *mocks.OrderRepositoryMock, productRepo *mocks.ProductRepositoryMock, paymentDispatcher *mocks.PaymentDispatcherMock, pushDispatcher *mocks.PushDispatcherMock) {
				paymentGateway.On("IsPaymentFailed", paymentID).Return(false, nil)
			},
			wantErr: false,
		},
		{
			name: "payment failed — order is PAID — late rollback",
			mockFn: func(paymentGateway *mocks.PaymentGatewayMock, orderRepo *mocks.OrderRepositoryMock, productRepo *mocks.ProductRepositoryMock, paymentDispatcher *mocks.PaymentDispatcherMock, pushDispatcher *mocks.PushDispatcherMock) {
				prodID := uuid.New()
				paymentGateway.On("IsPaymentFailed", paymentID).Return(true, nil)
				orderRepo.On("GetByID", mock.Anything, orderID).Return(&domain.Order{
					ID: orderID, UserID: userID, Status: domain.OrderStatusPaid, PaymentID: &paymentID,
					Items: []domain.OrderItem{{ProductID: prodID, Quantity: 1}},
				}, nil)
				productRepo.On("GetByIDForUpdate", mock.Anything, prodID).Return(&domain.Product{QuantityAvailable: 2}, nil)
				productRepo.On("UpdateStock", mock.Anything, prodID, 3).Return(nil)
				orderRepo.On("UpdateStatus", mock.Anything, orderID, domain.OrderStatusCancelled).Return(nil)
				orderRepo.On("AddStatusHistory", mock.Anything, mock.AnythingOfType("*domain.OrderStatusHistoryEntry")).Return(nil)
				paymentDispatcher.On("RefundAsync", mock.Anything, paymentID, orderID).Return(nil)
				pushDispatcher.On("SendAsync", mock.Anything, userID, mock.AnythingOfType("*domain.PushMessage")).Return(nil)
			},
			wantErr: false,
		},
		{
			name: "payment failed — order is COLLECTED — mark as UNPAID",
			mockFn: func(paymentGateway *mocks.PaymentGatewayMock, orderRepo *mocks.OrderRepositoryMock, productRepo *mocks.ProductRepositoryMock, paymentDispatcher *mocks.PaymentDispatcherMock, pushDispatcher *mocks.PushDispatcherMock) {
				paymentGateway.On("IsPaymentFailed", paymentID).Return(true, nil)
				orderRepo.On("GetByID", mock.Anything, orderID).Return(&domain.Order{
					ID: orderID, UserID: userID, Status: domain.OrderStatusCollected,
				}, nil)
				orderRepo.On("UpdateStatus", mock.Anything, orderID, domain.OrderStatusUnpaid).Return(nil)
				orderRepo.On("AddStatusHistory", mock.Anything, mock.AnythingOfType("*domain.OrderStatusHistoryEntry")).Return(nil)
				orderRepo.On("AddPaymentHistory", mock.Anything, mock.AnythingOfType("*domain.PaymentHistoryEntry")).Return(nil)
				pushDispatcher.On("SendAsync", mock.Anything, userID, mock.AnythingOfType("*domain.PushMessage")).Return(nil)
			},
			wantErr: false,
		},
		{
			name: "payment failed — order in safe state (CANCELLED) — skip",
			mockFn: func(paymentGateway *mocks.PaymentGatewayMock, orderRepo *mocks.OrderRepositoryMock, productRepo *mocks.ProductRepositoryMock, paymentDispatcher *mocks.PaymentDispatcherMock, pushDispatcher *mocks.PushDispatcherMock) {
				paymentGateway.On("IsPaymentFailed", paymentID).Return(true, nil)
				orderRepo.On("GetByID", mock.Anything, orderID).Return(&domain.Order{
					ID: orderID, UserID: userID, Status: domain.OrderStatusCancelled,
				}, nil)
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			_, orderRepo, productRepo, paymentGateway, paymentDispatcher, pushDispatcher, svc := setupTest()
			tt.mockFn(paymentGateway, orderRepo, productRepo, paymentDispatcher, pushDispatcher)

			err := svc.AuditPayment(context.Background(), orderID, paymentID)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			paymentGateway.AssertExpectations(t)
			orderRepo.AssertExpectations(t)
			productRepo.AssertExpectations(t)
			paymentDispatcher.AssertExpectations(t)
			pushDispatcher.AssertExpectations(t)
		})
	}
}

func TestOrderService_Getters(t *testing.T) {
	t.Run("GetByID", func(t *testing.T) {
		t.Parallel()
		orderID := uuid.New()
		_, orderRepo, _, _, _, _, svc := setupTest()

		orderRepo.On("GetByID", mock.Anything, orderID).Return(&domain.Order{ID: orderID}, nil).Once()
		res, err := svc.GetByID(context.Background(), orderID)
		assert.NoError(t, err)
		assert.Equal(t, orderID, res.ID)

		orderRepo.On("GetByID", mock.Anything, orderID).Return(nil, errors.New("db error")).Once()
		res, err = svc.GetByID(context.Background(), orderID)
		assert.Error(t, err)
		assert.Nil(t, res)
	})

	t.Run("GetByUser", func(t *testing.T) {
		t.Parallel()
		userID := uuid.New()
		_, orderRepo, _, _, _, _, svc := setupTest()

		orderRepo.On("GetByUser", mock.Anything, userID).Return([]*domain.Order{{UserID: userID}}, nil).Once()
		res, err := svc.GetByUser(context.Background(), userID)
		assert.NoError(t, err)
		assert.Len(t, res, 1)

		orderRepo.On("GetByUser", mock.Anything, userID).Return(nil, errors.New("db error")).Once()
		res, err = svc.GetByUser(context.Background(), userID)
		assert.Error(t, err)
		assert.Nil(t, res)
	})

	t.Run("GetByRestaurant", func(t *testing.T) {
		t.Parallel()
		restID := uuid.New()
		_, orderRepo, _, _, _, _, svc := setupTest()

		orderRepo.On("GetByRestaurant", mock.Anything, restID).Return([]*domain.Order{{RestaurantID: restID}}, nil).Once()
		res, err := svc.GetByRestaurant(context.Background(), restID)
		assert.NoError(t, err)
		assert.Len(t, res, 1)

		orderRepo.On("GetByRestaurant", mock.Anything, restID).Return(nil, errors.New("db error")).Once()
		res, err = svc.GetByRestaurant(context.Background(), restID)
		assert.Error(t, err)
		assert.Nil(t, res)
	})

	t.Run("GetActiveByRestaurant", func(t *testing.T) {
		t.Parallel()
		restID := uuid.New()
		_, orderRepo, _, _, _, _, svc := setupTest()

		orderRepo.On("GetActiveByRestaurant", mock.Anything, restID).Return([]*domain.Order{{RestaurantID: restID}}, nil).Once()
		res, err := svc.GetActiveByRestaurant(context.Background(), restID)
		assert.NoError(t, err)
		assert.Len(t, res, 1)
	})

	t.Run("GetByPickupCode", func(t *testing.T) {
		t.Parallel()
		restID := uuid.New()
		code := "123456"
		_, orderRepo, _, _, _, _, svc := setupTest()

		orderRepo.On("GetByPickupCode", mock.Anything, restID, code).Return(&domain.Order{PickupCode: code}, nil).Once()
		res, err := svc.GetByPickupCode(context.Background(), restID, code)
		assert.NoError(t, err)
		assert.Equal(t, code, res.PickupCode)
	})
}

func TestOrderService_MarkAsAccepted(t *testing.T) {
	orderID := uuid.New()
	restID := uuid.New()
	userID := uuid.New()
	staff := &domain.RestaurantStaff{UserID: uuid.New(), RestaurantID: restID}

	tests := []struct {
		name    string
		mockFn  func(orderRepo *mocks.OrderRepositoryMock, pushDispatcher *mocks.PushDispatcherMock)
		wantErr bool
	}{
		{
			name: "permission denied (wrong restaurant)",
			mockFn: func(orderRepo *mocks.OrderRepositoryMock, pushDispatcher *mocks.PushDispatcherMock) {
				orderRepo.On("GetByID", mock.Anything, orderID).Return(&domain.Order{RestaurantID: uuid.New()}, nil)
			},
			wantErr: true,
		},
		{
			name: "invalid status transition (from PENDING implies fail)",
			mockFn: func(orderRepo *mocks.OrderRepositoryMock, pushDispatcher *mocks.PushDispatcherMock) {
				orderRepo.On("GetByID", mock.Anything, orderID).Return(&domain.Order{
					ID: orderID, RestaurantID: restID, Status: domain.OrderStatusPending,
				}, nil)
			},
			wantErr: true,
		},
		{
			name: "success",
			mockFn: func(orderRepo *mocks.OrderRepositoryMock, pushDispatcher *mocks.PushDispatcherMock) {
				orderRepo.On("GetByID", mock.Anything, orderID).Return(&domain.Order{
					ID: orderID, UserID: userID, RestaurantID: restID, Status: domain.OrderStatusPaid,
				}, nil)
				orderRepo.On("UpdateStatus", mock.Anything, orderID, domain.OrderStatusAccepted).Return(nil)
				orderRepo.On("AddStatusHistory", mock.Anything, mock.AnythingOfType("*domain.OrderStatusHistoryEntry")).Return(nil)
				pushDispatcher.On("SendAsync", mock.Anything, userID, mock.AnythingOfType("*domain.PushMessage")).Return(nil)
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			_, orderRepo, _, _, _, pushDispatcher, svc := setupTest()
			tt.mockFn(orderRepo, pushDispatcher)

			err := svc.MarkAsAccepted(context.Background(), staff, orderID)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			orderRepo.AssertExpectations(t)
			pushDispatcher.AssertExpectations(t)
		})
	}
}

func TestOrderService_MarkAsReady(t *testing.T) {
	orderID := uuid.New()
	restID := uuid.New()
	userID := uuid.New()
	staff := &domain.RestaurantStaff{UserID: uuid.New(), RestaurantID: restID}

	tests := []struct {
		name    string
		mockFn  func(orderRepo *mocks.OrderRepositoryMock, pushDispatcher *mocks.PushDispatcherMock)
		wantErr bool
	}{
		{
			name: "order not found",
			mockFn: func(orderRepo *mocks.OrderRepositoryMock, pushDispatcher *mocks.PushDispatcherMock) {
				orderRepo.On("GetByID", mock.Anything, orderID).Return(nil, domain.ErrNotFound)
			},
			wantErr: true,
		},
		{
			name: "permission denied (wrong restaurant)",
			mockFn: func(orderRepo *mocks.OrderRepositoryMock, pushDispatcher *mocks.PushDispatcherMock) {
				orderRepo.On("GetByID", mock.Anything, orderID).Return(&domain.Order{RestaurantID: uuid.New()}, nil)
			},
			wantErr: true,
		},
		{
			name: "invalid status transition",
			mockFn: func(orderRepo *mocks.OrderRepositoryMock, pushDispatcher *mocks.PushDispatcherMock) {
				orderRepo.On("GetByID", mock.Anything, orderID).Return(&domain.Order{
					ID: orderID, RestaurantID: restID, Status: domain.OrderStatusPending,
				}, nil)
			},
			wantErr: true,
		},
		{
			name: "success",
			mockFn: func(orderRepo *mocks.OrderRepositoryMock, pushDispatcher *mocks.PushDispatcherMock) {
				orderRepo.On("GetByID", mock.Anything, orderID).Return(&domain.Order{
					ID: orderID, UserID: userID, RestaurantID: restID, Status: domain.OrderStatusAccepted,
				}, nil)
				orderRepo.On("UpdateStatus", mock.Anything, orderID, domain.OrderStatusReady).Return(nil)
				orderRepo.On("AddStatusHistory", mock.Anything, mock.AnythingOfType("*domain.OrderStatusHistoryEntry")).Return(nil)
				pushDispatcher.On("SendAsync", mock.Anything, userID, mock.AnythingOfType("*domain.PushMessage")).Return(nil)
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			_, orderRepo, _, _, _, pushDispatcher, svc := setupTest()
			tt.mockFn(orderRepo, pushDispatcher)

			err := svc.MarkAsReady(context.Background(), staff, orderID)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			orderRepo.AssertExpectations(t)
			pushDispatcher.AssertExpectations(t)
		})
	}
}

func TestOrderService_MarkAsCollected(t *testing.T) {
	orderID := uuid.New()
	restID := uuid.New()
	userID := uuid.New()
	staff := &domain.RestaurantStaff{UserID: uuid.New(), RestaurantID: restID}

	tests := []struct {
		name    string
		mockFn  func(orderRepo *mocks.OrderRepositoryMock, pushDispatcher *mocks.PushDispatcherMock)
		wantErr bool
	}{
		{
			name: "permission denied (wrong restaurant)",
			mockFn: func(orderRepo *mocks.OrderRepositoryMock, pushDispatcher *mocks.PushDispatcherMock) {
				orderRepo.On("GetByID", mock.Anything, orderID).Return(&domain.Order{RestaurantID: uuid.New()}, nil)
			},
			wantErr: true,
		},
		{
			name: "success",
			mockFn: func(orderRepo *mocks.OrderRepositoryMock, pushDispatcher *mocks.PushDispatcherMock) {
				orderRepo.On("GetByID", mock.Anything, orderID).Return(&domain.Order{
					ID: orderID, UserID: userID, RestaurantID: restID, Status: domain.OrderStatusReady,
				}, nil)
				orderRepo.On("UpdateStatus", mock.Anything, orderID, domain.OrderStatusCollected).Return(nil)
				orderRepo.On("AddStatusHistory", mock.Anything, mock.AnythingOfType("*domain.OrderStatusHistoryEntry")).Return(nil)
				pushDispatcher.On("SendAsync", mock.Anything, userID, mock.AnythingOfType("*domain.PushMessage")).Return(nil)
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			_, orderRepo, _, _, _, pushDispatcher, svc := setupTest()
			tt.mockFn(orderRepo, pushDispatcher)

			err := svc.MarkAsCollected(context.Background(), staff, orderID, "", true)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			orderRepo.AssertExpectations(t)
			pushDispatcher.AssertExpectations(t)
		})
	}
}

func TestOrderService_GetStatusHistoryByUser(t *testing.T) {
	orderID := uuid.New()
	userID := uuid.New()

	tests := []struct {
		name    string
		mockFn  func(orderRepo *mocks.OrderRepositoryMock)
		wantErr bool
	}{
		{
			name: "wrong user (permission denied)",
			mockFn: func(orderRepo *mocks.OrderRepositoryMock) {
				orderRepo.On("GetByID", mock.Anything, orderID).Return(&domain.Order{UserID: uuid.New()}, nil)
			},
			wantErr: true,
		},
		{
			name: "success",
			mockFn: func(orderRepo *mocks.OrderRepositoryMock) {
				orderRepo.On("GetByID", mock.Anything, orderID).Return(&domain.Order{UserID: userID}, nil)
				orderRepo.On("GetStatusHistory", mock.Anything, orderID).Return([]*domain.OrderStatusHistoryEntry{{OrderID: orderID}}, nil)
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			_, orderRepo, _, _, _, _, svc := setupTest()
			tt.mockFn(orderRepo)

			res, err := svc.GetStatusHistoryByUser(context.Background(), userID, orderID)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, res)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, res)
			}
			orderRepo.AssertExpectations(t)
		})
	}
}

func TestOrderService_GetStatusHistoryByRestaurant(t *testing.T) {
	orderID := uuid.New()
	restID := uuid.New()
	staff := &domain.RestaurantStaff{UserID: uuid.New(), RestaurantID: restID}

	tests := []struct {
		name    string
		mockFn  func(orderRepo *mocks.OrderRepositoryMock)
		wantErr bool
	}{
		{
			name: "permission denied",
			mockFn: func(orderRepo *mocks.OrderRepositoryMock) {
				orderRepo.On("GetByID", mock.Anything, orderID).Return(&domain.Order{RestaurantID: uuid.New()}, nil)
			},
			wantErr: true,
		},
		{
			name: "success",
			mockFn: func(orderRepo *mocks.OrderRepositoryMock) {
				orderRepo.On("GetByID", mock.Anything, orderID).Return(&domain.Order{RestaurantID: restID}, nil)
				orderRepo.On("GetStatusHistory", mock.Anything, orderID).Return([]*domain.OrderStatusHistoryEntry{{OrderID: orderID}}, nil)
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			_, orderRepo, _, _, _, _, svc := setupTest()
			tt.mockFn(orderRepo)

			res, err := svc.GetStatusHistoryByRestaurant(context.Background(), staff, orderID)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, res)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, res)
			}
			orderRepo.AssertExpectations(t)
		})
	}
}

func TestOrderService_CancelByUser(t *testing.T) {
	orderID := uuid.New()
	userID := uuid.New()
	paymentID := "pi_123"
	reason := domain.OrderReason("customer request")

	tests := []struct {
		name    string
		mockFn  func(orderRepo *mocks.OrderRepositoryMock, productRepo *mocks.ProductRepositoryMock, paymentDispatcher *mocks.PaymentDispatcherMock)
		wantErr bool
	}{
		{
			name: "permission denied",
			mockFn: func(orderRepo *mocks.OrderRepositoryMock, productRepo *mocks.ProductRepositoryMock, paymentDispatcher *mocks.PaymentDispatcherMock) {
				orderRepo.On("GetByID", mock.Anything, orderID).Return(&domain.Order{UserID: uuid.New()}, nil)
			},
			wantErr: true,
		},
		{
			name: "success cancel PENDING order",
			mockFn: func(orderRepo *mocks.OrderRepositoryMock, productRepo *mocks.ProductRepositoryMock, paymentDispatcher *mocks.PaymentDispatcherMock) {
				orderRepo.On("GetByID", mock.Anything, orderID).Return(&domain.Order{
					ID: orderID, UserID: userID, Status: domain.OrderStatusPending, PaymentID: &paymentID,
				}, nil)
				orderRepo.On("UpdateStatus", mock.Anything, orderID, domain.OrderStatusCancelled).Return(nil)
				orderRepo.On("AddStatusHistory", mock.Anything, mock.AnythingOfType("*domain.OrderStatusHistoryEntry")).Return(nil)
				paymentDispatcher.On("CancelAsync", mock.Anything, paymentID, orderID).Return(nil)
			},
			wantErr: false,
		},
		{
			name: "success cancel PAID order (refund and restore stock)",
			mockFn: func(orderRepo *mocks.OrderRepositoryMock, productRepo *mocks.ProductRepositoryMock, paymentDispatcher *mocks.PaymentDispatcherMock) {
				prodID := uuid.New()
				orderRepo.On("GetByID", mock.Anything, orderID).Return(&domain.Order{
					ID: orderID, UserID: userID, Status: domain.OrderStatusPaid, PaymentID: &paymentID,
					Items: []domain.OrderItem{{ProductID: prodID, Quantity: 2}},
				}, nil)
				productRepo.On("GetByIDForUpdate", mock.Anything, prodID).Return(&domain.Product{QuantityAvailable: 10}, nil)
				productRepo.On("UpdateStock", mock.Anything, prodID, 12).Return(nil)
				orderRepo.On("UpdateStatus", mock.Anything, orderID, domain.OrderStatusCancelled).Return(nil)
				orderRepo.On("AddStatusHistory", mock.Anything, mock.AnythingOfType("*domain.OrderStatusHistoryEntry")).Return(nil)
				paymentDispatcher.On("RefundAsync", mock.Anything, paymentID, orderID).Return(nil)
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			_, orderRepo, productRepo, _, paymentDispatcher, _, svc := setupTest()
			tt.mockFn(orderRepo, productRepo, paymentDispatcher)

			err := svc.CancelByUser(context.Background(), userID, orderID, &reason)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			orderRepo.AssertExpectations(t)
			productRepo.AssertExpectations(t)
			paymentDispatcher.AssertExpectations(t)
		})
	}
}

func TestOrderService_CancelByRestaurant(t *testing.T) {
	orderID := uuid.New()
	restID := uuid.New()
	userID := uuid.New()
	staff := &domain.RestaurantStaff{UserID: uuid.New(), RestaurantID: restID}
	reason := domain.OrderReason("restaurant requested")
	paymentID := "pi_123"

	tests := []struct {
		name    string
		mockFn  func(orderRepo *mocks.OrderRepositoryMock, productRepo *mocks.ProductRepositoryMock, paymentDispatcher *mocks.PaymentDispatcherMock, pushDispatcher *mocks.PushDispatcherMock)
		wantErr bool
	}{
		{
			name: "permission denied",
			mockFn: func(orderRepo *mocks.OrderRepositoryMock, productRepo *mocks.ProductRepositoryMock, paymentDispatcher *mocks.PaymentDispatcherMock, pushDispatcher *mocks.PushDispatcherMock) {
				orderRepo.On("GetByID", mock.Anything, orderID).Return(&domain.Order{RestaurantID: uuid.New()}, nil)
			},
			wantErr: true,
		},
		{
			name: "success cancel PAID order (refund and restore stock)",
			mockFn: func(orderRepo *mocks.OrderRepositoryMock, productRepo *mocks.ProductRepositoryMock, paymentDispatcher *mocks.PaymentDispatcherMock, pushDispatcher *mocks.PushDispatcherMock) {
				prodID := uuid.New()
				orderRepo.On("GetByID", mock.Anything, orderID).Return(&domain.Order{
					ID: orderID, UserID: userID, RestaurantID: restID,
					Status: domain.OrderStatusPaid, PaymentID: &paymentID,
					Items: []domain.OrderItem{{ProductID: prodID, Quantity: 1}},
				}, nil)
				productRepo.On("GetByIDForUpdate", mock.Anything, prodID).Return(&domain.Product{QuantityAvailable: 1}, nil)
				productRepo.On("UpdateStock", mock.Anything, prodID, 2).Return(nil)
				orderRepo.On("UpdateStatus", mock.Anything, orderID, domain.OrderStatusCancelled).Return(nil)
				orderRepo.On("AddStatusHistory", mock.Anything, mock.AnythingOfType("*domain.OrderStatusHistoryEntry")).Return(nil)
				paymentDispatcher.On("RefundAsync", mock.Anything, paymentID, orderID).Return(nil)
				pushDispatcher.On("SendAsync", mock.Anything, userID, mock.AnythingOfType("*domain.PushMessage")).Return(nil)
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			_, orderRepo, productRepo, _, paymentDispatcher, pushDispatcher, svc := setupTest()
			tt.mockFn(orderRepo, productRepo, paymentDispatcher, pushDispatcher)

			err := svc.CancelByRestaurant(context.Background(), staff, orderID, &reason)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			orderRepo.AssertExpectations(t)
			productRepo.AssertExpectations(t)
			paymentDispatcher.AssertExpectations(t)
			pushDispatcher.AssertExpectations(t)
		})
	}
}

func TestOrderService_CancelExpiredPendingOrders(t *testing.T) {
	tests := []struct {
		name    string
		mockFn  func(orderRepo *mocks.OrderRepositoryMock, paymentDispatcher *mocks.PaymentDispatcherMock)
		wantCnt int
		wantErr bool
	}{
		{
			name: "no expired orders",
			mockFn: func(orderRepo *mocks.OrderRepositoryMock, paymentDispatcher *mocks.PaymentDispatcherMock) {
				orderRepo.On("GetExpiredPendingOrders", mock.Anything, mock.AnythingOfType("time.Time")).Return([]*domain.Order{}, nil)
			},
			wantCnt: 0,
			wantErr: false,
		},
		{
			name: "cancel multiple expired orders successfully",
			mockFn: func(orderRepo *mocks.OrderRepositoryMock, paymentDispatcher *mocks.PaymentDispatcherMock) {
				pID1, pID2 := "pi_1", "pi_2"
				id1, id2 := uuid.New(), uuid.New()
				orders := []*domain.Order{
					{ID: id1, UserID: uuid.New(), Status: domain.OrderStatusPending, PaymentID: &pID1},
					{ID: id2, UserID: uuid.New(), Status: domain.OrderStatusPending, PaymentID: &pID2},
				}
				orderRepo.On("GetExpiredPendingOrders", mock.Anything, mock.AnythingOfType("time.Time")).Return(orders, nil)
				orderRepo.On("UpdateStatus", mock.Anything, mock.AnythingOfType("uuid.UUID"), domain.OrderStatusCancelled).Return(nil).Twice()
				orderRepo.On("AddStatusHistory", mock.Anything, mock.AnythingOfType("*domain.OrderStatusHistoryEntry")).Return(nil).Twice()
				paymentDispatcher.On("CancelAsync", mock.Anything, "pi_1", id1).Return(nil).Once()
				paymentDispatcher.On("CancelAsync", mock.Anything, "pi_2", id2).Return(nil).Once()
			},
			wantCnt: 2,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			_, orderRepo, _, _, paymentDispatcher, _, svc := setupTest()
			tt.mockFn(orderRepo, paymentDispatcher)

			cnt, err := svc.CancelExpiredPendingOrders(context.Background(), 15*time.Minute)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.wantCnt, cnt)
			}
			orderRepo.AssertExpectations(t)
			paymentDispatcher.AssertExpectations(t)
		})
	}
}

func TestOrderService_GetPaymentHistoryByOrder(t *testing.T) {
	orderID := uuid.New()
	userID := uuid.New()

	tests := []struct {
		name    string
		mockFn  func(orderRepo *mocks.OrderRepositoryMock)
		wantErr bool
	}{
		{
			name: "order not found",
			mockFn: func(orderRepo *mocks.OrderRepositoryMock) {
				orderRepo.On("GetByID", mock.Anything, orderID).Return(nil, domain.ErrNotFound)
			},
			wantErr: true,
		},
		{
			name: "permission denied (wrong user)",
			mockFn: func(orderRepo *mocks.OrderRepositoryMock) {
				orderRepo.On("GetByID", mock.Anything, orderID).Return(&domain.Order{UserID: uuid.New()}, nil)
			},
			wantErr: true,
		},
		{
			name: "success",
			mockFn: func(orderRepo *mocks.OrderRepositoryMock) {
				orderRepo.On("GetByID", mock.Anything, orderID).Return(&domain.Order{UserID: userID}, nil)
				orderRepo.On("GetPaymentHistory", mock.Anything, orderID).Return([]*domain.PaymentHistoryEntry{{OrderID: orderID}}, nil)
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			_, orderRepo, _, _, _, _, svc := setupTest()
			tt.mockFn(orderRepo)

			res, err := svc.GetPaymentHistoryByOrder(context.Background(), userID, orderID)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, res)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, res)
			}
			orderRepo.AssertExpectations(t)
		})
	}
}

func TestOrderService_GetPaymentHistoryByUser(t *testing.T) {
	userID := uuid.New()

	tests := []struct {
		name    string
		mockFn  func(orderRepo *mocks.OrderRepositoryMock)
		wantErr bool
	}{
		{
			name: "db error",
			mockFn: func(orderRepo *mocks.OrderRepositoryMock) {
				orderRepo.On("GetPaymentHistoryByUser", mock.Anything, userID).Return(nil, errors.New("db error"))
			},
			wantErr: true,
		},
		{
			name: "success",
			mockFn: func(orderRepo *mocks.OrderRepositoryMock) {
				orderRepo.On("GetPaymentHistoryByUser", mock.Anything, userID).Return([]*domain.PaymentHistoryEntry{{UserID: userID}}, nil)
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			_, orderRepo, _, _, _, _, svc := setupTest()
			tt.mockFn(orderRepo)

			res, err := svc.GetPaymentHistoryByUser(context.Background(), userID)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, res)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, res)
			}
			orderRepo.AssertExpectations(t)
		})
	}
}

func TestOrderService_RecordPaymentEvent(t *testing.T) {
	orderID := uuid.New()
	userID := uuid.New()
	paymentID := "pi_123"

	tests := []struct {
		name    string
		mockFn  func(orderRepo *mocks.OrderRepositoryMock)
		wantErr bool
	}{
		{
			name: "order not found",
			mockFn: func(orderRepo *mocks.OrderRepositoryMock) {
				orderRepo.On("GetByID", mock.Anything, orderID).Return(nil, domain.ErrNotFound)
			},
			wantErr: true,
		},
		{
			name: "success",
			mockFn: func(orderRepo *mocks.OrderRepositoryMock) {
				orderRepo.On("GetByID", mock.Anything, orderID).Return(&domain.Order{
					ID: orderID, UserID: userID, TotalAmount: 1500,
				}, nil)
				orderRepo.On("AddPaymentHistory", mock.Anything, mock.AnythingOfType("*domain.PaymentHistoryEntry")).Return(nil)
			},
			wantErr: false,
		},
		{
			name: "success with error message",
			mockFn: func(orderRepo *mocks.OrderRepositoryMock) {
				orderRepo.On("GetByID", mock.Anything, orderID).Return(&domain.Order{
					ID: orderID, UserID: userID, TotalAmount: 1500,
				}, nil)
				orderRepo.On("AddPaymentHistory", mock.Anything, mock.AnythingOfType("*domain.PaymentHistoryEntry")).Return(nil)
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			_, orderRepo, _, _, _, _, svc := setupTest()
			tt.mockFn(orderRepo)

			errMsg := "card_declined"
			var errMsgPtr *string
			if tt.name == "success with error message" {
				errMsgPtr = &errMsg
			}

			err := svc.RecordPaymentEvent(context.Background(), orderID, paymentID, domain.PaymentEventSucceeded, errMsgPtr)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			orderRepo.AssertExpectations(t)
		})
	}
}
