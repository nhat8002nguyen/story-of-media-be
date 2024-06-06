package services

import (
	"database/sql"
	"fmt"

	"github.com/nhat8002nguyen/story-of-media-be/story-service/src/models"
	"golang.org/x/crypto/bcrypt"
)

func GetUser(email string) (*models.User, error) {
	stmt := "SELECT * FROM users WHERE users.email=$1"
	u := models.User{}
	row := models.Db.QueryRow(stmt, email)
	switch err := row.Scan(&u.ID, &u.Name, &u.Email, &u.Password); err {
	case sql.ErrNoRows:
		return &u, fmt.Errorf("no record found for email: %s", email)
	case nil:
		return &u, nil
	default:
		return &u, err
	}
}

func AddUser(u *models.User) (*models.User, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), 10)
	if err != nil {
		return u, fmt.Errorf("error hashing password: %w", err)
	}

	statement := `
		INSERT INTO users (name, email, password)
		VALUES ($1, $2, $3)
		RETURNING id`

	id := ""
	err = models.Db.QueryRow(statement, u.Name, u.Email, hashedPassword).Scan(&id)
	if err != nil {
		return u, err
	}

	u.ID = id
	return u, nil
}
