package user

import (
	"authx-go/domain/user"
	"authx-go/internal/repository"
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"regexp"
	"strings"

	"github.com/LullNil/go-http-utils/apperr"
	"golang.org/x/crypto/bcrypt"
)

type service struct {
	userRepo user.Repository
	logger   *slog.Logger
}

// NewService returns a new user service.
func NewService(userRepo user.Repository, logger *slog.Logger) user.Service {
	return &service{
		userRepo: userRepo,
		logger:   logger,
	}
}

var (
	usernameRegexp = regexp.MustCompile(`^[a-z0-9_]+$`)
	emailRegexp    = regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,}$`)
)

// RegisterUser creates a new user with validations.
func (s *service) RegisterUser(ctx context.Context, req user.RegisterUserRequest) (int64, error) {
	const op = "service.user.RegisterUser"

	// normalize
	username := strings.ToLower(strings.TrimSpace(req.Username))
	email := strings.ToLower(strings.TrimSpace(req.Email))

	// validate username
	if len(username) < 3 {
		return 0, apperr.New(http.StatusBadRequest, "username is too short")
	}
	if len(username) > 25 {
		return 0, apperr.New(http.StatusBadRequest, "username is too long")
	}
	if !usernameRegexp.MatchString(username) {
		return 0, apperr.New(http.StatusBadRequest, "invalid username format")
	}

	// validate email
	if !emailRegexp.MatchString(email) {
		return 0, apperr.New(http.StatusBadRequest, "invalid email format")
	}

	// validate password
	if len(req.Password) < 6 {
		return 0, apperr.New(http.StatusBadRequest, "password is too weak")
	}

	// check conflicts
	if u, _ := s.userRepo.GetByEmail(ctx, email); u != nil {
		return 0, apperr.New(http.StatusConflict, "user already exists")
	}
	if u, _ := s.userRepo.GetByUsername(ctx, username); u != nil {
		return 0, apperr.New(http.StatusConflict, "user already exists")
	}

	// hash password
	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	id, err := s.userRepo.Save(ctx, &user.User{
		Email:    email,
		Username: username,
		Password: string(hash),
	})
	if err != nil {
		if errors.Is(err, repository.ErrConflict) {
			return 0, apperr.New(http.StatusConflict, "user already exists")
		}
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return id, nil
}

// LoginUser checks credentials and returns JWT stub.
func (s *service) LoginUser(ctx context.Context, req user.LoginRequest) (string, error) {
	const op = "service.user.LoginUser"

	// normalize email
	email := strings.ToLower(strings.TrimSpace(req.Email))
	if !emailRegexp.MatchString(email) {
		return "", apperr.New(http.StatusBadRequest, "invalid email format")
	}

	// get user
	u, err := s.userRepo.GetByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return "", apperr.New(http.StatusNotFound, "user not found")
		}
		return "", fmt.Errorf("%s: %w", op, err)
	}

	// compare password
	if err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(req.Password)); err != nil {
		return "", apperr.New(http.StatusBadRequest, "invalid login or password")
	}

	// JWT stub
	return "fake-jwt-token", nil
}

// GetUserByID retrieves an user by ID from the database.
func (s *service) GetUserByID(ctx context.Context, id int64) (*user.User, error) {
	const op = "service.user.GetUserByID"

	user, err := s.userRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return user, nil
}

// GetUserByEmail retrieves an user by email from the database.
// func (s *service) GetUserByEmail(ctx context.Context, email string) (*user.User, error) {
// 	const op = "service.user.GetUserByEmail"

// 	email = strings.ToLower(strings.TrimSpace(email))

// 	user, err := s.userRepo.GetByEmail(ctx, email)
// 	if err != nil {
// 		return nil, fmt.Errorf("%s: %w", op, err)
// 	}

// 	return user, nil
// }
