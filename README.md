# Seaport
Serve all your docker containers on port 80, using nginx.

### Overview
Seaport allows you to run many docker services on the same domain, without having to worry about setting up a (reverse proxy)[https://www.nginx.com/resources/admin-guide/reverse-proxy/]. It will scan all your docker containers for exposed tcp ports, and forward all requests to `/my_container_name/` to the corresponding exposed port, with `/my_container_name/` stripped from the path. Your containers will be none the wiser, and you're free to host as many services as you want from one domain.

### Quickstart
Make sure you have go, nginx, and docker installed, then run the following commands:
```
	go get github.com/drhurd/seaport # install seaport
```