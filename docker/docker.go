package docker

import (
	log "github.com/Sirupsen/logrus"
	dockerClient "github.com/fsouza/go-dockerclient"
	"strings"
)

const (
	dockerEndpoint = "unix:///var/run/docker.sock"
)

//Container is used pass around only the relevant docker container information
type Container struct {
	Name string
	Port int64
}

// ListContainers returns a list of containers running with exposed TCP ports
func ListContainers() []Container {
	// Connect to docker
	client, err := dockerClient.NewClient(dockerEndpoint)
	if err != nil {
		log.Fatal(err)
	}

	containers, err := client.ListContainers(dockerClient.ListContainersOptions{All: false})
	if err != nil {
		log.Fatal(err)
	}

	var validContainers []Container
	for _, container := range containers {
		// check for a valid port to map
		var portNumber int64
		validPortFound := false
		for _, port := range container.Ports {
			if port.Type == "tcp" {
				portNumber = port.PublicPort
				validPortFound = true
				break
			}
		}
		// If no valid port, skip this container
		if !validPortFound {
			continue
		}

		// Check that container has at least one name
		// TODO: do I even need to check for container name existence?
		if len(container.Names) == 0 {
			continue
		}

		// Strip the initial slash to set our name
		name := strings.Replace(container.Names[0], "/", "", 1)

		validContainers = append(validContainers, Container{name, portNumber})
	}

	return validContainers
}
