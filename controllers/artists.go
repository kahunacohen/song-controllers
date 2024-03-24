package controllers

import (
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"github.com/kahunacohen/songctls/models"
)

func ListArtists(conn *pgx.Conn, responder ListResponder) gin.HandlerFunc {
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
		songs, totalCount, err := models.SearchArtists(conn, userIDAsInt, q, pageInt)
		if err != nil {
			fmt.Println(err)
		}
		responder(c, userID, songs, totalCount, pageInt, c.Query("q"), content == "partial")
	}
}
