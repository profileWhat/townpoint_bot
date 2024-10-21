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
	TGbot  TGbot  `toml:"tgbot"`
	Links  Links  `toml:"links"`
	Source Source `toml:"source"`
}

type TGbot struct {
	Token string `toml:"token"`
}

type Links struct {
	Remarks string `toml:"remarks"`
	Kate    string `toml:"kate"`
}

type AbstractPicture struct {
	Name        string `toml:"name"`
	Path        string `toml:"path"`
	Description string `toml:"description"`
}

type Source struct {
	KateRitsonAbout  string            `toml:"kate_ritson_about"`
	AbstractPictures []AbstractPicture `toml:"abstract_pictures"`
}
