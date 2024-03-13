package models

import (
	"context"
	"fmt"

	"github.com/jackc/pgx"
)

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
