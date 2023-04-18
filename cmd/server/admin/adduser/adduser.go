package main

import (
	"context"
	"encoding/base64"
	"fmt"
	"os"

	"github.com/cloudflare/circl/sign/eddilithium2"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

func main() {
	conn, err := pgx.Connect(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer conn.Close(context.Background())

	pk, sk, err := eddilithium2.GenerateKey(nil)
	if err != nil {
		panic(err)
	}
	pkBytes, err := pk.MarshalBinary()
	if err != nil {
		panic(err)
	}

	newId := uuid.NewString()
	_, err = conn.Exec(context.Background(), "insert into users (id, pubkey) values ($1, $2)", newId, pkBytes)
	if err != nil {
		panic(err)
	}

	pkString := base64.RawURLEncoding.EncodeToString(pk.Bytes())
	skString := base64.RawURLEncoding.EncodeToString(sk.Bytes())
	fmt.Printf("\nid: %s\n\npk: %s\n\nsk: %s\n", newId, pkString, skString)
}
