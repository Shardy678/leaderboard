package repositories

import (
	"context"
	"database/sql"
	"fmt"
	"leaderboard-system/internal/models"

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
	case float64:
		return int(v), nil
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
	if score.GameID == "" || score.UserID == "" {
		return "", fmt.Errorf("game ID and user ID are required")
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
			for rank, score := range scores {
				allScores = append(allScores, map[string]interface{}{
					"game_id": key,
					"user_id": score.Member,
					"rank":    rank + 1,
					"score":   score.Score,
				})
			}
		}

		if cursor == 0 {
			break
		}
	}

	return allScores, nil
}

func (repo *ScoreRepository) GetScores(scoreID string) ([]map[string]interface{}, error) {
	scoreID = "score:" + scoreID
	scores, err := repo.redisClient.ZRevRangeWithScores(repo.ctx, scoreID, 0, -1).Result()
	if err != nil {
		return nil, err
	}

	if len(scores) == 0 {
		return nil, nil
	}

	var result []map[string]interface{}
	for _, score := range scores {
		result = append(result, map[string]interface{}{
			"score_id": scoreID[6:],
			"score":    score.Score,
			"member":   score.Member,
		})
	}
	return result, nil
}

func (repo *ScoreRepository) GetRank(scoreID, userID string) (int, error) {
	rank, err := repo.redisClient.ZRevRank(repo.ctx, scoreID, userID).Result()
	if err != nil {
		return 0, err
	}
	return int(rank) + 1, nil
}
