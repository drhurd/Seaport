package nginx

import (
	log "github.com/Sirupsen/logrus"

	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

const (
	configText = `error_log  log/nginx.error;
pid        run/nginx.pid;

events {
	worker_connections  1024;  ## Default: 1024
}

http {
	server {
		listen 80;
		server_name %s;

%s

		location / {
			%s
		}
	}
}
`
)

// Config encapsulates configuration options for the nginx.conf file
type Config struct {
	ServerName  string // will default to localhost
	DefaultPort int    // where to forward default requests (e.g. GET /)
}

// WriteConfigFile wirtes a valid nginx config to f
// Proxies are indicated in routes
// Other configuration specified in config
func WriteConfigFile(f *os.File, routes map[string]int, config Config) {
	var routesBuffer bytes.Buffer
	for path, port := range routes {
		locationStr := "\t\tlocation /%s {\n\t\t\treturn 302 /%s/;\n\t\t}\n\n\t\tlocation /%s/ {\n\t\t\tproxy_pass http://127.0.0.1:%d/;\n\t\t}\n"

		fmt.Fprintf(&routesBuffer, locationStr, path, path, path, port)
	}

	serverName := config.ServerName
	if serverName == "" {
		serverName = "localhost"
	}

	var rootRule string
	defaultPort := config.DefaultPort
	if defaultPort == 0 {
		rootRule = "return 404;"
	} else {
		rootRule = "proxy_pass http://127.0.0.1:" + strconv.Itoa(defaultPort) + ";"
	}

	fmt.Fprintf(f, configText, serverName, routesBuffer.String(), rootRule)

	err := f.Sync()
	if err != nil {
		log.Fatal("Error saving config file: ", err)
	}
}

// StartNginx starts nginx
func StartNginx(path string) error {
	cmd := exec.Command("nginx", "-p", "/var/", "-c", path)
	log.WithFields(log.Fields{
		"args": cmd.Args,
	}).Debug("Nginx Command created")

	var out bytes.Buffer
	cmd.Stdout = &out

	var errOut bytes.Buffer
	cmd.Stderr = &errOut

	err := cmd.Run()
	if err != nil {
		log.Error("Nginx: ", errOut.String())
	}
	return err
}

// StopNginx stops nginx
func StopNginx() error {
	cmd := exec.Command("nginx", "-p", "/var/", "-s", "stop")

	return cmd.Run()
}

// Status returns bool indicating if Nginx is currently running
func Status() bool {
	cmd := exec.Command("service", "nginx", "status")

	var out bytes.Buffer
	cmd.Stdout = &out

	var errOut bytes.Buffer
	cmd.Stderr = &errOut

	cmd.Run()
	// Service will return exit code 3 if nginx is not running.
	// Not sure how to handle this right now.
	/*	if err != nil && !exec.ExitError(err).ProcessState.Success() {
			log.Warn("Status error: ", errOut.String())
			log.Warn(out.String())
			log.Fatal("Error getting status: ", err)
		}
	*/

	return strings.Index(out.String(), "not") == -1
}
