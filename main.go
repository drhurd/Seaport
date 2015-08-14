package main

import (
	log "github.com/Sirupsen/logrus"
	"gopkg.in/alecthomas/kingpin.v2"

	"bytes"
	"net"
	"os"
	"strconv"
	"strings"

	"github.com/drhurd/seaport/docker"
	"github.com/drhurd/seaport/nginx"
	"github.com/drhurd/seaport/seaport"
)

const (
	logLevel = log.DebugLevel // set logging level here
)

var (
	serverPort  = kingpin.Flag("port", "Specify the port seaport should use. Defaults to 80.").Default("80").Short('p').Int()
	seaportMode = kingpin.Flag("forward", "Use Seaport instead of Nginx to forward to docker containers. Currently under development, not recommended.").Short('f').Bool()
	stop        = kingpin.Flag("stop", "Stop the Nginx server").Short('s').Bool()
	nginxFile   = kingpin.Flag("config", "Nginx config file, Default = /etc/nginx/nginx.conf").Short('c').Default("/etc/nginx/nginx.conf").String()
)

func main() {
	// Configure the logger
	log.SetLevel(logLevel)

	// Parse the input
	kingpin.Parse()

	if *seaportMode {
		runSeaportServer(*serverPort)
	} else if *stop {
		cmd := nginx.StopCommand()

		err := cmd.Run()
		if err != nil {
			log.Fatal(err)
		}
	} else {
		runNginxServer()
	}
}

func runSeaportServer(port int) {
	routes := makeRouteMap(docker.ListContainers())
	listener, err := net.Listen("tcp", ":" + strconv.Itoa(port))
	if err != nil {
		log.Fatal(err)
	}


	s := seaport.Seaport{routes, listener}

	s.Listen()
}

func runNginxServer() {
	routes := makeRouteMap(docker.ListContainers())

	// Use nginx to forward
	file, err := os.Create(*nginxFile)
	if err != nil {
		log.Fatal("Couldn't open nginx file: ", err)
	}

	nginx.WriteConfigFile(file, routes, nginx.Config{ServerName: "localhost"})

	if getNginxStatus() {
		cmd := nginx.StopCommand()

		err := cmd.Run()
		if err != nil {
			log.Fatal("Couldn't stop nginx: ", err)
		}
	}

	cmd := nginx.StartCommand(*nginxFile)
	err = cmd.Run()
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

func getNginxStatus() bool {
	cmd := nginx.StatusCommand()

	data, _ := cmd.Output() // exit codes are frustrating, will handle errors later

	output := bytes.NewBuffer(data).String()

	return strings.Index(output, "not") == -1
}