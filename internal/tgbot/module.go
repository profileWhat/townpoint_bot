package tgbot

import (
	"context"

	"go.uber.org/fx"
)

// Module provides tgbot and adds lifecycle hooks for graceful shutdown
var Module = fx.Module("tgbot",
	fx.Provide(New),
	fx.Invoke(
		Invoke,
	),
)

// Invoke tgbot module and appends fx.Lifecycle hooks for graceful shutdown
func Invoke(lc fx.Lifecycle, tgbot *TGBot) {
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			go func() {
				tgbot.Start()
			}()
			return nil
		},

		OnStop: func(ctx context.Context) error {
			return tgbot.Stop()
		},
	})
}
