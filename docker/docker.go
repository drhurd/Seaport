package docker

import (
	docker_client "github.com/fsouza/go-dockerclient"
	log "github.com/Sirupsen/logrus"
	"strings"
)

const (
	docker_endpoint = "unix:///var/run/docker.sock"
)

type Container struct {
	Name string
	Port int64
}

func ListContainers() []Container {
	// Connect to docker
	client, err := docker_client.NewClient(docker_endpoint)
	if err != nil {
		log.Fatal(err) 
	}

	containers, err := client.ListContainers(docker_client.ListContainersOptions {All: false})
	if err != nil {
		log.Fatal(err)
	}

	valid_containers := make([]Container, 0)
	for _, container := range containers {
		// check for a valid port to map
		var port_number int64
		valid_port_found := false
		for _, port := range container.Ports {
			if port.Type == "tcp" {
				port_number = port.PublicPort
				valid_port_found = true
				break
			}
		}
		// If no valid port, skip this container
		if !valid_port_found {
			continue
		}

		// Check that container has at least one name
		// TODO: do I even need to check for container name existence?
		if len(container.Names) == 0 {
			continue
		}

		// Strip the initial slash to set our name
		name := strings.Replace(container.Names[0], "/", "", 1)

		valid_containers = append(valid_containers, Container{name, port_number})
	}

	return valid_containers
}