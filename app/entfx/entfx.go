package entfx

import (
	"context"
	"database/sql"
	"townpoint_bot/config"
	ent "townpoint_bot/ent/generated"

	entsql "entgo.io/ent/dialect/sql"
	_ "github.com/lib/pq"
	"go.uber.org/fx"
)

// Module provides ent.Client instantiated with values from config
var Module = fx.Module("ent",
	fx.Provide(New),
) //nolint:gochecknoglobals // DI Container

// New returns ent.Client with open postgres connection
func New(lc fx.Lifecycle, cfg *config.Config) (*ent.Client, error) {
	client, err := Open(cfg)
	if err != nil {
		return nil, err
	}

	lc.Append(fx.Hook{
		OnStop: func(ctx context.Context) error {
			return client.Close()
		},
	})

	return client, nil
}

// Open new postgres connection and return ent.Client
func Open(cfg *config.Config) (*ent.Client, error) {
	db, err := sql.Open("postgres", cfg.Postgres.URL)
	if err != nil {
		return nil, err
	}

	// CreateOrder an ent.Driver from `db`.
	drv := entsql.OpenDB("postgres", db)
	return ent.NewClient(ent.Driver(drv)), nil
}