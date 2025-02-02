package yadiskclient

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

type API interface {
	GetUploadLink(filepath string) (*OperationResponse, error)
	GetOperationStatus(href string) (*OperationStatus, error)
	UploadByURL(filepath string, url string) (*OperationResponse, error)
	Publish(filepath string) (*OperationResponse, error)
	PublishInfo(href string) (*PublishInfoResponse, error)
}

type APIService struct {
	client *APIClient
	apiURL string
}

// NewAPIService make new api service implements API
func NewAPIService(client *APIClient) *APIService {
	return &APIService{
		client: client,
		apiURL: "",
	}
}

func (s *APIService) GetUploadLink(filepath string) (*OperationResponse, error) {
	path := strings.ReplaceAll(filepath, "/", "%2F")
	urlAPIMethod := fmt.Sprintf("%s/disk/resources/upload%s?path=%s", s.apiURL, path)
	req, err := s.client.createRequest("GET", urlAPIMethod, nil)
	if err != nil {
		return nil, err
	}

	data, _, err := s.client.doRequest(req)
	if err != nil {
		return nil, err
	}

	var resp *OperationResponse
	err = json.Unmarshal(data, resp)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (s *APIService) UploadByURL(filepath string, url string) (*OperationResponse, error) {
	path := strings.ReplaceAll(filepath, "/", "%2F")
	urlAPIMethod := fmt.Sprintf("%s/disk/resources/upload?path=%s&url=%s", s.apiURL, path, url)
	req, err := s.client.createRequest("POST", urlAPIMethod, nil)
	if err != nil {
		return nil, err
	}

	data, _, err := s.client.doRequest(req)
	if err != nil {
		return nil, err
	}

	resp := OperationResponse{}
	err = json.Unmarshal(data, &resp)
	if err != nil {
		return nil, err
	}

	return &resp, nil
}

func (s *APIService) Publish(filepath string) (*OperationResponse, error) {
	path := strings.ReplaceAll(filepath, "/", "%2F")
	urlAPIMethod := fmt.Sprintf("%s/disk/resources/publish?path=%s", s.apiURL, path)
	req, err := s.client.createRequest("PUT", urlAPIMethod, nil)
	if err != nil {
		return nil, err
	}

	data, _, err := s.client.doRequest(req)
	if err != nil {
		return nil, err
	}

	resp := OperationResponse{}
	err = json.Unmarshal(data, &resp)
	if err != nil {
		return nil, err
	}

	return &resp, nil
}

func (s *APIService) PublishInfo(href string) (*PublishInfoResponse, error) {
	req, err := http.NewRequestWithContext(context.Background(), "GET", href, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Authorization", s.client.cfg.Token)
	data, _, err := s.client.doRequest(req)
	if err != nil {
		return nil, err
	}

	resp := PublishInfoResponse{}
	err = json.Unmarshal(data, &resp)
	if err != nil {
		return nil, err
	}

	return &resp, nil
}

func (s *APIService) GetOperationStatus(href string) (*OperationStatus, error) {
	req, err := http.NewRequestWithContext(context.Background(), "GET", href, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Authorization", s.client.cfg.Token)
	data, _, err := s.client.doRequest(req)
	if err != nil {
		return nil, err
	}

	resp := OperationStatus{}
	err = json.Unmarshal(data, &resp)
	if err != nil {
		return nil, err
	}

	return &resp, nil
}
