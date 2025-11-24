package repository

import (
	"database/sql"
	"fmt"

	"github.com/forumapp/backend/internal/models"
)

type PostRepository struct {
	db *sql.DB
}

func NewPostRepository(db *sql.DB) *PostRepository {
	return &PostRepository{db: db}
}

func (r *PostRepository) Create(post *models.Post, content string) error {
	tx, err := r.db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	// Insert post metadata
	query := `
		INSERT INTO posts (thread_id, user_id, post_number, created_at)
		VALUES (?, ?, ?, NOW())
	`
	result, err := tx.Exec(query, post.ThreadID, post.UserID, post.PostNumber)
	if err != nil {
		return fmt.Errorf("failed to create post: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("failed to get last insert id: %w", err)
	}
	post.ID = uint64(id)

	// Insert post content
	contentQuery := `
		INSERT INTO post_contents (post_id, content_markdown, content_html)
		VALUES (?, ?, ?)
	`
	_, err = tx.Exec(contentQuery, post.ID, content, content) // For now, HTML is same as markdown
	if err != nil {
		return fmt.Errorf("failed to create post content: %w", err)
	}

	return tx.Commit()
}

func (r *PostRepository) FindByThreadID(threadID uint64, page, limit int) ([]models.Post, int, error) {
	offset := (page - 1) * limit

	// Get total count
	countQuery := `SELECT COUNT(*) FROM posts WHERE thread_id = ?`
	var totalCount int
	err := r.db.QueryRow(countQuery, threadID).Scan(&totalCount)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count posts: %w", err)
	}

	// Get posts
	query := `
		SELECT p.id, p.thread_id, p.user_id, p.post_number, p.created_at, p.updated_at,
		       pc.content_markdown,
		       u.id, u.username, u.email, u.avatar_url, u.about_text, u.created_at, u.updated_at
		FROM posts p
		LEFT JOIN post_contents pc ON p.id = pc.post_id
		LEFT JOIN users u ON p.user_id = u.id
		WHERE p.thread_id = ?
		ORDER BY p.post_number ASC
		LIMIT ? OFFSET ?
	`
	rows, err := r.db.Query(query, threadID, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to query posts: %w", err)
	}
	defer rows.Close()

	var posts []models.Post
	for rows.Next() {
		post := models.Post{
			Author: &models.User{},
		}
		err := rows.Scan(
			&post.ID,
			&post.ThreadID,
			&post.UserID,
			&post.PostNumber,
			&post.CreatedAt,
			&post.UpdatedAt,
			&post.Content,
			&post.Author.ID,
			&post.Author.Username,
			&post.Author.Email,
			&post.Author.AvatarURL,
			&post.Author.AboutText,
			&post.Author.CreatedAt,
			&post.Author.UpdatedAt,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan post: %w", err)
		}
		post.Author.PasswordHash = ""
		posts = append(posts, post)
	}

	return posts, totalCount, nil
}

func (r *PostRepository) FindByID(id uint64) (*models.Post, error) {
	query := `
		SELECT p.id, p.thread_id, p.user_id, p.post_number, p.created_at, p.updated_at,
		       pc.content_markdown
		FROM posts p
		LEFT JOIN post_contents pc ON p.id = pc.post_id
		WHERE p.id = ?
	`
	post := &models.Post{}
	err := r.db.QueryRow(query, id).Scan(
		&post.ID,
		&post.ThreadID,
		&post.UserID,
		&post.PostNumber,
		&post.CreatedAt,
		&post.UpdatedAt,
		&post.Content,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to find post by id: %w", err)
	}
	return post, nil
}

func (r *PostRepository) Update(post *models.Post) error {
	tx, err := r.db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	// Update post metadata
	query := `UPDATE posts SET updated_at = NOW() WHERE id = ?`
	_, err = tx.Exec(query, post.ID)
	if err != nil {
		return fmt.Errorf("failed to update post: %w", err)
	}

	// Update post content
	contentQuery := `
		UPDATE post_contents 
		SET content_markdown = ?, content_html = ?
		WHERE post_id = ?
	`
	_, err = tx.Exec(contentQuery, post.Content, post.Content, post.ID)
	if err != nil {
		return fmt.Errorf("failed to update post content: %w", err)
	}

	return tx.Commit()
}

func (r *PostRepository) Delete(id uint64) error {
	query := `DELETE FROM posts WHERE id = ?`
	_, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete post: %w", err)
	}
	return nil
}

func (r *PostRepository) GetNextPostNumber(threadID uint64) (uint, error) {
	query := `SELECT COALESCE(MAX(post_number), 0) + 1 FROM posts WHERE thread_id = ?`
	var nextNumber uint
	err := r.db.QueryRow(query, threadID).Scan(&nextNumber)
	if err != nil {
		return 0, fmt.Errorf("failed to get next post number: %w", err)
	}
	return nextNumber, nil
}

