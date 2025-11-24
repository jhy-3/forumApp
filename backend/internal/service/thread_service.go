package service

import (
	"fmt"
	"strings"

	"github.com/forumapp/backend/internal/models"
	"github.com/forumapp/backend/internal/repository"
)

type ThreadService struct {
	threadRepo *repository.ThreadRepository
	forumRepo  *repository.ForumRepository
	postRepo   *repository.PostRepository
}

func NewThreadService(threadRepo *repository.ThreadRepository) *ThreadService {
	return &ThreadService{
		threadRepo: threadRepo,
	}
}

func (s *ThreadService) SetForumRepo(forumRepo *repository.ForumRepository) {
	s.forumRepo = forumRepo
}

func (s *ThreadService) SetPostRepo(postRepo *repository.PostRepository) {
	s.postRepo = postRepo
}

func (s *ThreadService) CreateThread(userID uint64, req *models.ThreadCreateRequest) (*models.Thread, error) {
	// Check if forum exists
	if s.forumRepo != nil {
		forum, err := s.forumRepo.FindByID(req.ForumID)
		if err != nil {
			return nil, fmt.Errorf("failed to check forum: %w", err)
		}
		if forum == nil {
			return nil, fmt.Errorf("forum not found")
		}
	}

	// Generate slug from title
	slug := generateSlug(req.Title)

	thread := &models.Thread{
		ForumID:  req.ForumID,
		UserID:   userID,
		Title:    req.Title,
		Slug:     slug,
		IsLocked: false,
		IsPinned: false,
	}

	if err := s.threadRepo.Create(thread); err != nil {
		return nil, fmt.Errorf("failed to create thread: %w", err)
	}

	// Create first post (original content)
	if s.postRepo != nil {
		post := &models.Post{
			ThreadID:   thread.ID,
			UserID:     userID,
			PostNumber: 1,
			Content:    req.Content,
		}
		if err := s.postRepo.Create(post, req.Content); err != nil {
			return nil, fmt.Errorf("failed to create first post: %w", err)
		}
	}

	// Increment forum thread count
	if s.forumRepo != nil {
		s.forumRepo.IncrementThreadCount(req.ForumID)
		s.forumRepo.IncrementPostCount(req.ForumID)
	}

	// Fetch full thread with relations
	fullThread, err := s.threadRepo.FindByID(thread.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch created thread: %w", err)
	}

	return fullThread, nil
}

func (s *ThreadService) GetThreadByID(id uint64) (*models.Thread, error) {
	thread, err := s.threadRepo.FindByID(id)
	if err != nil {
		return nil, fmt.Errorf("failed to get thread: %w", err)
	}
	if thread == nil {
		return nil, fmt.Errorf("thread not found")
	}

	// Increment view count
	s.threadRepo.IncrementViewCount(id)

	return thread, nil
}

func (s *ThreadService) GetThreadsByForumSlug(slug string, page, limit int) (*models.ThreadListResponse, error) {
	threads, totalCount, err := s.threadRepo.FindByForumSlug(slug, page, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to get threads: %w", err)
	}

	totalPages := (totalCount + limit - 1) / limit

	return &models.ThreadListResponse{
		Threads: threads,
		Pagination: models.Pagination{
			Page:       page,
			Limit:      limit,
			TotalPages: totalPages,
			TotalItems: totalCount,
		},
	}, nil
}

func generateSlug(title string) string {
	slug := strings.ToLower(title)
	slug = strings.ReplaceAll(slug, " ", "-")
	// Remove special characters (simple version)
	var result strings.Builder
	for _, r := range slug {
		if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') || r == '-' {
			result.WriteRune(r)
		}
	}
	return result.String()
}

