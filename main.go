package main

import (
	log "github.com/Sirupsen/logrus"
	"gopkg.in/alecthomas/kingpin.v2"

	"os"
	"time"

	"github.com/drhurd/seaport/docker"
	"github.com/drhurd/seaport/nginx"
	"github.com/drhurd/seaport/seaport"
)

const (
	logLevel = log.DebugLevel // set logging level here
)

var (
	server_port = kingpin.Flag("port", "Specify the port seaport should use. Defaults to 80.").Default("80").Short('p').Int()
	nginx_mode = kingpin.Flag("nginx", "Enable nginx forwarding").Short('n').Bool()
	nginx_file = kingpin.Flag("config", "Nginx config file, Default = /etc/nginx/nginx.conf").Short('c').Default("/etc/nginx/nginx.conf").String()
)

func main() {
	// Configure the logger
	log.SetLevel(logLevel)

	// Parse the input
	kingpin.Parse()

	containers := docker.ListContainers()

	routes := make(map[string]int)
	for _, container := range containers {
		routes[container.Name] = int(container.Port)
	}

	if *nginx_mode {
		file, err := os.Open(*nginx_file)
		if err != nil {
			log.Fatal("Couldn't open nginx file: ", err)
		}

		nginx.WriteConfigFile(file, routes, nginx.Config{ServerName:"localhost"})

		if nginx.NginxStatus() {
			err = nginx.StopNginx()
			if err != nil {
				log.Fatal("Couldn't stop nginx: ", err)
			}
		}

		err = nginx.StartNginx(*nginx_file)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		s := seaport.NewSeaport(routes)

		go func() {
			time.Sleep(10000 * time.Millisecond)
			log.Debug("Closing")
			s.Close()
		}()

		log.Debug("Server port flag: ", *server_port)
		s.Listen(80)
	}
}
