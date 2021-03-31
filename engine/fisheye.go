package engine

import (
	"context"
	"errors"
	"fmt"
	"image"

	"giautm.dev/captcha/fisheye"
)

type FisheyePreprocessor struct{}

func NewFisheyePreprocessor() *FisheyePreprocessor {
	return &FisheyePreprocessor{}
}

func (p *FisheyePreprocessor) Preprocess(_ context.Context, img image.Image) (image.Image, error) {
	result, distance := fisheye.FindDistance(img)
	if distance < 0 {
		return nil, errors.New("preprocess: can not detect distance")
	}

	fmt.Printf("FISHEYE PREPROCESS DISTANCE: %d\n", distance)
	return result, nil
}
