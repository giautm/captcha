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

func ReportResult(eng CaptchaResolver, result *CaptchaResult) func(ctx context.Context, err error) {
	return func(ctx context.Context, err error) {
		if e, ok := (eng).(ResultReporter); ok {
			if err == nil {
				e.Report(ctx, result, true)
			} else if err == ErrCaptchaInvalid {
				e.Report(ctx, result, false)
			}
		}
	}
}
