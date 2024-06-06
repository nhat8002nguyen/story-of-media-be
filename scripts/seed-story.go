package main

import (
	"context"
	"fmt"
	"log"

	"github.com/jackc/pgx/v4/pgxpool"
)

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

	if err := seedStories(pool); err != nil {
		log.Fatalf("Error seeding users: %v", err)
	}

	log.Println("Database seeding completed successfully!")
}
