package repository

import (
	"database/sql"
	"fmt"

	"github.com/forumapp/backend/internal/models"
)

type ForumRepository struct {
	db *sql.DB
}

func NewForumRepository(db *sql.DB) *ForumRepository {
	return &ForumRepository{db: db}
}

func (r *ForumRepository) Create(forum *models.Forum) error {
	query := `
		INSERT INTO forums (name, description, slug, thread_count, post_count)
		VALUES (?, ?, ?, 0, 0)
	`
	result, err := r.db.Exec(query, forum.Name, forum.Description, forum.Slug)
	if err != nil {
		return fmt.Errorf("failed to create forum: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("failed to get last insert id: %w", err)
	}

	forum.ID = uint(id)
	return nil
}

func (r *ForumRepository) FindAll() ([]models.Forum, error) {
	query := `
		SELECT id, name, description, slug, thread_count, post_count
		FROM forums
		ORDER BY id ASC
	`
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to query forums: %w", err)
	}
	defer rows.Close()

	var forums []models.Forum
	for rows.Next() {
		var forum models.Forum
		err := rows.Scan(
			&forum.ID,
			&forum.Name,
			&forum.Description,
			&forum.Slug,
			&forum.ThreadCount,
			&forum.PostCount,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan forum: %w", err)
		}
		forums = append(forums, forum)
	}

	return forums, nil
}

func (r *ForumRepository) FindByID(id uint) (*models.Forum, error) {
	query := `
		SELECT id, name, description, slug, thread_count, post_count
		FROM forums
		WHERE id = ?
	`
	forum := &models.Forum{}
	err := r.db.QueryRow(query, id).Scan(
		&forum.ID,
		&forum.Name,
		&forum.Description,
		&forum.Slug,
		&forum.ThreadCount,
		&forum.PostCount,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to find forum by id: %w", err)
	}
	return forum, nil
}

func (r *ForumRepository) FindBySlug(slug string) (*models.Forum, error) {
	query := `
		SELECT id, name, description, slug, thread_count, post_count
		FROM forums
		WHERE slug = ?
	`
	forum := &models.Forum{}
	err := r.db.QueryRow(query, slug).Scan(
		&forum.ID,
		&forum.Name,
		&forum.Description,
		&forum.Slug,
		&forum.ThreadCount,
		&forum.PostCount,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to find forum by slug: %w", err)
	}
	return forum, nil
}

func (r *ForumRepository) IncrementThreadCount(forumID uint) error {
	query := `UPDATE forums SET thread_count = thread_count + 1 WHERE id = ?`
	_, err := r.db.Exec(query, forumID)
	return err
}

func (r *ForumRepository) IncrementPostCount(forumID uint) error {
	query := `UPDATE forums SET post_count = post_count + 1 WHERE id = ?`
	_, err := r.db.Exec(query, forumID)
	return err
}

