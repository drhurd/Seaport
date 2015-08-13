package server

import (
	log "github.com/Sirupsen/logrus"
	"github.com/drhurd/seaport/seaport"
	"net/http"
	"strconv"
)

var (
	seaport seaport.Seaport
)

// StartDaemonServer does nothing right now
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
			"name": container.Name,
			"port": container.Port,
		}).Debug("Container added")
	}

	s := seaport.NewSeaport(routes)

	if portStr, ok := r.Form["port"]; ok {
		port, err := strconv.Atoi(portStr)

		if err != nil {

		}
	} else {
		go s.Listen(80)
	}

}
