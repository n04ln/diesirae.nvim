package util

import "errors"

var (
	ErrCouldNotSetConfig = errors.New("could not set config")
	ErrResponseIsNotOK   = errors.New("response is not ok")
	ErrCookieIsNotFound  = errors.New("cookie is not found")
	ErrInvalidCookie     = errors.New("invalid cookie")
	ErrInvalidArgs       = errors.New("invalid args")
	ErrInvalidJudgeId    = errors.New("invalid judge_id")
)
