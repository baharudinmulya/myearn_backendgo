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

type Akun struct {
	IDAkun   int    `json:"id_akun"`
	NamaAkun string `json:"nama_akun"`
}

func GetAkun(c *gin.Context, db *sql.DB) {
	rows, err := db.Query("SELECT * FROM akun")
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve transactions"})
		return
	}
	defer rows.Close()

	var akun []Akun
	for rows.Next() {
		var akun1 Akun
		err := rows.Scan(&akun1.IDAkun, &akun1.NamaAkun)
		if err != nil {
			log.Println(err)
			continue
		}

		akun = append(akun, akun1)
	}

	c.JSON(http.StatusOK, akun)
}

func TambahAkun(c *gin.Context, db *sql.DB) {
	var form Akun
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
	query := "INSERT INTO akun (nama_akun) VALUES (?)"
	_, err = db.Exec(query, form.NamaAkun)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create akun"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Akun created successfully"})
}

func EditAkun(c *gin.Context, db *sql.DB) {
	var form Akun
	if err := c.ShouldBindJSON(&form); err != nil {
		log.Println("Failed to bind JSON:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	// Update transaksi in the database
	query := "UPDATE akun SET nama_akun=? WHERE id_akun=?"
	_, err := db.Exec(query, form.NamaAkun, form.IDAkun)
	if err != nil {
		log.Println("Failed to update AKUN:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update AKUN"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "AKUN updated successfully"})
}

func DeleteAkun(c *gin.Context, db *sql.DB) {
	var form Akun
	if err := c.ShouldBindJSON(&form); err != nil {
		log.Println("Failed to bind JSON:", err)
		c.JSON(http.StatusBadRequest, gin.H{"Error": "Invalid Request"})
		return
	}

	// Delete transaksi in the database
	query := "DELETE FROM akun WHERE id_akun = ?"

	// Prepare the statement
	stmt, err := db.Prepare(query)
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	// Execute the query
	result, err := stmt.Exec(form.IDAkun)
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
