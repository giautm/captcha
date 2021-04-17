package engine

import (
	"context"
	"errors"
)

type ResultReporter interface {
	Report(ctx context.Context, result *CaptchaResult, correct bool) error
}

var (
	ErrCaptchaInvalid = errors.New("captcha is invalid")
)
