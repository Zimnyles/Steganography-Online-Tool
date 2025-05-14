package transcribe

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
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

func (r *TranscribeRepository) userTranscribeCounterPlus(userLogin string) error {
	query := `
		UPDATE users
		SET transcribed = transcribed + 1
		WHERE login = @login;`

	args := pgx.NamedArgs{
		"login": userLogin,
	}
	_, err := r.Dbpool.Exec(context.Background(), query, args)
	if err != nil {
		return fmt.Errorf("failed to update encrypted counter: %w", err)
	}
	return nil

}

func (r *TranscribeRepository) AllUsersCounterPlus() error {
	query := `
		UPDATE usersactions
		SET actions = actions + 1
		WHERE id = 1
	`
	_, err := r.Dbpool.Exec(context.Background(), query)
	if err != nil {
		return fmt.Errorf("failed to update allusers actions counter: %w", err)
	}
	return nil
}
