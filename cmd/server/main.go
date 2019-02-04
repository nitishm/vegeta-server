package main

import (
	"fmt"
	"os"
	"os/signal"
	"vegeta-server/internal/app/attacker"
	"vegeta-server/internal/app/server/endpoints"

	log "github.com/sirupsen/logrus"

	"gopkg.in/alecthomas/kingpin.v2"
)

var (
	ip   = kingpin.Arg("ip", "Server IP Address.").Default("localhost").String()
	port = kingpin.Arg("port", "Server Port.").Default("8000").String()
)

func main() {
	kingpin.Parse()

	sig := make(chan os.Signal, 1)
	quit := make(chan struct{})

	signal.Notify(sig, os.Interrupt)
	go func() {
		for {
			select {
			case <-sig:
				quit <- struct{}{}
			}
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
