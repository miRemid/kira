package response

type StatusCode int

const (
	_ StatusCode = iota + 1000

	StatusOK
	StatusInternalError
	StatusUnauthorized
	StatusBadParams
)
