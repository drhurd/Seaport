package main

import (
	"github.com/fsouza/go-dockerclient"
	"fmt"
	"log"
	"strings"
)

func main() {

	endpoint := "unix:///var/run/docker.sock"
	client, err := docker.NewClient(endpoint)

	if err != nil {
		log.Fatal("Unable to connect to docker client:", err)
	}

	routes := make(map[string]int)

	// constructing the list of containers
	containers, _ := client.ListContainers(docker.ListContainersOptions {All: false})
	for _, container := range containers {
		inspected_container, err := client.InspectContainer(container.ID)
		if err != nil {
			log.Fatal(err)
		}

		if len(container.Ports) == 0 {
			continue
		}
		port := int(container.Ports[0].PublicPort)
		fmt.Println("Port: ", port)

		name := inspected_container.Name
		name = strings.Replace(name, "/", "", 1)
		fmt.Println("Name: ", name)

		routes[name] = port
	}

	s := NewSeaport(80, routes)

	s.Listen()
}
