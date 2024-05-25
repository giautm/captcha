package engine

type EngineOption struct {
	captchaLen   int
	binaryWidth  int
	preprocessor Preprocessor
	symbol       SymbolResolver
	stats        bool
}

type Option func(*EngineOption) error

func WithBinaryWidth(width int) Option {
	return func(opt *EngineOption) error {
		opt.binaryWidth = width
		return nil
	}
}

func WithCaptchaLength(len int) Option {
	return func(opt *EngineOption) error {
		opt.captchaLen = len
		return nil
	}
}

func WithPreprocessor(preprocessor Preprocessor) Option {
	return func(opt *EngineOption) error {
		if opt.stats {
			preprocessor = &StatsPreprocessor{preprocessor}
		}
		opt.preprocessor = preprocessor
		return nil
	}
}

func WithSymbolResolver(sr SymbolResolver) Option {
	return func(opt *EngineOption) error {
		if opt.stats {
			sr = &StatsSymbolResolver{sr}
		}
		opt.symbol = sr
		return nil
	}
}
