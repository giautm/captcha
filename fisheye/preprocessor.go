package fisheye

import (
	"context"
	"errors"
	"image"
	"os"
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

// Transform implements the Preprocessor interface.
func (p *FisheyePreprocessor) Transform(_ context.Context, img image.Image) (image.Image, error) {
	if result, distance := FindDistance(img, p.TestRowIndex); distance >= 0 {
		return result, nil
	}
	return nil, ErrDetectDistance
}

func GenFile(path string, testRowIndex int) (image.Image, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	img, _, err := image.Decode(file)
	if err != nil {
		return nil, err
	}
	result, distance := FindDistance(img, testRowIndex)
	if distance < 0 {
		return nil, errors.New("can not detect distance")
	}
	return result, nil
}
