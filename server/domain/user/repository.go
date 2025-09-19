package user

import (
	"context"
)

type Saver interface {
	Save(ctx context.Context, u *User) (int64, error)
}

type Getter interface {
	GetByEmail(ctx context.Context, email string) (*User, error)
	GetByID(ctx context.Context, id int64) (*User, error)
	GetByUsername(ctx context.Context, username string) (*User, error)
}

type Repository interface {
	Saver
	Getter
}
