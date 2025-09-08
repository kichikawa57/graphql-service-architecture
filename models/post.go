package models

import (
	"database/sql"
	"fmt"
	"time"
)

type Post struct {
	ID        int       `json:"id" db:"id"`
	UserID    int       `json:"user_id" db:"user_id"`
	Title     string    `json:"title" db:"title"`
	Content   string    `json:"content" db:"content"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

type PostRepository struct {
	db *sql.DB
}

func NewPostRepository(db *sql.DB) *PostRepository {
	return &PostRepository{db: db}
}

func (r *PostRepository) GetPostsByUserID(userID int) ([]*Post, error) {
	query := "SELECT id, user_id, title, content, created_at, updated_at FROM posts WHERE user_id = ? ORDER BY id"
	rows, err := r.db.Query(query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to query users: %w", err)
	}
	defer rows.Close()

	var users []*Post
	for rows.Next() {
		user := &Post{}
		err := rows.Scan(&user.ID, &user.UserID, &user.Title, &user.Content, &user.CreatedAt, &user.UpdatedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan post: %w", err)
		}
		users = append(users, user)
	}

	return users, nil
}

func (r *PostRepository) GetPostsByUserIDs(userIDs []int) ([]*Post, error) {
	query := "SELECT id, user_id, title, content, created_at, updated_at FROM posts WHERE user_id IN (?)"
	rows, err := r.db.Query(query, userIDs)
	if err != nil {
		return nil, fmt.Errorf("failed to query posts: %w", err)
	}
	defer rows.Close()

	var posts []*Post
	for rows.Next() {
		post := &Post{}
		err := rows.Scan(&post.ID, &post.UserID, &post.Title, &post.Content, &post.CreatedAt, &post.UpdatedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan post: %w", err)
		}
		posts = append(posts, post)
	}

	return posts, nil
}
