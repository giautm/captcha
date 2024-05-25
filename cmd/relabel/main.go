package main

import (
	"errors"
	"flag"
	"fmt"
	"image"
	"image/png"
	"io"
	"os"
	"path/filepath"
	"strings"

	"giautm.dev/captcha/binimg"
	"giautm.dev/captcha/fisheye"
	"giautm.dev/captcha/labeled"
	"github.com/manifoldco/promptui"
)

var (
	dir       = flag.String("dir", "./data/incorrect", "Input directory with labeled images")
	doneDir   = flag.String("done", "./data/done", "Input directory with labeled images")
	outputDir = flag.String("output", "./data/labeled", "Output directory with labeled images")
)

const (
	captchaLen   = 5
	binWidth     = 10
	testRowIndex = 42
)

func main() {
	flag.Parse()
	if *dir == "" || *outputDir == "" {
		flag.Usage()
		return
	}
	prompt := promptui.Prompt{
		Label: "What's the correct captcha?",
		Validate: func(input string) error {
			if len(input) != captchaLen {
				return errors.New("Cần nhập đủ 5 ký tự")
			}
			return nil
		},
	}
	err := filepath.Walk(*dir, func(path string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() {
			return err
		}
		return relabelFile(path, info, prompt.Run)
	})
	if err != nil {
		fmt.Printf("error walking the path %q: %v\n", *dir, err)
		return
	}
}

func relabelFile(path string, info os.FileInfo, askCaptcha func() (string, error)) error {
	fmt.Printf("-> Process: %s\n", path)
	file, err := os.OpenFile(path, os.O_RDONLY, 0)
	if err != nil {
		return err
	}
	defer file.Close()
	preview, err := os.Create("./data/captcha.png")
	if err != nil {
		return err
	}
	defer preview.Close()
	if _, err = io.Copy(preview, file); err != nil {
		return err
	}
	result, err := askCaptcha()
	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		return err
	}
	name := fmt.Sprintf("%s-%d.png", result, info.ModTime().Unix())
	if err = os.Rename(path, filepath.Join(*doneDir, name)); err != nil {
		return err
	}
	captcha := []rune(labeled.CaptchaFromName(info.Name()))
	wrongs := map[int]string{}
	wrongChars := []string{}
	for idx, c := range result {
		if captcha[idx] != c {
			wrongs[idx] = (string)(c)
			wrongChars = append(wrongChars, wrongs[idx])
		}
	}
	if c := len(wrongs); c > 0 {
		fmt.Printf("Found %d wrong character: %s\n", c, strings.Join(wrongChars, ","))
		// Reset to head for decode img
		file.Seek(0, 0)
		img, _, err := image.Decode(file)
		if err != nil {
			return err
		}
		img, distance := fisheye.FindDistance(img, testRowIndex)
		if distance < 0 {
			return errors.New("preprocess: can not detect distance")
		}
		for pos, label := range wrongs {
			if err = genBinFile(img, pos, label, name); err != nil {
				return err
			}
		}
	}
	return nil
}

func genBinFile(img image.Image, pos int, label, name string) error {
	file, err := labeled.CreateLabeledFile(*outputDir, label, name)
	if err != nil {
		return err
	}
	defer file.Close()
	bImg := binimg.MarkPosition(binimg.ExpandLeft(img, binWidth*captchaLen), binWidth, pos)
	return png.Encode(file, bImg)
}
