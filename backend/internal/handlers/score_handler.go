package handlers

import (
	"encoding/json"
	"leaderboard-system/internal/models"
	"leaderboard-system/internal/repositories"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type ScoreHandler struct {
	repo *repositories.ScoreRepository
}

func NewScoreHandler(repo *repositories.ScoreRepository) *ScoreHandler {
	return &ScoreHandler{repo: repo}
}

func (h *ScoreHandler) AddScore(w http.ResponseWriter, r *http.Request) {
	var score models.Score
	if err := json.NewDecoder(r.Body).Decode(&score); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if scoreID, err := h.repo.CreateScore(score); err != nil {
		http.Error(w, "Failed to create score", http.StatusInternalServerError)
		log.Printf("Error creating score: %v", err)
	} else {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(map[string]string{"message": "Score created successfully", "id": scoreID})
	}
}

func (h *ScoreHandler) GetAllScores(w http.ResponseWriter, r *http.Request) {
	scores, err := h.repo.GetAllScores()
	if err != nil {
		http.Error(w, "Failed to get all scores", http.StatusInternalServerError)
		log.Printf("Error getting all scores: %v", err)
	} else {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(scores)
	}
}

func (h *ScoreHandler) GetScore(w http.ResponseWriter, r *http.Request) {
	scoreID := mux.Vars(r)["id"]
	if scoreID == "" {
		http.Error(w, "Missing id parameter", http.StatusBadRequest)
		return
	}

	scores, err := h.repo.GetScores(scoreID)
	if err != nil {
		http.Error(w, "Failed to get score", http.StatusInternalServerError)
		log.Printf("Error getting score: %v", err)
	} else {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(scores)
	}
}

func (h *ScoreHandler) DeleteScore(w http.ResponseWriter, r *http.Request) {
	scoreID := mux.Vars(r)["score_id"]
	userID := mux.Vars(r)["user_id"]

	err := h.repo.DeleteScore(scoreID, userID)
	if err != nil {
		http.Error(w, "Failed to delete score", http.StatusInternalServerError)
		log.Printf("Error deleting score: %v", err)
	} else {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"message": "Score deleted successfully"})
	}
}

func (h *ScoreHandler) GetRank(w http.ResponseWriter, r *http.Request) {
	scoreID := mux.Vars(r)["score_id"]
	userID := mux.Vars(r)["user_id"]

	rank, err := h.repo.GetRank(scoreID, userID)
	if err != nil {
		http.Error(w, "Failed to get rank", http.StatusInternalServerError)
		log.Printf("Error getting rank: %v", err)
	} else {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]int{"rank": rank})
	}
}
