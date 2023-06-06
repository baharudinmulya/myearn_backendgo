package models

import (
	"database/sql"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type User struct {
	ID       int       `json:"user_id"`
	Name     string    `json:"username"`
	Password string    `json:"password"`
	Created  time.Time `json:"created"`
}

func GetUser(c *gin.Context, db *sql.DB) {
	rows, err := db.Query("SELECT * FROM users")
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve users"})
		return
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		var user User
		var createdStr string
		err := rows.Scan(&user.ID, &user.Name, &user.Password, &createdStr)
		if err != nil {
			log.Println(err)
			continue
		}

		user.Created, err = time.Parse("2006-01-02 15:04:05", createdStr)
		if err != nil {
			log.Println(err)
			continue
		}

		users = append(users, user)
	}

	c.JSON(http.StatusOK, users)
}
