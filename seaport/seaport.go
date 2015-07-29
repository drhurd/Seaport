package seaport

import (
	"bytes"
	"io"
	log "github.com/Sirupsen/logrus"
	"net"
	"strconv"
	"strings"
)

type Seaport struct {
	routes map[string]int // map of routes -> ports
	listener net.Listener
}

func NewSeaport(routes map[string]int) *Seaport {
	s := Seaport{routes, nil}
	return &s
}

func (s *Seaport) Listen(port int) {
	port_str := ":" + strconv.Itoa(port)
	listener, err := net.Listen("tcp", port_str)
	if err != nil {
		log.Fatal("Unable to listen: ", err)
	}

	s.listener = listener

	for {
		// Wait for a connection.
		conn, err := listener.Accept()
		if err != nil {
			log.Warn(err)
			break
		}

		// Handle the request
		go forward(conn, s.routes)
	}
}

func (s *Seaport) Close() {
	s.listener.Close()
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
	if key, ok := routeMatch(buf, routes); ok {
		log.WithFields(log.Fields{
			"name" : key,
		}).Debug("Forwarding to container")

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

func routeMatch(buf *bytes.Buffer, routes map[string]int) (string, bool) {
	// Extract the key
	// Expects: "[Method] /key/rest/of/path [Protocol]"
	data_str := buf.String()
	key := strings.Split(data_str, " ")[1]
	key = strings.Split(key, "/")[1]

	if _, ok := routes[key]; ok {
		return key, true
	} else {
		return "", false
	}
}

func findAndRemove(buf *bytes.Buffer, pattern string, n int) {
	data_str := buf.String()
	data_str = strings.Replace(data_str, pattern, "", n)
	data_str = strings.Replace(data_str, "//", "/", 1)

	buf.Truncate(0)
	buf.Write([]byte(data_str))
}



