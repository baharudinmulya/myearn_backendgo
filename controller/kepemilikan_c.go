package controller

import (
	"database/sql"
	"myearn/models"

	"github.com/gin-gonic/gin"
)

func GetKepemilikan(c *gin.Context, db *sql.DB) {
	models.GetKepemilikan(c, db)
}

func AddKepemilikan(c *gin.Context, db *sql.DB) {
	models.TambahKepemilikan(c, db)
}

func EditKepemilikan(c *gin.Context, db *sql.DB) {
	models.EditKepemilikan(c, db)
}

func DeleteKepemilikan(c *gin.Context, db *sql.DB) {
	models.DeleteKepemilikan(c, db)
}
