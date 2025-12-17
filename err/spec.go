package err

type Error interface {
	HttpCode() int
	Error() string
	Code() string
	Message() string
	Err() error
}
