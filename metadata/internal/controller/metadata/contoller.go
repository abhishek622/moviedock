package metadata

import (
	"context"
	"errors"
	"log"

	"github.com/abhishek622/moviedock/metadata/pkg/model"
)

// ErrNotFound is returned when a requested record is not found.
var ErrNotFound = errors.New("not found")

type metadataRepository interface {
	Get(ctx context.Context, id int32) (*model.Metadata, error)
	Put(ctx context.Context, id int32, m *model.Metadata) error
	Create(ctx context.Context, m *model.Metadata) (*model.Metadata, error)
	Delete(ctx context.Context, id int32) error
	List(ctx context.Context, limit, offset int) ([]*model.Metadata, error)
}

// Controller defines a metadata service controller.
type Controller struct {
	repo metadataRepository
}

// New creates a metadata service controller.
func New(repo metadataRepository) *Controller {
	return &Controller{repo}
}

// Get returns movie metadata by id.
func (c *Controller) Get(ctx context.Context, id int32) (*model.Metadata, error) {
	res, err := c.repo.Get(ctx, id)
	if err != nil {
		log.Printf("Failed to get metadata: %v", err)
		return nil, err
	}
	return res, err
}

// Create creates new movie metadata.
func (c *Controller) Create(ctx context.Context, metadata *model.Metadata) (*model.Metadata, error) {
	// Add validation if needed
	res, err := c.repo.Create(ctx, metadata)
	if err != nil {
		log.Printf("Failed to create metadata: %v", err)
		return nil, err
	}

	// Get the created metadata to return
	return res, nil
}

// Update updates movie metadata.
func (c *Controller) Update(ctx context.Context, id int32, metadata *model.Metadata) (*model.Metadata, error) {
	// Check if exists
	if _, err := c.Get(ctx, id); err != nil {
		return nil, err
	}

	// Update
	err := c.repo.Put(ctx, id, metadata)
	if err != nil {
		log.Printf("Failed to update metadata: %v", err)
		return nil, err
	}

	return c.Get(ctx, id)
}

// Delete deletes movie metadata.
func (c *Controller) Delete(ctx context.Context, id int32) error {
	return c.repo.Delete(ctx, id)
}

// List returns all metadata.
func (c *Controller) List(ctx context.Context, limit, offset int) ([]*model.Metadata, error) {
	return c.repo.List(ctx, limit, offset)
}
