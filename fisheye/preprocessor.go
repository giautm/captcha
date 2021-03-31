package fisheye

import (
	"context"
	"errors"
	"image"
)

var (
	ErrDetectDistance = errors.New("fisheye: can not detect distance")
)

type FisheyePreprocessor struct{}

func NewPreprocessor() *FisheyePreprocessor {
	return &FisheyePreprocessor{}
}

func (p *FisheyePreprocessor) Preprocess(_ context.Context, img image.Image) (image.Image, error) {
	result, distance := FindDistance(img)
	if distance < 0 {
		return nil, ErrDetectDistance
	}

	return result, nil
}
