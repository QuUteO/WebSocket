package service

import (
	"context"
	"log/slog"

	"github.com/QuUteO/video-communication/internal/model"
	"github.com/QuUteO/video-communication/internal/user/repository"
	"github.com/gofrs/uuid"
)

type Service interface {
	CreateUser(ctx context.Context, email string, password string) (uuid.UUID, error)
	DeleteUser(ctx context.Context, id string) error
	UpdateUser(ctx context.Context, id string, email string, password string) error
	FindAllUser(ctx context.Context) ([]model.DTOResponse, error)
	FindUserById(ctx context.Context, id string) (*model.User, error)
}

type service struct {
	repository repository.Repository
	logger     *slog.Logger
}

func (s *service) CreateUser(ctx context.Context, email string, password string) (uuid.UUID, error) {
	const op = "./internal/server/service.CreateUser"
	s.logger.With("op:", op)

	user := &model.User{
		Email:    email,
		Password: password,
	}

	id, err := s.repository.Create(ctx, user)
	if err != nil {
		s.logger.Error("Failed to create server", "error:", err, "email:", email)
		return uuid.Nil, err
	}

	s.logger.Info("Created server")
	return id, nil
}

func (s *service) DeleteUser(ctx context.Context, id string) error {
	const op = "./internal/server/service.DeleteUser"
	s.logger.With("op:", op)

	err := s.repository.Delete(ctx, id)
	if err != nil {
		s.logger.Error("Failed to delete server", "error:", err, "id", id)
		return err
	}

	s.logger.Info("Deleted server")
	return nil
}

func (s *service) UpdateUser(ctx context.Context, id string, email string, password string) error {
	const op = "./internal/server/service.UpdateUser"
	s.logger.With("op:", op)

	user, err := s.repository.FindByID(ctx, id)
	if err != nil {
		s.logger.Error("Failed to find server", "error:", err, "id", id)
		return err
	}

	s.logger.Info("Found server and updating server email:", email, "password:", password)
	user.Email = email
	user.Password = password

	err = s.repository.Update(ctx, user)
	if err != nil {
		s.logger.Error("Failed to update server", "error:", err, "id", id)
		return err
	}

	s.logger.Info("Updated server")
	return nil
}

func (s *service) FindAllUser(ctx context.Context) ([]model.DTOResponse, error) {
	const op = "./internal/server/service.FindAllUser"
	s.logger.With("op:", op)

	users, err := s.repository.FindAll(ctx)
	if err != nil {
		s.logger.Error("Failed to find all users", "error:", err, "users", users)
		return nil, err
	}

	s.logger.Info("Found all users")
	return users, nil
}

func (s *service) FindUserById(ctx context.Context, id string) (*model.User, error) {
	const op = "./internal/server/service.FindUserById"
	s.logger.With("op:", op)

	user, err := s.repository.FindByID(ctx, id)
	if err != nil {
		s.logger.Error("Failed to find server", "error:", err, "id", id)
	}

	s.logger.Info("Found server")
	return user, nil
}

func NewService(repository repository.Repository, logger *slog.Logger) Service {
	return &service{
		repository: repository,
		logger:     logger,
	}
}
