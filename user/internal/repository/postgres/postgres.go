package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"time"

	"github.com/abhishek622/moviedock/user/internal/repository"
	"github.com/abhishek622/moviedock/user/pkg/model"
	_ "github.com/jackc/pgx/v5/stdlib"
)

type Repository struct {
	db *sql.DB
}

func New() (*Repository, error) {
	// Read connection info from env with sensible defaults
	user := os.Getenv("POSTGRES_USER")
	pass := os.Getenv("POSTGRES_PASSWORD")
	dbname := os.Getenv("POSTGRES_DB")
	host := os.Getenv("POSTGRES_HOST")
	port := os.Getenv("POSTGRES_PORT")
	sslmode := os.Getenv("PGSSLMODE")

	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s", user, pass, host, port, dbname, sslmode)

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

func (r *Repository) Get(ctx context.Context, id string) (*model.User, error) {
	var user_id, full_name, email, password_hash, role, timezone string
	var created_at, updated_at, last_login time.Time
	var metadata map[string]interface{}
	var is_active bool

	row := r.db.QueryRowContext(ctx, "SELECT user_id, full_name, email, password_hash, role, is_active, timezone, last_login, metadata, created_at, updated_at WHERE user_id= $1", id)
	if err := row.Scan(&user_id, &full_name, &email, &password_hash, &role, &is_active, &timezone, &last_login, &metadata, &created_at, &updated_at); err != nil {
		if err == sql.ErrNoRows {
			return nil, repository.ErrNotFound
		}
		return nil, err
	}
	return &model.User{
		UserID:       user_id,
		FullName:     full_name,
		Email:        email,
		PasswordHash: password_hash,
		Role:         model.Role(role),
		IsActive:     is_active,
		Timezone:     timezone,
		LastLogin:    last_login,
		Metadata:     metadata,
		CreatedAt:    created_at,
		UpdatedAt:    updated_at,
	}, nil
}

func (r *Repository) Put(ctx context.Context, id string, user *model.User) error {
	_, err := r.db.ExecContext(ctx,
		`INSERT INTO users (user_id, full_name, email, password_hash, role, is_active, timezone, last_login, metadata, created_at, updated_at)
         VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
         ON CONFLICT (user_id) DO UPDATE
           SET full_name = EXCLUDED.full_name,
               email = EXCLUDED.email,
               password_hash = EXCLUDED.password_hash,
               role = EXCLUDED.role,
               is_active = EXCLUDED.is_active,
               timezone = EXCLUDED.timezone,
               last_login = EXCLUDED.last_login,
               metadata = EXCLUDED.metadata,
               updated_at = EXCLUDED.updated_at`,
		id, user.FullName, user.Email, user.PasswordHash, user.Role, user.IsActive, user.Timezone, user.LastLogin, user.Metadata, user.CreatedAt, user.UpdatedAt,
	)
	return err
}

func (r *Repository) Delete(ctx context.Context, id string) error {
	_, err := r.db.ExecContext(ctx, "DELETE FROM users WHERE user_id = $1", id)
	return err
}
