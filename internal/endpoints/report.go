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

	format := vegeta.NewFormat(c.DefaultQuery("format", "json"))
	bucket := c.DefaultQuery("bucket", vegeta.DefaultBucketString)
	format.SetMeta("bucket", bucket)

	resp, err := e.reporter.GetInFormat(id, format)
	if err != nil {
		ginErrNotFound(c, err)
		return
	}

	switch format.String() {
	case vegeta.JSONFormatString:
		c.Header("Content-Type", "application/json")
		var jsonReport models.JSONReportResponse
		err = json.Unmarshal(resp, &jsonReport)
		if err != nil {
			ginErrInternalServerError(c, err)
		}
		c.JSON(http.StatusOK, jsonReport)
	case vegeta.TextFormatString:
		c.Header("Content-Type", "text/plain")
		c.String(http.StatusOK, "%s", resp)
	case vegeta.BinaryFormatString:
		c.Header("Content-Type", "application/octet-stream")
		c.Data(http.StatusOK, "application/octet-stream", resp)
	case vegeta.HistogramFormatString:
		c.Header("Content-Type", "text/plain")
		c.String(http.StatusOK, "%s", resp)
	}
}
