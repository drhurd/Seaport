package main

import (
	// "github.com/fsouza/go-dockerclient"
)

func main() {
/*
	endpoint := "unix:///var/run/docker.sock"
	client, err := docker.NewClient(endpoint)

	if err != nil {
		log.Fatal("Unable to connect to docker client:", err)
	}

	imgs, _ := client.ListImages(docker.ListImagesOptions{All: false})
	for _, img := range imgs {
		fmt.Println("ID: ", img.ID)
		//		fmt.Println("P: ", img.
		fmt.Println("RepoTags: ", img.RepoTags)
		fmt.Println("Created: ", img.Created)
		fmt.Println("Size: ", img.Size)
		fmt.Println("VirtualSize: ", img.VirtualSize)
		fmt.Println("ParentId: ", img.ParentID)
	}
*/

	routes := map[string]int {
		"test" : 8080,
	}

	s := NewSeaport(3000, routes)

	s.Listen()
}
