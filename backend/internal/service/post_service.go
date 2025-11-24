package service

import (
	"fmt"

	"github.com/forumapp/backend/internal/models"
	"github.com/forumapp/backend/internal/repository"
)

type PostService struct {
	postRepo   *repository.PostRepository
	threadRepo *repository.ThreadRepository
	forumRepo  *repository.ForumRepository
}

func NewPostService(postRepo *repository.PostRepository) *PostService {
	return &PostService{
		postRepo: postRepo,
	}
}

func (s *PostService) SetThreadRepo(threadRepo *repository.ThreadRepository) {
	s.threadRepo = threadRepo
}

func (s *PostService) SetForumRepo(forumRepo *repository.ForumRepository) {
	s.forumRepo = forumRepo
}

func (s *PostService) CreatePost(userID uint64, threadID uint64, req *models.PostCreateRequest) (*models.Post, error) {
	// Check if thread exists
	if s.threadRepo != nil {
		thread, err := s.threadRepo.FindByID(threadID)
		if err != nil {
			return nil, fmt.Errorf("failed to check thread: %w", err)
		}
		if thread == nil {
			return nil, fmt.Errorf("thread not found")
		}
		if thread.IsLocked {
			return nil, fmt.Errorf("thread is locked")
		}
	}

	// Get next post number
	nextPostNumber, err := s.postRepo.GetNextPostNumber(threadID)
	if err != nil {
		return nil, fmt.Errorf("failed to get next post number: %w", err)
	}

	post := &models.Post{
		ThreadID:   threadID,
		UserID:     userID,
		PostNumber: nextPostNumber,
		Content:    req.Content,
	}

	if err := s.postRepo.Create(post, req.Content); err != nil {
		return nil, fmt.Errorf("failed to create post: %w", err)
	}

	// Increment thread reply count
	if s.threadRepo != nil {
		s.threadRepo.IncrementReplyCount(threadID)
	}

	// Increment forum post count
	if s.forumRepo != nil && s.threadRepo != nil {
		thread, _ := s.threadRepo.FindByID(threadID)
		if thread != nil {
			s.forumRepo.IncrementPostCount(thread.ForumID)
		}
	}

	// Fetch full post with author
	fullPost, err := s.postRepo.FindByID(post.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch created post: %w", err)
	}

	return fullPost, nil
}

func (s *PostService) GetPostsByThreadID(threadID uint64, page, limit int) (*models.PostListResponse, error) {
	posts, totalCount, err := s.postRepo.FindByThreadID(threadID, page, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to get posts: %w", err)
	}

	totalPages := (totalCount + limit - 1) / limit

	return &models.PostListResponse{
		Posts: posts,
		Pagination: models.Pagination{
			Page:       page,
			Limit:      limit,
			TotalPages: totalPages,
			TotalItems: totalCount,
		},
	}, nil
}

func (s *PostService) UpdatePost(postID, userID uint64, req *models.PostUpdateRequest) (*models.Post, error) {
	post, err := s.postRepo.FindByID(postID)
	if err != nil {
		return nil, fmt.Errorf("failed to find post: %w", err)
	}
	if post == nil {
		return nil, fmt.Errorf("post not found")
	}

	// Check if user owns the post
	if post.UserID != userID {
		return nil, fmt.Errorf("unauthorized")
	}

	post.Content = req.Content
	if err := s.postRepo.Update(post); err != nil {
		return nil, fmt.Errorf("failed to update post: %w", err)
	}

	return post, nil
}

func (s *PostService) DeletePost(postID, userID uint64) error {
	post, err := s.postRepo.FindByID(postID)
	if err != nil {
		return fmt.Errorf("failed to find post: %w", err)
	}
	if post == nil {
		return fmt.Errorf("post not found")
	}

	// Check if user owns the post
	if post.UserID != userID {
		return fmt.Errorf("unauthorized")
	}

	// Don't allow deleting the first post (post_number = 1)
	if post.PostNumber == 1 {
		return fmt.Errorf("cannot delete the first post of a thread")
	}

	if err := s.postRepo.Delete(postID); err != nil {
		return fmt.Errorf("failed to delete post: %w", err)
	}

	return nil
}

