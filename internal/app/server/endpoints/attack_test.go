package endpoints

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
	"vegeta-server/internal/app/attacker"
	"vegeta-server/internal/app/server/models"
)

func setupAttackResponses(n int) map[string]*models.AttackResponse {
	r := make(map[string]*models.AttackResponse)
	for i := 0; i < n; i++ {
		id := strconv.Itoa(i)
		r[id] = &models.AttackResponse{
			ID:     id,
			Status: models.AttackResponseStatusCompleted,
		}
	}
	return r
}

func setupAttackReqBody() (io.Reader, error) {
	body := map[string]interface{}{
		"rate":     5,
		"duration": "20s",
		"target": map[string]interface{}{
			"method": "GET",
			"URL":    "http://localhost:8000/api/v1/attack",
			"scheme": "http",
		},
	}

	req, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

	return bytes.NewReader(req), nil
}

func setupAttackCancelReqBody(id string) (io.Reader, error) {
	body := map[string]interface{}{
		"cancel": true,
	}

	req, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

	return bytes.NewReader(req), nil
}

func performRequest(r http.Handler, method, path string, body io.Reader) *httptest.ResponseRecorder {
	req, _ := http.NewRequest(method, path, body)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}

func TestEndpoints_PostAttackEndpoint(t *testing.T) {
	router := SetupRouter(
		&attacker.MockScheduler{
			Responses: setupAttackResponses(1),
		},
	)

	expectedResponse := models.AttackResponse{
		Status: models.AttackResponseStatusScheduled,
	}

	body, err := setupAttackReqBody()
	if err != nil {
		t.Fatal(err)
	}

	w := performRequest(
		router,
		"POST",
		"/api/v1/attack",
		body,
	)

	assert.Equal(t, http.StatusOK, w.Code)

	var response models.AttackResponse
	err = json.Unmarshal([]byte(w.Body.String()), &response)
	if err != nil {
		t.Fatal(err)
	}

	assert.Nil(t, err)
	assert.Equal(t, expectedResponse.Status, response.Status)
}

func TestEndpoints_GetAttackByIDEndpoint(t *testing.T) {
	expectedID := "0"

	router := SetupRouter(
		&attacker.MockScheduler{
			Responses: setupAttackResponses(1),
		},
	)

	w := performRequest(
		router,
		"GET",
		fmt.Sprintf("/api/v1/attack/%s", expectedID),
		nil,
	)

	assert.Equal(t, http.StatusOK, w.Code)

	var getResponse models.AttackResponse
	err := json.Unmarshal([]byte(w.Body.String()), &getResponse)
	if err != nil {
		t.Fatal(err)
	}

	assert.Nil(t, err)
	assert.Equal(t, expectedID, getResponse.ID)
}

func TestEndpoints_GetAttackByIDEndpoint_NotFound(t *testing.T) {
	expectedID := "100"

	router := SetupRouter(
		&attacker.MockScheduler{
			Responses: setupAttackResponses(1),
		},
	)

	w := performRequest(
		router,
		"GET",
		fmt.Sprintf("/api/v1/attack/%s", expectedID),
		nil,
	)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestEndpoints_GetAttackEndpoint(t *testing.T) {
	num := 5

	router := SetupRouter(
		&attacker.MockScheduler{
			Responses: setupAttackResponses(num),
		},
	)

	w := performRequest(
		router,
		"GET",
		fmt.Sprintf("/api/v1/attack"),
		nil,
	)

	assert.Equal(t, http.StatusOK, w.Code)

	var getResponse []*models.AttackResponse
	err := json.Unmarshal([]byte(w.Body.String()), &getResponse)
	if err != nil {
		t.Fatal(err)
	}

	assert.Nil(t, err)
	// FIXME : Compare only the values
	//assert.Equal(t, expectedResponse, getResponse)
}

func TestEndpoints_PostAttackByIDCancelEndpoint(t *testing.T) {
	id := "12345"
	responses := make(map[string]*models.AttackResponse)
	responses[id] = &models.AttackResponse{
		ID:     id,
		Status: models.AttackResponseStatusRunning,
	}

	router := SetupRouter(
		&attacker.MockScheduler{
			Responses: responses,
		},
	)

	body, err := setupAttackCancelReqBody(id)
	if err != nil {
		t.Fatal(err)
	}

	w := performRequest(
		router,
		"POST",
		fmt.Sprintf(fmt.Sprintf("/api/v1/attack/%s/cancel", id)),
		body,
	)

	assert.Equal(t, http.StatusOK, w.Code)
	var response models.AttackResponse
	err = json.Unmarshal([]byte(w.Body.String()), &response)
	if err != nil {
		t.Fatal(err)
	}

	assert.Nil(t, err)
	assert.Equal(t, models.AttackResponse{
		ID:     id,
		Status: models.AttackResponseStatusCanceled,
	}, response)
}
