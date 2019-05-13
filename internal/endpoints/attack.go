package endpoints

import (
	"net/http"
	"net/url"
	"vegeta-server/models"

	"github.com/gin-gonic/gin"
)

// PostAttackEndpoint implements a handler for the POST /api/v1/attack endpoint
func (e *Endpoints) PostAttackEndpoint(c *gin.Context) {
	var attackParams models.AttackParams
	if err := c.ShouldBindJSON(&attackParams); err != nil {
		ginErrBadRequest(c, err)
		return
	}

	// Submit the attack
	resp, err := e.dispatcher.Dispatch(attackParams)
	if err != nil {
		ginErrInternalServerError(c, err)
		return
	}

	c.JSON(http.StatusOK, resp)
}

// GetAttackByIDEndpoint implements a handler for the GET /api/v1/attack/<attackID> endpoint
func (e *Endpoints) GetAttackByIDEndpoint(c *gin.Context) {
	id := c.Param("attackID")
	resp, err := e.dispatcher.Get(id)
	if err != nil {
		ginErrNotFound(c, err)
		return
	}

	c.JSON(http.StatusOK, resp)
}

// GetAttackEndpoint implements a handler for the GET /api/v1/attack endpoint
func (e *Endpoints) GetAttackEndpoint(c *gin.Context) {
	var err error
	filterMap := make(models.FilterParams)
	filterMap["status"] = c.DefaultQuery("status", "")
	b := c.DefaultQuery("created_before", "")
	filterMap["created_before"], err = url.QueryUnescape(b)
	if err != nil {
		ginErrBadRequest(c, err)
		return
	}

	filterMap["created_after"] = c.DefaultQuery("created_after", "")
	resp := e.dispatcher.List(
		//models.StatusFilter(status),
		filterMap,
	)

	c.JSON(http.StatusOK, resp)
}

// PostAttackByIDCancelEndpoint implements a handler for the POST /api/v1/attack/<attackID>/cancel endpoint
func (e *Endpoints) PostAttackByIDCancelEndpoint(c *gin.Context) {
	id := c.Param("attackID")
	var attackCancelParams models.AttackCancel
	if err := c.ShouldBindJSON(&attackCancelParams); err != nil {
		ginErrBadRequest(c, err)
		return
	}

	_, err := e.dispatcher.Get(id)
	if err != nil {
		ginErrNotFound(c, err)
		return
	}

	err = e.dispatcher.Cancel(id, attackCancelParams.Cancel)
	if err != nil {
		ginErrInternalServerError(c, err)
		return
	}

	c.Status(http.StatusOK)
}
