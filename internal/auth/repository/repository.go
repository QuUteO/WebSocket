package authrepository

import (
	"context"
	"log/slog"

	"github.com/QuUteO/video-communication/internal/model"
	postgres "github.com/QuUteO/video-communication/pkg/db"
)

type Repository interface {
	Register(ctx context.Context, user *model.User) error
	FindByEmail(ctx context.Context, email string) (*model.User, error)
}

type repository struct {
	db     postgres.Client
	logger *slog.Logger
}

func (r *repository) Register(ctx context.Context, user *model.User) error {
	const op = "./internal/auth/repository.Register"
	log := r.logger.With("op: ", op)

	q := `INSERT INTO users (id, email, password, created_at) VALUES ($1, $2, $3, $4)`

	if err := r.db.QueryRow(ctx, q, user.Id, user.Email, user.Password, user.CreatedAt).Scan(&user.Id); err != nil {
		log.Error("Error to insert user", slog.Any("err", err))
		return err
	}

	return nil
}

func (r *repository) FindByEmail(ctx context.Context, email string) (*model.User, error) {
	const op = "./internal/auth/repository.FindByEmail"
	log := r.logger.With("op: ", op)

	q := `
		SELECT id, email, password, created_at 
		FROM users 
		WHERE email = $1
			`

	var user model.User
	if err := r.db.QueryRow(ctx, q, email).Scan(&user.Id, &user.Email, &user.Password); err != nil {
		log.Error("Error to find user by email", slog.Any("err", err))
		return nil, err
	}

	return &user, nil
}

func New(db postgres.Client, logger *slog.Logger) Repository {
	return &repository{
		db:     db,
		logger: logger,
	}
}
