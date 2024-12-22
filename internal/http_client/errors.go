package http_client

import "errors"

var (
	ErrRetryExceededMaximum = errors.New("maximum retry iterations exceeded")
)
