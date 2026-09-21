package payment

type Config struct {
	Currency         string `envconfig:"PAYMENT_CURRENCY" required:"true"`
	AmountMultiplier int    `envconfig:"PAYMENT_AMOUNT_MULTIPLIER" required:"true"`
}
