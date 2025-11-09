package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"time"

	"github.com/abhishek622/moviedock/rating/pkg/model"
	_ "github.com/jackc/pgx/v5/stdlib"
)

// Repository defines a MySQL-based rating repository.
type Repository struct {
	db *sql.DB
}

// New creates a new MySQL-based rating repository.
func New() (*Repository, error) {
	// Read connection info from env with sensible defaults
	user := os.Getenv("POSTGRES_USER")
	pass := os.Getenv("POSTGRES_PASSWORD")
	dbname := os.Getenv("POSTGRES_DB")
	host := os.Getenv("POSTGRES_HOST")
	port := os.Getenv("POSTGRES_PORT")
	sslmode := os.Getenv("PGSSLMODE")

	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?search_path=rating_service&sslmode=%s", user, pass, host, port, dbname, sslmode)

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

// Get retrieves all ratings for a given record.
func (r *Repository) Get(ctx context.Context, recordID model.RecordID, recordType model.RecordType) ([]model.Rating, error) {
	rows, err := r.db.QueryContext(ctx, "SELECT user_id, value FROM ratings WHERE record_id = $1 AND record_type = $2", recordID, recordType)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var res []model.Rating
	for rows.Next() {
		var userID string
		var value int32
		if err := rows.Scan(&userID, &value); err != nil {
			return nil, err
		}
		res = append(res, model.Rating{
			UserID: model.UserID(userID),
			Value:  model.RatingValue(value),
		})
	}
	return res, nil
}

// Put adds a rating for a given record.
func (r *Repository) Put(ctx context.Context, recordID model.RecordID, recordType model.RecordType, rating *model.Rating) error {
	_, err := r.db.ExecContext(ctx, `INSERT INTO ratings (record_id, record_type, user_id, value) VALUES ($1, $2, $3, $4)
	ON CONFLICT (record_id, user_id) DO UPDATE
	SET value = EXCLUDED.value, record_type = EXCLUDED.record_type`,
		recordID, recordType, rating.UserID, rating.Value)
	return err
}

// Delete rating by user_id
func (r *Repository) Delete(ctx context.Context, userID model.UserID) error {
	_, err := r.db.ExecContext(ctx, "DELETE FROM ratings WHERE user_id = $1", userID)
	return err
}
