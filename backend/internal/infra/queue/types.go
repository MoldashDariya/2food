package queue

const (
	TypePushSend           string = "notification:push_send"
	TypeOrderCancelExpired string = "cron:order:cancel_expired"
	TypeRoleCleanup        string = "cron:role:cleanup"

	TypePaymentConfirm string = "payment:confirm"
	TypePaymentCancel  string = "payment:cancel"
	TypePaymentRefund  string = "payment:refund"
	TypePaymentAudit   string = "payment:audit"
)
