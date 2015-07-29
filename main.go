package main

import (
	log "github.com/Sirupsen/logrus"
	"gopkg.in/alecthomas/kingpin.v2"
	"time"

	"github.com/drhurd/seaport/docker"
	"github.com/drhurd/seaport/seaport"
)

const (
	logLevel = log.DebugLevel // set logging level here
)

var (
	server_port = kingpin.Flag("port", "Specify the port seaport should use. Defaults to 80.").Default("80").Short('p').Int()
)

func main() {
	// Configure the logger
	log.SetLevel(logLevel)

	go func() {
		time.Sleep(10000 * time.Millisecond)
		log.Debug("Closing")
		s.Close()
	}()

	log.Debug("Server port flag: ", *server_port)

	s.Listen(80)
}
