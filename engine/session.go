package engine

import (
	"context"
	"errors"
	"net/http"
	"net/http/cookiejar"
	"time"

	"giautm.dev/captcha/srcimg"
)

var (
	ErrLimitExceeded = errors.New("session: retry limit exceeded")
)

type IDGenerator interface {
	NewID() string
}

type CaptchaSession struct {
	CaptchaURL  string
	Engine      CaptchaResolver
	IDGenerator IDGenerator
	RetryCount  int
	Transport   http.RoundTripper
}

type NoopIDGen struct{}

func (s NoopIDGen) NewID() string {
	return ""
}

func NewCaptchaSession(captchaURL string, eng CaptchaResolver) *CaptchaSession {
	return &CaptchaSession{
		CaptchaURL:  captchaURL,
		Engine:      eng,
		IDGenerator: &NoopIDGen{},
		RetryCount:  5,
	}
}

func (h *CaptchaSession) Fetch(ctx context.Context, fn func(c *http.Client, captcha string) error) error {
	jar, err := cookiejar.New(nil)
	if err != nil {
		return err
	}

	client := &http.Client{
		Jar:       jar,
		Timeout:   time.Second * 10,
		Transport: h.Transport,
	}
	for i := 0; i < h.RetryCount; i++ {
		err = func(ctx context.Context) error {
			file, err := srcimg.DownloadReader(ctx, client, h.CaptchaURL)
			if err != nil {
				return err
			}
			defer file.Close()

			result, err := h.Engine.ResolveFile(ctx, file)
			if err != nil {
				return err
			}

			err = fn(client, result.Captcha)

			if e, ok := h.Engine.(ResultReporter); ok {
				if err == nil {
					e.Report(ctx, result, true)
				} else if errors.Is(err, ErrCaptchaInvalid) {
					e.Report(ctx, result, false)
				}
			}

			return err
		}(WithSessionID(ctx, h.IDGenerator))
		if err == ErrCaptchaInvalid {
			continue
		}

		return err
	}

	return ErrLimitExceeded
}

type sessionID struct{}

var sessionIDKey sessionID

func SessionIDFromContext(ctx context.Context) string {
	if v := ctx.Value(sessionIDKey); v != nil {
		if id, ok := v.(string); ok {
			return id
		}
	}

	return ""
}

func WithSessionID(ctx context.Context, gen IDGenerator) context.Context {
	if _, ok := (gen).(*NoopIDGen); ok {
		return ctx
	}

	return context.WithValue(ctx, sessionIDKey, gen.NewID())
}
