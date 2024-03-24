package users

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v4"

	"gitlab.com/robotomize/gb-golang/homework/03-01-umanager/internal/database"
)

func New(userDB *pgx.Conn, timeout time.Duration) *Repository {
	return &Repository{userDB: userDB, timeout: timeout}
}

type Repository struct {
	userDB  *pgx.Conn
	timeout time.Duration
}

func (r *Repository) Create(ctx context.Context, req CreateUserReq) (database.User, error) {
    var u database.User

    ctx, cancel := context.WithTimeout(ctx, r.timeout)
    defer cancel()

    err := r.userDB.QueryRow(ctx, `
        INSERT INTO users (id, username, password, created_at, updated_at)
        VALUES ($1, $2, $3, $4, $5)
        RETURNING id, username, created_at, updated_at
    `, uuid.New(), req.Username, req.Password, time.Now(), time.Now()).Scan(&u.ID, &u.Username, &u.CreatedAt, &u.UpdatedAt)
    if err != nil {
        return u, err
    }

    return u, nil
}

func (r *Repository) FindByID(ctx context.Context, userID uuid.UUID) (database.User, error) {
    var u database.User

    ctx, cancel := context.WithTimeout(ctx, r.timeout)
    defer cancel()

    err := r.userDB.QueryRow(ctx, `
        SELECT id, username, created_at, updated_at
        FROM users
        WHERE id = $1
    `, userID).Scan(&u.ID, &u.Username, &u.CreatedAt, &u.UpdatedAt)
    if err != nil {
        return u, err
    }

    return u, nil
}

func (r *Repository) FindByUsername(ctx context.Context, username string) (database.User, error) {
    var u database.User

    ctx, cancel := context.WithTimeout(ctx, r.timeout)
    defer cancel()

    err := r.userDB.QueryRow(ctx, `
        SELECT id, username, created_at, updated_at
        FROM users
        WHERE username = $1
    `, username).Scan(&u.ID, &u.Username, &u.CreatedAt, &u.UpdatedAt)
    if err != nil {
        return u, err
    }

    return u, nil
}