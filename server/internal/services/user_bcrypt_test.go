package services

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
)

func TestBcryptPasswordHashing(t *testing.T) {
	password := "testpassword123"

	// 测试密码哈希生成
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	assert.NoError(t, err)
	assert.NotEqual(t, password, string(hashedPassword))

	// 测试正确密码验证
	err = bcrypt.CompareHashAndPassword(hashedPassword, []byte(password))
	assert.NoError(t, err)

	// 测试错误密码验证应该失败
	err = bcrypt.CompareHashAndPassword(hashedPassword, []byte("wrongpassword"))
	assert.Error(t, err)
	assert.Equal(t, bcrypt.ErrMismatchedHashAndPassword, err)
}
