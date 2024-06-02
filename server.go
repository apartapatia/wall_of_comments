package main

import (
	"flag"
	"net/http"
	"os"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/apartapatia/wall_of_comments/graph"
	"github.com/apartapatia/wall_of_comments/internal/config"
	"github.com/apartapatia/wall_of_comments/internal/database"
	"github.com/apartapatia/wall_of_comments/internal/database/pq"
	"github.com/apartapatia/wall_of_comments/internal/database/redis"
	"github.com/sirupsen/logrus"
)

const defaultPort = "8090"

func main() {
	dbtype := flag.String("db", "redis", "Database type to use (redis or postgres)")
	flag.Parse()

	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	conf, err := config.GetConfig()
	if err != nil {
		logrus.Fatalf("failed to get config: %v", err)
	}

	var repo database.Repo
	switch *dbtype {
	case "postgres":
		repo, err = pq.GetRepo(conf.PostgresConfig)
		if err != nil {
			logrus.Fatalf("failed to get postgres repo: %v", err)
		}
	case "redis":
		repo, err = redis.GetRepo(conf.RedisConfig)
		if err != nil {
			logrus.Fatalf("failed to get redis repo: %v", err)
		}
	default:
		logrus.Fatalf("unsupported database type: %s", *dbtype)
	}

	srv := handler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{Resolvers: &graph.Resolver{Repo: repo}}))

	http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", srv)

	logrus.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	logrus.Fatal(http.ListenAndServe(":"+port, nil))
}
