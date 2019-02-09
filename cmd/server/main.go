package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"os"
	"os/signal"
	"runtime"
	"vegeta-server/internal/dispatcher"
	"vegeta-server/internal/endpoints"
	"vegeta-server/pkg/vegeta"

	log "github.com/sirupsen/logrus"

	"gopkg.in/alecthomas/kingpin.v2"
)

var (
	commit  = "N/A"
	date    = "N/A"
	version = "N/A"
)

var (
	ip    = kingpin.Flag("ip", "Server IP Address.").Default("localhost").String()
	port  = kingpin.Flag("port", "Server Port.").Default("8000").String()
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

	d := dispatcher.NewDispatcher(
		vegeta.VegetaAttackFn,
	)

	go d.Run(quit)

	engine := endpoints.SetupRouter(d)

	sig := make(chan os.Signal, 1)

	signal.Notify(sig, os.Interrupt)
	go func() {
		for {
			select {
			case <-sig:
				quit <- struct{}{}
			}
			os.Exit(0)
		}
	}()

	// start server
	log.Fatal(engine.Run(fmt.Sprintf("%s:%s", *ip, *port)))
}
