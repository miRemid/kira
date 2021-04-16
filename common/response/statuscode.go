package response

type StatusCode int

const (
	_ StatusCode = iota + 1000

	StatusOK
	StatusInternalError
	StatusUnauthorized
	StatusExpired
	StatusBadParams
	StatusForbidden
	StatusPingError
	StatusRedisCheck
	StatusNeedToken
	StatusUserSuspend
	StatusAlreadyLike
)
