package store

import (
	"database/sql"
	"fmt"

	"github.com/zhetkerbaevan/personal-blog/internal/models"
)

type UserStore struct {
	db *sql.DB
}

func NewUserStore(db *sql.DB) *UserStore {
	return &UserStore{db : db}
}

func (s *UserStore) CreateUser(user models.User) error {
	_, err := s.db.Exec("INSERT INTO users (email, password, name, surname, age) VALUES($1, $2, $3, $4, $5)", user.Email, user.Password, user.Name, user.Surname, user.Age)
	if err != nil {
		return err
	}
	return nil
}

func (s *UserStore) GetUserByEmail(email string) (*models.User, error) {
	rows, err := s.db.Query("SELECT * FROM users WHERE email = $1", email)
	if err != nil {
		return nil, err
	}

	user := new(models.User)
	for rows.Next() {
		user, err = scanIntoUser(rows)
		if err != nil {
			return nil, err
		}
	}

	//check if user exists
	if user.Id == 0 {
		return nil, fmt.Errorf("USER NOT FOUND")
	}
	return user, nil

}

func scanIntoUser(rows *sql.Rows) (*models.User, error) {
	user := new(models.User)
	err := rows.Scan(
		&user.Id,
		&user.Email, 
		&user.Password, 
		&user.Name,
		&user.Surname,
		&user.Age)

	if err != nil {
		return nil, err
	}
	return user, nil
}