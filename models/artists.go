package models

import (
	"context"
	"fmt"
	"html"

	"github.com/jackc/pgx/v5"
)

type Artist struct {
	Id     int    `form:"id" json:"id"`
	Name   string `form:"name" binding:"required" json:"name"`
	UserID int    `form:"user_id" json:"user_id"`
}

func SearchArtists(conn *pgx.Conn, userID int, q string, page int) ([]Artist, int, error) {
	offset := (page - 1) * 10
	var query string
	var artists []Artist
	var rows pgx.Rows
	var totalCount int
	var err error
	if q == "" {
		query = "SELECT id, name, user_id FROM artists WHERE user_id = $1 ORDER BY name LIMIT 10 OFFSET $2;"
		rows, err = conn.Query(context.Background(), query, userID, offset)
		if err != nil {
			fmt.Printf("error: %v\n", err)
		}
	} else {
		query = "SELECT artist_id, name, user_id FROM artists WHERE user_id = $1 AND nane ILIKE '%' || $2 || '%' ORDER BY name LIMIT 10 OFFSET $3;"
		rows, err = conn.Query(context.Background(), query, userID, q, offset)
		if err != nil {
			fmt.Printf("error: %v\n", err)
		}
	}
	for rows.Next() {
		var artist Artist
		if err := rows.Scan(&artist.Id, &artist.Name, &artist.UserID); err != nil {
			return nil, 0, fmt.Errorf("error scanning row: %v", err)
		}
		artists = append(artists, artist)
	}
	if q == "" {
		if err := conn.QueryRow(context.Background(), "SELECT COUNT(*) FROM artists WHERE user_id = $1", userID).Scan(&totalCount); err != nil {
			return nil, 0, fmt.Errorf("error fetching total count: %v", err)
		}
	} else {
		if err := conn.QueryRow(context.Background(), "SELECT COUNT(*) FROM artists WHERE user_id = $1 AND name ILIKE '%' || $2 || '%';", userID, q).Scan(&totalCount); err != nil {
			return nil, 0, fmt.Errorf("error fetching total count: %v", err)
		}
	}
	return artists, totalCount, nil
}
func CreateArtist(conn *pgx.Conn, artist *Artist) error {
	var id int
	artist.Name = html.EscapeString(artist.Name)
	query := "INSERT INTO artists (name, user_id) VALUES($1, $2) RETURNING id"
	err := conn.QueryRow(context.Background(), query, artist.Name, artist.UserID).Scan(&id)
	if err != nil {
		return fmt.Errorf("error creating artist: %v", err)
	}
	artist.Id = id
	return nil
}
