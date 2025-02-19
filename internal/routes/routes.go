package routes

import (
	"context"
	"leaderboard-system/internal/handlers"
	"leaderboard-system/internal/repositories"

	"github.com/gorilla/mux"
	"github.com/redis/go-redis/v9"
)

func RegisterRoutes(r *mux.Router, client *redis.Client, ctx context.Context) {
	gameRepo := repositories.NewGameRepository(client, ctx)
	gameHandler := handlers.NewGameHandler(gameRepo)

	r.HandleFunc("/games", gameHandler.GetAllGames).Methods("GET")
	r.HandleFunc("/games", gameHandler.AddGame).Methods("POST")
	r.HandleFunc("/games/{id}", gameHandler.GetGame).Methods("GET")

	userRepo := repositories.NewUserRepository(client, ctx)
	userHandler := handlers.NewUserHandler(userRepo)

	r.HandleFunc("/users", userHandler.AddUser).Methods("POST")
	r.HandleFunc("/users/{id}", userHandler.GetUser).Methods("GET")
	r.HandleFunc("/users", userHandler.GetAllUsers).Methods("GET")

	scoreRepo := repositories.NewScoreRepository(client, ctx)
	scoreHandler := handlers.NewScoreHandler(scoreRepo)

	r.HandleFunc("/scores", scoreHandler.GetAllScores).Methods("GET")
	r.HandleFunc("/scores", scoreHandler.AddScore).Methods("POST")
	r.HandleFunc("/scores/{id}", scoreHandler.GetScore).Methods("GET")
	r.HandleFunc("/scores/{score_id}/{user_id}", scoreHandler.GetRank).Methods("GET")
}
