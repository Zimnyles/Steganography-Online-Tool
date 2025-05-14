package transcribe

import (
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog"
)

type TranscribeRepository struct {
	Dbpool       *pgxpool.Pool
	CustomLogger *zerolog.Logger
}

func NewEncryptRepository(dbpool *pgxpool.Pool, customLogger *zerolog.Logger) *TranscribeRepository {
	return &TranscribeRepository{
		Dbpool:       dbpool,
		CustomLogger: customLogger,
	}
}
