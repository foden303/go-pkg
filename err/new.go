package err

type impl struct {
	err      error
	httpCode int

	ErrorCode    string `json:"error_code"`
	ErrorMessage string `json:"error_message"`
}

func New(code, message string) Error {
	return &impl{
		ErrorCode:    code,
		ErrorMessage: message,
	}
}
