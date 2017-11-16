package util

import "errors"

var (
	ErrCouldNotSetConfig = errors.New("could not set config")
	ErrResponseIsNotOK   = errors.New("response is not ok")
	ErrCookieIsNotFound  = errors.New("cookie is not found")
	ErrInvalidArgs       = errors.New("invalid args")
)
