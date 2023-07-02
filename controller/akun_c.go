package controller

import (
	"database/sql"
	"myearn/models"

	"github.com/gin-gonic/gin"
)

func GetAkun(c *gin.Context, db *sql.DB) {
	models.GetAkun(c, db)
}
func AddAkun(c *gin.Context, db *sql.DB) {
	models.TambahAkun(c, db)
}

func EditAkun(c *gin.Context, db *sql.DB) {
	models.EditAkun(c, db)
}

func DeleteAkun(c *gin.Context, db *sql.DB) {
	models.DeleteAkun(c, db)
}
