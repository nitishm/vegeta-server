package endpoints

import (
	"net/http"
	"vegeta-server/models"

	"github.com/gin-gonic/gin"
)

func (e *Endpoints) PostAttackEndpoint(c *gin.Context) {
	var attackParams models.AttackParams
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
	resp := e.dispatcher.Dispatch(attackParams)

	c.JSON(http.StatusOK, resp)
}

func (e *Endpoints) GetAttackByIDEndpoint(c *gin.Context) {
	id := c.Param("attackID")
	resp, err := e.dispatcher.Get(id)
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

	c.JSON(http.StatusOK, resp)
}

func (e *Endpoints) GetAttackEndpoint(c *gin.Context) {
	resp := e.dispatcher.List()

	c.JSON(http.StatusOK, resp)
}

func (e *Endpoints) PostAttackByIDCancelEndpoint(c *gin.Context) {
	id := c.Param("attackID")
	var attackCancelParams models.AttackCancel
	if err := c.ShouldBindJSON(&attackCancelParams); err != nil {
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

	_, err := e.dispatcher.Get(id)
	if err != nil {
		c.JSON(
			http.StatusNotFound,
			gin.H{
				"message": "Not Found",
				"code":    http.StatusNotFound,
				"error":   err.Error(),
			},
		)
		return
	}

	resp, err := e.dispatcher.Cancel(id, attackCancelParams.Cancel)
	if err != nil {
		c.JSON(
			http.StatusInternalServerError,
			gin.H{
				"message": "Internal server error",
				"code":    http.StatusInternalServerError,
				"error":   err.Error(),
			},
		)
		return
	}

	c.JSON(http.StatusOK, resp)
}
