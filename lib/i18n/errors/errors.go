package errors

var (
	ErrBadRequest     = NewI18nError("err_bad_request")
	ErrInternalServer = NewI18nError("err_internal_server")
	ErrUnauthorized   = NewI18nError("err_unauthorized")
)
