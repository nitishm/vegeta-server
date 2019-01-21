// This file is safe to edit. Once it exists it will not be overwritten

package restapi

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	rtime "runtime"
	"vegeta-server/restapi/operations/report"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime"
	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/swag"

	"vegeta-server/internal/vegeta"
	"vegeta-server/models"
	"vegeta-server/restapi/operations"
	"vegeta-server/restapi/operations/attack"

	log "github.com/sirupsen/logrus"
)

var (
	commit  = "N/A"
	date    = "N/A"
	version = "N/A"
)

//go:generate swagger generate server --target ../../vegeta-server --name Vegeta --spec ../spec/swagger.yaml
type VersionFlag struct {
	Version bool `long:"version" description:"Show vegeta-server version details"`
}

func configureFlags(api *operations.VegetaAPI) {
	var vFlags = VersionFlag{}

	api.CommandLineOptionsGroups = []swag.CommandLineOptionsGroup{
		{
			ShortDescription: "Version",
			LongDescription:  "",
			Options:          &vFlags,
		},
	}
}

//nolint:lll
func configureAPI(api *operations.VegetaAPI) http.Handler {
	// FIXME: Find a better way to do this ?
	// Print version info if --version flag is set
	if vf, ok := api.CommandLineOptionsGroups[0].Options.(*VersionFlag); ok {
		if vf.Version {
			// Set at linking time
			fmt.Println("=======")
			fmt.Println("VERSION")
			fmt.Println("=======")
			fmt.Println("Version\t", version)
			fmt.Println("Commit \t", commit)
			fmt.Println("Runtime\t", fmt.Sprintf("%s %s/%s", rtime.Version(), rtime.GOOS, rtime.GOARCH))
			fmt.Println("Date   \t", date)

			os.Exit(0)
		}
	}

	// Initialize the attacker
	at := vegeta.NewAttacker()

	// Initialize the reporter with the attack results channel
	rp := vegeta.NewReporter(at.Results)

	// configure the api here
	api.ServeError = errors.ServeError

	api.Logger = log.Infof

	api.JSONConsumer = runtime.JSONConsumer()

	api.JSONProducer = runtime.JSONProducer()

	api.AttackGetAttackByIDHandler = attack.GetAttackByIDHandlerFunc(func(params attack.GetAttackByIDParams) middleware.Responder {
		if !at.Exists(params.AttackID) {
			code := http.StatusText(http.StatusNotFound)
			message := fmt.Sprintf("Attack by ID %v Not Found", params.AttackID)
			e := &models.Error{Code: &code, Message: &message}
			return attack.NewGetAttackByIDNotFound().WithPayload(e)
		}

		status, err := at.Status(params.AttackID)
		if err != nil {
			code := http.StatusText(http.StatusInternalServerError)
			message := err.Error()
			e := &models.Error{Code: &code, Message: &message}
			return attack.NewGetAttackByIDInternalServerError().WithPayload(e)
		}

		return attack.NewGetAttackByIDOK().WithPayload(&models.AttackResponse{
			ID:     params.AttackID,
			Status: status,
		})
	})

	api.AttackGetAttacksHandler = attack.GetAttacksHandlerFunc(func(params attack.GetAttacksParams) middleware.Responder {
		var attackResponseList models.AttackResponseList
		attacks := at.List()
		for uuid, status := range attacks {
			attackResponse := models.AttackResponse{
				ID:     uuid,
				Status: status,
			}
			attackResponseList = append(attackResponseList, &attackResponse)
		}
		return attack.NewGetAttacksOK().WithPayload(attackResponseList)
	})

	api.AttackPostAttackHandler = attack.PostAttackHandlerFunc(func(params attack.PostAttackParams) middleware.Responder {
		attackID := at.Schedule(params.Body)
		return attack.NewPostAttackOK().WithPayload(&models.AttackResponse{
			ID:     attackID,
			Status: models.AttackResponseStatusScheduled,
		})
	})

	api.AttackPutAttackByIDCancelHandler = attack.PutAttackByIDCancelHandlerFunc(func(params attack.PutAttackByIDCancelParams) middleware.Responder {
		if !at.Exists(params.AttackID) {
			code := http.StatusText(http.StatusNotFound)
			message := fmt.Sprintf("Attack by ID %v Not Found", params.AttackID)
			e := &models.Error{Code: &code, Message: &message}
			return attack.NewPutAttackByIDCancelNotFound().WithPayload(e)
		}

		status, err := at.Cancel(params.AttackID, *params.Body.IsCanceled)
		if err != nil {
			code := http.StatusText(http.StatusInternalServerError)
			message := err.Error()
			e := &models.Error{Code: &code, Message: &message}
			return attack.NewPutAttackByIDCancelInternalServerError().WithPayload(e)
		}

		return attack.NewPutAttackByIDCancelOK().WithPayload(&models.AttackResponse{
			ID:     params.AttackID,
			Status: status,
		})
	})

	api.ReportGetReportByIDHandler = report.GetReportByIDHandlerFunc(func(params report.GetReportByIDParams) middleware.Responder {
		if !at.Exists(params.AttackID) {
			code := http.StatusText(http.StatusNotFound)
			message := fmt.Sprintf("Attack by ID %v Not Found", params.AttackID)
			e := &models.Error{Code: &code, Message: &message}
			return report.NewGetReportByIDNotFound().WithPayload(e)
		}

		r, err := rp.Get(params.AttackID)
		if err != nil {
			code := http.StatusText(http.StatusInternalServerError)
			message := err.Error()
			e := &models.Error{Code: &code, Message: &message}
			return report.NewGetReportByIDInternalServerError().WithPayload(e)
		}

		rep, err := adaptReportToSpec(r)
		if err != nil {
			code := http.StatusText(http.StatusInternalServerError)
			message := err.Error()
			e := &models.Error{Code: &code, Message: &message}
			return report.NewGetReportByIDInternalServerError().WithPayload(e)
		}

		return report.NewGetReportByIDOK().WithPayload(&models.ReportResponse{
			ID:     params.AttackID,
			Report: rep,
		})
	})

	api.ReportGetReportsHandler = report.GetReportsHandlerFunc(func(params report.GetReportsParams) middleware.Responder {
		var reportResponseList models.ReportResponseList
		reports := rp.List()
		for uuid, rep := range reports {
			report, err := adaptReportToSpec(rep)
			if err != nil {
				// FIXME: Should we just skip or panic if this happens ?
				continue
			}

			reportResponse := models.ReportResponse{
				ID:     uuid,
				Report: report,
			}
			reportResponseList = append(reportResponseList, &reportResponse)
		}
		return report.NewGetReportsOK().WithPayload(reportResponseList)
	})

	api.ServerShutdown = func() {}

	return setupGlobalMiddleware(api.Serve(setupMiddlewares))
}

// The TLS configuration before HTTPS server starts.
func configureTLS(tlsConfig *tls.Config) {
	// Make all necessary changes to the TLS configuration here.
}

// As soon as server is initialized but not run yet, this function will be called.
// If you need to modify a config, store server instance to stop it individually later, this is the place.
// This function can be called multiple times, depending on the number of serving schemes.
// scheme value will be set accordingly: "http", "https" or "unix"
func configureServer(s *http.Server, scheme, addr string) {
}

// The middleware configuration is for the handler executors. These do not apply to the swagger.json document.
// The middleware executes after routing but before authentication, binding and validation
func setupMiddlewares(handler http.Handler) http.Handler {
	return handler
}

// The middleware configuration happens before anything, this middleware also
// applies to serving the swagger.json document.
// So this is a good place to plug in a panic handling middleware, logging and metrics
func setupGlobalMiddleware(handler http.Handler) http.Handler {
	return handler
}

func adaptReportToSpec(r string) (*models.Report, error) {
	report := &models.Report{}
	err := json.Unmarshal([]byte(r), report)
	if err != nil {
		return nil, err
	}
	return report, nil
}
