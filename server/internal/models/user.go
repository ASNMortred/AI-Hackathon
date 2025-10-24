package models

import (
	"database/sql"
	"time"
)

// User 用户模型
type User struct {
	Uid       int          `json:"uid" db:"uid"`
	Username  string       `json:"username" db:"username"`
	Password  string       `json:"password" db:"password"`
	CreatedAt time.Time    `json:"created_at" db:"created_at"`
	UpdatedAt sql.NullTime `json:"updated_at" db:"updated_at"`
}
