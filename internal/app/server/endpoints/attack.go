package endpoints

import (
	"net/http"
	"vegeta-server/internal/app/server/models"

	"github.com/gin-gonic/gin"
)

func (e *Endpoints) PostAttackEndpoint(c *gin.Context) {
	var attackParams models.Attack
	if err := c.ShouldBindJSON(&attackParams); err != nil {
		c.JSON(
			http.StatusBadRequest,
			gin.H{
				"message": "Bad request params",
				"code":    http.StatusBadRequest,
				"error":   err.Error(),
			},
		)
		return
	}

	// Submit the attack
	resp := e.attacker.Submit(attackParams)

	c.JSON(http.StatusOK, resp)
}

func (e *Endpoints) GetAttackEndpoint(c *gin.Context) {
	c.JSON(
		http.StatusOK,
		gin.H{
			"code":     http.StatusOK,
			"endpoint": "GET /attack",
		},
	)
}

func (e *Endpoints) GetAttackByIDEndpoint(c *gin.Context) {
	c.JSON(
		http.StatusOK,
		gin.H{
			"code":     http.StatusOK,
			"endpoint": "GET /attack/" + c.Param("attackID"),
		},
	)
}

func (e *Endpoints) PostAttackByIDCancelEndpoint(c *gin.Context) {
	c.JSON(
		http.StatusOK,
		gin.H{
			"code":     http.StatusOK,
			"endpoint": "GET /attack/" + c.Param("attackID") + "/cancel",
		},
	)
}
