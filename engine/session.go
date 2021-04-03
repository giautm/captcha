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

type CaptchaSession struct {
	CaptchaURL string
	Engine     CaptchaResolver
	RetryCount int
	Transport  http.RoundTripper
}

func NewCaptchaSession(captchaURL string, eng CaptchaResolver) *CaptchaSession {
	return &CaptchaSession{
		CaptchaURL: captchaURL,
		Engine:     eng,
		RetryCount: 5,
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
		err = func() error {
			file, err := srcimg.DownloadReader(ctx, client, h.CaptchaURL)
			if err != nil {
				return err
			}
			defer file.Close()

			result, err := h.Engine.ResolveFile(ctx, file)
			if err != nil {
				return err
			}
			defer ReportResult(h.Engine, result)(ctx, err)

			return fn(client, result.Captcha)
		}()
		if err == ErrCaptchaInvalid {
			continue
		}

		return err
	}

	return ErrLimitExceeded
}
