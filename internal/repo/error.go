package repo

import "errors"

var (
	ErrDatabaseError = errors.New("database error")
	ErrAddJTI        = errors.New("failed to add jti to auth record")
)
