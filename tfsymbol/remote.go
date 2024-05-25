package tfsymbol

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"image"
	"net/http"
	"net/url"
)

type (
	// RemoteResolver is a struct that resolves the symbol
	// of the image using the TensorFlow Serving server.
	RemoteResolver struct {
		client     HTTPDoer
		predictURL string
		labels     LabelLookup
	}
	// LabelLookup is an interface that provides a way to
	// look up the label of the symbol.
	LabelLookup interface {
		// BestMatch returns the label with the highest probability.
		BestMatch(probabilities []float32) (string, error)
	}
	// Option is a function that sets an option on the RemoteResolver.
	Option func(*options) error
	// HTTPDoer is an interface that provides a way to make HTTP requests.
	HTTPDoer interface {
		Do(*http.Request) (*http.Response, error)
	}
)

// WithHTTPClient sets the HTTP client to use for making requests to the TensorFlow Serving server.
// The default is http.DefaultClient.
func WithHTTPClient(c HTTPDoer) Option {
	return func(o *options) error {
		o.client = c
		return nil
	}
}

// WithBaseURL sets the base URL of the TensorFlow Serving server.
// The default is http://localhost:8501.
func WithBaseURL(u string) Option {
	return func(o *options) (err error) {
		o.baseURL, err = url.Parse(u)
		return err
	}
}

// WithModelName sets the name of the model to use for prediction.
//
// The default is "resnet".
func WithModelName(name string) Option {
	return func(o *options) error {
		o.model = name
		return nil
	}
}

// NewRemoteResolver creates a new RemoteResolver that uses the TensorFlow Serving server to resolve symbols.
func NewRemoteResolver(labels LabelLookup, opt ...Option) (*RemoteResolver, error) {
	if labels == nil {
		return nil, fmt.Errorf("labels must not be nil")
	}
	b, _ := url.Parse("http://:8501")
	opts := &options{
		baseURL: b,
		client:  http.DefaultClient,
		model:   "resnet",
	}
	for _, o := range opt {
		if err := o(opts); err != nil {
			return nil, err
		}
	}
	predictURL, err := opts.baseURL.Parse(
		fmt.Sprintf("/v1/models/%s:predict", opts.model))
	if err != nil {
		return nil, err
	}
	return &RemoteResolver{
		predictURL: predictURL.String(),
		client:     http.DefaultClient,
		labels:     labels,
	}, nil
}

// SymbolResolve resolves the symbol of the image using the TensorFlow Serving server.
func (s *RemoteResolver) SymbolResolve(ctx context.Context, img image.Image) (string, error) {
	input := predictRequest{Instances: ImageToTensorValue(img)}
	buf := bytes.NewBuffer(nil)
	if err := json.NewEncoder(buf).Encode(input); err != nil {
		return "", err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, s.predictURL, buf)
	if err != nil {
		return "", err
	}
	res, err := s.client.Do(req)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()
	var data predictResponse
	if err = json.NewDecoder(res.Body).Decode(&data); err != nil {
		return "", err
	}
	return s.labels.BestMatch(data.Predictions[0])
}

type (
	predictRequest struct {
		Instances [][][][]float32 `json:"instances"`
	}
	predictResponse struct {
		Predictions [][]float32 `json:"predictions"`
	}
	options struct {
		client  HTTPDoer
		baseURL *url.URL
		model   string
	}
)
