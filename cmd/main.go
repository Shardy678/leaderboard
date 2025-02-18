package main

import (
	"context"
	"encoding/json"
	"leaderboard-system/internal/handlers"
	"leaderboard-system/internal/repositories"
	"log"
	"net/http"

	"github.com/redis/go-redis/v9"
)

var (
	ctx = context.Background()
	rdb *redis.Client
)

func main() {
	rdb = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	_, err := rdb.Ping(ctx).Result()
	if err != nil {
		log.Fatalf("Could not connect to Redis: %v", err)
	}

	userRepo := repositories.NewUserRepository(rdb, ctx)
	userHandler := handlers.NewUserHandler(userRepo)

	gameRepo := repositories.NewGameRepository(rdb, ctx)
	gameHandler := handlers.NewGameHandler(gameRepo)

	http.HandleFunc("/add_user", userHandler.AddUser)
	http.HandleFunc("/get_user", userHandler.GetUser)
	http.HandleFunc("/get_all_users", userHandler.GetAllUsers)
	http.HandleFunc("/add_game", gameHandler.AddGame)
	http.HandleFunc("/get_game", gameHandler.GetGame)
	log.Println("Starting server on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func getAllGames(w http.ResponseWriter, r *http.Request) {
	games, err := rdb.Keys(ctx, "*").Result()
	if err != nil {
		http.Error(w, "Could not retrieve games", http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"games": games,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
