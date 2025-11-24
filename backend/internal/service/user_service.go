package service

import (
	"fmt"

	"github.com/forumapp/backend/internal/models"
	"github.com/forumapp/backend/internal/repository"
)

type UserService struct {
	userRepo *repository.UserRepository
}

func NewUserService(userRepo *repository.UserRepository) *UserService {
	return &UserService{
		userRepo: userRepo,
	}
}

func (s *UserService) GetUserByID(id uint64) (*models.User, error) {
	user, err := s.userRepo.FindByID(id)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}
	if user == nil {
		return nil, fmt.Errorf("user not found")
	}
	// Clear sensitive data
	user.PasswordHash = ""
	return user, nil
}

func (s *UserService) UpdateUser(id uint64, req *models.UserUpdateRequest) (*models.User, error) {
	user, err := s.userRepo.FindByID(id)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}
	if user == nil {
		return nil, fmt.Errorf("user not found")
	}

	// Update fields
	if req.AvatarURL != nil {
		user.AvatarURL = req.AvatarURL
	}
	if req.AboutText != nil {
		user.AboutText = req.AboutText
	}

	if err := s.userRepo.Update(user); err != nil {
		return nil, fmt.Errorf("failed to update user: %w", err)
	}

	// Clear sensitive data
	user.PasswordHash = ""
	return user, nil
}

