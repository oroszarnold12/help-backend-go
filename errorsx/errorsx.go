package errorsx

type StatusCoder interface {
	error
	StatusCode() int
}
