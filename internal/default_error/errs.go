package default_error

import "errors"

var (
	ErrTokenNotFound   = errors.New("auth token not found")
	ErrTokenExpired    = errors.New("auth token expired")
	ErrInvalidToken    = errors.New("invalid auth token")
	ErrLogoutFault     = errors.New("logout failed")
	ErrUnauthorized    = errors.New("unauthorized")
	ErrInvalidUserData = errors.New("invalid user data")
	ErrUpdatingAvatar  = errors.New("error updating avatar")
	ErrUserNotFound    = errors.New("user not found")
)

const (
	ErrInvalidEnv      = "is not set in .env"
	ErrDatabaseError   = "database error"
	ErrSigningJWT      = "failed to sign jwt"
	ErrAddJTI          = "failed to add jti to auth record"
	ErrSigningMethod   = "unexpected signing method"
	ErrParseJWT        = "failed  to parse jwt"
	ErrInvalidIDType   = "userID is not of type float64"
	ErrEmptyClaims     = "token claims doesnt contain required fields"
	ErrorInvalidJWT    = "invalid jwt"
	ErrParsingJWT      = "failed to parse jwt"
	ErrInvalidClaims   = "invalid token claims"
	ErrStartingBot     = "Error starting Telegram bot"
	ErrSendingMsg      = "Error sending message"
	ErrGeneratingLink  = "Error generating auth link"
	ErrDownloadingPic  = "Error downloading picture"
	ErrUploadingPic    = "error uploading picture to imgur"
	ErrInvalidPic      = "Invalid picture"
	ExErrTokenNotFound = "auth token record not found"
	WarnInvalidLink    = "is trying to use invalid link"
)
