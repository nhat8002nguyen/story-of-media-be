package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"
)

var dbURL string

type User struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

// Placeholder data as per your JavaScript file
var users = []User{
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

func SeedUserService() {
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

func seedStories(pool *pgxpool.Pool) error {
	ctx := context.Background()
	_, err := pool.Exec(ctx, `CREATE EXTENSION IF NOT EXISTS "uuid-ossp"`)
	if err != nil {
		return fmt.Errorf("error creating extension: %w", err)
	}

	sqlStmt := `CREATE TABLE IF NOT EXISTS chat_sessions (
		id SERIAL PRIMARY KEY,
		user_id TEXT NOT NULL,
		session_id TEXT NOT NULL,
		message TEXT NOT NULL,
		sender TEXT NOT NULL,
		timestamp TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);`

	_, err = pool.Exec(ctx, sqlStmt)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Seeded stories data.")
	return nil
}

func SeedStoryService() {
	poolConfig, err := pgxpool.ParseConfig(dbURL)
	if err != nil {
		log.Fatalf("Unable to parse database URL: %v", err)
	}

	pool, err := pgxpool.ConnectConfig(context.Background(), poolConfig)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v", err)
	}
	defer pool.Close()

	if err := seedStories(pool); err != nil {
		log.Fatalf("Error seeding users: %v", err)
	}

	if err := seedUploads(pool); err != nil {
		log.Fatalf("Error seeding users: %v", err)
	}

	log.Println("Database seeding completed successfully!")
}

func seedUploads(pool *pgxpool.Pool) error {
	ctx := context.Background()
	_, err := pool.Exec(ctx, `CREATE EXTENSION IF NOT EXISTS "uuid-ossp"`)
	if err != nil {
		return fmt.Errorf("error creating extension: %w", err)
	}

	sqlStmt := `CREATE TABLE IF NOT EXISTS session_files (
		id SERIAL PRIMARY KEY,
		user_id TEXT NOT NULL,
		session_id TEXT NOT NULL,
		filename VARCHAR(255) NOT NULL,
		content_type TEXT NOT NULL,
		file_data BYTEA NOT NULL,
		upload_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);`

	_, err = pool.Exec(ctx, sqlStmt)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Seeded stories data.")
	return nil
}

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file.")
	} else {
		log.Println("Loaded .env")
	}

	dbURL = os.Getenv("dbURL")

	SeedUserService()
	SeedStoryService()
}
