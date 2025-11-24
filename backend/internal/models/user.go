package models

import "time"

type User struct {
	ID           uint64     `json:"id"`
	Username     string     `json:"username"`
	Email        string     `json:"email"`
	PasswordHash string     `json:"-"` // Never send to client
	AvatarURL    *string    `json:"avatarUrl,omitempty"`
	AboutText    *string    `json:"aboutText,omitempty"`
	CreatedAt    time.Time  `json:"createdAt"`
	UpdatedAt    time.Time  `json:"updatedAt"`
}

type UserRegistrationRequest struct {
	Username string `json:"username" validate:"required,min=3,max=50"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
}

type UserLoginRequest struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type UserUpdateRequest struct {
	AvatarURL *string `json:"avatarUrl,omitempty"`
	AboutText *string `json:"aboutText,omitempty"`
}

type AuthResponse struct {
	AccessToken string `json:"accessToken"`
	User        *User  `json:"user"`
}

