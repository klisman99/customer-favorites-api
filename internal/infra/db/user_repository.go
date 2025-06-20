package db

import (
	"app/internal/domain/model"
	"database/sql"
)

type UserRepository struct {
	DB *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{DB: db}
}

func (r *UserRepository) Create(user model.User) (*model.User, error) {
	query := `
		INSERT INTO users (id, username, password)
		VALUES ($1, $2, $3)
		RETURNING id, username, password, created_at
	`

	row := r.DB.QueryRow(query, user.ID, user.Username, user.Password)

	var createdUser model.User
	if err := row.Scan(&createdUser.ID, &createdUser.Username, &createdUser.Password, &createdUser.CreatedAt); err != nil {
		return nil, err
	}

	return &createdUser, nil
}

func (r *UserRepository) FindByUsername(username string) (*model.User, error) {
	query := `
		SELECT id, username, password, created_at
		FROM users
		WHERE username = $1
	`

	row := r.DB.QueryRow(query, username)

	var user model.User
	err := row.Scan(&user.ID, &user.Username, &user.Password, &user.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &user, nil
}
