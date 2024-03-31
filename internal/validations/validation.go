package validations

type FormValidator struct {
	FieldErrors map[string]string
	Flash       string
}

func (fv *FormValidator) Valid() bool {
	return len(fv.FieldErrors) == 0
}

func (fv *FormValidator) AddFieldError(field, message string) {
	if fv.FieldErrors == nil {
		fv.FieldErrors = make(map[string]string)
	}
	fv.FieldErrors[field] = message
}
