package fisheye

import (
	"image"
	"image/draw"
	"math"
)

func Formula(s float64) float64 {
	if s < 0 {
		return 0
	} else if s > 1 {
		return s
	}

	return -0.75*s*s*s + 1.5*s*s + 0.25*s
}

func DistanceRange(dx int) (int, int) {
	return dx / 4, dx/3 + 1
}

func FindDistance(src image.Image, testRowIndex int) (image.Image, int) {
	results := make(map[int]image.Image)

	b := src.Bounds()
	dx := b.Dx()
	score, distance := dx, 0
	for d, maxDistance := DistanceRange(dx); d < maxDistance; d++ {
		results[d] = RemoveFishEye(image.NewRGBA(b), src, d)
		if s := WhitePoints(results[d], testRowIndex); s < score {
			score = s
			distance = d
		}
	}

	return results[distance], distance
}

func WhitePoints(src image.Image, y int) int {
	score := 0

	b := src.Bounds()
	for x := b.Min.X; x < b.Max.X; x++ {
		_, _, _, a := src.At(x, y).RGBA()
		if a == 0 {
			score++
		}
	}

	return score
}

func RemoveFishEye(dest draw.Image, src image.Image, distance int) draw.Image {
	dis := float64(distance)

	b := src.Bounds()
	midY, midX := b.Dy()/2, b.Dx()/2

	for y := b.Min.Y; y < b.Max.Y; y++ {
		relY := float64(y - midY)
		for x := b.Min.X; x < b.Max.X; x++ {
			relX := float64(x - midX)
			d := math.Sqrt(relX*relX + relY*relY)

			y2, x2 := y, x
			if d < dis {
				tmp := Formula(d/dis) * dis / d
				x2 = midX + (int)(tmp*relX)
				y2 = midY + (int)(tmp*relY)
			}

			dest.Set(x2, y2, src.At(x, y))
		}
	}

	return dest
}

func ApplyFishEye(dest draw.Image, src image.Image, distance int) draw.Image {
	dis := float64(distance)

	b := dest.Bounds()
	midY, midX := b.Dy()/2, b.Dx()/2

	for y := b.Min.Y; y < b.Max.Y; y++ {
		for x := b.Min.X; x < b.Max.X; x++ {
			relY := float64(y - midY)
			relX := float64(x - midX)

			d := math.Sqrt(relX*relX + relY*relY)
			if d < dis {
				tmp := Formula(d/dis) * dis / d
				x2 := midX + (int)(tmp*relX)
				y2 := midY + (int)(tmp*relY)

				dest.Set(x, y, src.At(x2, y2))
			} else {
				dest.Set(x, y, src.At(x, y))
			}
		}
	}

	return dest
}

func MakeNoise(img image.Image, distance int, factory func() draw.Image) draw.Image {
	fe := ApplyFishEye(CloneImage(factory(), img), img, distance)
	return RemoveFishEye(factory(), fe, distance)
}
