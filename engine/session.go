package engine

import (
	"context"
	"errors"
	"io"
	"net/http"
	"net/http/cookiejar"
	"time"
)

type (
	CaptchaSession struct {
		// Captcha is the function to fetch the captcha image, it should use
		// the HTTPDoer to make requests.
		Captcha func(_ context.Context, c HTTPDoer) (io.ReadCloser, error)
		// Main is the function will be invoked after the captcha is resolved.
		Main func(_ context.Context, c HTTPDoer, captcha string) (any, error)
		// Engine is the captcha resolver.
		Engine CaptchaResolver
		// IDGenerator generates a new session ID.
		IDGenerator IDGenerator
		// RetryCount is the number of retries, 0 means no retry.
		RetryCount int
		// Transport is the HTTP transport used to make requests.
		Transport http.RoundTripper
	}
	HTTPDoer interface {
		Do(*http.Request) (*http.Response, error)
	}
	IDGenerator interface {
		NewID() string
	}
	IDFunc    func() string
	sessionID struct{}
)

var (
	ErrLimitExceeded = errors.New("session: retry limit exceeded")
)

// Start starts the captcha session, it will fetch the captcha image,
// resolve the captcha, and then invoke the main function.
func (h *CaptchaSession) Start(ctx context.Context) (any, error) {
	jar, err := cookiejar.New(nil)
	if err != nil {
		return nil, err
	}
	client := &http.Client{
		Jar:       jar,
		Timeout:   time.Second * 10,
		Transport: h.Transport,
	}
	handler := func(ctx context.Context) (any, error) {
		file, err := h.Captcha(ctx, client)
		if err != nil {
			return nil, err
		}
		defer file.Close()
		result, err := h.Engine.ResolveFile(ctx, file)
		if err != nil {
			return nil, err
		}
		t, err := h.Main(ctx, client, result.Captcha)
		if e, ok := h.Engine.(ResultReporter); ok {
			if err == nil {
				e.Report(ctx, result, true)
			} else if errors.Is(err, ErrCaptchaInvalid) {
				e.Report(ctx, result, false)
			}
		}
		return t, err
	}
	for i := 0; i <= h.RetryCount; i++ {
		switch t, err := handler(WithSessionID(ctx, h.IDGenerator)); {
		case err == nil:
			return t, nil
		case errors.Is(err, ErrLimitExceeded):
			continue
		default:
			return t, err
		}
	}
	return nil, ErrLimitExceeded
}

// NewID returns a new session ID.
func (fn IDFunc) NewID() string {
	return fn()
}

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
	if gen == nil {
		return ctx
	}
	return context.WithValue(ctx, sessionIDKey, gen.NewID())
}
