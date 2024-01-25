package apperror

type StatusError struct {
	error
	status int
}

func (se StatusError) StatusCode() int {
	return se.status
}

func WithStatus(err error, status int) error {
	return StatusError{
		error:  err,
		status: status,
	}
}
