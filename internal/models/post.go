package models

import "time"

type Post struct {
	Id          int `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	CreatedAt   time.Time `json:"createdAt"`
	UserId      int `json:"userId"`
}

type CreatePost struct {
	Title       string `json:"title" validate:"required"`
	Description string `json:"description" validate:"required"`
	UserId      int `json:"userId"`
}

type UpdatePost struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}

type PostStoreInterface interface {
	CreatePost(CreatePost) error
	GetPostsByUserIds(int) ([]Post, error)
	UpdatePost(int, Post) error
	DeletePost(int) error
	GetPostById(int) (*Post, error)
}