package main

import (
	"context"
	"flag"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"os"
	"path/filepath"

	"giautm.dev/captcha/binimg"
	"giautm.dev/captcha/engine"
	"giautm.dev/captcha/fisheye"
	"giautm.dev/captcha/labeled"
	"github.com/gammazero/workerpool"
)

var (
	dir          = flag.String("dir", "./labeled", "Input directory with labeled images")
	outputDir    = flag.String("output", "./labeled", "Output directory with labeled images")
	outputFormat = flag.String("outputFormat", "png", "Format of output images: png/jpeg")
	workers      = flag.Int("workers", 100, "Number of workers")
	processor    = flag.String("processor", "", "Pre-processor for images")
)

func main() {
	flag.Parse()
	if *dir == "" || *outputDir == "" || *workers < 1 {
		flag.Usage()
		return
	}

	var p engine.Preprocessor = &noopProcess{}
	if *processor == "fisheye" {
		p = fisheye.NewPreprocessor()
	}

	wp := workerpool.New(*workers)
	err := filepath.Walk(*dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			fmt.Printf("prevent panic by handling failure accessing a path %q: %v\n", path, err)
			return err
		}
		if info.IsDir() {
			return nil
		}

		wp.Submit(func() {
			defer fmt.Printf("processed file: %q\n", path)
			processFile(p, path, info.Name())
		})

		return nil
	})
	wp.StopWait()

	if err != nil {
		fmt.Printf("error walking the path %q: %v\n", *dir, err)
		return
	}
}

type noopProcess struct{}

func (*noopProcess) Preprocess(_ context.Context, img image.Image) (image.Image, error) {
	return img, nil
}

func processFile(p engine.Preprocessor, path, name string) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	img, _, err := image.Decode(file)
	if err != nil {
		return err
	}

	result, err := p.Preprocess(context.Background(), img)
	if err != nil {
		return err
	}

	captcha := labeled.CaptchaFromName(name)
	imgs := binimg.AttachBinaryImages(result, len(captcha), 10)
	for idx, bimg := range imgs {
		err = saveImage(bimg, string(captcha[idx]), name)
		if err != nil {
			return err
		}
	}

	return nil
}

func saveImage(img image.Image, label, name string) error {
	f, err := labeled.CreateLabeledFile(*outputDir, label, name)
	if err != nil {
		return err
	}
	defer f.Close()

	if *outputFormat == "png" {
		err = png.Encode(f, img)
	} else {
		err = jpeg.Encode(f, img, &jpeg.Options{Quality: 100})
	}
	return err
}
