package main

import (
	"bytes"
	"io"
	"log"
	"net"
	"strconv"
	"strings"
	"fmt"
)

type Seaport struct {
	port int // listening port
	routes map[string]int // map of routes -> ports
}

func NewSeaport(port int, routes map[string]int) *Seaport {
	s := Seaport{port, routes}
	return &s
}

func (s *Seaport) Listen() {
	port_str := ":" + strconv.Itoa(s.port)
	log.Print("port: ", port_str)
	listener, err := net.Listen("tcp", port_str)

	if err != nil {
		log.Fatal("Unable to listen: ", err)
	}

	for {
		// Wait for a connection.
		conn, err := listener.Accept()
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println("accepting connection")

		// Handle the request
		go forward(conn, s.routes)
	}
}

func forward(in net.Conn, routes map[string]int) {	
	buf := bytes.NewBuffer(make([]byte, 1024)) // connection data buffer

	// read in from the connection
	_, err := in.Read(buf.Bytes())
	if err != nil {
		// exit once we are done reading from the connection
		log.Fatal(err)
	}

	// Search for a match the first time through
	// TODO: wait for minimum buffer size
	key := routeMatch(buf, routes)
	if key != "" {
		findAndRemove(buf, key, 1)

		addr := ":" + strconv.Itoa(routes[key]) // port number -> address
		out, err := net.Dial("tcp", addr)
		if err != nil {
			log.Fatal(err)
		}

		io.Copy(out, buf)
		
		pipe(in, out)
	} else {
		// TODO: return 404 error message
		in.Close()
	}
}

func pipe(conn1, conn2 net.Conn) {
	go func() {
		defer conn1.Close()
		defer conn2.Close()
		io.Copy(conn1, conn2)
	}()

	go func() {
		defer conn1.Close()
		defer conn2.Close()
		io.Copy(conn2, conn1)
	}()
}

func routeMatch(buf *bytes.Buffer, routes map[string]int) string {
	// Extract the key
	// Expects: "[Method] /key/rest_of_path [Protocol]"
	data_str := buf.String()
	key := strings.Split(data_str, " ")[1]
	key = strings.Split(key, "/")[1]

	if routes[key] != 0 {
		return key
	} else {
		return ""
	}
}

func findAndRemove(buf *bytes.Buffer, pattern string, n int) {
	data_str := buf.String()
	data_str = strings.Replace(data_str, pattern, "", n)
	data_str = strings.Replace(data_str, "//", "/", 1)

	buf.Truncate(0)
	buf.Write([]byte(data_str))
}



