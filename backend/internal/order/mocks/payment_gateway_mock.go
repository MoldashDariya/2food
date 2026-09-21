package mocks

import (
	"github.com/stretchr/testify/mock"
)

type PaymentGatewayMock struct {
	mock.Mock
}

func (m *PaymentGatewayMock) CreatePaymentIntent(amount float64, metadata map[string]string) (string, string, error) {
	args := m.Called(amount, metadata)
	return args.String(0), args.String(1), args.Error(2)
}

func (m *PaymentGatewayMock) VerifyWebhookSignature(payload []byte, signature string) (string, string, error) {
	args := m.Called(payload, signature)
	return args.String(0), args.String(1), args.Error(2)
}

func (m *PaymentGatewayMock) IsPaymentFailed(paymentID string) (bool, error) {
	args := m.Called(paymentID)
	return args.Bool(0), args.Error(1)
}
