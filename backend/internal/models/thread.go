package models

import "time"

type Thread struct {
	ID         uint64     `json:"id"`
	ForumID    uint       `json:"forumId"`
	UserID     uint64     `json:"userId"`
	Title      string     `json:"title"`
	Slug       string     `json:"slug"`
	IsLocked   bool       `json:"isLocked"`
	IsPinned   bool       `json:"isPinned"`
	ViewCount  uint       `json:"viewCount"`
	ReplyCount uint       `json:"replyCount"`
	LastPostAt *time.Time `json:"lastPostAt,omitempty"`
	CreatedAt  time.Time  `json:"createdAt"`
	Author     *User      `json:"author,omitempty"`
	Forum      *Forum     `json:"forum,omitempty"`
}

type ThreadCreateRequest struct {
	ForumID uint   `json:"forumId" validate:"required"`
	Title   string `json:"title" validate:"required,min=3,max=255"`
	Content string `json:"content" validate:"required,min=10"`
}

type ThreadListResponse struct {
	Threads    []Thread   `json:"threads"`
	Pagination Pagination `json:"pagination"`
}

