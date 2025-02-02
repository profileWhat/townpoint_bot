package app

import (
	"townpoint_bot/app/entfx"
	"townpoint_bot/config"
	"townpoint_bot/internal/services"
	"townpoint_bot/internal/tgbot"

	"go.uber.org/fx"
)

// New returns constructed backend as fx.App (DI Bundle)
func New() *fx.App {
	return fx.New(
		fx.Provide(
			config.New,
			services.NewYadisk,
		),

		entfx.Module,
		tgbot.Module,
	)
}
