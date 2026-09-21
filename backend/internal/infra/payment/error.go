package payment

type Error struct {
	Code    string // payment provider codes
	Message string
}

func (e *Error) Error() string {
	return e.Code + ": " + e.Message
}
