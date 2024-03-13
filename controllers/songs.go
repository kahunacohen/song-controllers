package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx"
)

func ListSongs(conn *pgx.Conn) gin.HandlerFunc {
	return func(c *gin.Context) {

		// fmt.Println(songctls.ListSongs)
		// userID := c.Param("user_id")
		// userIDAsInt, _ := strconv.Atoi(userID)
		// q := c.Query("q")
		// page := c.Query("page")
		// content := c.Query("ct")
		// fmt.Println(content)
		// pageInt, err := strconv.Atoi(page)
		// if err != nil {
		// 	pageInt = 1
		// }
		// songs, totalCount, err := models.SearchSongs(conn, userIDAsInt, q, pageInt)
		// if err != nil {
		// 	fmt.Println(err)
		// }
	}
}
