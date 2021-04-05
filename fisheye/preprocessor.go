package fisheye

import (
	"context"
	"errors"
	"image"
)

var (
	ErrDetectDistance = errors.New("fisheye: can not detect distance")
)

type FisheyePreprocessor struct {
	TestRowIndex int
}

func NewPreprocessor() *FisheyePreprocessor {
	return &FisheyePreprocessor{
		// We test pixels at line 42 to check correct distance
		TestRowIndex: 42,
	}
}

func (p *FisheyePreprocessor) Preprocess(_ context.Context, img image.Image) (image.Image, error) {
	result, distance := FindDistance(img, p.TestRowIndex)
	if distance < 0 {
		return nil, ErrDetectDistance
	}

	return result, nil
}
