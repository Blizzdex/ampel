default: docker

docker:
	docker-compose up --build

protoc:
	# get the generated protoc files from the docker environment
	docker build --target proto -t ampel2:proto .
	# docker run  test-project
	docker run -v $(shell pwd):/tmp ampel2:proto cp -r servis /tmp

