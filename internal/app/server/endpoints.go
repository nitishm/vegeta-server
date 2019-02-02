package server

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func PostAttackEndpoint(c *gin.Context) {
	c.JSON(
		http.StatusOK,
		gin.H{
			"code":     http.StatusOK,
			"endpoint": "POST /attack",
		},
	)
}

func GetAttackEndpoint(c *gin.Context) {
	c.JSON(
		http.StatusOK,
		gin.H{
			"code":     http.StatusOK,
			"endpoint": "GET /attack",
		},
	)
}

func GetAttackByIDEndpoint(c *gin.Context) {
	c.JSON(
		http.StatusOK,
		gin.H{
			"code":     http.StatusOK,
			"endpoint": "GET /attack/" + c.Param("id"),
		},
	)
}
