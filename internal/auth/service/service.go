package service

import (
	"context"
	"log/slog"
	"time"

	"github.com/QuUteO/video-communication/internal/auth/jwt"
	"github.com/QuUteO/video-communication/internal/auth/repository"
	"github.com/QuUteO/video-communication/internal/model"
	uuid2 "github.com/gofrs/uuid"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type Service interface {
	Register(ctx context.Context, req *model.AuthRequest) (string, error)
	Login(ctx context.Context, req *model.LoginRequest) (string, error)
}

type AuthService struct {
	repo   repository.Repository
	jwt    *jwt.Manager
	logger *slog.Logger
}

func (a *AuthService) Register(ctx context.Context, req *model.AuthRequest) (string, error) {
	const op = "internal/auth/service.Create"
	log := a.logger.With("op: ", op)

	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Error("Error hashing password", slog.Any("error", err))
		return "", err
	}

	user := &model.User{
		Id:        uuid2.UUID(uuid.New()),
		Email:     req.Email,
		Password:  string(hash),
		CreatedAt: time.Now(),
	}

	if err := a.repo.Register(ctx, user); err != nil {
		log.Error("Error registering user", slog.Any("error", err))
		return "", err
	}

	return a.jwt.Generate(user.Id.String())
}

func (a *AuthService) Login(ctx context.Context, req *model.LoginRequest) (string, error) {
	const op = "internal/auth/service.Create"
	log := a.logger.With("op: ", op)

	user, err := a.repo.FindByEmail(ctx, req.Email)
	if err != nil {
		log.Error("Error finding user", slog.Any("error", err))
		return "", err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		log.Error("Error comparing password", slog.Any("error", err))
		return "", err
	}

	return a.jwt.Generate(user.Id.String())
}

func NewAuthService(repo repository.Repository, jwt *jwt.Manager, logger *slog.Logger) Service {
	return &AuthService{
		repo:   repo,
		jwt:    jwt,
		logger: logger,
	}
}
