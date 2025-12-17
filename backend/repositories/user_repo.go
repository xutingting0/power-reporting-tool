// file: repositories/user_repo.go
package repositories

import (
	"backend/models"
	"database/sql"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) GetUserByUsername(username string) (*models.User, error) {
	var user models.User
	query := "SELECT id, username, password, real_name, phone, position, role, unit FROM users WHERE username = ?"
	err := r.db.QueryRow(query, username).Scan(
		&user.ID, &user.Username, &user.Password, &user.RealName,
		&user.Phone, &user.Position, &user.Role, &user.Unit,
	)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) GetUserByID(id int) (*models.User, error) {
	var user models.User
	query := "SELECT id, username, real_name, phone, position, role, unit FROM users WHERE id = ?"
	err := r.db.QueryRow(query, id).Scan(
		&user.ID, &user.Username, &user.RealName,
		&user.Phone, &user.Position, &user.Role, &user.Unit,
	)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) GetAllUsers() ([]models.User, error) {
	rows, err := r.db.Query("SELECT id, username, real_name, phone, position, role, unit FROM users ORDER BY id")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []models.User
	for rows.Next() {
		var user models.User
		err := rows.Scan(&user.ID, &user.Username, &user.RealName,
			&user.Phone, &user.Position, &user.Role, &user.Unit)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	return users, nil
}

func (r *UserRepository) CreateUser(user *models.User) error {
	_, err := r.db.Exec(
		"INSERT INTO users (username, password, real_name, phone, position, role, unit) VALUES (?, ?, ?, ?, ?, ?, ?)",
		user.Username, user.Password, user.RealName, user.Phone,
		user.Position, user.Role, user.Unit,
	)
	return err
}

func (r *UserRepository) UpdateUser(user *models.User) error {
	_, err := r.db.Exec(`
		UPDATE users 
		SET real_name = ?, phone = ?, position = ?, role = ?, unit = ?
		WHERE username = ?`,
		user.RealName, user.Phone, user.Position, user.Role, user.Unit, user.Username,
	)
	return err
}

func (r *UserRepository) DeleteUser(username string) error {
	_, err := r.db.Exec("DELETE FROM users WHERE username = ?", username)
	return err
}

func (r *UserRepository) UpdateUserProfile(user *models.User) error {
	_, err := r.db.Exec(
		"UPDATE users SET real_name = ?, phone = ?, position = ?, unit = ? WHERE id = ?",
		user.RealName, user.Phone, user.Position, user.Unit, user.ID,
	)
	return err
}
