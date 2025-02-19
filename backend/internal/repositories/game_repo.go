package repositories

import (
	"context"
	"fmt"
	"leaderboard-system/internal/models"

	"database/sql"

	"github.com/redis/go-redis/v9"
)

type GameRepository struct {
	redisClient *redis.Client
	ctx         context.Context
	db          *sql.DB
}

func NewGameRepository(redisClient *redis.Client, ctx context.Context, db *sql.DB) *GameRepository {
	return &GameRepository{
		redisClient: redisClient,
		ctx:         ctx,
		db:          db,
	}
}

func (repo *GameRepository) CreateGame(game models.Game) (string, error) {
	if game.Name == "" {
		return "", fmt.Errorf("game name cannot be empty")
	}

	query := "INSERT INTO games (name) VALUES ($1) RETURNING id"
	var gameID string
	err := repo.db.QueryRowContext(repo.ctx, query, game.Name).Scan(&gameID)
	if err != nil {
		return "", err
	}

	return gameID, nil
}

func (repo *GameRepository) GetGame(gameID string) (models.Game, error) {
	if !isValidGameID(gameID) {
		return models.Game{}, fmt.Errorf("invalid game ID format")
	}

	query := "SELECT id, name FROM games WHERE id = $1"
	var game models.Game
	err := repo.db.QueryRowContext(repo.ctx, query, gameID).Scan(&game.ID, &game.Name)
	if err != nil {
		if err == sql.ErrNoRows {
			return models.Game{}, fmt.Errorf("game not found")
		}
		return models.Game{}, err
	}

	return game, nil
}

func isValidGameID(gameID string) bool {
	return len(gameID) == 14 && gameID[:5] == "game:"
}

func (repo *GameRepository) GetAllGames() ([]models.Game, error) {
	query := "SELECT id, name FROM games"
	rows, err := repo.db.QueryContext(repo.ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var games []models.Game
	for rows.Next() {
		var game models.Game
		if err := rows.Scan(&game.ID, &game.Name); err != nil {
			return nil, err
		}
		games = append(games, game)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return games, nil
}
