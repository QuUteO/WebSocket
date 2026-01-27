package repository

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/QuUteO/video-communication/internal/model"
	"github.com/QuUteO/video-communication/pkg/db"
	"github.com/gofrs/uuid"
	"github.com/jackc/pgconn"
)

type Repository interface {
	Create(ctx context.Context, user *model.User) (uuid.UUID, error)
	FindAll(ctx context.Context) ([]model.DTOResponse, error)
	FindByID(ctx context.Context, id string) (*model.User, error)
	Update(ctx context.Context, user *model.User) error
	Delete(ctx context.Context, id string) error

	SaveMsg(ctx context.Context, msg model.Message) error
	GetMessagesByChannel(ctx context.Context, channel string) ([]model.Message, error)
}

type repository struct {
	client postgres.Client
	logger *slog.Logger
}

func (r *repository) GetMessagesByChannel(ctx context.Context, channel string) ([]model.Message, error) {
	const op = "./internal/server/repository/GetMessagesByChannel"
	r.logger.With("op: ", op)

	q := `SELECT text, channel, username, created_at 
		FROM message 
		WHERE channel = $1
		ORDER BY created_at DESC
		LIMIT 10
`

	var messages []model.Message

	rows, err := r.client.Query(ctx, q, channel)
	if err != nil {
		r.logger.Error("error querying message: ", slog.String("error", err.Error()))
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var msg model.Message

		if err := rows.Scan(
			&msg.User,
			&msg.Msg,
			&msg.Channel,
			&msg.Time,
		); err != nil {
			r.logger.Error("error scanning message: ", slog.String("error", err.Error()))
			return nil, err
		}

		messages = append(messages, msg)
	}

	return messages, nil
}

func (r *repository) SaveMsg(ctx context.Context, msg model.Message) error {
	const op = "./internal/server/repository/SaveMsg"
	r.logger.With("op:", op)

	q := `
		INSERT INTO message (text, channel, username, created_at)
		VALUES ($1, $2, $3, $4)
	`

	if _, err := r.client.Exec(ctx, q, msg.Msg, msg.Channel, msg.User, msg.Time); err != nil {
		r.logger.Info("Error saving message", slog.String("error", err.Error()))
		return fmt.Errorf("%w: %s", err, msg)
	}

	return nil
}

func (r *repository) Create(ctx context.Context, user *model.User) (uuid.UUID, error) {
	const op = "./internal/server/repository/Create"
	r.logger.With("op:", op)

	q := `
		INSERT INTO users (email, password)
		VALUES ($1, $2)
		RETURNING id
	`

	if err := r.client.QueryRow(ctx, q, user.Email, user.Password).Scan(&user.Id); err != nil {
		var PGerr *pgconn.PgError
		if errors.As(err, &PGerr) {
			r.logger.Error(fmt.Sprintf("QueryRow failed: %s Code: %s Where: %s SQL State: %s", PGerr.Message, PGerr.Code, PGerr.Where, PGerr.SQLState()))
		}
		return uuid.UUID{}, err
	}

	return user.Id, nil
}

func (r *repository) FindAll(ctx context.Context) ([]model.DTOResponse, error) {
	const op = "./internal/server/repository/FindAll"
	r.logger.With("op:", op)

	q := `
	SELECT id, email FROM users
	`

	var users []model.DTOResponse
	rows, err := r.client.Query(ctx, q)
	if err != nil {
		r.logger.Error("Error querying users: ", slog.String("error", err.Error()))
		return nil, err
	}

	for rows.Next() {
		var user model.DTOResponse
		if err := rows.Scan(&user.ID, &user.Email); err != nil {
			r.logger.Error(err.Error())
			return nil, err
		}
		users = append(users, user)
	}

	return users, nil
}

func (r *repository) FindByID(ctx context.Context, id string) (*model.User, error) {
	const op = "./internal/server/repository/FindByID"
	r.logger.With("op:", op)

	q := `
	SELECT id, email FROM users WHERE id = $1
	`

	var user model.User
	if err := r.client.QueryRow(ctx, q, id).Scan(&user.Id, &user.Email); err != nil {
		r.logger.Info("Error querying user: ", slog.String("error", err.Error()))
		return nil, err
	}

	return &user, nil
}

func (r *repository) Update(ctx context.Context, user *model.User) error {
	const op = "./internal/server/repository/Update"
	r.logger.With("op:", op)

	q := `
	UPDATE users 
	SET email = $1, password = $2 
	WHERE id = $3
	`

	if _, err := r.client.Exec(ctx, q, user.Email, user.Password, user.Id); err != nil {
		r.logger.Error("Error updating user: ", slog.String("error", err.Error()))
		return err
	}
	return nil
}

func (r *repository) Delete(ctx context.Context, id string) error {
	const op = "./internal/server/repository/Delete"
	r.logger.With("op:", op)

	q := `
		DELETE FROM users WHERE id = $1	
	`

	if _, err := r.client.Exec(ctx, q, id); err != nil {
		r.logger.Error("Error deleting user: ", slog.String("error", err.Error()))
		return err
	}

	return nil
}

func NewRepository(client postgres.Client, logger *slog.Logger) Repository {
	return &repository{
		client: client,
		logger: logger,
	}
}
