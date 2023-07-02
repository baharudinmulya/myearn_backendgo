package controller

import (
	"database/sql"
	"myearn/models"

	"github.com/gin-gonic/gin"
)

func GetTransaksi(c *gin.Context, db *sql.DB) {
	models.GetTransaksi(c, db)
}

func CountTransaksi(c *gin.Context, db *sql.DB) {
	models.CountTransaksi(c, db)
}

func AddTransaksi(c *gin.Context, db *sql.DB) {
	models.TambahTransaksi(c, db)
}

func EditTransaksi(c *gin.Context, db *sql.DB) {
	models.EditTransaksi(c, db)
}

func DeleteTransaksi(c *gin.Context, db *sql.DB) {
	models.DeleteTransaksi(c, db)
}
