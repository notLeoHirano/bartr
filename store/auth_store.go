package store

import (
	"database/sql"

	"github.com/notLeoHirano/bartr/models"
)

func (r *Store) CreateUser(user *models.User) error {
	result, err := r.db.Exec(
		"INSERT INTO users (name, email, password_hash) VALUES (?, ?, ?)",
		user.Name, user.Email, user.PasswordHash,
	)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}

	user.ID = int(id)
	return nil
}

func (r *Store) GetUserByEmail(email string) (*models.User, error) {
	var user models.User
	err := r.db.QueryRow(
		"SELECT id, name, email, password_hash, created_at FROM users WHERE email = ?",
		email,
	).Scan(&user.ID, &user.Name, &user.Email, &user.PasswordHash, &user.CreatedAt)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *Store) GetUserByID(id int) (*models.User, error) {
	var user models.User
	err := r.db.QueryRow(
		"SELECT id, name, email, created_at FROM users WHERE id = ?",
		id,
	).Scan(&user.ID, &user.Name, &user.Email, &user.CreatedAt)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &user, nil
}