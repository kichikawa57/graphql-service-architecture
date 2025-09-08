package graph

import (
	"database/sql"
	"graphql-backend/models"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	userRepo *models.UserRepository
	postRepo *models.PostRepository
}

func NewResolver(db *sql.DB) *Resolver {
	return &Resolver{
		userRepo: models.NewUserRepository(db),
		postRepo: models.NewPostRepository(db),
	}
}
