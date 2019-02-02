package main

import (
	"fmt"
	"vegeta-server/internal/app/server"

	log "github.com/sirupsen/logrus"

	"github.com/gin-gonic/gin"
	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

var (
	verbose = kingpin.Flag("verbose", "Verbose mode.").Short('v').Bool()
	ip      = kingpin.Arg("ip", "Server IP Address.").Default("localhost").String()
	port    = kingpin.Arg("port", "Server Port.").Default("8000").String()
)

func main() {
	kingpin.Parse()

	router := gin.Default()

	// api/v1 router group
	v1 := router.Group("/api/v1")
	{
		v1.POST("/attack", server.PostAttackEndpoint)
		v1.GET("/attack", server.GetAttackEndpoint)
		v1.GET("/attack/:id", server.GetAttackByIDEndpoint)
	}

	// start server
	log.Fatal(router.Run(fmt.Sprintf("%s:%s", *ip, *port)))
}
