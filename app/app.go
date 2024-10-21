package app

import (
	"kate_ritson_art_bot/config"
	"kate_ritson_art_bot/internal/tgbot"

	"go.uber.org/fx"
)

// New returns constructed backend as fx.App (DI Bundle)
func New() *fx.App {
	return fx.New(
		fx.Provide(
			config.New,
		),

		tgbot.Module,
	)
}
