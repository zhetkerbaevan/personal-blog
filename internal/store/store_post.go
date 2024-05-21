package store

import (
	"database/sql"
	"fmt"

	"github.com/zhetkerbaevan/personal-blog/internal/models"
)

type PostStore struct {
	db *sql.DB
}

func NewPostStore(db *sql.DB) *PostStore {
	return &PostStore{db : db}
}

func (s *PostStore) CreatePost(post models.CreatePost) error {
	_, err := s.db.Exec("INSERT INTO posts (title, description, userId) VALUES ($1, $2, $3)", post.Title, post.Description, post.UserId)
	if err != nil {
		return err
	}
	return nil
}

func (s *PostStore) GetPostsByUserIds(userId int) ([]models.Post, error) {
	rows, err := s.db.Query("SELECT * FROM posts WHERE userId = $1", userId)
	if err != nil {
		return nil, err
	}

	posts := make([]models.Post, 0)

	for rows.Next() {
		p, err := ScanIntoPost(rows)
		if err != nil {
			return nil, err
		}
		posts = append(posts, *p)
	}

	return posts, nil
} 

func (s *PostStore) UpdatePost(postId int, post models.Post) error {
	res, err := s.db.Exec("UPDATE posts SET title = $1, description = $2 WHERE id = $3", post.Title, post.Description, postId)
	if err != nil {
		return err
	}

	//check if any row was updated
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return fmt.Errorf("NO ROWS WITH POSTID %d WERE AFFECTED", postId)
	}
	return nil
}

func (s *PostStore) DeletePost(postId int) error {
	_, err := s.db.Exec("DELETE FROM posts WHERE id = $1", postId)
	if err != nil {
		return err
	}
	return nil
}

func (s *PostStore) GetPostById(postId int) (*models.Post, error) {
	rows, err := s.db.Query("SELECT * FROM posts WHERE id = $1", postId)
	if err != nil {
		return nil, err
	}

	post := new(models.Post)
	for rows.Next() {
		post, err = ScanIntoPost(rows)
		if err != nil {
			return nil, err
		}
	}
	if post.Id == 0 {
		return nil, fmt.Errorf("POST NOT FOUND")
	}
	return post, nil
}

func ScanIntoPost(rows *sql.Rows) (*models.Post, error) {
	post := new(models.Post)

	err := rows.Scan(
		&post.Id,
		&post.Title,
		&post.Description,
		&post.CreatedAt,
		&post.UserId)
	if err != nil {
		return nil, err
	}
	
	return post, nil
}