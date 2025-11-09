package postgres

import (
	"context"
	"database/sql"
	"errors"
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

	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?search_path=user_service&sslmode=%s", user, pass, host, port, dbname, sslmode)

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

func (r *Repository) GetByEmail(ctx context.Context, email string) (*model.User, error) {
	var user model.User
	var lastLogin, createdAt, updatedAt sql.NullTime
	var metadata map[string]interface{}

	query := `SELECT user_id, full_name, email, encrypted_password, role, is_active, timezone, 
             last_login, metadata, created_at, updated_at FROM users WHERE email = $1`

	err := r.db.QueryRowContext(ctx, query, email).Scan(
		&user.UserID, &user.FullName, &user.Email, &user.EncryptedPassword, &user.Role,
		&user.IsActive, &user.Timezone, &lastLogin, &metadata, &createdAt, &updatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, repository.ErrNotFound
		}
		return nil, fmt.Errorf("error getting user by email: %w", err)
	}

	if lastLogin.Valid {
		user.LastLogin = &lastLogin.Time
	}
	if createdAt.Valid {
		user.CreatedAt = createdAt.Time
	}
	if updatedAt.Valid {
		user.UpdatedAt = updatedAt.Time
	}

	return &user, nil
}

func (r *Repository) GetUser(ctx context.Context, user_id string) (*model.UserResponse, error) {
	var user model.UserResponse

	err := r.db.QueryRowContext(ctx, `SELECT user_id, full_name, email, role FROM users WHERE user_id = $1`,
		user_id).Scan(&user.UserId, &user.FullName, &user.Email, &user.Role)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, repository.ErrNotFound
		}
		return nil, fmt.Errorf("error getting user by id: %w", err)
	}

	return &user, nil
}

func (r *Repository) Put(ctx context.Context, id string, user *model.User) error {
	_, err := r.db.ExecContext(ctx,
		`INSERT INTO users (user_id, full_name, email, role, is_active, timezone, last_login, metadata, created_at, updated_at)
         VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
         ON CONFLICT (user_id) DO UPDATE
           SET full_name = EXCLUDED.full_name,
               email = EXCLUDED.email,
               role = EXCLUDED.role,
               is_active = EXCLUDED.is_active,
               timezone = EXCLUDED.timezone,
               last_login = EXCLUDED.last_login,
               metadata = EXCLUDED.metadata,
               updated_at = EXCLUDED.updated_at`,
		id, user.FullName, user.Email, user.Role, user.IsActive, user.Timezone, user.LastLogin, user.Metadata, user.CreatedAt, user.UpdatedAt,
	)
	return err
}

func (r *Repository) RegisterUser(ctx context.Context, user *model.User) (*model.User, error) {
	query := `
		INSERT INTO users (
			full_name, email, encrypted_password, role, is_active, timezone
		) VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING user_id, created_at, updated_at
	`

	err := r.db.QueryRowContext(
		ctx,
		query,
		user.FullName,
		user.Email,
		user.EncryptedPassword,
		user.Role,
		user.IsActive,
		user.Timezone,
	).Scan(&user.UserID, &user.CreatedAt, &user.UpdatedAt)

	if err != nil {
		return nil, fmt.Errorf("error creating user: %w", err)
	}

	return user, nil
}

func (r *Repository) Delete(ctx context.Context, id string) error {
	_, err := r.db.ExecContext(ctx, "DELETE FROM users WHERE user_id = $1", id)
	return err
}

// LoginUser retrieves a user by email for login purposes
func (r *Repository) LoginUser(ctx context.Context, user *model.User) (*model.User, error) {
	query := `
		SELECT user_id, full_name, email, encrypted_password, role, is_active, timezone, created_at, updated_at
		FROM users
		WHERE email = $1
	`

	var u model.User
	err := r.db.QueryRowContext(ctx, query, user.Email).Scan(
		&u.UserID,
		&u.FullName,
		&u.Email,
		&u.EncryptedPassword,
		&u.Role,
		&u.IsActive,
		&u.Timezone,
		&u.CreatedAt,
		&u.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("user not found")
		}
		return nil, fmt.Errorf("error retrieving user: %w", err)
	}

	return &u, nil
}
