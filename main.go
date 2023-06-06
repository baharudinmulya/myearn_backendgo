package main

import (
	"database/sql"
	"log"
	"net/http"
	"time"
	"os"

	"myearn/controller"
	"myearn/middleware"
	"github.com/joho/godotenv"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

type User struct {
	ID       int       `json:"user_id"`
	Name     string    `json:"username"`
	Password sql.NullString `json:"password"`
	Created  time.Time `json:"created"`
}

func main() {
	// Load environment variables from .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Database connection
	db, err := sql.Open(os.Getenv("DB_TYPE"), os.Getenv("DB_CONNECTION"))
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Gin router
	r := gin.Default()

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	r.POST("/login", controller.LoginHandler)

	auth := r.Group("/api")
	auth.Use(middleware.AuthMiddleware())

	auth.GET("/users", func(c *gin.Context) {
		controller.GetUsersHandler(c, db)
	})
	
	r.Run(":8080")
}
