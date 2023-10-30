package errors

import "errors"

var (
	ErrMissingActivationCodeInQueryString = errors.New("missing activation code in query string")
)
