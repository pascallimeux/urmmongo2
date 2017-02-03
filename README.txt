
Installation on host
	install mongodb on server,
	update configuration (ip address for server in ..src/github.com/pascallimeux/urmmongo/server/config.json)
	update script (user name in ..src/github.com/pascallimeux/urmmongo/scripts/start_urmmongo.sh)
	https_proxy=proxy:8080 go get gopkg.in/mgo.v2
	https_proxy=proxy:8080 go get github.com/gorilla/mux
	cd ..src/github.com/pascallimeux/urmmongo/server
	./build_urmmongo.sh
	cd ../dist
	./start_urmmongo.sh
	./stop_urmmongo.sh

Installation in docker container
	update Dockerfile (proxy in ..src/github.com/pascallimeux/urmmongo/Dockerfile)
	cd ..src/github.com/pascallimeux/urmmongo/server
	./build_urmmongo.sh
	cd ..
	sudo docker-compose build
	cd ../dist
	./start_docker.sh
	./stop_docker.sh


