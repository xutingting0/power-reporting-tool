// file: models/user_model.go
package models

import (
	"backend/db"
)

// User 用户模型
type User struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Password string `json:"-"`
	RealName string `json:"real_name"`
	Phone    string `json:"phone"`
	Position string `json:"position"`
	Role     string `json:"role"`
	Unit     string `json:"unit"`
}

// GetUserByID 根据ID获取用户
func GetUserByID(id int) (*User, error) {
	var user User
	err := db.DB.QueryRow(`
		SELECT id, username, password, real_name, phone, position, role, unit
		FROM users
		WHERE id = ?`, id).Scan(
		&user.ID, &user.Username, &user.Password, &user.RealName,
		&user.Phone, &user.Position, &user.Role, &user.Unit)
	if err != nil {
		return nil, err
	}
	return &user, nil
}
