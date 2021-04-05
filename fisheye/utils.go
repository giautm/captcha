package fisheye

import (
	"errors"
	"image"
	"image/draw"
	"os"
)

func CloneImage(dst draw.Image, src image.Image) draw.Image {
	b := src.Bounds()
	draw.Draw(dst, b, src, b.Min, draw.Src)
	return dst
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
