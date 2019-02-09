package endpoints

import (
	"github.com/gin-gonic/gin"
	"vegeta-server/internal/dispatcher"
)

type Endpoints struct {
	dispatcher dispatcher.IDispatcher
}

func NewEndpoints(d dispatcher.IDispatcher) *Endpoints {
	return &Endpoints{
		d,
	}
}

func SetupRouter(d dispatcher.IDispatcher) *gin.Engine {
	router := gin.Default()

	e := NewEndpoints(d)

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
