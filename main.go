package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"time"

	"myearn/controller"
	"myearn/middleware"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

type User struct {
	ID       int            `json:"user_id"`
	Name     string         `json:"username"`
	Password sql.NullString `json:"password"`
	Created  time.Time      `json:"created"`
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
		os.Exit(1)
	}
	defer db.Close()

	// Check the database connection
	err = db.Ping()
	if err != nil {
		log.Fatal(err)
		os.Exit(1) // Exit the script with exit code 1
	}

	// Gin router
	r := gin.Default()

	// Apply CORS middleware
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"*"}
	config.AllowMethods = []string{"GET", "POST", "OPTIONS"}
	config.AllowHeaders = []string{"Authorization", "Content-Type"}
	// config.AllowPreflight = false // Disable handling preflight OPTIONS requests
	r.Use(cors.New(config))

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	r.POST("/login", controller.LoginHandler)
	r.POST("/signup", controller.SignUpHandler)

	auth := r.Group("/api")
	auth.Use(middleware.AuthMiddleware())
	auth.Use(cors.Default())

	auth.GET("/users", func(c *gin.Context) {
		controller.GetUsersHandler(c, db)
	})

	// transaksi
	auth.GET("/transaksi", func(c *gin.Context) {
		controller.GetTransaksi(c, db)
	})
	auth.POST("/addtransaksi", func(c *gin.Context) {
		controller.AddTransaksi(c, db)
	})
	auth.POST("/edittransaksi", func(c *gin.Context) {
		controller.EditTransaksi(c, db)
	})
	auth.DELETE("/deletetransaksi", func(c *gin.Context) {
		controller.DeleteTransaksi(c, db)
	})

	// akun
	auth.GET("/akun", func(c *gin.Context) {
		controller.GetAkun(c, db)
	})
	auth.POST("/addakun", func(c *gin.Context) {
		controller.AddAkun(c, db)
	})
	auth.POST("/editakun", func(c *gin.Context) {
		controller.EditAkun(c, db)
	})
	auth.DELETE("/deleteakun", func(c *gin.Context) {
		controller.DeleteAkun(c, db)
	})

	//kepemilikan
	auth.GET("/milik", func(c *gin.Context) {
		controller.GetKepemilikan(c, db)
	})
	auth.POST("/addmilik", func(c *gin.Context) {
		controller.AddKepemilikan(c, db)
	})
	auth.POST("/editmilik", func(c *gin.Context) {
		controller.EditKepemilikan(c, db)
	})
	auth.DELETE("/deletemilik", func(c *gin.Context) {
		controller.DeleteKepemilikan(c, db)
	})

	//allmoney
	auth.GET("/counttransaksi", func(c *gin.Context) {
		controller.CountTransaksi(c, db)
	})

	r.Run(":8080")
}
