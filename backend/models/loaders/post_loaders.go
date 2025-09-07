// pkg/loaders/loaders.go
package loaders

import (
	"context"
	"fmt"

	"graphql-backend/graph/model"
	"graphql-backend/models"

	"github.com/graph-gophers/dataloader/v7"
)

type Loaders struct {
	// key: int64(user_id), value: []*model.Post
	PostsByUserID *dataloader.Loader[int, []*model.Post]
}

func NewPostLoaders(postRepo *models.PostRepository) *Loaders {
	batch := func(ctx context.Context, keys []int) []*dataloader.Result[[]*model.Post] {
		// 1) 一括取得
		rows, err := postRepo.GetPostsByUserIDs(keys)
		if err != nil {
			res := make([]*dataloader.Result[[]*model.Post], len(keys))
			for i := range res {
				res[i] = &dataloader.Result[[]*model.Post]{Error: err}
			}
			return res
		}

		// 2) user_id -> []*model.Post にグルーピング
		group := make(map[int][]*model.Post, len(keys))
		for _, r := range rows {
			group[r.UserID] = append(group[r.UserID], &model.Post{
				ID:        fmt.Sprint(r.ID),
				Title:     r.Title,
				Content:   r.Content,
				CreatedAt: r.CreatedAt.Format("2006-01-02T15:04:05Z"),
				UpdatedAt: r.UpdatedAt.Format("2006-01-02T15:04:05Z"),
			})
		}

		// 3) keys の順序で結果を返す
		res := make([]*dataloader.Result[[]*model.Post], len(keys))
		for i, k := range keys {
			res[i] = &dataloader.Result[[]*model.Post]{Data: group[k]}
		}
		return res
	}

	return &Loaders{
		PostsByUserID: dataloader.NewBatchedLoader(batch),
	}
}
