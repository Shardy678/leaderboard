package repositories

import (
	"context"
	"crypto/md5"
	"fmt"
	"leaderboard-system/internal/models"

	"database/sql"

	"github.com/redis/go-redis/v9"
)

type UserRepository struct {
	redisClient *redis.Client
	ctx         context.Context
	db          *sql.DB
}

func NewUserRepository(redisClient *redis.Client, ctx context.Context, db *sql.DB) *UserRepository {
	return &UserRepository{
		redisClient: redisClient,
		ctx:         ctx,
		db:          db,
	}
}

func (repo *UserRepository) CreateUser(user models.User) (string, error) {
	hash := md5.Sum([]byte(user.Username))
	user.ID = "user:" + fmt.Sprintf("%x", hash)
	user.ID = user.ID[:14]

	err := repo.redisClient.HSet(repo.ctx, user.ID, map[string]interface{}{
		"username": user.Username,
		"password": user.Password,
	}).Err()

	if err != nil {
		return "", err
	}

	return user.ID, nil
}

func (repo *UserRepository) GetUser(userID string) (models.User, error) {
	userData, err := repo.redisClient.HGetAll(repo.ctx, userID).Result()
	if err != nil {
		return models.User{}, err
	}

	if len(userData) == 0 {
		return models.User{}, nil // User not found
	}

	return models.User{
		ID:       userID,
		Username: userData["username"],
		Password: userData["password"],
	}, nil
}

func (repo *UserRepository) GetAllUsers() ([]string, error) {
	users, err := repo.redisClient.Keys(repo.ctx, "*").Result()
	return users, err
}
