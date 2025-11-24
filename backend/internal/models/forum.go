package models

type Forum struct {
	ID          uint   `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Slug        string `json:"slug"`
	ThreadCount uint   `json:"threadCount"`
	PostCount   uint   `json:"postCount"`
}

type ForumCreateRequest struct {
	Name        string `json:"name" validate:"required,min=3,max=100"`
	Description string `json:"description"`
	Slug        string `json:"slug" validate:"required,min=3,max=100"`
}

