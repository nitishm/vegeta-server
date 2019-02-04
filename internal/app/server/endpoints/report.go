package endpoints

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (e *Endpoints) GetReportEndpoint(c *gin.Context) {
	c.JSON(
		http.StatusOK,
		gin.H{
			"code":     http.StatusOK,
			"endpoint": "GET /report",
		},
	)
}

func (e *Endpoints) GetReportByIDEndpoint(c *gin.Context) {
	c.JSON(
		http.StatusOK,
		gin.H{
			"code":     http.StatusOK,
			"endpoint": "GET /report/" + c.Param("attackID"),
		},
	)
}
