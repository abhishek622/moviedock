package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"time"

	"github.com/abhishek622/moviedock/metadata/internal/repository"
	"github.com/abhishek622/moviedock/metadata/pkg/model"
	_ "github.com/jackc/pgx/v5/stdlib"
)

// Repository defines a Postgres-backed movie metadata repository.
type Repository struct {
	db *sql.DB
}

// New creates a new Postgres-based repository.
func New() (*Repository, error) {
	// Read connection info from env with sensible defaults
	user := os.Getenv("POSTGRES_USER")
	pass := os.Getenv("POSTGRES_PASSWORD")
	dbname := os.Getenv("POSTGRES_DB")
	host := os.Getenv("POSTGRES_HOST")
	port := os.Getenv("POSTGRES_PORT")
	sslmode := os.Getenv("PGSSLMODE")

	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?search_path=metadata_service&sslmode=%s", user, pass, host, port, dbname, sslmode)

	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}

	// // set sensible pool config
	// db.SetMaxOpenConns(25)
	// db.SetMaxIdleConns(5)
	// db.SetConnMaxIdleTime(5 * time.Minute)
	// db.SetConnMaxLifetime(60 * time.Minute)

	// ping with timeout / retry in case DB isn't ready
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := db.PingContext(ctx); err != nil {
		return nil, err
	}

	return &Repository{db: db}, nil
}

// Get retrieves movie metadata for by movie id.
func (r *Repository) Get(ctx context.Context, id int32) (*model.Metadata, error) {
	var title, description, director string
	var runtime int32
	// Postgres uses $1 style placeholders
	row := r.db.QueryRowContext(ctx, "SELECT title, description, director, runtime FROM movies WHERE metadata_id = $1", id)
	if err := row.Scan(&title, &description, &director, &runtime); err != nil {
		if err == sql.ErrNoRows {
			return nil, repository.ErrNotFound
		}
		return nil, err
	}
	return &model.Metadata{
		MetadataID:  id,
		Title:       title,
		Description: description,
		Director:    director,
		Runtime:     runtime,
	}, nil
}

// Put adds or updates movie metadata for a given movie id.
func (r *Repository) Put(ctx context.Context, id int32, metadata *model.Metadata) error {
	_, err := r.db.ExecContext(ctx,
		`INSERT INTO movies (metadata_id, title, description, director, runtime)
         VALUES ($1, $2, $3, $4, $5)
         ON CONFLICT (metadata_id) DO UPDATE
           SET title = EXCLUDED.title,
               description = EXCLUDED.description,
               director = EXCLUDED.director,
               runtime = EXCLUDED.runtime`,
		id, metadata.Title, metadata.Description, metadata.Director, metadata.Runtime,
	)
	return err
}

func (r *Repository) Delete(ctx context.Context, id int32) error {
	_, err := r.db.ExecContext(ctx, "DELETE FROM movies WHERE metadata_id = $1", id)
	return err
}

func (r *Repository) Create(ctx context.Context, metadata *model.Metadata) (*model.Metadata, error) {
	err := r.db.QueryRowContext(ctx,
		`INSERT INTO movies (title, description, director, runtime) 
         VALUES ($1, $2, $3, $4) 
         RETURNING metadata_id`,
		metadata.Title, metadata.Description, metadata.Director, metadata.Runtime).
		Scan(&metadata.MetadataID)
	if err != nil {
		return nil, err
	}
	return metadata, nil
}

func (r *Repository) List(ctx context.Context, limit, offset int) ([]*model.Metadata, error) {
	rows, err := r.db.QueryContext(ctx,
		`SELECT metadata_id, title, description, director, runtime
         FROM movies
         LIMIT $1 OFFSET $2`,
		limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var metadatas []*model.Metadata
	for rows.Next() {
		var metadata model.Metadata
		if err := rows.Scan(&metadata.MetadataID, &metadata.Title, &metadata.Description, &metadata.Director, &metadata.Runtime); err != nil {
			return nil, err
		}
		metadatas = append(metadatas, &metadata)
	}
	return metadatas, nil
}
