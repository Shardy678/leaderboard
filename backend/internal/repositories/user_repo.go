package repositories

import (
	"context"
	"fmt"
	"leaderboard-system/internal/models"
	"log"

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
	if user.Username == "" || user.Password == "" {
		return "", fmt.Errorf("username and password cannot be empty")
	}

	query := "INSERT INTO users (username, password) VALUES ($1, $2) RETURNING id"
	var userID string
	err := repo.db.QueryRowContext(repo.ctx, query, user.Username, user.Password).Scan(&userID)
	if err != nil {
		log.Println("Error creating user:", err)
		return "", err
	}

	return userID, nil
}

func (repo *UserRepository) GetUser(userID string) (models.User, error) {
	query := "SELECT username, password FROM users WHERE id = $1"
	var user models.User
	err := repo.db.QueryRowContext(repo.ctx, query, userID).Scan(&user.Username, &user.Password)
	if err != nil {
		if err == sql.ErrNoRows {
			return models.User{}, nil
		}
		return models.User{}, err
	}
	user.ID = userID
	return user, nil
}

func (repo *UserRepository) GetAllUsers() ([]interface{}, error) {
	query := "SELECT id, username FROM users"
	rows, err := repo.db.QueryContext(repo.ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []interface{}
	for rows.Next() {
		var user struct {
			ID       string
			Username string
		}
		if err := rows.Scan(&user.ID, &user.Username); err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}
