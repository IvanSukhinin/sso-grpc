package postgresql

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/lib/pq"
	_ "github.com/lib/pq"
	"sso/internal/config"
	"sso/internal/domain/models"
	"sso/internal/storage"
)

const UniqueViolationErrCode = pq.ErrorCode("23505")

type Storage struct {
	db *sql.DB
}

// New creates database connection.
func New(cfg config.Db) (*Storage, error) {
	const op = "storage.postgresql.New"
	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.Db)
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	err = db.Ping()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &Storage{db: db}, nil
}

// SaveUser saves user to db.
func (s *Storage) SaveUser(ctx context.Context, email string, passHash []byte) error {
	const op = "storage.postgresql.SaveUser"
	_, err := s.db.Query("INSERT INTO users(email, pass_hash) VALUES($1, $2)", email, passHash)
	if err != nil {
		var pqErr *pq.Error
		if errors.As(err, &pqErr) && pqErr.Code == UniqueViolationErrCode {
			return storage.ErrEmailExists
		}
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}

// Close closes database connection.
func (s *Storage) Close() error {
	return s.db.Close()
}

// User returns user by email.
func (s *Storage) User(ctx context.Context, email string) (models.User, error) {
	const op = "storage.postgresql.User"

	row := s.db.QueryRow("SELECT id, email, pass_hash FROM users WHERE email = $1", email)
	var user models.User
	err := row.Scan(&user.Id, &user.Email, &user.PassHash)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.User{}, fmt.Errorf("%s: %w", op, storage.ErrUserNotFound)
		}
		return models.User{}, fmt.Errorf("%s: %w", op, err)
	}

	return user, nil
}

// App returns app by app id.
func (s *Storage) App(ctx context.Context, appId int64) (models.App, error) {
	const op = "storage.postgresql.App"

	row := s.db.QueryRow("SELECT id, name, secret FROM apps WHERE id = $1", appId)

	var app models.App
	err := row.Scan(&app.Id, &app.Name, &app.Secret)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.App{}, fmt.Errorf("%s: %w", op, storage.ErrAppNotFound)
		}
		return models.App{}, fmt.Errorf("%s: %w", op, err)
	}

	return app, nil
}

// IsAdmin checks user is admin.
func (s *Storage) IsAdmin(ctx context.Context, userId uuid.UUID) (bool, error) {
	const op = "storage.postgresql.IsAdmin"

	row := s.db.QueryRow("SELECT is_admin FROM users WHERE id = $1", userId)
	var isAdmin bool
	err := row.Scan(&isAdmin)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, fmt.Errorf("%s: %w", op, storage.ErrUserNotFound)
		}
		return false, fmt.Errorf("%s: %w", op, err)
	}

	return isAdmin, nil
}
