package fisheye

import (
	"context"
	"errors"
	"image"
)

type FisheyePreprocessor struct{}

func NewPreprocessor() *FisheyePreprocessor {
	return &FisheyePreprocessor{}
}

func (p *FisheyePreprocessor) Preprocess(_ context.Context, img image.Image) (image.Image, error) {
	result, distance := FindDistance(img)
	if distance < 0 {
		return nil, errors.New("preprocess: can not detect distance")
	}

	return result, nil
}
