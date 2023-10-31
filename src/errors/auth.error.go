package errors

import "errors"

var (
	ErrMissingActivationCodeInQueryString = errors.New("missing activation code in query string")
	ErrMissingBearerToken                 = errors.New("missing bearer token in auth header")
)
