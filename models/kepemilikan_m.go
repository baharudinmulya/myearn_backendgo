package models

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql" // Import the MariaDB driver
)

type Kepemilikan struct {
	IDMilik   int    `json:"id_milik"`
	NamaMilik string `json:"nama_milik"`
}

func GetKepemilikan(c *gin.Context, db *sql.DB) {
	rows, err := db.Query("SELECT * FROM kepemilikan")
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve transactions"})
		return
	}
	defer rows.Close()

	var kepemilikan []Kepemilikan
	for rows.Next() {
		var kepemilikan1 Kepemilikan
		err := rows.Scan(&kepemilikan1.IDMilik, &kepemilikan1.NamaMilik)
		if err != nil {
			log.Println(err)
			continue
		}

		kepemilikan = append(kepemilikan, kepemilikan1)
	}

	c.JSON(http.StatusOK, kepemilikan)
}

func TambahKepemilikan(c *gin.Context, db *sql.DB) {
	var form Kepemilikan
	if err := c.ShouldBindJSON(&form); err != nil {
		log.Println("Failed to bind JSON:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	// Database connection
	db, err := sql.Open(os.Getenv("DB_TYPE"), os.Getenv("DB_CONNECTION"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to connect to the database"})
		return
	}
	defer db.Close()

	// Insert user into the database
	query := "INSERT INTO Kepemilikan (nama_milik) VALUES (?)"
	_, err = db.Exec(query, form.NamaMilik)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create Kepemilikan"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Kepemilikan created successfully"})
}

func EditKepemilikan(c *gin.Context, db *sql.DB) {
	var form Kepemilikan
	if err := c.ShouldBindJSON(&form); err != nil {
		log.Println("Failed to bind JSON:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	// Update transaksi in the database
	query := "UPDATE Kepemilikan SET nama_milik=? WHERE id_milik=?"
	_, err := db.Exec(query, form.NamaMilik, form.IDMilik)
	if err != nil {
		log.Println("Failed to update Kepemilikan:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update Kepemilikan"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Kepemilikan updated successfully"})
}

func DeleteKepemilikan(c *gin.Context, db *sql.DB) {
	var form Kepemilikan
	if err := c.ShouldBindJSON(&form); err != nil {
		log.Println("Failed to bind JSON:", err)
		c.JSON(http.StatusBadRequest, gin.H{"Error": "Invalid Request"})
		return
	}

	// Delete transaksi in the database
	query := "DELETE FROM Kepemilikan WHERE id_milik = ?"

	// Prepare the statement
	stmt, err := db.Prepare(query)
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	// Execute the query
	result, err := stmt.Exec(form.IDMilik)
	if err != nil {
		log.Fatal(err)
	}

	// Check if any rows were affected
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Fatal(err)
	}

	// Check if the record was deleted successfully
	if rowsAffected > 0 {
		fmt.Println("Record deleted successfully")
		c.JSON(http.StatusOK, gin.H{"message": "Record deleted successfully"})
	} else {
		fmt.Println("No record found with the specified ID")
		c.JSON(http.StatusNotFound, gin.H{"message": "No record found with the specified ID"})
	}
}
