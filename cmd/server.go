package main

import (
	"github.com/KilyakArata/gqlgen-todos/postgres"
	"github.com/go-pg/pg/v10"
	"log"
	"net/http"
	"os"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/KilyakArata/gqlgen-todos/graph"
)

const defaultPort = "8080"

func main() {
	DB := postgres.New(&pg.Options{
		User:     "postgres",
		Password: "Robotesf5",
		Database: "posts_dev",
		Addr:     "ps-psql:5432",
	})

	defer DB.Close()

	DB.AddQueryHook(postgres.DBLogger{})

	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	c := graph.Config{Resolvers: &graph.Resolver{
		PostsRepo:    postgres.PostsRepo{DB: DB},
		CommentsRepo: postgres.CommentsRepo{DB: DB},
	}}

	srv := handler.NewDefaultServer(graph.NewExecutableSchema(c))

	http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", srv)

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
