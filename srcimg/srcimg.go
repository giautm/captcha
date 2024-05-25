package srcimg

import (
	"context"
	"image"
	"io"
	"net/http"
)

type HTTPDoer interface {
	Do(*http.Request) (*http.Response, error)
}

func DownloadReader(ctx context.Context, c HTTPDoer, url string) (io.ReadCloser, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	res, err := c.Do(req)
	if err != nil {
		return nil, err
	}
	return res.Body, nil
}

func DownloadImage(ctx context.Context, c HTTPDoer, url string) (image.Image, error) {
	file, err := DownloadReader(ctx, c, url)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	img, _, err := image.Decode(file)
	if err != nil {
		return nil, err
	}
	return img, nil
}
