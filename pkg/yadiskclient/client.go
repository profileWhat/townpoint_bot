package yadiskclient

import (
	"context"
	"fmt"
	"io"
	"net/http"
)

type Configuration struct {
	Token      string `json:"token,omitempty"`
	HTTPClient *http.Client
}

type APIClient struct {
	cfg *Configuration

	API API
}

func NewAPIClient(cfg *Configuration) *APIClient {
	api := &APIClient{
		cfg: cfg,
	}

	api.API = NewAPIService(api)

	return api
}

func (s *APIClient) createRequest(httpMethod string, urlAPIMethod string, reqBody io.Reader) (*http.Request, error) {
	ctx := context.Background()
	url := fmt.Sprintf("https://cloud-api.yandex.net/v1%s", urlAPIMethod)
	req, err := http.NewRequestWithContext(ctx, httpMethod, url, reqBody)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Authorization", s.cfg.Token)

	return req, nil
}

func (s *APIClient) doRequest(req *http.Request) ([]byte, *http.Response, error) {
	resp, err := s.cfg.HTTPClient.Do(req)
	if err != nil {
		return nil, resp, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, resp, err
	}

	return body, resp, nil
}
