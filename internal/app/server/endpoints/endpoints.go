package endpoints

import (
	"github.com/gin-gonic/gin"
	"vegeta-server/internal/app/attacker"
)

type Endpoints struct {
	attacker attacker.IAttacker
}

func NewEndpoints(a attacker.IAttacker) *Endpoints {
	return &Endpoints{
		a,
	}
}

func SetupRouter(a attacker.IAttacker) *gin.Engine {
	router := gin.Default()

	e := NewEndpoints(a)

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
