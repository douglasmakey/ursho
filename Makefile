# Usage:

release: clean build

build:
	docker build -t ursho-builder .
	docker run ursho-builder | docker build -t ursho -

# remove previous images and containers
clean:
	docker rm -f ursho-builder 2> /dev/null || true
	docker rmi -f ursho-builder || true
	docker rm -f ursho 2> /dev/null || true
	docker rmi -f ursho || true
