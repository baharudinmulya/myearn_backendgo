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

type Transaksi struct {
	ID        int    `json:"id_trx"`
	NamaTrx   string `json:"nama_trx"`
	IDMilik   int    `json:"id_milik"`
	IDAkun    int    `json:"id_akun"`
	TglTrx    string `json:"tgl_trx"`
	Value     string `json:"value"`
	Created   string `json:"created"`
	Updated   string `json:"updated"`
	NamaMilik string `json:"nama_milik"`
	NamaAkun  string `json:"nama_akun"`
}

type Count_Transaksi struct {
	Hasil string `json:"difference"`
}

type AddTransaksi struct {
	NamaTransaksi string `json:"nama_trx"`
	IDMilik       int    `json:"id_milik"`
	IDAkun        int    `json:"id_akun"`
	TglTrx        string `json:"tgl_trx"`
	Value         int    `json:"value"`
}

type Edit_Transaksi struct {
	NamaTransaksi string `json:"nama_trx"`
	IDMilik       int    `json:"id_milik"`
	TglTrx        string `json:"tgl_trx"`
	Value         int    `json:"value"`
	IDTrx         int    `json:"id_trx"`
}

type Delete_Transaksi struct {
	ID string `json:"id_trx"`
}

func GetTransaksi(c *gin.Context, db *sql.DB) {
	rows, err := db.Query("SELECT transaksi.*, kepemilikan.nama_milik AS nama_milik, akun.nama_akun AS nama_akun FROM transaksi JOIN kepemilikan ON transaksi.id_milik = kepemilikan.id_milik JOIN akun on akun.id_akun = transaksi.id_akun ORDER BY transaksi.tgl_trx DESC")
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve transactions"})
		return
	}
	defer rows.Close()

	var transaksi []Transaksi
	for rows.Next() {
		var transaksi1 Transaksi
		err := rows.Scan(&transaksi1.ID, &transaksi1.NamaTrx, &transaksi1.IDMilik, &transaksi1.IDAkun, &transaksi1.TglTrx, &transaksi1.Value, &transaksi1.Created, &transaksi1.Updated, &transaksi1.NamaMilik, &transaksi1.NamaAkun)
		if err != nil {
			log.Println(err)
			continue
		}

		transaksi = append(transaksi, transaksi1)
	}

	c.JSON(http.StatusOK, transaksi)
}

func TambahTransaksi(c *gin.Context, db *sql.DB) {
	var form AddTransaksi
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
	query := "INSERT INTO transaksi (nama_trx, id_milik, tgl_trx, value,id_akun) VALUES (?, ?, ?, ?,?)"
	_, err = db.Exec(query, form.NamaTransaksi, form.IDMilik, form.TglTrx, form.Value, form.IDAkun)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create transaksi"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Transaksi created successfully"})
}

func EditTransaksi(c *gin.Context, db *sql.DB) {
	var form Edit_Transaksi
	if err := c.ShouldBindJSON(&form); err != nil {
		log.Println("Failed to bind JSON:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	// Update transaksi in the database
	query := "UPDATE transaksi SET nama_trx=?, id_milik=?, tgl_trx=?, value=? WHERE id_trx=?"
	_, err := db.Exec(query, form.NamaTransaksi, form.IDMilik, form.TglTrx, form.Value, form.IDTrx)
	if err != nil {
		log.Println("Failed to update transaksi:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update transaksi"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Transaksi updated successfully"})
}

func DeleteTransaksi(c *gin.Context, db *sql.DB) {
	var form Delete_Transaksi
	if err := c.ShouldBindJSON(&form); err != nil {
		log.Println("Failed to bind JSON:", err)
		c.JSON(http.StatusBadRequest, gin.H{"Error": "Invalid Request"})
		return
	}

	// Delete transaksi in the database
	query := "DELETE FROM transaksi WHERE id_trx = ?"

	// Prepare the statement
	stmt, err := db.Prepare(query)
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	// Execute the query
	result, err := stmt.Exec(form.ID)
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

func CountTransaksi(c *gin.Context, db *sql.DB) {
	rows, err := db.Query("SELECT (SELECT SUM(value) FROM transaksi WHERE id_akun = 1) - (SELECT SUM(value) FROM transaksi WHERE id_akun = 2) AS difference;")
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve transactions"})
		return
	}
	defer rows.Close()

	var transaksi []Count_Transaksi
	for rows.Next() {
		var transaksi1 Count_Transaksi
		err := rows.Scan(&transaksi1.Hasil)
		if err != nil {
			log.Println(err)
			continue
		}

		transaksi = append(transaksi, transaksi1)
	}

	c.JSON(http.StatusOK, transaksi)
}
