package endpoints

import (
	"encoding/json"
	"net/http"
	"vegeta-server/models"
	"vegeta-server/pkg/vegeta"

	"github.com/gin-gonic/gin"
)

func (e *Endpoints) GetReportEndpoint(c *gin.Context) {
	format := c.DefaultQuery("format", "json")
	resp := e.reporter.GetAllInFormat(vegeta.Format(format))

	if format == "json" {
		c.Header("Content-Type", "application/json")
		jsonReports := make([]models.JSONReportResponse, 0)

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
	} else if format == "text" {
		c.Header("Content-Type", "text/plain")
		c.JSON(http.StatusOK, resp)
	}
}

func (e *Endpoints) GetReportByIDEndpoint(c *gin.Context) {
	id := c.Param("attackID")
	format := c.DefaultQuery("format", "json")
	resp, err := e.reporter.GetInFormat(id, vegeta.Format(format))
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

	if format == "json" {
		c.Header("Content-Type", "application/json")
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
	} else if format == "text" {
		c.Header("Content-Type", "text/plain")
		c.JSON(http.StatusOK, resp)
	}
}
