package main

import (
	"errors"
	"flag"
	"fmt"
	"image"
	"image/png"
	"os"
	"path/filepath"

	"giautm.dev/captcha/binimg"
	"giautm.dev/captcha/fisheye"
	"giautm.dev/captcha/labeled"
	"github.com/gammazero/workerpool"
)

var (
	dir       = flag.String("dir", "./labeled", "Input directory with labeled images")
	outputDir = flag.String("output", "./labeled", "Output directory with labeled images")
	workers   = flag.Int("workers", 100, "Number of workers")
)

const (
	binWidth     = 10
	testRowIndex = 42
)

func main() {
	flag.Parse()
	if *dir == "" || *outputDir == "" || *workers < 1 {
		flag.Usage()
		return
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
			processFile(path, info.Name())
		})
		return nil
	})
	wp.StopWait()
	if err != nil {
		fmt.Printf("error walking the path %q: %v\n", *dir, err)
		return
	}
}

func processFile(path, name string) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	img, _, err := image.Decode(file)
	if err != nil {
		return err
	}
	result, distance := fisheye.FindDistance(img, testRowIndex)
	if distance < 0 {
		return errors.New("can not detect distance")
	}
	captcha := labeled.CaptchaFromName(name)
	images := binimg.GenImages(result, len(captcha), binWidth)
	for idx, bimg := range images {
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
	return png.Encode(f, img)
}
