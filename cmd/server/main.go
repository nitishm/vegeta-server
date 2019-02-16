package main

import (
	"fmt"
	"os"
	"os/signal"
	"runtime"
	"vegeta-server/internal/dispatcher"
	"vegeta-server/internal/endpoints"
	"vegeta-server/internal/reporter"
	"vegeta-server/models"
	"vegeta-server/pkg/vegeta"

	gin "github.com/gin-gonic/gin"

	log "github.com/sirupsen/logrus"

	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

var (
	commit  = "N/A"
	date    = "N/A"
	version = "N/A"
)

var (
	ip    = kingpin.Flag("ip", "Server IP Address.").Default("0.0.0.0").String()
	port  = kingpin.Flag("port", "Server Port.").Default("80").String()
	v     = kingpin.Flag("version", "Version Info").Short('v').Bool()
	debug = kingpin.Flag("debug", "Enabled Debug").Bool()
)

func main() {
	kingpin.Parse()

	if *v {
		// Set at linking time
		fmt.Println("Version\t", version)
		fmt.Println("Commit \t", commit)
		fmt.Println("Runtime\t", fmt.Sprintf("%s %s/%s", runtime.Version(), runtime.GOOS, runtime.GOARCH))
		fmt.Println("Date   \t", date)

		os.Exit(0)
		return
	}

	if *debug {
		log.SetLevel(log.DebugLevel)
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
		gin.DisableConsoleColor()
	}

	quit := make(chan struct{})
	defer close(quit)

	db := models.NewTaskMap()

	d := dispatcher.NewDispatcher(
		db,
		vegeta.AttackFn,
	)

	r := reporter.NewReporter(db)

	go d.Run(quit)

	engine := endpoints.SetupRouter(d, r)

	sig := make(chan os.Signal, 1)

	signal.Notify(sig, os.Interrupt)
	go func() {
		for {
			select { //nolint: megacheck
			case <-sig:
				quit <- struct{}{}
			}
			os.Exit(0)
		}
	}()

	// start server
	log.Fatal(engine.Run(fmt.Sprintf("%s:%s", *ip, *port)))
}
