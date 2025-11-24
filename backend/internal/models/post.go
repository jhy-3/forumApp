package models

import "time"

type Post struct {
	ID         uint64    `json:"id"`
	ThreadID   uint64    `json:"threadId"`
	UserID     uint64    `json:"userId"`
	PostNumber uint      `json:"postNumber"`
	Content    string    `json:"content"`
	CreatedAt  time.Time `json:"createdAt"`
	UpdatedAt  *time.Time `json:"updatedAt,omitempty"`
	Author     *User     `json:"author,omitempty"`
}

type PostCreateRequest struct {
	Content string `json:"content" validate:"required,min=1"`
}

type PostUpdateRequest struct {
	Content string `json:"content" validate:"required,min=1"`
}

type PostListResponse struct {
	Posts      []Post     `json:"posts"`
	Pagination Pagination `json:"pagination"`
}

type Pagination struct {
	Page       int `json:"page"`
	Limit      int `json:"limit"`
	TotalPages int `json:"totalPages"`
	TotalItems int `json:"totalItems"`
}

