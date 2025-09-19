package user

import (
	"context"
)

type Service interface {
	RegisterUser(ctx context.Context, req RegisterUserRequest) (int64, error)
	LoginUser(ctx context.Context, req LoginRequest) (string, error)
	GetUserByID(ctx context.Context, id int64) (*User, error)
	// GetUserByEmail(ctx context.Context, email string) (*User, error)
}

type RegisterUserRequest struct {
	Email    string `json:"email" validate:"required"`
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}
