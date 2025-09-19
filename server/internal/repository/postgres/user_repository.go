package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/LullNil/authx-go/domain/user"
	"github.com/LullNil/authx-go/internal/repository"

	"github.com/lib/pq"
)

type userRepo struct {
	db *sql.DB
}

// NewUserRepository creates a new user repository.
func NewUserRepository(db *sql.DB) *userRepo {
	return &userRepo{
		db: db,
	}
}

// Save saves a new user to the database.
func (r *userRepo) Save(ctx context.Context, user *user.User) (int64, error) {
	const op = "repository.postgres.user.Save"

	query := `
		INSERT INTO users (email, username, password)
		VALUES ($1, $2, $3)
		RETURNING id
	`

	var id int64
	err := r.db.QueryRowContext(
		ctx,
		query,
		user.Email,
		user.Username,
		user.Password,
	).Scan(&id)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			if pqErr.Code == "23505" { // unique_violation
				return 0, repository.ErrConflict
			}
		}
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return id, nil
}

// GetByEmail retrieves an user by email from the database.
func (r *userRepo) GetByEmail(ctx context.Context, email string) (*user.User, error) {
	const op = "repository.postgres.user.GetByEmail"

	query := `
		SELECT id, email, username, password
		FROM users
		WHERE email = $1
	`

	var u user.User
	err := r.db.QueryRowContext(ctx, query, email).Scan(
		&u.ID,
		&u.Email,
		&u.Username,
		&u.Password,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, repository.ErrNotFound
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &u, nil
}

// GetByID retrieves an user by ID from the database.
func (r *userRepo) GetByID(ctx context.Context, id int64) (*user.User, error) {
	const op = "repository.postgres.user.GetByID"

	query := `
		SELECT id, email, username
		FROM users
		WHERE id = $1
	`

	var u user.User
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&u.ID,
		&u.Email,
		&u.Username,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, repository.ErrNotFound
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &u, nil
}

// GetByUsername retrieves user by username
func (r *userRepo) GetByUsername(ctx context.Context, username string) (*user.User, error) {
	const op = "repository.postgres.user.GetByUsername"

	query := `
		SELECT id, email, username, password
		FROM users
		WHERE username = $1
	`

	var u user.User
	err := r.db.QueryRowContext(ctx, query, username).Scan(
		&u.ID,
		&u.Email,
		&u.Username,
		&u.Password,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, repository.ErrNotFound
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &u, nil
}
