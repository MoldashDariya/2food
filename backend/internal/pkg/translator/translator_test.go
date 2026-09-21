package translator

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTranslator_New(t *testing.T) {
	assert.NotPanics(t, func() {
		tr := New()
		assert.NotNil(t, tr)
	})
}

func TestTranslator_T(t *testing.T) {
	tr := New()

	tests := []struct {
		name      string
		lang      string
		messageID string
		data      map[string]any
		expected  string
	}{
		{
			name:      "English OrderStatusPaid_Title",
			lang:      "en",
			messageID: "OrderStatusPaid_Title",
			expected:  "Order Paid",
		},
		{
			name:      "Russian OrderStatusPaid_Title",
			lang:      "ru",
			messageID: "OrderStatusPaid_Title",
			expected:  "Заказ оплачен",
		},
		{
			name:      "Kazakh OrderStatusPaid_Title",
			lang:      "kk",
			messageID: "OrderStatusPaid_Title",
			expected:  "Тапсырыс төленді",
		},
		{
			name:      "Chinese OrderStatusPaid_Title",
			lang:      "zh",
			messageID: "OrderStatusPaid_Title",
			expected:  "订单已支付",
		},
		{
			name:      "English City Almaty",
			lang:      "en",
			messageID: "City_almaty",
			expected:  "Almaty",
		},
		{
			name:      "Chinese City Almaty",
			lang:      "zh",
			messageID: "City_almaty",
			expected:  "阿拉木图",
		},
		{
			name:      "Template data replacement - EN",
			lang:      "en",
			messageID: "SMS_OTP_Text",
			data:      map[string]any{"code": "123456"},
			expected:  "Your verification code: 123456",
		},
		{
			name:      "Template data replacement - ZH",
			lang:      "zh",
			messageID: "SMS_OTP_Text",
			data:      map[string]any{"code": "123456"},
			expected:  "您的验证码为：123456",
		},
		{
			name:      "Non-existent message returns message ID",
			lang:      "en",
			messageID: "NonExistentKey",
			expected:  "NonExistentKey",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := tr.T(tt.lang, tt.messageID, tt.data)
			assert.Equal(t, tt.expected, actual)
		})
	}
}

func TestTranslator_TFallback(t *testing.T) {
	tr := New()

	tests := []struct {
		name      string
		lang      string
		messageID string
		fallback  string
		expected  string
	}{
		{
			name:      "Existent key returns translation",
			lang:      "zh",
			messageID: "OrderStatusPaid_Title",
			fallback:  "Fallback Title",
			expected:  "订单已支付",
		},
		{
			name:      "Non-existent key returns fallback",
			lang:      "zh",
			messageID: "NonExistentKey",
			fallback:  "Fallback Title",
			expected:  "Fallback Title",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := tr.TFallback(tt.lang, tt.messageID, tt.fallback, nil)
			assert.Equal(t, tt.expected, actual)
		})
	}
}
