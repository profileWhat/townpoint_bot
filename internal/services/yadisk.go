package services

import (
	"net/http"
	"townpoint_bot/config"
	"townpoint_bot/pkg/yadiskclient"
)

type Yadisk struct {
	Client *yadiskclient.APIClient
}

func NewYadisk(cfg *config.Config) *Yadisk {
	cl := yadiskclient.NewAPIClient(&yadiskclient.Configuration{
		Token:      cfg.Yadisk.Token,
		HTTPClient: &http.Client{},
	})

	return &Yadisk{
		Client: cl,
	}
}
