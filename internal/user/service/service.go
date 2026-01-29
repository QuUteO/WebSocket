package service

import (
	"context"
	"log/slog"
	"time"

	"github.com/QuUteO/video-communication/internal/model"
	"github.com/QuUteO/video-communication/internal/user/repository"
	"github.com/gofrs/uuid"
	"golang.org/x/crypto/bcrypt"
)

type Service interface {
	CreateUser(ctx context.Context, email string, password string) (uuid.UUID, error)
	DeleteUser(ctx context.Context, id string) error
	UpdateUser(ctx context.Context, id string, email string, password string) error
	FindAllUser(ctx context.Context) ([]model.DTOResponse, error)
	FindUserById(ctx context.Context, id string) (*model.User, error)

	SaveMsg(ctx context.Context, msg model.Message) error
	GetMessageByChannel(ctx context.Context, channel string) ([]model.Message, error)
}

type service struct {
	repository repository.Repository
	logger     *slog.Logger
}

func (s *service) CreateUser(ctx context.Context, email string, password string) (uuid.UUID, error) {
	const op = "./internal/server/service.CreateUser"
	log := s.logger.With("op:", op)

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return uuid.Nil, err
	}

	user := &model.User{
		Email:    email,
		Password: string(hash),
	}

	id, err := s.repository.Create(ctx, user)
	if err != nil {
		log.Error("Failed to create server", "error:", err, "email:", email)
		return uuid.Nil, err
	}

	log.Info("Created server")
	return id, nil
}

func (s *service) DeleteUser(ctx context.Context, id string) error {
	const op = "./internal/server/service.DeleteUser"
	log := s.logger.With("op:", op)

	err := s.repository.Delete(ctx, id)
	if err != nil {
		log.Error("Failed to delete server", "error:", err, "id", id)
		return err
	}

	log.Info("Deleted server")
	return nil
}

func (s *service) UpdateUser(ctx context.Context, id string, email string, password string) error {
	const op = "./internal/server/service.UpdateUser"
	log := s.logger.With("op:", op)

	user, err := s.repository.FindByID(ctx, id)
	if err != nil {
		log.Error("Failed to find server", "error:", err, "id", id)
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

	log.Info("Updated server")
	return nil
}

func (s *service) FindAllUser(ctx context.Context) ([]model.DTOResponse, error) {
	const op = "./internal/server/service.FindAllUser"
	log := s.logger.With("op:", op)

	users, err := s.repository.FindAll(ctx)
	if err != nil {
		log.Error("Failed to find all users", "error:", err, "users", users)
		return nil, err
	}

	log.Info("Found all users")
	return users, nil
}

func (s *service) FindUserById(ctx context.Context, id string) (*model.User, error) {
	const op = "./internal/server/service.FindUserById"
	log := s.logger.With("op:", op)

	user, err := s.repository.FindByID(ctx, id)
	if err != nil {
		log.Error("Failed to find server", "error:", err, "id", id)
	}

	log.Info("Found server")
	return user, nil
}

func (s *service) GetMessageByChannel(ctx context.Context, channel string) ([]model.Message, error) {
	const op = "./internal/user/service.GetMessageByChannel"
	log := s.logger.With("op: ", op)

	message, err := s.repository.GetMessagesByChannel(ctx, channel)
	if err != nil {
		log.Error("error inserting insertMessage: ", slog.String("error", err.Error()))
		return nil, err
	}

	return message, nil
}

func (s *service) SaveMsg(ctx context.Context, msg model.Message) error {
	const op = "./internal/server/repository/SaveMsg"
	log := s.logger.With("op: ", op)

	if msg.Time.IsZero() {
		msg.Time = time.Now()
	}

	if err := s.repository.SaveMsg(ctx, msg); err != nil {
		log.Error("Error saving message: ", slog.Any("err", err))
		return err
	}

	return nil
}

func NewService(repository repository.Repository, logger *slog.Logger) Service {
	return &service{
		repository: repository,
		logger:     logger,
	}
}
