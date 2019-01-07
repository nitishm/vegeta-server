// This file is safe to edit. Once it exists it will not be overwritten

package restapi

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"os"

	rtime "runtime"

	errors "github.com/go-openapi/errors"
	runtime "github.com/go-openapi/runtime"
	middleware "github.com/go-openapi/runtime/middleware"
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
		swag.CommandLineOptionsGroup{
			ShortDescription: "Version",
			LongDescription:  "",
			Options:          &vFlags,
		},
	}
}

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

	// configure the api here
	api.ServeError = errors.ServeError

	api.Logger = log.Infof

	api.JSONConsumer = runtime.JSONConsumer()

	api.JSONProducer = runtime.JSONProducer()

	api.AttackGetAttackByIDHandler = attack.GetAttackByIDHandlerFunc(func(params attack.GetAttackByIDParams) middleware.Responder {
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
		return middleware.NotImplemented("operation attack.GetAttacks has not yet been implemented")
	})
	api.AttackPostAttackHandler = attack.PostAttackHandlerFunc(func(params attack.PostAttackParams) middleware.Responder {
		attackID := at.Schedule(params.Body)
		return attack.NewPostAttackOK().WithPayload(&models.AttackResponse{
			ID:     attackID,
			Status: models.AttackResponseStatusScheduled,
		})
	})
	api.AttackPutAttackByIDCancelHandler = attack.PutAttackByIDCancelHandlerFunc(func(params attack.PutAttackByIDCancelParams) middleware.Responder {
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

// The middleware configuration happens before anything, this middleware also applies to serving the swagger.json document.
// So this is a good place to plug in a panic handling middleware, logging and metrics
func setupGlobalMiddleware(handler http.Handler) http.Handler {
	return handler
}
