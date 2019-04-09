package validation

type Validation struct {
	Message string
	Payload interface{}
}

func New(message string, payload interface{}) *Validation {
	return &Validation{Message: message, Payload: payload}
}
