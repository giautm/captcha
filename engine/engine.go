package engine

import (
	"context"
	"image"
	"io"

	"giautm.dev/captcha/binimg"
)

type CaptchaResolver interface {
	ResolveFile(ctx context.Context, r io.Reader) (*CaptchaResult, error)
	ResolveImage(ctx context.Context, img image.Image) (*CaptchaResult, error)
}

type Preprocessor interface {
	Preprocess(ctx context.Context, img image.Image) (image.Image, error)
}

type SymbolResolver interface {
	SymbolResolve(ctx context.Context, img image.Image) (string, error)
}

type CaptchaResolveEngine struct {
	binaryWidth  int
	captchaLen   int
	preprocessor Preprocessor
	symResolver  SymbolResolver
}

type CaptchaResult struct {
	ID      string `json:"id"`
	Captcha string `json:"captcha"`
}

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

func (e *CaptchaResolveEngine) ResolveFile(ctx context.Context, r io.Reader) (*CaptchaResult, error) {
	img, _, err := image.Decode(r)
	if err != nil {
		return nil, err
	}

	return e.ResolveImage(ctx, img)
}

func (e *CaptchaResolveEngine) ResolveImage(ctx context.Context, img image.Image) (*CaptchaResult, error) {
	var err error
	if e.preprocessor != nil {
		if img, err = e.preprocessor.Preprocess(ctx, img); err != nil {
			return nil, err
		}
	}

	captcha := ""
	binImgs := binimg.AttachBinaryImages(img, e.captchaLen, e.binaryWidth)
	for _, binImg := range binImgs {
		symbol, err := e.symResolver.SymbolResolve(ctx, binImg)
		if err != nil {
			return nil, err
		}

		captcha += symbol
	}

	return &CaptchaResult{
		Captcha: captcha,
	}, nil
}
