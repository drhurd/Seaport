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
	config_text = 
`error_log  log/nginx.error;
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

type Config struct {
	ServerName string // will default to localhost
	DefaultPort int // where to forward default requests (e.g. GET /)
}

/* Writes a valid nginx config file to f, with the proxies indicated by routes and any other configuration specified in config
*/
func WriteConfigFile(f *os.File, routes map[string]int, config Config) {
	var routes_buffer bytes.Buffer
	for path, port := range routes {
		location_str := "\t\tlocation /%s {\n\t\t\treturn 302 /%s/;\n\t\t}\n\n\t\tlocation /%s/ {\n\t\t\tproxy_pass http://127.0.0.1:%d/;\n\t\t}\n"

		fmt.Fprintf(&routes_buffer, location_str, path, path, path, port)
	}

	server_name := config.ServerName
	if server_name == "" {
		server_name = "localhost"
	}

	var root_rule string
	default_port := config.DefaultPort
	if default_port == 0 {
		root_rule = "return 404;"
	} else {
		root_rule = "proxy_pass http://127.0.0.1:" + strconv.Itoa(default_port) + ";"
	}

	fmt.Fprintf(f, config_text, server_name, routes_buffer.String(), root_rule)

	err := f.Sync()
	if err != nil {
		log.Fatal("Error saving config file: ", err)
	}
}


/* Runs the command to start nginx
*/
func StartNginx(path string) error {
	cmd := exec.Command("nginx", "-p", "/var/", "-c", path)
	log.WithFields(log.Fields{
		"args" : cmd.Args,
	}).Debug("Nginx Command created")

	var out bytes.Buffer
	cmd.Stdout = &out

	var err_out bytes.Buffer
	cmd.Stderr = &err_out

	err := cmd.Run()
	if err != nil {
		log.Error("Nginx: ", err_out.String())
	}
	return err
}

func StopNginx() error {
	cmd := exec.Command("nginx", "-p", "/var/", "-s", "stop")

	return cmd.Run()
}

/* Returns bool indicating if Nginx is currently running
*/
func NginxStatus() bool {
	cmd := exec.Command("service", "nginx", "status")

	var out bytes.Buffer
	cmd.Stdout = &out

	var err_out bytes.Buffer
	cmd.Stderr = &err_out

	cmd.Run()
	// Service will return exit code 3 if nginx is not running.
	// Not sure how to handle this right now.
/*	if err != nil && !exec.ExitError(err).ProcessState.Success() {
		log.Warn("Status error: ", err_out.String())
		log.Warn(out.String())
		log.Fatal("Error getting status: ", err)
	}
*/

	status_msg := out.String()

	return strings.Index(status_msg, "not") == -1
}