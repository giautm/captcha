package fisheye

import (
	"image"
	"image/draw"
	"math"
)

// MakeNoise applies the fish eye effect to the image,
// and then tries to remove the fish eye effect from the image.
//
// The image will be used as the input for training the model.
func MakeNoise(img image.Image, distance int, factory func() draw.Image) draw.Image {
	fe := ApplyFishEye(clone(factory(), img), img, distance)
	return RemoveFishEye(factory(), fe, distance)
}

// ApplyFishEye applies the fish eye effect to the image.
func ApplyFishEye(dest draw.Image, src image.Image, distance int) draw.Image {
	b := dest.Bounds()
	midX, midY, radius := b.Dx()/2, b.Dy()/2, float64(distance)
	for y := b.Min.Y; y < b.Max.Y; y++ {
		for x := b.Min.X; x < b.Max.X; x++ {
			relX, relY := float64(x-midX), float64(y-midY)
			if r := math.Sqrt(relX*relX + relY*relY); r < radius {
				tmp := formula(r/radius) * radius / r
				x2, y2 := midX+(int)(tmp*relX), midY+(int)(tmp*relY)
				dest.Set(x, y, src.At(x2, y2))
			} else {
				dest.Set(x, y, src.At(x, y))
			}
		}
	}
	return dest
}

// RemoveFishEye try to remove the fish eye effect from the image.
// The image still be loss some pixels at the center
func RemoveFishEye(dest draw.Image, src image.Image, distance int) draw.Image {
	b := src.Bounds()
	midX, midY, radius := b.Dx()/2, b.Dy()/2, float64(distance)
	for y := b.Min.Y; y < b.Max.Y; y++ {
		for x, relY := b.Min.X, float64(y-midY); x < b.Max.X; x++ {
			x2, y2, relX := x, y, float64(x-midX)
			if r := math.Sqrt(relX*relX + relY*relY); r < radius {
				tmp := formula(r/radius) * radius / r
				x2, y2 = midX+(int)(tmp*relX), midY+(int)(tmp*relY)
			}
			dest.Set(x2, y2, src.At(x, y))
		}
	}
	return dest
}

// FindDistance finds the distance of the fish eye effect.
func FindDistance(src image.Image, testRowIndex int) (image.Image, int) {
	results := make(map[int]image.Image)
	b := src.Bounds()
	dx := b.Dx()
	score, distance := dx, 0
	for d, maxDistance := distanceRange(dx); d < maxDistance; d++ {
		results[d] = RemoveFishEye(image.NewRGBA(b), src, d)
		if s := WhitePoints(results[d], testRowIndex); s < score {
			score = s
			distance = d
		}
	}
	return results[distance], distance
}

// WhitePoints returns the number of white points at the y-axis.
func WhitePoints(src image.Image, y int) int {
	b, score := src.Bounds(), 0
	for x := b.Min.X; x < b.Max.X; x++ {
		_, _, _, a := src.At(x, y).RGBA()
		if a == 0 {
			score++
		}
	}
	return score
}

func clone(dst draw.Image, src image.Image) draw.Image {
	b := src.Bounds()
	draw.Draw(dst, b, src, b.Min, draw.Src)
	return dst
}

// formula calculates the fish eye effect.
func formula(s float64) float64 {
	if s < 0 {
		return 0
	} else if s > 1 {
		return s
	}
	return -0.75*s*s*s + 1.5*s*s + 0.25*s
}

// distanceRange returns the range of the distance.
func distanceRange(dx int) (int, int) {
	return dx / 4, dx/3 + 1
}
