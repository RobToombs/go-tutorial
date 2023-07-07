// Package dataAccess - swap this to 'package main' to run alone
package dataAccess

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
	"log"
	"os"
)

type Album struct {
	ID     int64   `json:"id"`
	Title  string  `json:"title"`
	Artist string  `json:"artist"`
	Price  float32 `json:"price"`
}

var conn *pgx.Conn

func main() {
	EstablishConnection()
	defer CloseConnection()

	albums, err := AlbumsByArtist("John Coltrane")
	if err != nil {
		fmt.Printf("Error collecting albums by artist: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("Albums found: %v\n", albums)

	// Hard-code ID 2 here to test the query.
	alb, err := AlbumByID(2)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Album found: %v\n", alb)

	album := Album{
		Title:  "The Modern Sound of Betty Carter",
		Artist: "Betty Carter",
		Price:  49.99,
	}

	err = AddAlbumPtr(&album)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("ID of added album: %v\n", album.ID)
}

func EstablishConnection() {
	if conn == nil {
		var err error
		conn, err = pgx.Connect(context.Background(), "postgresql://postgres:postgres@localhost:5432/recordings")
		if err != nil {
			fmt.Printf("Unable to connect to database: %v\n", err)
			os.Exit(1)
		}
	}
}

func CloseConnection() {
	if conn != nil {
		conn.Close(context.Background())
	}
}

// AlbumsByArtist queries for albums that have the specified artist name.
func AlbumsByArtist(name string) ([]Album, error) {
	// An albums slice to hold data from returned rows.
	var albums []Album

	rows, err := conn.Query(context.Background(), "SELECT * FROM album WHERE artist = $1", name)
	if err != nil {
		return nil, fmt.Errorf("AlbumsByArtist %q: %v", name, err)
	}

	defer rows.Close()

	// Loop through rows, using Scan to assign column data to struct fields.
	for rows.Next() {
		var alb Album
		err = rows.Scan(&alb.ID, &alb.Title, &alb.Artist, &alb.Price)
		if err != nil {
			return nil, fmt.Errorf("AlbumsByArtist %q: %v", name, err)
		}
		albums = append(albums, alb)
	}

	err = rows.Err()
	if err != nil {
		return nil, fmt.Errorf("AlbumsByArtist %q: %v", name, err)
	}

	return albums, nil
}

func AlbumByID(id int64) (Album, error) {
	row := conn.QueryRow(context.Background(), "SELECT * FROM album WHERE id = $1", id)
	alb := rowToAlbum(row)

	return alb, nil
}

// AddAlbum adds the specified album to the database,
// returning the album ID of the new entry
func AddAlbum(alb Album) (int64, error) {
	//var result Album
	err := conn.QueryRow(context.Background(), "INSERT INTO album (title, artist, price) VALUES ($1, $2, $3) RETURNING id", alb.Title, alb.Artist, alb.Price).Scan(&alb.ID)
	if err != nil {
		return 0, fmt.Errorf("AddAlbum: %v", err)
	}

	return alb.ID, nil
}

// AddAlbumPtr adds the specified album to the database,
// but updates the provided album ID
func AddAlbumPtr(alb *Album) error {
	row := conn.QueryRow(context.Background(), "INSERT INTO album (title, artist, price) VALUES ($1, $2, $3) RETURNING id", alb.Title, alb.Artist, alb.Price)
	err := row.Scan(&alb.ID)
	if err != nil {
		return fmt.Errorf("AddAlbumPtr: %v", err)
	}

	return nil
}

func Albums() ([]Album, error) {
	rows, err := conn.Query(context.Background(), "SELECT * FROM album")
	if err != nil {
		return nil, fmt.Errorf("Albums: %v", err)
	}

	// An albums slice to hold data from returned rows.
	var albums []Album

	// Loop through rows, using Scan to assign column data to struct fields.
	for rows.Next() {
		alb := rowsToAlbum(rows)
		albums = append(albums, alb)
	}

	return albums, nil
}

func rowsToAlbum(row pgx.Rows) Album {
	var alb Album
	err := row.Scan(&alb.ID, &alb.Title, &alb.Artist, &alb.Price)
	if err != nil {
		log.Fatalf("Albums: converting to albums %v", err)
	}

	return alb
}

func rowToAlbum(row pgx.Row) Album {
	var alb Album
	err := row.Scan(&alb.ID, &alb.Title, &alb.Artist, &alb.Price)
	if err != nil {
		log.Fatalf("Albums: converting to albums %v", err)
	}

	return alb
}
