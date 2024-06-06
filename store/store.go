package store

import (
	"context"

	"github.com/hanzalahimran7/go-auth/model"
)

type DatabaseStore interface {
	Signup(ctx context.Context, user *model.User) error
	Login(ctx context.Context, username string, password string) (*model.User, error)
	GetUser(ctx context.Context, username string) (*model.User, error)
	DeleteUser(ctx context.Context, username string) error
	UpdateUser(ctx context.Context, user *model.User) error
	FindByEmail(ctx context.Context, email string) error
}
