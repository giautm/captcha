package binimg

import (
	"image"
	"image/draw"
)

// GenImages generates a list of images by attaching the binary image
// to the left of the source image.
func GenImages(src image.Image, count, width int) []image.Image {
	images := make([]image.Image, count)
	for pos := 0; pos < count; pos++ {
		img := ExpandLeft(src, width*count)
		images[pos] = MarkPosition(img, width, pos)
	}
	return images
}

// ExpandLeft expands the image to the left by width pixels.
// The new image will have the same height as the original image.
func ExpandLeft(src image.Image, width int) draw.Image {
	p := image.Pt(width, 0)
	b := src.Bounds().Add(p)
	b.Min = b.Min.Sub(p)
	dst := image.NewRGBA(b)
	draw.Draw(dst, src.Bounds().Add(p), src, image.Point{}, draw.Src)
	return dst
}

// MarkPosition marks the position of the image with a black line.
func MarkPosition(dst draw.Image, width, position int) draw.Image {
	b := dst.Bounds()
	x := b.Min.X + width*position
	r := image.Rect(x, b.Min.Y, x+width, b.Max.Y)
	draw.Draw(dst, r, image.Black, image.Point{}, draw.Src)
	return dst
}
