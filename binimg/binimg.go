package binimg

import (
	"image"
	"image/draw"
)

func AttachBinaryImages(src image.Image, count, width int) []image.Image {
	imgs := make([]image.Image, count)
	for pos := 0; pos < count; pos++ {
		imgs[pos] = AttachBinaryImage(src, count, width, pos)
	}

	return imgs
}

func AttachBinaryImage(src image.Image, count, width, pos int) draw.Image {
	return DrawBinary(ResizeImage(src, width*count), width, pos)
}

func ResizeImage(src image.Image, width int) draw.Image {
	p := image.Pt(width, 0)
	b := src.Bounds().Add(p)
	b.Min = b.Min.Sub(p)

	dst := image.NewRGBA(b)
	draw.Draw(dst, src.Bounds().Add(p), src, image.ZP, draw.Src)

	return dst
}

func DrawBinary(dst draw.Image, width, position int) draw.Image {
	b := dst.Bounds()
	x := b.Min.X + width*position
	r := image.Rect(x, b.Min.Y, x+width, b.Max.Y)

	draw.Draw(dst, r, image.Black, image.ZP, draw.Src)

	return dst
}
