package labeled

import (
	"os"
	"path/filepath"
	"strings"
)

func CaptchaFromName(name string) string {
	return strings.SplitN(name, "--", 2)[0]
}

func CreateLabeledFile(dir, label, name string) (*os.File, error) {
	labelDir := filepath.Join(dir, label)
	if err := os.MkdirAll(labelDir, os.ModePerm); err != nil {
		return nil, err
	}

	return os.Create(filepath.Join(labelDir, name))
}
