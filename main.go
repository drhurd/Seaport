package main

import (
	log "github.com/Sirupsen/logrus"
	"gopkg.in/alecthomas/kingpin.v2"

	"os"

	"github.com/drhurd/seaport/docker"
	"github.com/drhurd/seaport/nginx"
	"github.com/drhurd/seaport/seaport"
)

const (
	logLevel = log.DebugLevel // set logging level here
)

var (
	server_port = kingpin.Flag("port", "Specify the port seaport should use. Defaults to 80.").Default("80").Short('p').Int()
	seaport_mode = kingpin.Flag("forward", "Use Seaport instead of Nginx to forward to docker containers. Currently under development, not recommended.").Short('f').Bool()
	stop = kingpin.Flag("stop", "Stop the Nginx server").Short('s').Bool()
	nginx_file = kingpin.Flag("config", "Nginx config file, Default = /etc/nginx/nginx.conf").Short('c').Default("/etc/nginx/nginx.conf").String()
)

func main() {
	// Configure the logger
	log.SetLevel(logLevel)

	// Parse the input
	kingpin.Parse()

	if *seaport_mode {
		runSeaportServer(*server_port)
	} else if *stop {
		nginx.StopNginx()
	} else {
		runNginxServer()
	}
}

func runSeaportServer(port int) {
	routes := makeRouteMap(docker.ListContainers())
	s := seaport.NewSeaport(routes)
	s.Listen(port)
	return
}

func runNginxServer() {
	routes := makeRouteMap(docker.ListContainers())

	// Use nginx to forward
	file, err := os.Create(*nginx_file)
	if err != nil {
		log.Fatal("Couldn't open nginx file: ", err)
	}

	nginx.WriteConfigFile(file, routes, nginx.Config{ServerName:"localhost"})

	log.Debug("file written")

	if nginx.NginxStatus() {
		err = nginx.StopNginx()
		if err != nil {
			log.Fatal("Couldn't stop nginx: ", err)
		}
	}

	err = nginx.StartNginx(*nginx_file)
	if err != nil {
		log.Fatal("Couldn't start nginx: ", err)
	}
}

func makeRouteMap(containers []docker.Container) map[string]int {
	routes := make(map[string]int)
	for _, container := range containers {
		routes[container.Name] = int(container.Port)
	}
	return routes
}