package repositories

import (
	"context"
	"crypto/md5"
	"fmt"
	"leaderboard-system/internal/models"

	"github.com/redis/go-redis/v9"
)

type GameRepository struct {
	client *redis.Client
	ctx    context.Context
}

func NewGameRepository(client *redis.Client, ctx context.Context) *GameRepository {
	return &GameRepository{
		client: client,
		ctx:    ctx,
	}
}

func (repo *GameRepository) CreateGame(game models.Game) (string, error) {
	if game.Name == "" {
		return "", fmt.Errorf("game name cannot be empty")
	}

	hash := md5.Sum([]byte(game.Name))
	game.ID = fmt.Sprintf("%x", hash)
	game.ID = "game:" + game.ID[:10]

	err := repo.client.HSet(repo.ctx, game.ID, map[string]interface{}{
		"name": game.Name,
	}).Err()
	if err != nil {
		return "", err
	}

	return game.ID, nil
}

func (repo *GameRepository) GetGame(gameID string) (models.Game, error) {
	if !isValidGameID(gameID) {
		return models.Game{}, fmt.Errorf("invalid game ID format")
	}

	game, err := repo.client.HGetAll(repo.ctx, gameID).Result()
	if err != nil {
		return models.Game{}, err
	}

	if len(game) == 0 {
		return models.Game{}, fmt.Errorf("game not found")
	}

	return models.Game{
		ID:   gameID,
		Name: game["name"],
	}, nil
}

func isValidGameID(gameID string) bool {
	return len(gameID) == 14 && gameID[:5] == "game:"
}

func (repo *GameRepository) GetAllGames() ([]models.Game, error) {
	keys, err := repo.client.Keys(repo.ctx, "game:*").Result()
	if err != nil {
		return nil, err
	}

	var games []models.Game
	for _, key := range keys {
		gameData, err := repo.client.HGetAll(repo.ctx, key).Result()
		if err != nil {
			return nil, err
		}

		if len(gameData) > 0 {
			games = append(games, models.Game{
				ID:   key,
				Name: gameData["name"],
			})
		}
	}

	return games, nil
}
