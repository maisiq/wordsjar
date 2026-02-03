package user

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/maisiq/go-words-jar/internal/db"
	errx "github.com/maisiq/go-words-jar/internal/errors"
	"github.com/maisiq/go-words-jar/internal/models"
)

type IUserRepository interface {
	AddUser(ctx context.Context, user models.User) error
	User(ctx context.Context, username string) (models.User, error)
}

type UserRepository struct {
	db *db.DBClient
}

func NewUserRepository(db *db.DBClient) *UserRepository {
	return &UserRepository{
		db: db,
	}
}

func (r *UserRepository) AddUser(ctx context.Context, user models.User) error {
	builder := sq.Insert("users").
		Columns("id", "username", "password").
		Values(user.ID, user.Username, user.HashedPassword).
		PlaceholderFormat(sq.Dollar)
	query, args, err := builder.ToSql()

	if err != nil {
		return fmt.Errorf("failed to create query: %w", err)
	}
	_, execErr := r.db.DB.ExecContext(ctx, query, args...)
	if execErr != nil {
		var pgErr *pgconn.PgError
		if errors.As(execErr, &pgErr) {
			switch pgErr.Code {
			case "23505":
				return errx.ErrUserAlreadyExists
			}
		}
		return fmt.Errorf("failed to insert user: %w", execErr)
	}
	return nil
}

func (r *UserRepository) User(ctx context.Context, username string) (models.User, error) {
	builder := sq.Select("id", "username", "password", "is_admin").
		From("users").
		Where("username = $1", username).
		PlaceholderFormat(sq.Dollar)
	query, args, err := builder.ToSql()

	if err != nil {
		return models.User{}, fmt.Errorf("failed to create query: %w", err)
	}
	res := r.db.DB.QueryRowxContext(ctx, query, args...)

	var user models.User
	errScan := res.StructScan(&user)

	if errScan != nil {
		if errors.Is(errScan, sql.ErrNoRows) {
			return models.User{}, errx.ErrUserNotFound
		}
		return models.User{}, errScan
	}
	return user, nil
}
