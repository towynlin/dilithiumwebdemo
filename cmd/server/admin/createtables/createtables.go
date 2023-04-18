package main

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5"
)

func main() {
	conn, err := pgx.Connect(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer conn.Close(context.Background())

	_, err = conn.Exec(context.Background(), "create table if not exists users (id uuid primary key, pubkey bytea not null)")
	if err != nil {
		panic(err)
	}
	_, err = conn.Exec(context.Background(), "create table if not exists jobs (id uuid primary key, user_id uuid not null)")
	if err != nil {
		panic(err)
	}
}
