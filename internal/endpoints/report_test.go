package endpoints

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"vegeta-server/internal/reporter"
	"vegeta-server/models"
	"vegeta-server/pkg/vegeta"

	"github.com/gin-gonic/gin/json"

	rmock "vegeta-server/internal/reporter/mocks"

	assert "gopkg.in/go-playground/assert.v1"
)

func setupTestReporterRouter(r reporter.IReporter, req *http.Request) *httptest.ResponseRecorder {
	router := SetupRouter(nil, r)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)
	return w
}

type setupReporterFunc func() (reporter.IReporter, *http.Request)

func TestEndpoints_GetReportEndpoint(t *testing.T) {
	type params struct {
		setup    setupReporterFunc
		wantCode int
	}
	tests := []struct {
		name   string
		params params
	}{
		{
			name: "OK",
			params: params{
				setup: func() (reporter.IReporter, *http.Request) {
					reports := []models.JSONReportResponse{
						{
							ID: "1",
						},
					}
					bReports := make([][]byte, 0)
					for _, report := range reports {
						res, _ := json.Marshal(&report)
						bReports = append(bReports, res)
					}

					r := &rmock.IReporter{}

					r.
						On("GetAll").
						Return(bReports)

					// Setup router
					req, _ := http.NewRequest("GET", "/api/v1/report", nil)

					return r, req
				},
				wantCode: http.StatusOK,
			},
		},
		{
			name: "Internal Server Error",
			params: params{
				setup: func() (reporter.IReporter, *http.Request) {
					reports := [][]byte{[]byte("1234")}

					r := &rmock.IReporter{}

					r.
						On("GetAll").
						Return(reports)

					// Setup router
					req, _ := http.NewRequest("GET", "/api/v1/report", nil)

					return r, req
				},
				wantCode: http.StatusInternalServerError,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := setupTestReporterRouter(tt.params.setup())
			gotCode := w.Code
			assert.Equal(t, tt.params.wantCode, gotCode)
		})
	}
}

func TestEndpoints_GetReportByIDEndpoint(t *testing.T) {
	type params struct {
		setup    setupReporterFunc
		wantCode int
	}
	tests := []struct {
		name   string
		params params
	}{
		{
			name: "Not Found",
			params: params{
				setup: func() (reporter.IReporter, *http.Request) {
					r := &rmock.IReporter{}
					r.
						On("GetInFormat", "123", vegeta.JSONFormat).
						Return([]byte{}, fmt.Errorf("not found"))

					// Setup router
					req, _ := http.NewRequest("GET", "/api/v1/report/123?format=json", nil)

					return r, req
				},
				wantCode: http.StatusNotFound,
			},
		},
		{
			name: "OK - json (implicit)",
			params: params{
				setup: func() (reporter.IReporter, *http.Request) {
					resp := models.JSONReportResponse{}
					bResp, _ := json.Marshal(resp)
					r := &rmock.IReporter{}
					r.
						On("GetInFormat", "123", vegeta.JSONFormat).
						Return(bResp, nil)

					// Setup router
					req, _ := http.NewRequest("GET", "/api/v1/report/123", nil)

					return r, req
				},
				wantCode: http.StatusOK,
			},
		},
		{
			name: "OK - json (explicit)",
			params: params{
				setup: func() (reporter.IReporter, *http.Request) {
					resp := models.JSONReportResponse{}
					bResp, _ := json.Marshal(resp)
					r := &rmock.IReporter{}
					r.
						On("GetInFormat", "123", vegeta.JSONFormat).
						Return(bResp, nil)

					// Setup router
					req, _ := http.NewRequest("GET", "/api/v1/report/123?format=json", nil)

					return r, req
				},
				wantCode: http.StatusOK,
			},
		},
		{
			name: "Internal Server Error",
			params: params{
				setup: func() (reporter.IReporter, *http.Request) {
					r := &rmock.IReporter{}
					r.
						On("GetInFormat", "123", vegeta.JSONFormat).
						Return(nil, nil)

					// Setup router
					req, _ := http.NewRequest("GET", "/api/v1/report/123?format=json", nil)

					return r, req
				},
				wantCode: http.StatusInternalServerError,
			},
		},
		{
			name: "OK - text",
			params: params{
				setup: func() (reporter.IReporter, *http.Request) {
					r := &rmock.IReporter{}
					r.
						On("GetInFormat", "123", vegeta.TextFormat).
						Return([]byte{}, nil)

					// Setup router
					req, _ := http.NewRequest("GET", "/api/v1/report/123?format=text", nil)

					return r, req
				},
				wantCode: http.StatusOK,
			},
		},
		{
			name: "OK - binary",
			params: params{
				setup: func() (reporter.IReporter, *http.Request) {
					r := &rmock.IReporter{}
					r.
						On("GetInFormat", "123", vegeta.BinaryFormat).
						Return([]byte{}, nil)

					// Setup router
					req, _ := http.NewRequest("GET", "/api/v1/report/123?format=binary", nil)

					return r, req
				},
				wantCode: http.StatusOK,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := setupTestReporterRouter(tt.params.setup())
			gotCode := w.Code
			assert.Equal(t, tt.params.wantCode, gotCode)
		})
	}
}
