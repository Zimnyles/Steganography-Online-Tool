package database

import (
	"context"
	"stegano-webapp/steagano-webapp/config"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog"
)

func CreateDataBasePool(config *config.DataBaseConfig, logger *zerolog.Logger) *pgxpool.Pool {
	msg := "database connection: "
	dbpool, err := pgxpool.New(context.Background(), config.Url)
	if err != nil {
		logger.Error().Msg(msg + "NO")
		panic(err)
	}
	logger.Info().Msg(msg + "OK")
	return dbpool
}
