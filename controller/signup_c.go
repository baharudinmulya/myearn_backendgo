package controller

import (
	"myearn/models"

	"github.com/gin-gonic/gin"
)

func SignUpHandler(c *gin.Context) {
	models.SignUpHandler(c)
}
