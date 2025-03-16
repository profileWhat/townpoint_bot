package config

import (
	"log"
	"os"

	"github.com/BurntSushi/toml"
	"github.com/joho/godotenv"
)

func New() (*Config, error) {
	cfg := &Config{}

	err := godotenv.Load()
	if err != nil {
		log.Println("cant download env")
	}

	_, err = toml.DecodeFile("./config.toml", &cfg)
	if err != nil {
		return nil, err
	}

	tgbotToken := os.Getenv("TGBOT_TOKEN")
	if tgbotToken != "" {
		cfg.TGbot.Token = tgbotToken
	}

	yadiskToken := os.Getenv("YADISK_TOKEN")
	if yadiskToken != "" {
		cfg.Yadisk.Token = yadiskToken
	}

	yadiskPath := os.Getenv("YADISK_PATH")
	if yadiskPath != "" {
		cfg.Yadisk.Path = yadiskPath
	}

	postgresURL := os.Getenv("POSTGRES_URL")
	if postgresURL != "" {
		cfg.Postgres.URL = postgresURL
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
