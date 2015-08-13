package nginx

import (
	//log "github.com/Sirupsen/logrus"

	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

const (
	config_text = 
`error_log  logs/error.log;
pid        logs/nginx.pid;

events {
	worker_connections  1024;  ## Default: 1024
}

http {
	server {
		listen 80;
		server_names_hash_bucket_size 128;
		server_name %s;

%s

		location / {
			return 404;
		}
	}
}
`
)

type Config struct {
	ServerName string // will default to localhost
}

/* Writes a valid nginx config file to f, with the proxies indicated by routes and any other configuration specified in config
*/
func WriteConfigFile(f *os.File, routes map[string]int, config Config) {
	var locations bytes.Buffer
	for path, port := range routes {
		location_str := "\t\tlocation /" + path + " {\n\t\t\treturn 302 /" + path + "/;\n\t\t}\n\n\t\tlocation /" + path + "/ {\n\t\t\tproxy_pass 127.0.0.1:" + strconv.Itoa(port) + ";\n\t\t}\n"

		locations.WriteString(location_str)
	}

	server_name := config.ServerName
	if server_name == "" {
		server_name = "localhost"
	}

	fmt.Fprintf(f, config_text, server_name, locations)
}


/* Runs the command to start nginx
*/
func StartNginx(path string) error {
	cmd := exec.Command("nginx", "-c", path)
	
	var out bytes.Buffer
	cmd.Stdout = &out

	var err_out bytes.Buffer
	cmd.Stderr = &err_out

	return cmd.Run()
}

func StopNginx() error {
	cmd := exec.Command("nginx", "-s", "stop")

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