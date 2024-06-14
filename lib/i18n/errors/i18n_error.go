package errors

// Error that support internationalization. The error is store a key instead of message.
// Use errros.ErrorMessage() to get localization error message.
type I18nError interface {
	error
}

type i18nError struct {
	key string
}

func NewI18nError(key string) I18nError {
	return &i18nError{
		key: key,
	}
}

func (e *i18nError) Error() string {
	return e.key
}
