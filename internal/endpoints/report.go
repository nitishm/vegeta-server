package endpoints

import (
	"encoding/json"
	"net/http"
	"vegeta-server/models"
	"vegeta-server/pkg/vegeta"

	"github.com/gin-gonic/gin"
)

// GetReportEndpoint implements a handler for the GET /api/v1/report endpoint
func (e *Endpoints) GetReportEndpoint(c *gin.Context) {
	resp := e.reporter.GetAll()

	jsonReports := make([]models.JSONReportResponse, 0)

	for _, report := range resp {
		var jsonReport models.JSONReportResponse
		err := json.Unmarshal(report, &jsonReport)
		if err != nil {
			ginErrInternalServerError(c, err)
			return
		}
		jsonReports = append(jsonReports, jsonReport)
	}

	c.JSON(http.StatusOK, jsonReports)
}

// GetReportByIDEndpoint implements a handler for the GET /api/v1/report/<attackID> endpoint
func (e *Endpoints) GetReportByIDEndpoint(c *gin.Context) {
	id := c.Param("attackID")

	format := vegeta.Format(c.DefaultQuery("format", "json"))
	bucket := c.DefaultQuery("bucket", "0,500ms,1s,1.5s,2s,2.5s,3s")
	format = format.StringsToFormat(string(format), bucket)

	resp, err := e.reporter.GetInFormat(id, format)
	if err != nil {
		ginErrNotFound(c, err)
		return
	}

	switch format.GetFormat() {
	case vegeta.JSONFormat:
		c.Header("Content-Type", "application/json")
		var jsonReport models.JSONReportResponse
		err = json.Unmarshal(resp, &jsonReport)
		if err != nil {
			ginErrInternalServerError(c, err)
		}
		c.JSON(http.StatusOK, jsonReport)
	case vegeta.TextFormat:
		c.Header("Content-Type", "text/plain")
		c.String(http.StatusOK, "%s", resp)
	case vegeta.BinaryFormat:
		c.Header("Content-Type", "application/octet-stream")
		c.Data(http.StatusOK, "application/octet-stream", resp)
	case vegeta.HistogramFormat:
		c.Header("Content-Type", "text/plain")
		c.String(http.StatusOK, "%s", resp)
	}
}
