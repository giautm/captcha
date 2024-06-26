package engine

import (
	"context"
	"image"
	"io"

	"go.opencensus.io/trace"
)

type (
	StatsPreprocessor struct {
		Preprocessor
	}
	StatsSymbolResolver struct {
		SymbolResolver
	}
	StatsCaptchaResolver struct {
		CaptchaResolver
	}
)

func (s *StatsPreprocessor) Transform(ctx context.Context, img image.Image) (image.Image, error) {
	ctx, span := trace.StartSpan(ctx, "engine.Preprocess")
	defer span.End()
	return s.Preprocessor.Transform(ctx, img)
}

func (s *StatsSymbolResolver) SymbolResolve(ctx context.Context, img image.Image) (string, error) {
	ctx, span := trace.StartSpan(ctx, "engine.SymbolResolve")
	defer span.End()
	return s.SymbolResolver.SymbolResolve(ctx, img)
}

func (s *StatsCaptchaResolver) Report(ctx context.Context, captcha *CaptchaResult, correct bool) error {
	return nil
}

func (s *StatsCaptchaResolver) ResolveFile(ctx context.Context, r io.Reader) (*CaptchaResult, error) {
	ctx, span := trace.StartSpan(ctx, "engine.ResolveFile")
	defer span.End()
	return s.CaptchaResolver.ResolveFile(ctx, r)
}

func (s *StatsCaptchaResolver) ResolveImage(ctx context.Context, img image.Image) (*CaptchaResult, error) {
	ctx, span := trace.StartSpan(ctx, "engine.ResolveImage")
	defer span.End()
	return s.CaptchaResolver.ResolveImage(ctx, img)
}
