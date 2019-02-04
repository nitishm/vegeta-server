package main

import (
	"fmt"
	"os"
	"os/signal"
	"runtime"
	"vegeta-server/internal/app/attacker"
	"vegeta-server/internal/app/server/endpoints"

	log "github.com/sirupsen/logrus"

	"gopkg.in/alecthomas/kingpin.v2"
)

var (
	commit  = "N/A"
	date    = "N/A"
	version = "N/A"
)

var (
	ip   = kingpin.Flag("ip", "Server IP Address.").Default("localhost").String()
	port = kingpin.Flag("port", "Server Port.").Default("8000").String()
	v    = kingpin.Flag("version", "Version Info").Short('v').Bool()
)

func main() {
	kingpin.Parse()
	log.Infof("PID - %v", os.Getpid())

	if *v {
		// Set at linking time
		fmt.Println("=======")
		fmt.Println("VERSION")
		fmt.Println("=======")
		fmt.Println("Version\t", version)
		fmt.Println("Commit \t", commit)
		fmt.Println("Runtime\t", fmt.Sprintf("%s %s/%s", runtime.Version(), runtime.GOOS, runtime.GOARCH))
		fmt.Println("Date   \t", date)

		os.Exit(0)
		return
	}

	sig := make(chan os.Signal, 1)
	quit := make(chan struct{})

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

	attacker := attacker.NewAttacker(
		attacker.NewScheduler(
			attacker.NewDispatcher(
				attacker.DefaultAttackFn,
			),
			quit,
		),
	)

	engine := endpoints.SetupRouter(attacker)

	// start server
	log.Fatal(engine.Run(fmt.Sprintf("%s:%s", *ip, *port)))
}
