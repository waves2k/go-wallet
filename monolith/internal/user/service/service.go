package service

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/waves2k/go-wallet/monolith/internal/user/model"
	"github.com/waves2k/go-wallet/monolith/internal/user/repository"
)

type UserService interface {
	Register(ctx context.Context, req model.CreateUserRequest) (*model.User, error)
	GetProfile(ctx context.Context, id string) (*model.User, error)
	UpdateProfile(ctx context.Context, id string, req model.UpdateUserRequest) (*model.User, error)
}

type userService struct {
	repo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) UserService {
	return &userService{repo: repo}
}

func (s *userService) Register(ctx context.Context, req model.CreateUserRequest) (*model.User, error) {
	existing, _ := s.repo.FindByEmail(ctx, req.Email)
	if existing != nil {
		return nil, errors.New("email already registerd")
	}

	user := model.User{
		ID:           uuid.New().String(),
		FullName:     req.Username,
		Email:        req.Email,
		PasswordHash: req.Password,
	}

	if err := s.repo.Create(ctx, &user); err != nil {
		return nil, err
	}

	return &user, nil

}

func (s *userService) GetProfile(ctx context.Context, id string) (*model.User, error) {
	return s.repo.FindByID(ctx, id)
}

func (s *userService) UpdateProfile(ctx context.Context, id string, req model.UpdateUserRequest) (*model.User, error) {
	user, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	user.FullName = req.FullName
	if err := s.repo.Update(ctx, user); err != nil {
		return nil, err
	}

	return user, nil
}
