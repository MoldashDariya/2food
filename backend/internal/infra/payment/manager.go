package payment

type Provider interface {
	CreatePaymentIntent(amount float64, multiplier int, currency string, metadata map[string]string) (paymentID string, clientSecret string, err error)
	ConfirmPayment(paymentID string) error
	CancelPaymentIntent(paymentID string) error
	RefundPayment(paymentID string) error
	VerifyWebhookSignature(payload []byte, signature string) (paymentID string, orderID string, err error)
	IsPaymentFailed(paymentID string) (bool, error)
}

type Manager struct {
	cfg      Config
	provider Provider
}

func NewManager(cfg Config, provider Provider) *Manager {
	return &Manager{cfg: cfg, provider: provider}
}

func (m *Manager) CreatePaymentIntent(amount float64, metadata map[string]string) (string, string, error) {
	return m.provider.CreatePaymentIntent(amount, m.cfg.AmountMultiplier, m.cfg.Currency, metadata)
}

func (m *Manager) ConfirmPayment(paymentID string) error {
	return m.provider.ConfirmPayment(paymentID)
}

func (m *Manager) CancelPaymentIntent(paymentID string) error {
	return m.provider.CancelPaymentIntent(paymentID)
}

func (m *Manager) RefundPayment(paymentID string) error {
	return m.provider.RefundPayment(paymentID)
}

func (m *Manager) VerifyWebhookSignature(payload []byte, signature string) (string, string, error) {
	return m.provider.VerifyWebhookSignature(payload, signature)
}

func (m *Manager) IsPaymentFailed(paymentID string) (bool, error) {
	return m.provider.IsPaymentFailed(paymentID)
}
