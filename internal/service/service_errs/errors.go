package service_errs

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
	ErrSigningMethod   = errors.New("unexpected signing method")
	ErrEmptyClaims     = errors.New("token claims doesnt contain required fields")
	ErrorInvalidJWT    = errors.New("invalid jwt")
	ErrInvalidClaims   = errors.New("invalid token claims")
	ErrBindingJSON     = errors.New("error binding json body")
	ErrInvalidProduct  = errors.New("invalid product data")
	ErrAddingProduct   = errors.New("error adding product")
)
