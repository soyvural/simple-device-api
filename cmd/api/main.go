package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/soyvural/simple-device-api/service"

	"github.com/gin-gonic/gin"

	_ "github.com/soyvural/simple-device-api/docs"
)

var (
	port = flag.Uint("port", 8080, "http port to serve device API")
)

func main() {
	flag.Parse()
	g := gin.Default()
	svc := service.New(g)
	svc.SetRoute_v1()
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", *port),
		WriteTimeout: time.Second * 10,
		ReadTimeout:  time.Second * 5,
		IdleTimeout:  time.Second * 30,
		Handler:      g,
	}
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("Error occurred while the server is getting started, err: %v\n", err)
	}
}
