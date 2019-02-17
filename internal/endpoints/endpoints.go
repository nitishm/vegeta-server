package endpoints

import (
	"net/http"
	"vegeta-server/internal/dispatcher"
	"vegeta-server/internal/reporter"

	"github.com/gin-gonic/gin"
)

var (
	ginErrNotFound = func(c *gin.Context, err error) {
		c.JSON(
			http.StatusNotFound,
			gin.H{
				"message": "Not found",
				"code":    http.StatusNotFound,
				"error":   err.Error(),
			},
		)
	}

	ginErrBadRequest = func(c *gin.Context, err error) {
		c.JSON(
			http.StatusBadRequest,
			gin.H{
				"message": "Bad request params",
				"code":    http.StatusBadRequest,
				"error":   err.Error(),
			},
		)
	}

	ginErrInternalServerError = func(c *gin.Context, err error) {
		c.JSON(
			http.StatusInternalServerError,
			gin.H{
				"message": "Internal server error",
				"code":    http.StatusInternalServerError,
				"error":   err.Error(),
			},
		)
	}
)

// Endpoints provides an encapsulation for all dependencies required by the
// API handlers.
type Endpoints struct {
	dispatcher dispatcher.IDispatcher
	reporter   reporter.IReporter
}

// NewEndpoints returns an instance of the Endpoints object
func NewEndpoints(d dispatcher.IDispatcher, r reporter.IReporter) *Endpoints {
	return &Endpoints{
		d,
		r,
	}
}

// SetupRouter registers the endpoint handlers and returns a pointer to the
// server instance.
func SetupRouter(d dispatcher.IDispatcher, r reporter.IReporter) *gin.Engine {
	router := gin.Default()

	e := NewEndpoints(d, r)

	// api/v1 router group
	v1 := router.Group("/api/v1")
	{
		// Attack endpoints
		v1.POST("/attack", e.PostAttackEndpoint)
		v1.GET("/attack", e.GetAttackEndpoint)
		v1.GET("/attack/:attackID", e.GetAttackByIDEndpoint)
		v1.POST("/attack/:attackID/cancel", e.PostAttackByIDCancelEndpoint)

		// Report endpoints
		v1.GET("/report", e.GetReportEndpoint)
		v1.GET("/report/:attackID", e.GetReportByIDEndpoint)
	}

	return router
}
