package handlers

import (
	"encoding/json"
	"leaderboard-system/internal/models"
	"leaderboard-system/internal/repositories"
	"log"
	"net/http"
)

type GameHandler struct {
	repo *repositories.GameRepository
}

func NewGameHandler(repo *repositories.GameRepository) *GameHandler {
	return &GameHandler{repo: repo}
}

func (h *GameHandler) AddGame(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var game models.Game
	if err := json.NewDecoder(r.Body).Decode(&game); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	gameID, err := h.repo.CreateGame(game)
	if err != nil {
		http.Error(w, "Failed to create game", http.StatusInternalServerError)
		log.Printf("Error creating game: %v", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Game created successfully",
		"id":      gameID,
	})
}

func (h *GameHandler) GetGame(w http.ResponseWriter, r *http.Request) {
	gameID := r.URL.Query().Get("id")
	if gameID == "" {
		http.Error(w, "Missing id parameter", http.StatusBadRequest)
		return
	}

	game, err := h.repo.GetGame(gameID)
	if err != nil {
		http.Error(w, "Failed to get game", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(game)
}

func (h *GameHandler) GetAllGames(w http.ResponseWriter, r *http.Request) {
	games, err := h.repo.GetAllGames()
	if err != nil {
		http.Error(w, "Failed to get all games", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(games)
}
