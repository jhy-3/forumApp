package repository

import (
	"database/sql"
	"fmt"

	"github.com/forumapp/backend/internal/models"
)

type ThreadRepository struct {
	db *sql.DB
}

func NewThreadRepository(db *sql.DB) *ThreadRepository {
	return &ThreadRepository{db: db}
}

func (r *ThreadRepository) Create(thread *models.Thread) error {
	query := `
		INSERT INTO threads (forum_id, user_id, title, slug, is_locked, is_pinned, view_count, reply_count, created_at)
		VALUES (?, ?, ?, ?, ?, ?, 0, 0, NOW())
	`
	result, err := r.db.Exec(query, thread.ForumID, thread.UserID, thread.Title, thread.Slug, thread.IsLocked, thread.IsPinned)
	if err != nil {
		return fmt.Errorf("failed to create thread: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("failed to get last insert id: %w", err)
	}

	thread.ID = uint64(id)
	return nil
}

func (r *ThreadRepository) FindByID(id uint64) (*models.Thread, error) {
	query := `
		SELECT t.id, t.forum_id, t.user_id, t.title, t.slug, t.is_locked, t.is_pinned, 
		       t.view_count, t.reply_count, t.last_post_at, t.created_at,
		       u.id, u.username, u.email, u.avatar_url, u.about_text, u.created_at, u.updated_at,
		       f.id, f.name, f.description, f.slug, f.thread_count, f.post_count
		FROM threads t
		LEFT JOIN users u ON t.user_id = u.id
		LEFT JOIN forums f ON t.forum_id = f.id
		WHERE t.id = ?
	`
	thread := &models.Thread{
		Author: &models.User{},
		Forum:  &models.Forum{},
	}
	err := r.db.QueryRow(query, id).Scan(
		&thread.ID,
		&thread.ForumID,
		&thread.UserID,
		&thread.Title,
		&thread.Slug,
		&thread.IsLocked,
		&thread.IsPinned,
		&thread.ViewCount,
		&thread.ReplyCount,
		&thread.LastPostAt,
		&thread.CreatedAt,
		&thread.Author.ID,
		&thread.Author.Username,
		&thread.Author.Email,
		&thread.Author.AvatarURL,
		&thread.Author.AboutText,
		&thread.Author.CreatedAt,
		&thread.Author.UpdatedAt,
		&thread.Forum.ID,
		&thread.Forum.Name,
		&thread.Forum.Description,
		&thread.Forum.Slug,
		&thread.Forum.ThreadCount,
		&thread.Forum.PostCount,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to find thread by id: %w", err)
	}
	// Clear password hash
	thread.Author.PasswordHash = ""
	return thread, nil
}

func (r *ThreadRepository) FindByForumSlug(slug string, page, limit int) ([]models.Thread, int, error) {
	offset := (page - 1) * limit

	// Get total count
	countQuery := `
		SELECT COUNT(*)
		FROM threads t
		JOIN forums f ON t.forum_id = f.id
		WHERE f.slug = ?
	`
	var totalCount int
	err := r.db.QueryRow(countQuery, slug).Scan(&totalCount)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count threads: %w", err)
	}

	// Get threads
	query := `
		SELECT t.id, t.forum_id, t.user_id, t.title, t.slug, t.is_locked, t.is_pinned, 
		       t.view_count, t.reply_count, t.last_post_at, t.created_at,
		       u.id, u.username, u.email, u.avatar_url, u.about_text, u.created_at, u.updated_at
		FROM threads t
		JOIN forums f ON t.forum_id = f.id
		LEFT JOIN users u ON t.user_id = u.id
		WHERE f.slug = ?
		ORDER BY t.is_pinned DESC, t.last_post_at DESC, t.created_at DESC
		LIMIT ? OFFSET ?
	`
	rows, err := r.db.Query(query, slug, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to query threads: %w", err)
	}
	defer rows.Close()

	var threads []models.Thread
	for rows.Next() {
		thread := models.Thread{
			Author: &models.User{},
		}
		err := rows.Scan(
			&thread.ID,
			&thread.ForumID,
			&thread.UserID,
			&thread.Title,
			&thread.Slug,
			&thread.IsLocked,
			&thread.IsPinned,
			&thread.ViewCount,
			&thread.ReplyCount,
			&thread.LastPostAt,
			&thread.CreatedAt,
			&thread.Author.ID,
			&thread.Author.Username,
			&thread.Author.Email,
			&thread.Author.AvatarURL,
			&thread.Author.AboutText,
			&thread.Author.CreatedAt,
			&thread.Author.UpdatedAt,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan thread: %w", err)
		}
		thread.Author.PasswordHash = ""
		threads = append(threads, thread)
	}

	return threads, totalCount, nil
}

func (r *ThreadRepository) IncrementViewCount(threadID uint64) error {
	query := `UPDATE threads SET view_count = view_count + 1 WHERE id = ?`
	_, err := r.db.Exec(query, threadID)
	return err
}

func (r *ThreadRepository) IncrementReplyCount(threadID uint64) error {
	query := `UPDATE threads SET reply_count = reply_count + 1, last_post_at = NOW() WHERE id = ?`
	_, err := r.db.Exec(query, threadID)
	return err
}

