package services

import (
	"fmt"

	"github.com/ASNMortred/AI-Hackathon/server/internal/dao"
	"github.com/ASNMortred/AI-Hackathon/server/internal/models"
	"golang.org/x/crypto/bcrypt"
)

// UserService 用户服务
type UserService struct {
	userDAO *dao.UserDAO
}

// NewUserService 创建新的用户服务实例
func NewUserService() *UserService {
	return &UserService{
		userDAO: dao.NewUserDAO(),
	}
}

// RegisterUser 注册用户（使用 bcrypt 加密密码）
func (s *UserService) RegisterUser(username, password string) error {
	// 检查用户名是否已存在
	exists, err := s.userDAO.CheckUserExists(username)
	if err != nil {
		return fmt.Errorf("failed to check user existence: %w", err)
	}
	if exists {
		return fmt.Errorf("username already exists")
	}

	// 使用 bcrypt 加密密码
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}

	// 创建用户，存储加密后的密码
	err = s.userDAO.CreateUser(username, string(hashedPassword))
	if err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}

	return nil
}

// AuthenticateUser 验证用户身份（使用 bcrypt 验证密码）
func (s *UserService) AuthenticateUser(username, password string) (*models.User, error) {
	// 获取用户信息
	user, err := s.userDAO.GetUserByUsername(username)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}
	if user == nil {
		return nil, fmt.Errorf("user not found")
	}

	// 使用 bcrypt 验证密码
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return nil, fmt.Errorf("invalid password")
	}

	return user, nil
}
