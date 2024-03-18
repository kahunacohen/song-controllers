package models

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
)

type Song struct {
	CreatedAt time.Time `json:"created_at"`
	Capo      int
	Artist    string    `json:"artist"`
	Genre     string    `form:"genre" json:"genre"`
	Lyrics    string    `form:"lyrics" binding:"required" json:"lyrics"`
	Id        int       `form:"id" json:"id"`
	Title     string    `form:"title" binding:"required" json:"title"`
	UpdatedAt time.Time `json:"updated_at"`
	UserID    int       `form:"user_id" json:"user_id"`
}

func SearchSongs(conn *pgx.Conn, userID int, q string, page int) ([]Song, int, error) {
	offset := (page - 1) * 10
	var query string
	var songs []Song
	var rows pgx.Rows
	var totalCount int
	var err error
	if q == "" {
		query = "SELECT song_id, user_id, title, genre, artist_name FROM songs_by_user WHERE user_id = $1 ORDER BY title LIMIT 10 OFFSET $2;"
		rows, err = conn.Query(context.Background(), query, userID, offset)
		if err != nil {
			fmt.Printf("error: %v\n", err)
		}
	} else {
		query = "SELECT song_id, user_id, title, genre, artist_name FROM songs_by_user WHERE user_id = $1 AND CONCAT(title, ' ', artist_name) ILIKE '%' || $2 || '%' ORDER BY title LIMIT 10 OFFSET $3;"
		rows, err = conn.Query(context.Background(), query, userID, q, offset)
		if err != nil {
			fmt.Printf("error: %v\n", err)
		}
	}
	for rows.Next() {
		var song Song
		if err := rows.Scan(&song.Id, &song.UserID, &song.Title, &song.Genre, &song.Artist); err != nil {
			return nil, 0, fmt.Errorf("error scanning row: %v", err)
		}
		songs = append(songs, song)
	}
	if q == "" {
		if err := conn.QueryRow(context.Background(), "SELECT COUNT(*) FROM songs_by_user WHERE user_id = $1", userID).Scan(&totalCount); err != nil {
			return nil, 0, fmt.Errorf("error fetching total count: %v", err)
		}
	} else {
		if err := conn.QueryRow(context.Background(), "SELECT COUNT(*) FROM songs_by_user WHERE user_id = $1 AND CONCAT(title, ' ', artist_name) ILIKE '%' || $2 || '%';", userID, q).Scan(&totalCount); err != nil {
			return nil, 0, fmt.Errorf("error fetching total count: %v", err)
		}
	}
	return songs, totalCount, nil
}

func GetSongByID(conn *pgx.Conn, id int) (*Song, error) {
	query := "SELECT song_id, title, genre, lyrics, artist_name FROM songs_by_user WHERE song_id = $1;"
	row := conn.QueryRow(context.Background(), query, id)
	var song Song
	if err := row.Scan(&song.Id, &song.Title, &song.Genre, &song.Lyrics, &song.Artist); err != nil {
		return nil, fmt.Errorf("error scanning row: %v", err)
	}
	return &song, nil
}
func UpdateSong(conn *pgx.Conn, song *Song) error {
	query := "UPDATE songs SET title='$1', lyrics='$2' WHERE id=$3;"
	_, err := conn.Exec(context.Background(), query, song.Title, song.Lyrics, song.Id)
	if err != nil {
		return fmt.Errorf("error updating song: %v", err)
	}
	return nil
}
func CreateSong(conn *pgx.Conn, song *Song) error {
	var id int
	fmt.Println(song.Genre)
	query := "INSERT INTO songs (title, lyrics, user_id, genre_id, artist_id) VALUES($1, $2, $3, $4, $5) RETURNING id"
	err := conn.QueryRow(context.Background(), query, song.Title, song.Lyrics, song.UserID, 1, 2).Scan(&id)
	if err != nil {
		return fmt.Errorf("error creating song: %v", err)
	}
	song.Id = id
	return nil
}
func DeleteSong(conn *pgx.Conn, songID int) error {
	query := "DELETE FROM songs WHERE id=$1"
	_, err := conn.Exec(context.Background(), query, songID)
	if err != nil {
		fmt.Println("error here!")
		return err
	}
	return nil
}
