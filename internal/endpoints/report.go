package endpoints

import (
	"encoding/json"
	"net/http"
	"vegeta-server/models"

	"github.com/gin-gonic/gin"
)

func (e *Endpoints) GetReportEndpoint(c *gin.Context) {
	resp := e.reporter.GetAll()

	var jsonReports []models.JSONReportResponse

	for _, report := range resp {
		var jsonReport models.JSONReportResponse
		err := json.Unmarshal([]byte(report), &jsonReport)
		if err != nil {
			c.JSON(
				http.StatusInternalServerError,
				gin.H{
					"message": "Failed to decode",
					"code":    http.StatusInternalServerError,
					"error":   err.Error(),
				},
			)
			return
		}
		jsonReports = append(jsonReports, jsonReport)
	}

	c.JSON(http.StatusOK, jsonReports)
}

func (e *Endpoints) GetReportByIDEndpoint(c *gin.Context) {
	id := c.Param("attackID")
	resp, err := e.reporter.Get(id)
	if err != nil {
		c.JSON(
			http.StatusNotFound,
			gin.H{
				"message": "Not found",
				"code":    http.StatusNotFound,
				"error":   err.Error(),
			},
		)
		return
	}

	var jsonReport models.JSONReportResponse
	err = json.Unmarshal([]byte(resp), &jsonReport)
	if err != nil {
		c.JSON(
			http.StatusInternalServerError,
			gin.H{
				"message": "Failed to decode",
				"code":    http.StatusInternalServerError,
				"error":   err.Error(),
			},
		)
		return
	}

	c.JSON(http.StatusOK, jsonReport)
}
