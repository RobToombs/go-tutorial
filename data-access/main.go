package main

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/jackc/pgx/v5"
	"log"
	"os"
)

type Album struct {
	ID     int64
	Title  string
	Artist string
	Price  float32
}

var conn *pgx.Conn

func main() {
	var err error
	conn, err = pgx.Connect(context.Background(), "postgresql://postgres:postgres@localhost:5432/recordings")
	if err != nil {
		fmt.Printf("Unable to connect to database: %v\n", err)
		os.Exit(1)
	}

	defer conn.Close(context.Background())

	albums, err := albumsByArtist("John Coltrane")
	if err != nil {
		fmt.Printf("Error collecting albums by artist: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("Albums found: %v\n", albums)

	// Hard-code ID 2 here to test the query.
	alb, err := albumByID(2)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Album found: %v\n", alb)

	album := Album{
		Title:  "The Modern Sound of Betty Carter",
		Artist: "Betty Carter",
		Price:  49.99,
	}

	err = addAlbumPtr(&album)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("ID of added album: %v\n", album.ID)
}

// albumsByArtist queries for albums that have the specified artist name.
func albumsByArtist(name string) ([]Album, error) {
	// An albums slice to hold data from returned rows.
	var albums []Album

	rows, err := conn.Query(context.Background(), "SELECT * FROM album WHERE artist = $1", name)
	if err != nil {
		return nil, fmt.Errorf("albumsByArtist %q: %v", name, err)
	}

	defer rows.Close()

	// Loop through rows, using Scan to assign column data to struct fields.
	for rows.Next() {
		var alb Album
		if err := rows.Scan(&alb.ID, &alb.Title, &alb.Artist, &alb.Price); err != nil {
			return nil, fmt.Errorf("albumsByArtist %q: %v", name, err)
		}
		albums = append(albums, alb)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("albumsByArtist %q: %v", name, err)
	}

	return albums, nil
}

func albumByID(id int64) (Album, error) {
	// An album to hold data from the returned row.
	var alb Album

	row := conn.QueryRow(context.Background(), "SELECT * FROM album WHERE id = $1", id)
	err := row.Scan(&alb.ID, &alb.Title, &alb.Artist, &alb.Price)
	if err != nil {
		if err == sql.ErrNoRows {
			return alb, fmt.Errorf("albumsById %d: no such album", id)
		}

		return alb, fmt.Errorf("albumsById %d: %v", id, err)
	}

	return alb, nil
}

// addAlbum adds the specified album to the database,
// returning the album ID of the new entry
func addAlbum(alb Album) (int64, error) {
	//var result Album
	err := conn.QueryRow(context.Background(), "INSERT INTO album (title, artist, price) VALUES ($1, $2, $3) RETURNING id", alb.Title, alb.Artist, alb.Price).Scan(&alb.ID)
	if err != nil {
		return 0, fmt.Errorf("addAlbum: %v", err)
	}

	alb.Artist = "FARTIST"

	return alb.ID, nil
}

// addAlbumPtr adds the specified album to the database,
// but updates the provided album ID
func addAlbumPtr(alb *Album) error {
	row := conn.QueryRow(context.Background(), "INSERT INTO album (title, artist, price) VALUES ($1, $2, $3) RETURNING id", alb.Title, alb.Artist, alb.Price)
	err := row.Scan(&alb.ID)
	if err != nil {
		return fmt.Errorf("addAlbum: %v", err)
	}

	return nil
}
