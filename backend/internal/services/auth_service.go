package services

import (
	httpcommon "chi-mysql-boilerplate/internal/domain/http_common"
	"chi-mysql-boilerplate/internal/domain/models"
	"context"
	"database/sql"
	"errors"

	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	db *sql.DB
}

func NewAuthService(db *sql.DB) *AuthService {
	return &AuthService{db: db}
}

func (a *AuthService) Register(req models.AuthRequest) error {
	ctx, cancel := context.WithTimeout(context.Background(), httpcommon.DbConstants.Timeout)
	defer cancel()

	// check if the user already exists
	query := `
		SELECT id
		FROM users
		WHERE username = ?
	`
	row := a.db.QueryRowContext(ctx, query, req.Username)
	var id int64
	err := row.Scan(&id)
	if err != nil && err != sql.ErrNoRows {
		return err
	}
	if err == nil {
		return errors.New(httpcommon.ErrorMessage.ErrUserAlreadyExists)
	}

	// hash the password
	req.Password, err = HashPassword(req.Password)
	if err != nil {
		return err
	}

	query = `
		INSERT INTO users (username, password, refresh_token)
		VALUES (?, ?, ?)
	`

	_, err = a.db.ExecContext(ctx, query, req.Username, req.Password, nil)
	if err != nil {
		return err
	}

	return nil
}

func (a *AuthService) Login(req models.AuthRequest) (*models.AuthResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), httpcommon.DbConstants.Timeout)
	defer cancel()

	query := `
		SELECT id, password
		FROM users
		WHERE username = ?
	`
	row := a.db.QueryRowContext(ctx, query, req.Username)

	var res models.AuthResponse
	var password string
	if err := row.Scan(&res.ID, &password); err != nil {
		// if the user does not exist, return a bad credentials error
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New(httpcommon.ErrorMessage.ErrUserDoesNotExist)
		}

		return nil, err
	}

	// check if the password is correct
	if !IsCorrectPassword(req.Password, password) {
		return nil, errors.New(httpcommon.ErrorMessage.BadCredentials)
	}

	return &res, nil
}

func (a *AuthService) UpdateRefreshToken(userId uint64, newRefreshToken string) error {
	ctx, cancel := context.WithTimeout(context.Background(), httpcommon.DbConstants.Timeout)
	defer cancel()

	query := `
		UPDATE users
		SET
			refresh_token = ?
		WHERE id = ?
	`
	_, err := a.db.ExecContext(ctx, query, newRefreshToken, userId)
	if err != nil {
		return err
	}

	return nil
}

func (a *AuthService) RemoveRefreshToken(refreshToken string) error {
	ctx, cancel := context.WithTimeout(context.Background(), httpcommon.DbConstants.Timeout)
	defer cancel()

	query := `
		UPDATE users
		SET
			refresh_token = NULL
		WHERE refresh_token = ?
	`
	_, err := a.db.ExecContext(ctx, query, refreshToken)
	if err != nil {
		return err
	}

	return nil
}

func (a *AuthService) ValidateRefreshToken(userId uint64, refreshToken string) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), httpcommon.DbConstants.Timeout)
	defer cancel()

	query := `
		SELECT refresh_token
		FROM users
		WHERE id = ?
	`
	row := a.db.QueryRowContext(ctx, query, userId)

	var retrievedToken string
	if err := row.Scan(&retrievedToken); err != nil {
		return false, err
	}

	return retrievedToken == refreshToken && retrievedToken != "", nil
}

// helper function to hash the password
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

// helper function to check if the password in the request matches the hashed password in the database
func IsCorrectPassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
