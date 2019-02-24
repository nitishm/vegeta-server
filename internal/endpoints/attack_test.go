package endpoints

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"vegeta-server/internal/dispatcher"
	dmocks "vegeta-server/internal/dispatcher/mocks"
	"vegeta-server/models"

	"github.com/stretchr/testify/mock"

	assert "gopkg.in/go-playground/assert.v1"
)

func setupTestDispatcherRouter(d dispatcher.IDispatcher, req *http.Request) *httptest.ResponseRecorder {
	router := SetupRouter(d, nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)
	return w
}

type setupDispatcherFunc func() (dispatcher.IDispatcher, *http.Request)

func TestEndpoints_PostAttackEndpoint(t *testing.T) {
	type params struct {
		setup    setupDispatcherFunc
		wantCode int
	}
	tests := []struct {
		name   string
		params params
	}{
		{
			name: "Bad Request - Nil body",
			params: params{
				func() (dispatcher.IDispatcher, *http.Request) {
					req, _ := http.NewRequest("POST", "/api/v1/attack", strings.NewReader(""))
					return new(dmocks.IDispatcher), req
				},
				http.StatusBadRequest,
			},
		},
		{
			name: "Bad Request - Missing Duration",
			params: params{
				func() (dispatcher.IDispatcher, *http.Request) {
					attackParams := models.AttackParams{
						Rate: 1,
						Target: models.Target{
							Method: "GET",
							URL:    "localhost:80/api/v1/",
							Scheme: "http",
						},
					}
					d := new(dmocks.IDispatcher)

					d.
						On("Dispatch", attackParams).
						Return(nil, nil)
					bAttackParamsBody, _ := json.Marshal(attackParams)
					attackParamsBody := string(bAttackParamsBody)

					req, _ := http.NewRequest("POST", "/api/v1/attack", strings.NewReader(attackParamsBody))

					return d, req
				},
				http.StatusBadRequest,
			},
		},
		{
			name: "OK",
			params: params{
				func() (dispatcher.IDispatcher, *http.Request) {
					attackParams := models.AttackParams{
						Rate: 1,
						Target: models.Target{
							Method: "GET",
							URL:    "localhost:80/api/v1/",
							Scheme: "http",
						},
						Duration: "1s",
					}
					d := new(dmocks.IDispatcher)

					d.
						On("Dispatch", attackParams).
						Return(nil, nil)
					bAttackParamsBody, _ := json.Marshal(attackParams)
					attackParamsBody := string(bAttackParamsBody)

					req, _ := http.NewRequest("POST", "/api/v1/attack", strings.NewReader(attackParamsBody))

					return d, req
				},
				http.StatusOK,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := setupTestDispatcherRouter(tt.params.setup())
			gotCode := w.Code
			assert.Equal(t, tt.params.wantCode, gotCode)
		})
	}
}

func TestEndpoints_GetAttackByIDEndpoint(t *testing.T) {
	type params struct {
		setup    setupDispatcherFunc
		wantCode int
	}
	tests := []struct {
		name   string
		params params
	}{
		{
			name: "Not Found",
			params: params{
				func() (dispatcher.IDispatcher, *http.Request) {
					d := &dmocks.IDispatcher{}

					// Prepare mock
					wantErr := fmt.Errorf("not found")
					d.
						On("Get", mock.AnythingOfType("string")).
						Return(nil, wantErr)

					// Setup router
					req, _ := http.NewRequest("GET", "/api/v1/attack/123", nil)
					return d, req
				},
				http.StatusNotFound,
			},
		},
		{
			name: "OK",
			params: params{
				func() (dispatcher.IDispatcher, *http.Request) {
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

					return d, req
				},
				http.StatusOK,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := setupTestDispatcherRouter(tt.params.setup())
			gotCode := w.Code
			assert.Equal(t, tt.params.wantCode, gotCode)
		})
	}
}

func TestEndpoints_GetAttackEndpoint(t *testing.T) {
	type params struct {
		setup    setupDispatcherFunc
		wantCode int
	}
	tests := []struct {
		name   string
		params params
	}{
		{
			name: "OK",
			params: params{
				setup: func() (iDispatcher dispatcher.IDispatcher, request *http.Request) {
					d := &dmocks.IDispatcher{}
					d.
						On("List").
						Return([]*models.AttackResponse{})

					// Setup router
					req, _ := http.NewRequest("GET", "/api/v1/attack", nil)
					return d, req
				},
				wantCode: http.StatusOK,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := setupTestDispatcherRouter(tt.params.setup())
			gotCode := w.Code
			assert.Equal(t, tt.params.wantCode, gotCode)
		})
	}
}

func TestEndpoints_PostAttackByIDCancelEndpoint(t *testing.T) {
	type params struct {
		setup    setupDispatcherFunc
		wantCode int
	}
	tests := []struct {
		name   string
		params params
	}{
		{
			name: "Bad Request - nil body",
			params: params{
				setup: func() (iDispatcher dispatcher.IDispatcher, request *http.Request) {
					d := &dmocks.IDispatcher{}
					// Setup router
					req, _ := http.NewRequest("POST", "/api/v1/attack/123/cancel", strings.NewReader(""))
					return d, req
				},
				wantCode: http.StatusBadRequest,
			},
		},
		{
			name: "Not Found",
			params: params{
				setup: func() (iDispatcher dispatcher.IDispatcher, request *http.Request) {
					d := &dmocks.IDispatcher{}
					// Return valid response
					d.
						On("Get", "123").
						Return(nil, fmt.Errorf("not found"))

					bAttackCancelBody, _ := json.Marshal(&models.AttackCancel{
						Cancel: true,
					})
					attackCancelBody := string(bAttackCancelBody)
					// Setup router
					req, _ := http.NewRequest("POST", "/api/v1/attack/123/cancel", strings.NewReader(attackCancelBody))
					return d, req
				},
				wantCode: http.StatusNotFound,
			},
		},
		{
			name: "Internal Server Error",
			params: params{
				setup: func() (iDispatcher dispatcher.IDispatcher, request *http.Request) {
					d := &dmocks.IDispatcher{}
					// Return valid response
					d.
						On("Get", "123").
						Return(nil, nil)

					// Return error on Cancel
					d.
						On("Cancel", "123", true).
						Return(fmt.Errorf("internal server error"))

					bAttackCancelBody, _ := json.Marshal(&models.AttackCancel{
						Cancel: true,
					})
					attackCancelBody := string(bAttackCancelBody)
					// Setup router
					req, _ := http.NewRequest("POST", "/api/v1/attack/123/cancel", strings.NewReader(attackCancelBody))
					return d, req
				},
				wantCode: http.StatusInternalServerError,
			},
		},
		{
			name: "OK",
			params: params{
				setup: func() (iDispatcher dispatcher.IDispatcher, request *http.Request) {
					d := &dmocks.IDispatcher{}
					// Return valid response
					d.
						On("Get", "123").
						Return(nil, nil)

					// Return error on Cancel
					d.
						On("Cancel", "123", true).
						Return(nil)

					bAttackCancelBody, _ := json.Marshal(&models.AttackCancel{
						Cancel: true,
					})
					attackCancelBody := string(bAttackCancelBody)
					// Setup router
					req, _ := http.NewRequest("POST", "/api/v1/attack/123/cancel", strings.NewReader(attackCancelBody))
					return d, req
				},
				wantCode: http.StatusOK,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := setupTestDispatcherRouter(tt.params.setup())
			gotCode := w.Code
			assert.Equal(t, tt.params.wantCode, gotCode)
		})
	}
}
