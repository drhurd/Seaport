package nginx

import (
	log "github.com/Sirupsen/logrus"

	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
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

// StartNginx returns a command to start nginx
func StartCommand(path string) *exec.Cmd {
	return exec.Command("nginx", "-p", "/var/", "-c", path)
}

// StopNginx returns a command to stop nginx
func StopCommand() *exec.Cmd {
	return exec.Command("nginx", "-p", "/var/", "-s", "stop")
}

// Status returns bool indicating if Nginx is currently running
func StatusCommand() *exec.Cmd {
	return exec.Command("service", "nginx", "status")
}
