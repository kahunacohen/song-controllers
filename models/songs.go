package models

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
	"html"
	"regexp"
	"time"
)

type Song struct {
	CreatedAt time.Time `json:"created_at"`
	Capo      int
	Artist    string    `form:"artist" binding:"required" json:"artist"`
	Genre     string    `form:"genre" json:"genre"`
	Lyrics    string    `form:"lyrics" binding:"required" json:"lyrics"`
	Id        int       `form:"id" json:"id"`
	Title     string    `form:"title" binding:"required" json:"title"`
	UpdatedAt time.Time `json:"updated_at"`
	UserID    int       `form:"user_id" json:"user_id"`
}

func (s Song) Html() string {
	var ret string
	chordRe := regexp.MustCompile(`\[(.+?)\]`)
	ret = fmt.Sprintf(
		"<div style='position: relative;'>%s",
		chordRe.ReplaceAllString(html.EscapeString(s.Lyrics), "<span style='position:absolute;top:-12px;font-size:90%;font-weight:bold;'>$1</span>"))
	lineWrapperRe := regexp.MustCompile(`(?m)^.*$`)
	ret = lineWrapperRe.ReplaceAllString(ret, "<div style='position:relative;line-height:2.5;'>$0</div>")
	return ret
}
func SearchSongs(conn *pgx.Conn, userID int, q string, page int) ([]Song, int, error) {
	offset := (page - 1) * 10
	var query string
	var songs []Song
	var rows pgx.Rows
	var totalCount int
	var err error
	if q == "" {
		fmt.Println("blank")
		query = "SELECT song_id, user_id, title, genre, artist_name FROM songs_by_user WHERE user_id = $1 ORDER BY title LIMIT 10 OFFSET $2;"
		rows, err = conn.Query(context.Background(), query, userID, offset)
		if err != nil {
			fmt.Printf("error: %v\n", err)
		}
	} else {
		fmt.Println("SHIT")
		query = "SELECT song_id, user_id, title, genre, artist_name FROM songs_by_user WHERE user_id = $1 AND CONCAT(title, ' ', artist_name) ILIKE '%' || $2 || '%' ORDER BY title LIMIT 10 OFFSET $3;"
		rows, err = conn.Query(context.Background(), query, userID, q, offset)
		if err != nil {
			fmt.Printf("error: %v\n", err)
		}
	}
	for rows.Next() {
		var song Song
		if err := rows.Scan(&song.Id, &song.UserID, &song.Title, &song.Genre, &song.Artist); err != nil {
			return nil, 0, fmt.Errorf("error scanning row when getting songs: %v", err)
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
	query := "UPDATE songs SET title=$1, lyrics=$2 WHERE id=$3;"
	_, err := conn.Exec(context.Background(), query, song.Title, song.Lyrics, song.Id)
	if err != nil {
		return fmt.Errorf("error updating song: %v", err)
	}
	return nil
}
func CreateSong(conn *pgx.Conn, song *Song) error {
	var id int
	query := "INSERT INTO songs (title, lyrics, user_id, genre_id, artist_id) VALUES($1, $2, $3, $4, $5) RETURNING id"
	err := conn.QueryRow(context.Background(), query, song.Title, song.Lyrics, song.UserID, 1, song.Artist).Scan(&id)
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
		return err
	}
	return nil
}
