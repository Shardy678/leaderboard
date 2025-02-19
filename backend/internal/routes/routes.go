package routes

import (
	"context"
	"database/sql"
	"leaderboard-system/internal/handlers"
	"leaderboard-system/internal/repositories"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/redis/go-redis/v9"
)

func RegisterRoutes(r *mux.Router, redisClient *redis.Client, ctx context.Context, db *sql.DB) {
	r.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Access-Control-Allow-Origin", "http://localhost:5173")
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
			if r.Method == http.MethodOptions {
				w.WriteHeader(http.StatusOK)
				return
			}
			next.ServeHTTP(w, r)
		})
	})

	gameRepo := repositories.NewGameRepository(redisClient, ctx, db)
	gameHandler := handlers.NewGameHandler(gameRepo)

	r.HandleFunc("/games", gameHandler.GetAllGames).Methods("GET")
	r.HandleFunc("/games", gameHandler.AddGame).Methods("POST")
	r.HandleFunc("/games/{id}", gameHandler.GetGame).Methods("GET")

	userRepo := repositories.NewUserRepository(redisClient, ctx, db)
	userHandler := handlers.NewUserHandler(userRepo)

	r.HandleFunc("/users", userHandler.AddUser).Methods("POST", "OPTIONS")
	r.HandleFunc("/users/{id}", userHandler.GetUser).Methods("GET")
	r.HandleFunc("/users", userHandler.GetAllUsers).Methods("GET")

	scoreRepo := repositories.NewScoreRepository(redisClient, ctx, db)
	scoreHandler := handlers.NewScoreHandler(scoreRepo)

	r.HandleFunc("/scores", scoreHandler.GetAllScores).Methods("GET")
	r.HandleFunc("/scores", scoreHandler.AddScore).Methods("POST", "OPTIONS")
	r.HandleFunc("/scores/{id}", scoreHandler.GetScore).Methods("GET")
	r.HandleFunc("/scores/{score_id}/{user_id}", scoreHandler.GetRank).Methods("GET")
}
