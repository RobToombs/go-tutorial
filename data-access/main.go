package main

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
	"os"
)

func main() {
	conn, err := pgx.Connect(context.Background(), "postgresql://postgres:postgres@localhost:5432/recordings")
	if err != nil {
		fmt.Printf("Unable to connect to database: %v\n", err)
		os.Exit(1)
	}

	defer conn.Close(context.Background())

	rows, err := conn.Query(context.Background(), "select * from album")
	if err != nil {
		fmt.Printf("QueryRow failed: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Connected!")
	for rows.Next() {
		row, _ := rows.Values()
		fmt.Println(row)
	}
}
