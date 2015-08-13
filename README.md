# Seaport
Serve all your docker containers on port, using nginx.

### Overview
Seaport allows you to run many docker services on the same domain, without having to worry about setting up a (reverse proxy)[https://www.nginx.com/resources/admin-guide/reverse-proxy/]. It will scan all your docker containers for exposed tcp ports, and forward all requests to `/my_container_name/` to the corresponding exposed port, with `/my_container_name/` stripped from the path. Your containers will be none the wiser, and you're free to host as many services as you want from one domain.

### Quickstart
Make sure you have go, nginx, and docker installed, then run the following commands:
```
	go get github.com/drhurd/seaport # install seaport
	go install github.com/drhurd/seaport # compile seaport
	
	docker run -d -p 80 tutum/hello-world # run an example container in the background

	sudo seaport # run seaport -> configure and start nginx

	curl localhost:80/your_container_name/ # make the request!

	sudo seaport -s # stop nginx
```

### Development
Testing is on its way.