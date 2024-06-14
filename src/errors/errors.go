package errors

import (
	i18n_err "loverly/lib/i18n/errors"
)

var (
	// Auth
	ErrEmailUnregistered      = i18n_err.NewI18nError("err_email_unregistered")
	ErrInvalidEmailOrPassword = i18n_err.NewI18nError("err_invalid_email_or_password")
	ErrPasswordNotMatch       = i18n_err.NewI18nError("err_password_not_match")
	ErrInvalidEmailFormat     = i18n_err.NewI18nError("err_invalid_email_format")
	ErrInvalidUserId          = i18n_err.NewI18nError("err_invalid_user_id")
)
