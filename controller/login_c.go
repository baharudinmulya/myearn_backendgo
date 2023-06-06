package controller

import (
	"github.com/gin-gonic/gin"
	"myearn/models"
)


func LoginHandler(c *gin.Context) {
	models.LoginHandler(c)
}
