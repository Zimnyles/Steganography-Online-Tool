package profile

import (
	"context"
	"fmt"
	"stegano-webapp/steagano-webapp/internal/models"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog"
)

type ProfileRepository struct {
	Dbpool       *pgxpool.Pool
	CustomLogger *zerolog.Logger
}

func NewProfileRepository(dbpool *pgxpool.Pool, customLogger *zerolog.Logger) *ProfileRepository {
	return &ProfileRepository{
		Dbpool:       dbpool,
		CustomLogger: customLogger,
	}
}

func (r *ProfileRepository) IsLoginExistsForString(login string, logger *zerolog.Logger) (bool, error) {
	var exists bool
	err := r.Dbpool.QueryRow(
		context.Background(),
		"SELECT EXISTS(SELECT 1 FROM users WHERE login = $1)",
		login,
	).Scan(&exists)

	return exists, err
}

func (r *ProfileRepository) GetUserDataFromLogin(login string, logger *zerolog.Logger) (*models.ProfileCreditionals, error) {
	logger.Info().Msg("1")
	query := `
        SELECT 
			email,
            login,
			encrypted,
			transcribed,
			createdat 
        FROM users 
        WHERE login = @login`
	args := pgx.NamedArgs{
		"login": login,
	}

	row := r.Dbpool.QueryRow(context.Background(), query, args)
	var ProfileCredentials models.ProfileCreditionals

	err := row.Scan(&ProfileCredentials.Email, &ProfileCredentials.Login, &ProfileCredentials.Encrypted, &ProfileCredentials.Transcribed, &ProfileCredentials.Createdat	)

	if err != nil {
		logger.Info().Msg("2")
		return nil, fmt.Errorf("error scanning password s36 : %w", err)
	}
	return &ProfileCredentials, nil
}
