package main

import (
	"context"
	"leaderboard-system/internal/database"
	"leaderboard-system/internal/routes"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/redis/go-redis/v9"
)

var (
	ctx       = context.Background()
	rdb       *redis.Client
	dbService *database.Service
)

func main() {
	rdb = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
		PoolSize: 100,
	})

	_, err := rdb.Ping(ctx).Result()
	if err != nil {
		log.Fatalf("Could not connect to Redis: %v", err)
	}

	dbService = &database.Service{}
	db, err := dbService.Connect()
	if err != nil {
		log.Fatalf("Could not connect to database: %v", err)
	}
	log.Println("Connected to database")
	defer dbService.Close()

	router := mux.NewRouter()
	routes.RegisterRoutes(router, rdb, ctx, db)

	log.Println("Starting server on :8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}
