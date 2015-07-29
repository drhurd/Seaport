package server

import (
	"net/http"
	log "github.com/Sirupsen/logrus"
	"github.com/drhurd/seaport/seaport"
	"strconv"
)

var (
	seaport seaport.Seaport
)

func StartDaemonServer() {
	http.Handle("/listen", listenHandler)
}

func listenHandler(w http.ResponseWriter, r *http.Request) {
	// Construct the list of routes
	routes := make(map[string]int)
	containers := docker.ListContainers()
	
	for _, container := range containers {
		routes[container.Name] = int(container.Port)
		
		log.WithFields(log.Fields{
			"name" : container.Name,
			"port" : container.Port,
		}).Debug("Container added")
	}

	s := seaport.NewSeaport(routes)

	if port_str, ok := r.Form["port"]; ok {
		port, err := strconv.Atoi(port_str)
		
		if err != nil {
			
		}
	} else {
		go s.Listen(80)
	}

}