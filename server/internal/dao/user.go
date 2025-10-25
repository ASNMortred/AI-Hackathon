package dao

import (
	"database/sql"
	"fmt"

	"github.com/ASNMortred/AI-Hackathon/server/internal/database"
	"github.com/ASNMortred/AI-Hackathon/server/internal/models"
)

// UserDAO 用户数据访问对象
type UserDAO struct{}

// NewUserDAO 创建新的用户DAO实例
func NewUserDAO() *UserDAO {
	return &UserDAO{}
}

// CreateUser 创建新用户（密码应为加密后的哈希值）
func (dao *UserDAO) CreateUser(username, password string) error {
	query := "INSERT INTO users (username, password) VALUES (?, ?)"
	_, err := database.DB.Exec(query, username, password)
	if err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}
	return nil
}

// GetUserByUsername 根据用户名获取用户
func (dao *UserDAO) GetUserByUsername(username string) (*models.User, error) {
	query := "SELECT uid, username, password, created_at, updated_at FROM users WHERE username = ?"
	row := database.DB.QueryRow(query, username)

	user := &models.User{}
	err := row.Scan(&user.Uid, &user.Username, &user.Password, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // 用户不存在
		}
		return nil, fmt.Errorf("failed to get user by username: %w", err)
	}

	return user, nil
}

// CheckUserExists 检查用户是否存在
func (dao *UserDAO) CheckUserExists(username string) (bool, error) {
	query := "SELECT COUNT(*) FROM users WHERE username = ?"
	var count int
	err := database.DB.QueryRow(query, username).Scan(&count)
	if err != nil {
		return false, fmt.Errorf("failed to check user existence: %w", err)
	}
	return count > 0, nil
}
