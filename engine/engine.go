package engine

import (
	"context"
	"errors"
	"image"
	"io"

	"giautm.dev/captcha/binimg"
)

type (
	CaptchaResolveEngine struct {
		binaryWidth  int
		captchaLen   int
		preprocessor Preprocessor
		symResolver  SymbolResolver
	}
	CaptchaResult struct {
		Captcha string `json:"captcha"`
	}
	CaptchaResolver interface {
		ResolveFile(ctx context.Context, r io.Reader) (*CaptchaResult, error)
		ResolveImage(ctx context.Context, img image.Image) (*CaptchaResult, error)
	}
	Preprocessor interface {
		Transform(ctx context.Context, img image.Image) (image.Image, error)
	}
	SymbolResolver interface {
		SymbolResolve(ctx context.Context, img image.Image) (string, error)
	}
	ResultReporter interface {
		Report(ctx context.Context, result *CaptchaResult, correct bool) error
	}
)

var (
	ErrCaptchaInvalid = errors.New("captcha is invalid")
)

// NewCaptchaResolveEngine creates a new captcha resolve engine.
func NewCaptchaResolveEngine(opts ...Option) (*CaptchaResolveEngine, error) {
	opt := &EngineOption{
		binaryWidth: 10,
		captchaLen:  5,
	}
	for _, fn := range opts {
		if err := fn(opt); err != nil {
			return nil, err
		}
	}
	return &CaptchaResolveEngine{
		binaryWidth:  opt.binaryWidth,
		captchaLen:   opt.captchaLen,
		preprocessor: opt.preprocessor,
		symResolver:  opt.symbol,
	}, nil
}

// ResolveFile resolves the captcha from the file.
func (e *CaptchaResolveEngine) ResolveFile(ctx context.Context, r io.Reader) (*CaptchaResult, error) {
	img, _, err := image.Decode(r)
	if err != nil {
		return nil, err
	}
	return e.ResolveImage(ctx, img)
}

// ResolveImage resolves the captcha from the image.
func (e *CaptchaResolveEngine) ResolveImage(ctx context.Context, img image.Image) (*CaptchaResult, error) {
	var err error
	if e.preprocessor != nil {
		if img, err = e.preprocessor.Transform(ctx, img); err != nil {
			return nil, err
		}
	}
	captcha := ""
	binImages := binimg.GenImages(img, e.captchaLen, e.binaryWidth)
	for _, bimg := range binImages {
		symbol, err := e.symResolver.SymbolResolve(ctx, bimg)
		if err != nil {
			return nil, err
		}
		captcha += symbol
	}
	return &CaptchaResult{Captcha: captcha}, nil
}
