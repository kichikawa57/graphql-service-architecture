// pkg/loaders/loaders.go
package loaders

import (
	"context"
	"graphql-backend/models"

	"github.com/graph-gophers/dataloader/v7"
)

type postLoader = dataloader.Loader[int, []*models.Post]

func newPostLoaders(postRepo *models.PostRepository) *postLoader {
	batch := func(ctx context.Context, keys []int) []*dataloader.Result[[]*models.Post] {
		// 1) 一括取得
		rows, err := postRepo.GetPostsByUserIDs(keys)
		if err != nil {
			res := make([]*dataloader.Result[[]*models.Post], len(keys))
			for i := range res {
				res[i] = &dataloader.Result[[]*models.Post]{Error: err}
			}
			return res
		}

		// 2) user_id -> []*models.Post にグルーピング
		group := make(map[int][]*models.Post, len(keys))
		for _, r := range rows {
			group[r.UserID] = append(group[r.UserID], &models.Post{
				ID:        r.ID,
				Title:     r.Title,
				Content:   r.Content,
				CreatedAt: r.CreatedAt,
				UpdatedAt: r.UpdatedAt,
			})
		}

		// 3) keys の順序で結果を返す
		res := make([]*dataloader.Result[[]*models.Post], len(keys))
		for i, k := range keys {
			res[i] = &dataloader.Result[[]*models.Post]{Data: group[k]}
		}
		return res
	}

	return dataloader.NewBatchedLoader(batch)
}
