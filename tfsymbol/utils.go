package tfsymbol

import (
	"bufio"
	"image"
	"os"
)

type Labels []string

func (s Labels) BestMatch(probabilities []float32) string {
	bestIdx := 0
	for i, p := range probabilities {
		if p > probabilities[bestIdx] {
			bestIdx = i
		}
	}

	return s[bestIdx]
}

func ReadLabels(labelsFile string) (Labels, error) {
	// Read the string from labelsFile, which
	// contains one line per label.
	file, err := os.Open(labelsFile)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var labels []string
	for scanner.Scan() {
		labels = append(labels, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return labels, nil
}

// ImageToTensorValue return an array with
// shape [1,50,180,4] for tensor as an input
func ImageToTensorValue(img image.Image) [][][][]float32 {
	bounds := img.Bounds()

	result := make([][][]float32, bounds.Dy())
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		dy := y - bounds.Min.Y
		result[dy] = make([][]float32, bounds.Dx())
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			r, g, b, a := img.At(x, y).RGBA()

			dx := x - bounds.Min.X
			result[dy][dx] = []float32{
				(float32)(r), (float32)(g), (float32)(b), (float32)(a),
			}
		}
	}

	return [][][][]float32{result}
}
