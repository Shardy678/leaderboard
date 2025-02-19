package repositories

import (
	"context"
	"database/sql"
	"fmt"
	"leaderboard-system/internal/models"
	"strings"

	"github.com/redis/go-redis/v9"
)

type ScoreRepository struct {
	redisClient *redis.Client
	ctx         context.Context
	db          *sql.DB
}

func NewScoreRepository(redisClient *redis.Client, ctx context.Context, db *sql.DB) *ScoreRepository {
	return &ScoreRepository{
		redisClient: redisClient,
		ctx:         ctx,
		db:          db,
	}
}

func (repo *ScoreRepository) ScoreToInt(score models.Score) (int, error) {
	switch v := score.Score.(type) {
	case int:
		return v, nil
	case string:
		var score int
		_, err := fmt.Sscanf(v, "%d", &score)
		if err != nil {
			return 0, fmt.Errorf("invalid score value: %s", v)
		}
		return score, nil
	default:
		return 0, fmt.Errorf("score must be an integer or string representation of an integer")
	}
}

func (repo *ScoreRepository) CreateScore(score models.Score) (string, error) {
	if !isValidID(score.GameID) || !isValidID(score.UserID) {
		return "", fmt.Errorf("game ID and user ID are required and must start with 'game:' and 'user:' respectively")
	}

	scoreInt, err := repo.ScoreToInt(score)
	if err != nil {
		return "", err
	}

	if scoreInt < 0 {
		return "", fmt.Errorf("score cannot be negative")
	}

	score.ID = "score:" + score.GameID

	_, err = repo.redisClient.ZAdd(repo.ctx, score.ID,
		redis.Z{
			Score:  float64(scoreInt),
			Member: score.UserID,
		},
	).Result()

	if err != nil {
		return "", err
	}

	return score.ID, nil
}

func isValidID(id string) bool {
	return len(id) > 0 && (id[:5] == "game:" || id[:5] == "user:")
}

func (repo *ScoreRepository) GetAllScores() ([]map[string]interface{}, error) {
	var allScores []map[string]interface{}

	var cursor uint64
	for {
		keys, newCursor, err := repo.redisClient.Scan(repo.ctx, cursor, "score:*", 0).Result()
		if err != nil {
			return nil, err
		}
		cursor = newCursor

		for _, key := range keys {
			scores, err := repo.redisClient.ZRevRangeWithScores(repo.ctx, key, 0, -1).Result()
			if err != nil {
				return nil, err
			}
			gameID := strings.TrimPrefix(key, "score:")
			gameName, err := repo.redisClient.HGet(repo.ctx, gameID, "name").Result()
			if err != nil {
				return nil, err
			}
			for rank, score := range scores {
				userID := score.Member
				userIDStr := fmt.Sprintf("%v", userID)
				username, err := repo.redisClient.HGet(repo.ctx, userIDStr, "username").Result()
				fmt.Println(userIDStr)
				if err != nil {
					return nil, err
				}
				allScores = append(allScores, map[string]interface{}{
					"rank":      rank + 1,
					"score":     score.Score,
					"user_id":   userID,
					"username":  username,
					"game_id":   gameID,
					"game_name": gameName,
				})
			}
		}

		if cursor == 0 {
			break
		}
	}

	return allScores, nil
}

func (repo *ScoreRepository) GetScore(scoreID string) ([]redis.Z, error) {
	scores, err := repo.redisClient.ZRevRangeWithScores(repo.ctx, scoreID, 0, 1).Result()

	if err != nil {
		return nil, err
	}

	return scores, nil
}

func (repo *ScoreRepository) GetRank(scoreID, userID string) (int, error) {
	rank, err := repo.redisClient.ZRevRank(repo.ctx, scoreID, userID).Result()
	if err != nil {
		return 0, err
	}
	return int(rank) + 1, nil
}
