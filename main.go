package main

import (
	"log"
	"net/http"
	"os"

	"graphql-backend/database"
	"graphql-backend/graph"
	"graphql-backend/models"
	"graphql-backend/models/loaders"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/rs/cors"
)

const defaultPort = "8080"

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	db, err := database.NewDB()
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer db.Close()

	resolver := graph.NewResolver(db.DB)
	srv := handler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{Resolvers: resolver}))

	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		AllowCredentials: true,
	})

	// --- ここがポイント ---
	// 1) まず CORS で包む
	corsWrapped := c.Handler(srv)

	postRepo := models.NewPostRepository(db.DB)

	// 2) その上から DataLoader ミドルウェアで包む（リクエストごとにLoadersを注入）
	//    NewLoaders に必要な依存（repo/DB等）を渡してください
	loaderWrapped := loaders.Middleware(postRepo)(corsWrapped)
	

	http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", loaderWrapped)

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}