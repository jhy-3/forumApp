package service

import (
	"fmt"

	"github.com/forumapp/backend/internal/models"
	"github.com/forumapp/backend/internal/repository"
)

type ForumService struct {
	forumRepo *repository.ForumRepository
}

func NewForumService(forumRepo *repository.ForumRepository) *ForumService {
	return &ForumService{
		forumRepo: forumRepo,
	}
}

func (s *ForumService) GetAllForums() ([]models.Forum, error) {
	forums, err := s.forumRepo.FindAll()
	if err != nil {
		return nil, fmt.Errorf("failed to get forums: %w", err)
	}
	return forums, nil
}

func (s *ForumService) GetForumBySlug(slug string) (*models.Forum, error) {
	forum, err := s.forumRepo.FindBySlug(slug)
	if err != nil {
		return nil, fmt.Errorf("failed to get forum: %w", err)
	}
	if forum == nil {
		return nil, fmt.Errorf("forum not found")
	}
	return forum, nil
}

func (s *ForumService) CreateForum(req *models.ForumCreateRequest) (*models.Forum, error) {
	forum := &models.Forum{
		Name:        req.Name,
		Description: req.Description,
		Slug:        req.Slug,
	}

	if err := s.forumRepo.Create(forum); err != nil {
		return nil, fmt.Errorf("failed to create forum: %w", err)
	}

	return forum, nil
}

