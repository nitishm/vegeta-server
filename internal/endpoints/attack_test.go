package endpoints

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"vegeta-server/internal/dispatcher"
	dmocks "vegeta-server/internal/dispatcher/mocks"
	"vegeta-server/models"

	"github.com/stretchr/testify/mock"

	assert "gopkg.in/go-playground/assert.v1"
)

func setupTestRouter(d dispatcher.IDispatcher, req *http.Request) *httptest.ResponseRecorder {
	router := SetupRouter(d, nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)
	return w
}

func TestEndpoints_PostAttackEndpoint(t *testing.T) {

}

func TestEndpoints_GetAttackByIDEndpoint_NotFound(t *testing.T) {
	d := &dmocks.IDispatcher{}

	// Prepare mock
	wantErr := fmt.Errorf("not found")
	d.
		On("Get", mock.AnythingOfType("string")).
		Return(nil, wantErr)

	// Setup router
	req, _ := http.NewRequest("GET", "/api/v1/attack/123", nil)
	w := setupTestRouter(d, req)

	// Assert results
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestEndpoints_GetAttackByIDEndpoint(t *testing.T) {
	d := &dmocks.IDispatcher{}

	// Prepare mock
	wantBody := &models.AttackResponse{
		ID:     "123",
		Status: models.AttackResponseStatusScheduled,
	}

	d.
		On("Get", "123").
		Return(wantBody, nil)

	// Setup router
	req, _ := http.NewRequest("GET", "/api/v1/attack/123", nil)
	w := setupTestRouter(d, req)

	// Assert results
	assert.Equal(t, http.StatusOK, w.Code)
	gotBody := &models.AttackResponse{}
	_ = json.Unmarshal(w.Body.Bytes(), gotBody)
	assert.Equal(t, *wantBody, *gotBody)
}

func TestEndpoints_GetAttackEndpoint(t *testing.T) {

}

func TestEndpoints_PostAttackByIDCancelEndpoint(t *testing.T) {

}
