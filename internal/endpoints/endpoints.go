package endpoints

import (
	"github.com/gin-gonic/gin"
	"vegeta-server/internal/dispatcher"
	"vegeta-server/internal/reporter"
)

type Endpoints struct {
	dispatcher dispatcher.IDispatcher
	reporter   reporter.IReporter
}

func NewEndpoints(d dispatcher.IDispatcher, r reporter.IReporter) *Endpoints {
	return &Endpoints{
		d,
		r,
	}
}

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
