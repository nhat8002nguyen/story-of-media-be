package services

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/nhat8002nguyen/story-of-media-be/story-service/src/models"
	"golang.org/x/crypto/bcrypt"
)

var SecretKey []byte

type Claims struct {
	Email string
	jwt.Claims
}

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

func Authenticate(creds models.User) (string, error) {
	stmt := "SELECT * FROM users WHERE users.email=$1"
	u := models.User{}
	row := models.Db.QueryRow(stmt, creds.Email)
	err := row.Scan(&u.ID, &u.Name, &u.Email, &u.Password)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", &CustomError{
				Code:    ERROR_NOT_FOUND,
				Message: "no user found",
			}
		}
		return "", &CustomError{
			Code:    ERROR_INTERNAL_SERVER,
			Message: err.Error(),
		}
	}

	// Check password
	if err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(creds.Password)); err != nil {
		return "", &CustomError{
			Code:    ERROR_UNAUTHORIZED,
			Message: err.Error(),
		}
	}

	claims := Claims{
		Email: creds.Email,
		Claims: jwt.RegisteredClaims{
			ExpiresAt: &jwt.NumericDate{
				Time: time.Now().Add(30 * time.Minute),
			},
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString(SecretKey)
	if err != nil {
		return "", &CustomError{
			Code:    ERROR_INTERNAL_SERVER,
			Message: err.Error(),
		}
	}

	return signedToken, nil
}

func GetScretKey() []byte {
	return SecretKey
}
