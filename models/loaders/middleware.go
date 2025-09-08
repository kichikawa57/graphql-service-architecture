package loaders

import (
	"context"
	"graphql-backend/models"
	"net/http"
)

type ctxKey struct{}

func Middleware(repo *models.PostRepository) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			lds := NewLoaders(repo) // ★ リクエストごとに新しいLoaders
			ctx := context.WithValue(r.Context(), ctxKey{}, lds)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func FromContext(ctx context.Context) *Loaders {
	if v, ok := ctx.Value(ctxKey{}).(*Loaders); ok {
		return v
	}
	return nil
}
