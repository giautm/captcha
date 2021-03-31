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

type PredictRequest struct {
	Instances [][][][]float32 `json:"instances"`
}

type PredictReponse struct {
	Predictions []float32 `json:"predictions"`
}

type HttpClient interface {
	Do(*http.Request) (*http.Response, error)
}

type LabelLookup interface {
	BestMatch(probabilities []float32) string
}

type RemoteResolver struct {
	predictURL string
	client     HttpClient
	labels     LabelLookup
}

func NewRemoteResolver(
	baseURL, modelName string,
	client HttpClient,
	labels LabelLookup,
) (*RemoteResolver, error) {
	base, err := url.Parse(baseURL)
	if err != nil {
		return nil, err
	}

	predictURL, err := base.Parse(fmt.Sprintf("/v1/models/%s:predict", modelName))
	if err != nil {
		return nil, err
	}

	return &RemoteResolver{
		predictURL: predictURL.String(),
		client:     client,
		labels:     labels,
	}, nil
}

func (s *RemoteResolver) SymbolResolve(ctx context.Context, img image.Image) (string, error) {
	input := PredictRequest{
		Instances: ImageToTensorValue(img),
	}

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

	var data PredictReponse
	if err = json.NewDecoder(res.Body).Decode(&data); err != nil {
		return "", err
	}

	return s.labels.BestMatch(data.Predictions), nil
}
