package controller

import (
	"database/sql"

	"github.com/gin-gonic/gin"
	"myearn/models"
)

func GetUsersHandler(c *gin.Context, db *sql.DB) {
	models.GetUser(c, db)
}

// Other controller functions...

