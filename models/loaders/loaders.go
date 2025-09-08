package loaders

import (
	"graphql-backend/models"

	"github.com/graph-gophers/dataloader/v7"
)

type Loaders struct {
	// key: int64(user_id), value: []*model.Post
	PostsByUserID *dataloader.Loader[int, []*models.Post]
}


func NewLoaders(postRepo *models.PostRepository) *Loaders {
	postLoader :=  newPostLoaders(postRepo)

	return &Loaders{
		PostsByUserID: postLoader,
	}
}