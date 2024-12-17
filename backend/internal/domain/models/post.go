package models

import "time"

type Post struct {
	ID        uint64    `db:"id"`
	Content   string    `db:"content"`
	UserID    uint64    `db:"user_id"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

type PostRequest struct {
	Content string `json:"content" validate:"required"`
}

type PostResponse struct {
	ID        uint64    `json:"id"`
	Content   string    `json:"content"`
	UserID    uint64    `json:"userId"`
	UserName  string    `json:"userName"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}
