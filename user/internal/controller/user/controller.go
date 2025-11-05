package user

import (
	"context"
	"errors"

	"github.com/abhishek622/moviedock/user/internal/repository"
	"github.com/abhishek622/moviedock/user/pkg/model"
)

// ErrNotFound is returned when a requested record is not found.
var ErrNotFound = errors.New("not found")

type userRepository interface {
	Get(ctx context.Context, id string) (*model.User, error)
	Put(ctx context.Context, id string, user *model.User) error
	Delete(ctx context.Context, id string) error
}

type Controller struct {
	repo userRepository
}

func New(repo userRepository) *Controller {
	return &Controller{repo: repo}
}

func (c *Controller) Get(ctx context.Context, id string) (*model.User, error) {
	res, err := c.repo.Get(ctx, id)
	if err != nil && errors.Is(err, repository.ErrNotFound) {
		return nil, ErrNotFound
	}
	return res, err
}

func (c *Controller) Put(ctx context.Context, user *model.User) error {
	return c.repo.Put(ctx, user.UserID, user)
}

func (c *Controller) Delete(ctx context.Context, id string) error {
	return c.repo.Delete(ctx, id)
}
