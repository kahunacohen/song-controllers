package controllers

import (
	"fmt"
	"html"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"github.com/kahunacohen/songctls/models"
)

type ListResponder func(context *gin.Context, userID string, songs []models.Song, totalCount int, page int, searchTerm string, partial bool)

func ListSongs(conn *pgx.Conn, responder ListResponder) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.Param("user_id")
		userIDAsInt, _ := strconv.Atoi(userID)
		q := c.Query("q")
		page := c.Query("page")
		content := c.Query("ct")
		pageInt, err := strconv.Atoi(page)
		if err != nil {
			pageInt = 1
		}
		songs, totalCount, err := models.SearchSongs(conn, userIDAsInt, q, pageInt)
		if err != nil {
			fmt.Println(err)
		}
		responder(c, userID, songs, totalCount, pageInt, c.Query("q"), content == "partial")
	}
}

type ReadResponder func(context *gin.Context, mode string, song models.Song, artists []models.Artist, uri string, editModeUri string)

func ReadSong(conn *pgx.Conn, responder ReadResponder) gin.HandlerFunc {
	return func(c *gin.Context) {
		songID := c.Param("song_id")
		songIDAsInt, _ := strconv.Atoi(songID)
		userID := c.Param("user_id")
		userIdAsInt, _ := strconv.Atoi(userID)
		song, getSongErr := models.GetSongByID(conn, songIDAsInt)
		if getSongErr != nil {
			// templates.Render(c, templates.Base("Not found", templates.NotFound()))
			// return
			return
		}
		uri := fmt.Sprintf("/users/%s/songs/%d", userID, song.Id)
		editModeUri := fmt.Sprintf("%s?mode=edit", uri)
		mode := c.Query("mode")
		artists, _, _ := models.SearchArtists(conn, userIdAsInt, nil, nil)
		responder(c, mode, *song, artists, uri, editModeUri)
	}
}

type UpdateResponder func(context *gin.Context, song models.Song, artists []models.Artist)

func UpdateSong(conn *pgx.Conn, responder UpdateResponder) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.Param("user_id")
		var song models.Song
		c.Bind(&song)
		song.Title = html.EscapeString(song.Title)
		err := models.UpdateSong(conn, &song)
		if err != nil {
			// @TODO error handling.
			fmt.Printf("error updating song: %v\n", err)

		}
		uri := fmt.Sprintf("/users/%s/songs/%d?flashOn=true&flashMsg=Song%%20saved", userID, song.Id)
		if c.Request.Method == "POST" {
			// We are receiving from old-school form where method=POST
			// is not supported by browsers, so redirect to same page
			// with a GET.
			c.Redirect(http.StatusSeeOther, uri)
			return
		}
		var artists []models.Artist
		responder(c, song, artists)
	}
}
