package endpoints

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/tsenart/vegeta/lib"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"vegeta-server/internal/app/attacker"
	"vegeta-server/internal/app/server/models"
)

func setupDispatcher() attacker.IDispatcher{
	return attacker.NewDispatcher(func(*attacker.AttackOpts) <-chan *vegeta.Result {
		fmt.Printf("ATTACKING\n")
		return nil
	})
}

func setupScheduler() attacker.IScheduler {
	return attacker.NewScheduler(
		setupDispatcher(),
		make(chan struct{}),
	)
}

func setupEndpoints() *Endpoints {
	return &Endpoints{}
}

func setupAttackReqBody() (io.Reader, error) {
	body := map[string]interface{}{
		"rate": 5,
		"duration": "20s",
		"target": map[string]interface{}{
			"method": "GET",
			"URL": "http://localhost:8000/api/v1/attack",
			"scheme": "http",
		},
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
		attacker.NewAttacker(
			setupScheduler(),
		),
	)

	expectedResponse := models.AttackResponse {
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
