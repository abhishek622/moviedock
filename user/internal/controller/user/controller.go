package user

import (
	"context"
	"errors"

	"github.com/abhishek622/moviedock/user/pkg/model"
	"golang.org/x/crypto/bcrypt"
)

// ErrNotFound is returned when a requested record is not found.
var ErrNotFound = errors.New("not found")

func HashPassword(password string) (string, error) {
	HashPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return string(HashPassword), nil

}

type userRepository interface {
	RegisterUser(ctx context.Context, user *model.User) (*model.User, error)
	LoginUser(ctx context.Context, user *model.User) (*model.User, error)
	// LogoutUser(ctx context.Context, user *model.User) (*model.User, error)
	// RefreshToken(ctx context.Context, user *model.User) (*model.User, error)
}

type Controller struct {
	repo userRepository
}

func New(repo userRepository) *Controller {
	return &Controller{repo}
}

func (c *Controller) RegisterUser(ctx context.Context, user *model.User) (*model.User, error) {
	return c.repo.RegisterUser(ctx, user)
}

func (c *Controller) LoginUser(ctx context.Context, user *model.User) (*model.User, error) {
	return c.repo.LoginUser(ctx, user)
}

// func (c *Controller) LogoutUser(ctx context.Context, user *model.User) (*model.User, error) {
// 	return c.repo.LogoutUser(ctx, user)
// }

// func (c *Controller) RefreshToken(ctx context.Context, user *model.User) (*model.User, error) {
// 	return c.repo.RefreshToken(ctx, user)
// }
