package repositories

type RepositoryError struct {
	error
}

func newRepositoryError(err error) error {
	return &RepositoryError{error: err}
}
