package handlers

import (
	"encoding/json"
	"leaderboard-system/internal/models"
	"leaderboard-system/internal/repositories"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type GameHandler struct {
	repo *repositories.GameRepository
}

func NewGameHandler(repo *repositories.GameRepository) *GameHandler {
	return &GameHandler{repo: repo}
}

func (h *GameHandler) AddGame(w http.ResponseWriter, r *http.Request) {
	var game models.Game
	if err := json.NewDecoder(r.Body).Decode(&game); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	if gameID, err := h.repo.CreateGame(game); err != nil {
		http.Error(w, "Failed to create game", http.StatusInternalServerError)
		log.Printf("Error creating game: %v", err)
	} else {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(map[string]string{"message": "Game created successfully", "id": gameID})
	}
}

func (h *GameHandler) GetGame(w http.ResponseWriter, r *http.Request) {
	gameID := mux.Vars(r)["id"]
	if gameID == "" {
		http.Error(w, "Missing id parameter", http.StatusBadRequest)
		return
	}

	if game, err := h.repo.GetGame(gameID); err != nil {
		http.Error(w, "Failed to get game", http.StatusInternalServerError)
	} else {
		json.NewEncoder(w).Encode(game)
	}
}

func (h *GameHandler) GetAllGames(w http.ResponseWriter, r *http.Request) {
	if games, err := h.repo.GetAllGames(); err != nil {
		http.Error(w, "Failed to get all games", http.StatusInternalServerError)
		log.Printf("Error getting all games: %v", err)
	} else {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(games)
	}
}
