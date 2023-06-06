package models

import (
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"os"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/dgrijalva/jwt-go"
	"net/http"
	"time"
)

type Claims struct {
	UserID int `json:"user_id"`
	jwt.StandardClaims
}

type Credentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func LoginHandler(c *gin.Context) {
	var creds Credentials
	if err := c.ShouldBindJSON(&creds); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	// Hash the password
	hashedPassword := sha256.Sum256([]byte(creds.Password))
	hashedPasswordStr := hex.EncodeToString(hashedPassword[:])

	// Database connection
	db, err := sql.Open(os.Getenv("DB_TYPE"), os.Getenv("DB_CONNECTION"));
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to connect to database"})
		return
	}
	defer db.Close()

	// Fetch user from the database
	var user User
	query := "SELECT * FROM users WHERE username=? AND password=?"
	if err := db.QueryRow(query, creds.Username, hashedPasswordStr).Scan(&user.ID, &user.Name, &user.Password, &user.Created); err == nil {
		fmt.Println("Query Error:", err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials", "id":user.ID,"name":user.Name,"err":hashedPasswordStr})
		return
	}

	// Create the JWT token
	claims := &Claims{
		UserID: user.ID, // Example user ID
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 24).Unix(), // Token expiration time
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	jwtSecret := []byte(os.Getenv("JWT_TOKEN")) // Replace with your own secret key
	signedToken, err := token.SignedString(jwtSecret)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": signedToken})
}
