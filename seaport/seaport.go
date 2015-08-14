package seaport

import (
	"bytes"
	log "github.com/Sirupsen/logrus"
	"io"
	"net"
	"strconv"
	"strings"
)

// Seaport represents the core functionality of the port forwarding
// It forwards the paths specified by the routes keys to their corresponding port values
type Seaport struct {
	Routes   map[string]int // map of routes -> ports
	Listener net.Listener
}

// Listen starts the server
func (s Seaport) Listen() {	
	for {
		// Wait for a connection.
		conn, err := s.Listener.Accept()
		if err != nil {
			log.Warn(err)
			break
		}

		// Handle the request
		go forward(conn, s.Routes)
	}
}

// Close closes the listener connection, stopping seaport
func (s Seaport) Close() {
	s.Listener.Close()
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
			"name": key,
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
	str := buf.String()
	key := strings.Split(str, " ")[1]
	key = strings.Split(key, "/")[1]

	if _, ok := routes[key]; ok {
		return key, true
	} 
	
	return "", false
}

func findAndRemove(buf *bytes.Buffer, pattern string, n int) {
	str := buf.String()
	str = strings.Replace(str, pattern, "", n)
	str = strings.Replace(str, "//", "/", 1)

	buf.Truncate(0)
	buf.Write([]byte(str))
}
