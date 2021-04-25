package sendgrid

type ErrorResponse struct {
	Errors Errors `json:"errors"`
}

type Error struct {
	Message string `json:"message"`
	Field   string `json:"field"`
	Help    string `json:"help"`
}

type Errors []Error

func (e Error) Error() string {
	return e.Message
}

func (e Errors) Error() string {
	if len(e) > 0 {
		return e[0].Error()
	}
	return ""
}
