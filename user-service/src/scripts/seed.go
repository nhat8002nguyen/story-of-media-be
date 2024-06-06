package scripts

import (
	"context"
	"fmt"
	"log"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/nhat8002nguyen/story-of-media-be/story-service/src/models"
	"golang.org/x/crypto/bcrypt"
)

// Placeholder data as per your JavaScript file
var users = []models.User{
	{
		ID:       "13c1f7d1-8c16-4d76-944d-6f6ae930c99d",
		Name:     "Nhat Nguyen",
		Email:    "nv.nhat8002@gmail.com",
		Password: "12345678",
	},
	{
		ID:       "23c1f7d1-8c16-4d76-944d-6f6ae930c99d",
		Name:     "Nathan",
		Email:    "kagaminguyendu123@gmail.com",
		Password: "123456",
	},
}

func seedUsers(pool *pgxpool.Pool) error {
	ctx := context.Background()
	_, err := pool.Exec(ctx, `CREATE EXTENSION IF NOT EXISTS "uuid-ossp"`)
	if err != nil {
		return fmt.Errorf("error creating extension: %w", err)
	}

	_, err = pool.Exec(ctx, `
    CREATE TABLE IF NOT EXISTS users (
        id UUID DEFAULT uuid_generate_v4() PRIMARY KEY,
        name VARCHAR(255) NOT NULL,
				email VARCHAR(255) NOT NULL CHECK (email <> '') UNIQUE,
        password TEXT NOT NULL
    );`)
	if err != nil {
		return fmt.Errorf("error creating users table: %w", err)
	}

	for _, user := range users {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), 10)
		if err != nil {
			return fmt.Errorf("error hashing password: %w", err)
		}

		_, err = pool.Exec(ctx, `
        INSERT INTO users (id, name, email, password)
        VALUES ($1, $2, $3, $4)
        ON CONFLICT (id) DO NOTHING`, user.ID, user.Name, user.Email, string(hashedPassword))
		if err != nil {
			return fmt.Errorf("error inserting user: %w", err)
		}
	}

	log.Printf("Seeded %d users\n", len(users))
	return nil
}

func SeedData() {
	dbURL := "postgres://nhatnguyen@localhost:5432/postgres"

	poolConfig, err := pgxpool.ParseConfig(dbURL)
	if err != nil {
		log.Fatalf("Unable to parse database URL: %v", err)
	}

	pool, err := pgxpool.ConnectConfig(context.Background(), poolConfig)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v", err)
	}
	defer pool.Close()

	if err := seedUsers(pool); err != nil {
		log.Fatalf("Error seeding users: %v", err)
	}

	log.Println("Database seeding completed successfully!")
}
