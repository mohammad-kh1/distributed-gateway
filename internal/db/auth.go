package auth

import (
	"context"
	"database/sql"
	"log"

	_ "github.com/jackc/pgx/v5/stdlib"
)

type Repository struct {
	db *sql.DB
}

func NewRepository(conStr string) *Repository {
	db, err := sql.Open("pgx", conStr)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v", err)

	}
	return &Repository{db: db}
}

type UserData struct {
	UserID    string
	RateLimit int
	Tier      string
}

func (r *Repository) GetUserByToken(ctx context.Context, token string) (*UserData, error) {
	var user UserData
	query := `SELECT * FROM users WHERE token = $1`
	err := r.db.QueryRowContext(ctx, query, token).Scan(&user.UserID, &user.RateLimit, &user.Tier)
	if err != nil {
		return nil, err
	}
	return &user, nil
}
