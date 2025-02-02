package config

import (
	"github.com/BurntSushi/toml"
)

func New() (*Config, error) {
	cfg := &Config{}

	_, err := toml.DecodeFile("./config.toml", &cfg)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}

type Config struct {
	TGbot    TGbot    `toml:"tgbot"`
	Links    Links    `toml:"links"`
	Source   Source   `toml:"source"`
	Postgres Postgres `toml:"postgres"`
	Yadisk   Yadisk   `toml:"yadisk"`
}

type Yadisk struct {
	Token string `toml:"token"`
	Path  string `toml:"path"`
}

type TGbot struct {
	Token string `toml:"token"`
}

type Links struct {
	Remarks string `toml:"remarks"`
	Kate    string `toml:"kate"`
}

type Source struct {
	PointPicturesPath string `toml:"point_pictures"`
	PointVideosPath   string `toml:"point_videos"`
}

type Postgres struct {
	URL string `toml:"url" env:"POSTGRES_URL"`
}
