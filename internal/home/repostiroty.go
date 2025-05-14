package home

import (
	"context"
	"fmt"
	"stegano-webapp/steagano-webapp/internal/models"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog"
	"golang.org/x/crypto/bcrypt"
)

type UsersRepository struct {
	Dbpool       *pgxpool.Pool
	CustomLogger *zerolog.Logger
}

func NewUsersRepository(dbpool *pgxpool.Pool, customLogger *zerolog.Logger) *UsersRepository {
	return &UsersRepository{
		Dbpool:       dbpool,
		CustomLogger: customLogger,
	}
}

func (r *UsersRepository) getUserActionsCounter(login string, logger *zerolog.Logger) (*int, error) {
	query := `
    SELECT actions
    FROM usersactions
	WHERE id = 1
	`

	var actions int
	failCounter := 0
	row := r.Dbpool.QueryRow(context.Background(), query)
	err := row.Scan(&actions)
	fmt.Println(actions)
	if err != nil {
    	return &failCounter, fmt.Errorf("failed to get actions: %w", err)
	}

	return &actions, nil
}

func (r *UsersRepository) getTranscribedCounter(login string, logger *zerolog.Logger) (*int, error) {
	query := `
    SELECT transcribed 
    FROM users 
    WHERE login = @login`

	args := pgx.NamedArgs{
		"login": login, // предполагается, что userLogin передаётся в функцию
	}

	row := r.Dbpool.QueryRow(context.Background(), query, args)
	var transcribedCounter int

	err := row.Scan(&transcribedCounter)
	failCounter := 0

	if err != nil {
		return &failCounter, fmt.Errorf("error: %w", err)
	}

	return &transcribedCounter, nil

}

func (r *UsersRepository) getEncryptedCounter(login string, logger *zerolog.Logger) (*int, error) {
	query := `
    SELECT encrypted 
    FROM users 
    WHERE login = @login`

	args := pgx.NamedArgs{
		"login": login, // предполагается, что userLogin передаётся в функцию
	}

	row := r.Dbpool.QueryRow(context.Background(), query, args)
	var encryptCounter int

	err := row.Scan(&encryptCounter)
	failCounter := 0

	if err != nil {
		return &failCounter, fmt.Errorf("error: %w", err)
	}

	return &encryptCounter, nil

}

func (r *UsersRepository) getUserData(login string, logger *zerolog.Logger) (*models.UserData, error) {
	query := `
        SELECT 
            login,  
            email,
			encrypted,
			transcribed   
        FROM users 
        WHERE login = @login`
	args := pgx.NamedArgs{
		"login": login,
	}
	row := r.Dbpool.QueryRow(context.Background(), query, args)
	var UserData models.UserData
	err := row.Scan(&UserData.Login, &UserData.Email, &UserData.Encrypted, &UserData.Transcribed)
	if err != nil {
		return nil, fmt.Errorf("error: %w", err)
	}
	fmt.Println(UserData)
	fmt.Println(login)
	return &UserData, nil
}

func (r *UsersRepository) addUser(form UserCreateForm, logger *zerolog.Logger) (string, error) {

	emailIsExists, err := r.IsEmailExists(form, logger)
	if emailIsExists {
		return "Аккаунт с таким email уже существует", fmt.Errorf("невозможно зарегестрировать аккаунт : %w", err)
	}

	loginIsExists, err := r.IsLoginExists(form, logger)
	if loginIsExists {
		return "Аккаунт с таким логином уже существует", fmt.Errorf("невозможно зарегестрировать аккаунт : %w", err)
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(form.Password), bcrypt.DefaultCost)
	if err != nil {
		return "Ошибка сервера, попробуйте позже", fmt.Errorf("невозможно зарегестрировать аккаунт : %w", err)
	}

	query := `
		INSERT INTO users (email, login, password, createdat) 
		VALUES (@email, @login, @password, @createdat)
	`
	args := pgx.NamedArgs{
		"email":     form.Email,
		"login":     strings.ToLower(form.Login),
		"password":  hashedPassword,
		"createdat": time.Now(),
	}

	_, err = r.Dbpool.Exec(context.Background(), query, args)
	if err != nil {
		return "Ошибка сервера, попробуйте позже", fmt.Errorf("невозможно зарегестрировать аккаунт : %w", err)
	}
	logger.Info().Msg("зарегестрирован аккаунт")
	return "Аккаунт зарегестрирован", nil

}

func (r *UsersRepository) IsEmailExists(form UserCreateForm, logger *zerolog.Logger) (bool, error) {
	var exists bool
	err := r.Dbpool.QueryRow(
		context.Background(),
		"SELECT EXISTS(SELECT 1 FROM users WHERE email = $1)",
		form.Email,
	).Scan(&exists)

	return exists, err
}

func (r *UsersRepository) IsLoginExists(form UserCreateForm, logger *zerolog.Logger) (bool, error) {
	var exists bool
	err := r.Dbpool.QueryRow(
		context.Background(),
		"SELECT EXISTS(SELECT 1 FROM users WHERE login = $1)",
		form.Login,
	).Scan(&exists)

	return exists, err
}

func (r *UsersRepository) GetPasswordByEmail(form LoginForm, logger *zerolog.Logger) (*UserCredentials, error) {
	query := `
        SELECT 
            login,  
            password   
        FROM users 
        WHERE email = @email`
	args := pgx.NamedArgs{
		"email": form.Email,
	}
	row := r.Dbpool.QueryRow(context.Background(), query, args)
	var UserCredentials UserCredentials
	err := row.Scan(&UserCredentials.Login, &UserCredentials.PasswordHash)
	if err != nil {
		return nil, fmt.Errorf("error scanning password s36 : %w", err)
	}
	return &UserCredentials, nil

}

func (r *UsersRepository) IsEmailExistsForLogin(form LoginForm, logger *zerolog.Logger) (bool, error) {
	var exists bool
	err := r.Dbpool.QueryRow(
		context.Background(),
		"SELECT EXISTS(SELECT 1 FROM users WHERE email = $1)",
		form.Email,
	).Scan(&exists)

	return exists, err
}
