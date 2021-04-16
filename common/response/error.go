package response

import "errors"

var (
	ErrRecordNotFound = errors.New("record not found")
	ErrInternalBad    = errors.New("internal error")
)
